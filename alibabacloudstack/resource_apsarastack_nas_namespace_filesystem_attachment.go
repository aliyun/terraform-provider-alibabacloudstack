package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackNasNamespaceFilesystemAttachment() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"nas_namespace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_system_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mapped_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"storage_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_system_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackNasNamespaceFilesystemAttachmentCreate, resourceAlibabacloudStackNasNamespaceFilesystemAttachmentRead, resourceAlibabacloudStackNasNamespaceFilesystemAttachmentUpdate, resourceAlibabacloudStackNasNamespaceFilesystemAttachmentDelete)
	return resource
}

func resourceAlibabacloudStackNasNamespaceFilesystemAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"NasNamespaceId": d.Get("nas_namespace_id").(string),
		"FileSystemId":   d.Get("file_system_id").(string),
		"MappedPath":     d.Get("mapped_path").(string),
	}

	// Call the API to add file system to namespace
	_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "AddFileSystemToNamespace", "", nil, request, nil)
	if err != nil {
		return err
	}

	id := fmt.Sprintf("%s:%s", d.Get("nas_namespace_id").(string), d.Get("file_system_id").(string))

	d.SetId(id)
	return nil
}

func resourceAlibabacloudStackNasNamespaceFilesystemAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	nasService := NasService{client}

	object, err := nasService.DescribeNasNamespaceFilesystemAttachment(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_nas_namespace_filesystem_attachment nasService.DescribeNasNamespaceFilesystemAttachment Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource id format, expected NasNamespaceId:FileSystemId")
	}

	d.Set("nas_namespace_id", parts[0])
	d.Set("file_system_id", object["FileSystemId"])
	d.Set("mapped_path", object["MappedPath"])
	d.Set("create_time", object["CreateTime"])
	d.Set("storage_type", object["StorageType"])
	d.Set("file_system_type", object["FileSystemType"])

	return nil
}

func resourceAlibabacloudStackNasNamespaceFilesystemAttachmentUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}

	if d.HasChange("mapped_path") {
		request := map[string]interface{}{
			"NasNamespaceId": d.Get("nas_namespace_id").(string),
			"FileSystemId":   d.Get("file_system_id").(string),
			"MappedPath":     d.Get("mapped_path").(string),
		}

		// Call the RenameMappedPathInNamespace API
		_, err := client.DoTeaRequest("POST", "nas", "2017-06-26", "RenameMappedPathInNamespace", "", nil, request, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_nas_namespace_filesystem_attachment", "RenameMappedPathInNamespace", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackNasNamespaceFilesystemAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts := strings.SplitN(d.Id(), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid resource id format, expected NasNamespaceId:FileSystemId")
	}

	reqQuery := map[string]interface{}{
		"NasNamespaceId": parts[0],
		"FileSystemId":   parts[1],
	}

	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "RemoveFileSystemFromNamespace", "", nil, reqQuery, nil)
		if err != nil {
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveFileSystemFromNamespace", errmsgs.AlibabacloudStackSdkGoERROR, err.Error())
			return resource.RetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
