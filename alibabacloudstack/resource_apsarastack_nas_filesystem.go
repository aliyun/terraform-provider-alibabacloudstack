package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackNasFileSystem() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{"Capacity", "Performance", "standard", "advance"}, false),
			},
			"protocol_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"NFS", "SMB"}, false),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(2, 256),
			},
			"encrypt_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1, 2}),
				Default:      0,
			},
			"file_system_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"extreme", "standard"}, false),
				Default:      "standard",
			},
			"capacity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasFileSystemCreate,
		resourceAlibabacloudStackNasFileSystemRead, resourceAlibabacloudStackNasFileSystemUpdate, resourceAlibabacloudStackNasFileSystemDelete)
	return resource
}

func resourceAlibabacloudStackNasFileSystemCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateFileSystem"
	request := make(map[string]interface{})
	request["ProtocolType"] = d.Get("protocol_type")
	if v, ok := d.GetOk("file_system_type"); ok {
		request["FileSystemType"] = v
	}
	request["StorageType"] = d.Get("storage_type")
	request["ClusterId"] = d.Get("cluster_id")
	if v, ok := d.GetOk("encrypt_type"); ok {
		request["EncryptType"] = v
	}
	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}
	if v, ok := d.GetOk("capacity"); ok {
		request["Capacity"] = v
	}
	if v, ok := d.GetOk("kms_key_id"); ok {
		request["KmsKeyId"] = v
	}

	response, err = client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)

	if err != nil {
		return err
	}

	d.SetId(fmt.Sprint(response["FileSystemId"]))
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutRead), 3*time.Second, nasService.DescribeNasFileSystemStateRefreshFunc(d.Id(), "Pending", []string{"Stopped", "Stopping", "Deleting"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	//}
	return nil
}

func resourceAlibabacloudStackNasFileSystemUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"RegionId":     client.RegionId,
		"FileSystemId": d.Id(),
	}
	if d.HasChanges("description", "capacity") {
		request["Description"] = d.Get("description")
		request["Capacity"] = d.Get("capacity")
		request["newCapacity"] = d.Get("capacity")
		action := "ModifyFileSystem"
		_, err = client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackNasFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}
	object, err := nasService.DescribeNasFileSystem(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alicloud_nas_file_system nasService.DescribeNasFileSystem Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("description", object["Description"])
	d.Set("protocol_type", object["ProtocolType"])
	d.Set("storage_type", object["StorageType"])
	d.Set("encrypt_type", object["EncryptType"])
	d.Set("file_system_type", object["FileSystemType"])
	d.Set("capacity", object["Capacity"])
	d.Set("zone_id", object["ZoneId"])
	d.Set("kms_key_id", object["KMSKeyId"])
	if d.Get("cluster_id").(string) == "" {
		d.Set("cluster_id", object["Location"])
	}
	return nil
}

func resourceAlibabacloudStackNasFileSystemDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteFileSystem"
	request := map[string]interface{}{
		"FileSystemId": d.Id(),
	}
	_, err = client.DoTeaRequest("POST", "Nas", "2017-06-26", action, "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidFileSystem.NotFound", "Forbidden.NasNotFound"}) {
			return nil
		}
		return err
	}
	nasService := NasService{client}
	stateConf := BuildStateConf([]string{"Stopped", "Stopping", "Deleting"}, []string{}, d.Timeout(schema.TimeoutRead), 3*time.Second, nasService.DescribeNasFileSystemStateRefreshFunc(d.Id(), "Deleting", []string{"Running", "Pending"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
