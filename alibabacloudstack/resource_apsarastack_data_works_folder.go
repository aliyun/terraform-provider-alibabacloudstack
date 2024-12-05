package alibabacloudstack

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDataWorksFolderCreate,
		Read:   resourceAlibabacloudStackDataWorksFolderRead,
		Update: resourceAlibabacloudStackDataWorksFolderUpdate,
		Delete: resourceAlibabacloudStackDataWorksFolderDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"folder_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"folder_path": {
				Type:     schema.TypeString,
				Required: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"project_identifier": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackDataWorksFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateFolder"
	request := make(map[string]interface{})
	folderPath := ConvertDataWorksFrontEndFolderPathToBackEndFolderPath(d.Get("folder_path").(string))
	request["FolderPath"] = folderPath
	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v
	}
	if v, ok := d.GetOk("project_identifier"); ok {
		request["ProjectIdentifier"] = v
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	response, err := client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(response["Data"], ":", request["ProjectId"]))

	return resourceAlibabacloudStackDataWorksFolderRead(d, meta)
}

func resourceAlibabacloudStackDataWorksFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksFolder(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksFolder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("folder_id", parts[0])
	d.Set("project_id", parts[1])
	d.Set("folder_path", object["FolderPath"].(string))
	return nil
}

func resourceAlibabacloudStackDataWorksFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	request := map[string]interface{}{
		"FolderId":   parts[0],
		"ProjectId":  parts[1],
	}
	if d.HasChange("folder_path") {
		folderPath := ConvertDataWorksFrontEndFolderPathToBackEndFolderPath(d.Get("folder_path").(string))
		absolutePath := folderPath
		_, lastDir := path.Split(absolutePath)
		request["FolderName"] = lastDir
	}
	if v, ok := d.GetOk("project_identifier"); ok {
		request["ProjectIdentifier"] = v
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	action := "UpdateFolder"
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	return resourceAlibabacloudStackDataWorksFolderRead(d, meta)
}

func resourceAlibabacloudStackDataWorksFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteFolder"
	request := map[string]interface{}{
		"FolderId":   parts[0],
		"ProjectId":  parts[1],
	}
	if v, ok := d.GetOk("project_identifier"); ok {
		request["ProjectIdentifier"] = v
	}
	request["PageSize"] = 1
	request["PageNumber"] = 1
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	return nil
}

func getConvertMap() map[string]string {
	convertMap := make(map[string]string)
	convertMap["Business Flow"] = "业务流程"
	convertMap["folderAlgm"] = "算法"
	convertMap["folderCDH"] = "CDH"
	convertMap["folderDi"] = "数据集成"
	convertMap["folderFlink"] = "Flink"
	convertMap["folderGeneral"] = "通用"
	convertMap["folderHologres"] = "Hologres"
	convertMap["folderMaxCompute"] = "MaxCompute"
	convertMap["folderUserDefined"] = "自定义"
	convertMap["folderEMR"] = "EMR"
	convertMap["folderErd"] = "数据模型"
	convertMap["folderADB"] = "AnalyticDB for PostgreSQL"
	convertMap["folderJdbc"] = "数据库"
	return convertMap
}

func ConvertDataWorksFrontEndFolderPathToBackEndFolderPath(source string) string {
	result := source
	convertMap := getConvertMap()

	for convert := range convertMap {
		result = strings.Replace(result, convert, convertMap[convert], 1)
	}
	return result
}

func ConvertDataWorksBackEndFolderPathToFrontEndFolderPath(source string) string {
	result := source
	convertMap := getConvertMap()

	for convert := range convertMap {
		result = strings.Replace(result, convertMap[convert], convert, 1)
	}
	return result
}
