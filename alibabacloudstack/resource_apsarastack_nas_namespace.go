package alibabacloudstack

import (
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNasNamespace() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
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
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"protocol_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"encrypt_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_system_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"mount_target_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasNamespaceCreate, resourceAlibabacloudStackNasNamespaceRead, nil, resourceAlibabacloudStackNasNamespaceDelete)
	return resource
}

func resourceAlibabacloudStackNasNamespaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"ZoneId":         d.Get("zone_id").(string),
		"StorageType":    d.Get("storage_type").(string),
		"ProtocolType":   d.Get("protocol_type").(string),
		"ClusterId":      d.Get("cluster_id").(string),
		"Description":    d.Get("description").(string),
		"EncryptType":    d.Get("encrypt_type").(int),
		"FileSystemType": "standard",
	}

	raw, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "CreateNamespace", "", nil, request, nil)
	if err != nil {
		return err
	}

	// Get the NasNamespaceId from response
	nasNamespaceId := raw["NasNamespaceId"].(string)

	// Generate resource ID using the NasNamespaceId
	d.SetId(nasNamespaceId)

	return nil
}

func resourceAlibabacloudStackNasNamespaceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasNamespace(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_namespace nasService.DescribeNasNamespace Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("zone_id", object["ZoneId"])
	d.Set("storage_type", object["StorageType"])
	d.Set("protocol_type", object["ProtocolType"])
	d.Set("description", object["Description"])
	d.Set("create_time", object["CreateTime"])
	d.Set("file_system_type", object["FileSystemType"])
	d.Set("encrypt_type", object["EncryptType"])
	d.Set("mount_target_count", object["MountTargetCount"])

	return nil
}

func resourceAlibabacloudStackNasNamespaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := map[string]interface{}{
		"NasNamespaceId": d.Id(),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "DeleteNamespace", "", nil, reqQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "NamespaceNotFound", "Forbidden.NasNamespaceNotFound") {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteNamespace", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		if errmsgs.IsExpectedErrors(err, "NamespaceNotFound", "Forbidden.NasNamespaceNotFound") {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	return nil
}
