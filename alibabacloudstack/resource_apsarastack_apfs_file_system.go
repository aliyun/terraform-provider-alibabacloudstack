package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackApfsFileSystem() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"protocol_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "EFS",
			},
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"file_system_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "efs",
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"volume_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metered_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"quota_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackApfsFileSystemCreate, resourceAlibabacloudStackApfsFileSystemRead, resourceAlibabacloudStackApfsFileSystemUpdate, resourceAlibabacloudStackApfsFileSystemDelete)
	return resource
}

func resourceAlibabacloudStackApfsFileSystemCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := map[string]interface{}{
		"ProtocolType":   d.Get("protocol_type"),
		"FileSystemType": d.Get("file_system_type"),
		"StorageType":    d.Get("storage_type"),
		"ZoneId":         d.Get("zone_id"),
		"ClusterId":      d.Get("cluster_id"),
		"VolumeSize":     d.Get("volume_size"),
		"Description":    d.Get("description"),
	}

	resp, err := client.DoTeaRequest("POST", "EFS", "2017-06-26", "CreateFileSystem", "", nil, request, nil)
	if err != nil {
		return fmt.Errorf("failed to create APFS file system: %v", err)
	}

	// Extract FileSystemId from response
	fileSystemId, ok := resp["FileSystemId"].(string)
	if !ok || fileSystemId == "" {
		return fmt.Errorf("failed to retrieve FileSystemId after creation")
	}

	// Save the ID temporarily
	d.SetId(fileSystemId)

	efsService := EfsService{client}
	stateConf := BuildStateConf([]string{"Pending"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, efsService.ApfsFileSystemStateRefreshFunc(d.Id(), []string{""}))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for APFS file system to be ready failed: %v", err)
	}

	return nil
}

func resourceAlibabacloudStackApfsFileSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	efsService := EfsService{client}
	fs, err := efsService.DescribeFileSystem(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_apfs_file_system", "DescribeFileSystem", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("file_system_id", fs["FileSystemId"])
	d.Set("description", fs["Description"])
	d.Set("zone_id", fs["ZoneId"])
	d.Set("create_time", fs["CreateTime"])
	d.Set("status", fs["Status"])
	d.Set("metered_size", fs["MeteredSize"])
	d.Set("quota_size", fs["QuotaSize"])
	d.Set("volume_size", fs["VolumeSize"])
	d.Set("storage_type", fs["StorageType"])
	d.Set("protocol_type", fs["ProtocolType"])
	d.Set("file_system_type", fs["FileSystemType"])
	d.Set("cluster_id", fs["Location"])

	return nil
}

func resourceAlibabacloudStackApfsFileSystemUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("volume_size", "description") {
		query := make(map[string]interface{})
		query["FileSystemId"] = d.Id()
		query["Description"] = d.Get("description")
		query["VolumeSize"] = d.Get("volume_size")
		_, err := client.DoTeaRequest("POST", "EFS", "2017-06-26", "ModifyFileSystem", "", nil, query, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
				"alibabacloudstack_apfs_file_system", "ModifyFileSystem", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		efsService := EfsService{client}
		stateConf := BuildStateConf([]string{"MODIFYING"}, []string{"Running"}, d.Timeout(schema.TimeoutDelete), 10*time.Second, efsService.ApfsFileSystemStateRefreshFunc(d.Id(), []string{""}))
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for APFS file system to be ready failed: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackApfsFileSystemDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	efsService := EfsService{client}

	requestQuery := map[string]interface{}{
		"FileSystemId": d.Id(),
	}

	err := resource.Retry(10*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EFS", "2017-06-26", "DeleteFileSystem", "", nil, requestQuery, nil)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"InvalidFileSystemId.NotFound"}) {
				return resource.NonRetryableError(err)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteFileSystem", errmsgs.AlibabacloudStackSdkGoERROR, "")
			return resource.RetryableError(err)
		}
		addDebug("DeleteFileSystem", raw)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"InvalidFileSystemId.NotFound"}) {
			return nil
		}
		return errmsgs.WrapError(err)
	}

	stateConf := BuildStateConf([]string{"Pending", "Running", "Deleting"}, []string{}, d.Timeout(schema.TimeoutDelete), 3*time.Second, efsService.ApfsFileSystemStateRefreshFunc(d.Id(), []string{}))
	_, err = stateConf.WaitForState()
	return errmsgs.WrapError(err)
}
