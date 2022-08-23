package apsarastack

import (
	"fmt"
	"log"
	"path"
	"strings"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackDataWorksFolder() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackDataWorksFolderCreate,
		Read:   resourceApsaraStackDataWorksFolderRead,
		Update: resourceApsaraStackDataWorksFolderUpdate,
		Delete: resourceApsaraStackDataWorksFolderDelete,
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

func resourceApsaraStackDataWorksFolderCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateFolder"
	request := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	folderPath := ConvertDataWorksFrontEndFolderPathToBackEndFolderPath(d.Get("folder_path").(string))
	request["FolderPath"] = folderPath
	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v
	}
	if v, ok := d.GetOk("project_identifier"); ok {
		request["ProjectIdentifier"] = v
	}
	request["RegionId"] = "default"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_data_works_folder", action, ApsaraStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"], ":", request["ProjectId"]))

	return resourceApsaraStackDataWorksFolderRead(d, meta)
}
func resourceApsaraStackDataWorksFolderRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksFolder(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource apsarastack_data_works_folder dataworksPublicService.DescribeDataWorksFolder Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	d.Set("folder_id", parts[0])
	d.Set("project_id", parts[1])
	d.Set("folder_path", object["FolderPath"].(string))
	return nil
}
func resourceApsaraStackDataWorksFolderUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FolderId":  parts[0],
		"ProjectId": parts[1],
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
	request["RegionId"] = "default"
	action := "UpdateFolder"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return resourceApsaraStackDataWorksFolderRead(d, meta)
}
func resourceApsaraStackDataWorksFolderDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteFolder"
	var response map[string]interface{}
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"FolderId":  parts[0],
		"ProjectId": parts[1],
	}

	if v, ok := d.GetOk("project_identifier"); ok {
		request["ProjectIdentifier"] = v
	}
	request["RegionId"] = "default"
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
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
