package alibabacloudstack

import (
	"fmt"
	"log"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksFile() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"file_type": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"file_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksFileCreate, resourceAlibabacloudStackDataWorksFileRead, resourceAlibabacloudStackDataWorksFileUpdate, resourceAlibabacloudStackDataWorksFileDelete)
	return resource
}

func resourceAlibabacloudStackDataWorksFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	var response map[string]interface{}
	action := "CreateFile"
	path, err := dataworksPublicService.GetFolderPath(fmt.Sprintf("%s:%s", d.Get("project_id"), d.Get("folder_id")))
	if err != nil {
		return err
	}
	request := map[string]interface{}{
		"ProjectId":       d.Get("project_id"),
		"FileFolderPath":  path,
		"FileName":        d.Get("file_name"),
		"FileType":        d.Get("file_type"),
		"Content":         d.Get("content"),
		"FileDescription": d.Get("description"),
	}

	response, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("project_id"), response["Data"]))
	return nil
}

func resourceAlibabacloudStackDataWorksFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	object, err := dataworksPublicService.DescribeDataWorksFile(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksFolder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, _ := ParseResourceId(d.Id(), 2)
	d.Set("project_id", parts[0])
	d.Set("file_id", parts[1])
	file := object["File"].(map[string]interface{})
	d.Set("file_name", file["FileName"])
	d.Set("file_type", file["FileType"])
	d.Set("description", file["FileDescription"])
	d.Set("content", file["Content"])
	d.Set("folder_id", file["FileFolderId"])

	return nil
}

func resourceAlibabacloudStackDataWorksFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("folder_id", "file_name", "file_type", "file_id", "description", "content") {
		parts, err := ParseResourceId(d.Id(), 2)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		path, err := dataworksPublicService.GetFolderPath(fmt.Sprintf("%s:%s", d.Get("project_id"), d.Get("folder_id")))
		if err != nil {
			return err
		}
		request := map[string]interface{}{
			"FileId":          parts[1],
			"ProjectId":       parts[0],
			"FileFolderPath":  path,
			"FileName":        d.Get("file_name"),
			"FileType":        d.Get("file_type"),
			"Content":         d.Get("content"),
			"FileDescription": d.Get("description"),
		}

		_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "UpdateFile", "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackDataWorksFileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteFile"
	request := map[string]interface{}{
		"FileId":    parts[1],
		"ProjectId": parts[0],
	}
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
