package alibabacloudstack

import (
	"fmt"
	"log"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackDataWorksFolder() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"business_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"engine_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Data Integration", "Hologres", "Algorithm", "Database", "General", "UserDefined"}, false),
			},
			"folder_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"folder_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackDataWorksFolderCreate, resourceAlibabacloudStackDataWorksFolderRead, resourceAlibabacloudStackDataWorksFolderUpdate, resourceAlibabacloudStackDataWorksFolderDelete)
	return resource
}

func getDataworksFolderPath(d *schema.ResourceData) string {
	return fmt.Sprintf("业务流程/%v/%v/%v", d.Get("business_name"), d.Get("engine_type"), d.Get("folder_path"))
}

func resourceAlibabacloudStackDataWorksFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateFolder"
	request := map[string]interface{}{
		"ProjectId":  d.Get("project_id"),
		"FolderPath": getDataworksFolderPath(d),
	}
	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("project_id"), response["Data"]))
	return nil
}

func resourceAlibabacloudStackDataWorksFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksService{client}
	object, err := dataworksPublicService.DescribeDataWorksFolder(d.Id())
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
	d.Set("folder_id", parts[1])
	parts = strings.SplitN(object["FolderPath"].(string), "/", 4)

	if len(parts) != 4 {
		return errmsgs.WrapError(fmt.Errorf("Invalid Folder Path %s.", object["FolderPath"].(string)))
	}
	d.Set("business_name", parts[1])
	if strings.HasPrefix(parts[2], "folder") {
		d.Set("engine_type", strings.TrimPrefix(parts[2], "folder"))
	} else {
		d.Set("engine_type", parts[2])
	}
	d.Set("folder_path", parts[3])

	return nil
}

func resourceAlibabacloudStackDataWorksFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := noUpdatesAllowedCheck(d, []string{"business_name", "engine_type", "folder_path"}); err != nil {
		return err 
	}
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if d.HasChanges("business_name", "engine_type", "folder_path") {
		request := map[string]interface{}{
			"ProjectId":  parts[0],
			"FolderId":   parts[1],
			"FolderName": getDataworksFolderPath(d),
		}
		_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", "UpdateFolder", "", nil, nil, request)
		if err != nil {
			return err
		}
	}
	return nil
}

func resourceAlibabacloudStackDataWorksFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteFolder"
	request := map[string]interface{}{
		"FolderId":  parts[1],
		"ProjectId": parts[0],
	}
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	return nil
}
