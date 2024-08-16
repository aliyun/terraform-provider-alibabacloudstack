package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDataWorksConnection() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDataWorksConnectionCreate,
		Read:   resourceAlibabacloudStackDataWorksConnectionRead,
		Update: resourceAlibabacloudStackDataWorksConnectionUpdate,
		Delete: resourceAlibabacloudStackDataWorksConnectionDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"connection_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Computed: true,
			},
			"connection_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"odps", "mysql", "rds", "oss", "sqlserver", "polardb", "oracle", "mongodb", "emr", "postgresql", "analyticdb_for_mysql", "hybriddb_for_postgresql", "holo"}, false),
			},
			"content": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"env_type": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntInSlice([]int{0, 1}),
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sub_type": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: rdsDiffSuppressFunc,
				ValidateFunc:     validation.StringInSlice([]string{"mysql", "sqlserver", "postgresql"}, false),
			},
		},
	}
}

func resourceAlibabacloudStackDataWorksConnectionCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateConnection"
	request := make(map[string]interface{})
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("connection_type"); ok {
		request["ConnectionType"] = v.(string)
	}

	if v, ok := d.GetOk("content"); ok {
		content, _ := json.Marshal(v)
		request["Content"] = string(content)
	}

	if v, ok := d.GetOk("env_type"); ok {
		request["EnvType"] = v.(int)
	}

	if v, ok := d.GetOk("name"); ok {
		request["Name"] = v.(string)
	}

	if v, ok := d.GetOk("project_id"); ok {
		request["ProjectId"] = v.(string)
	}

	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	if v, ok := d.GetOk("sub_type"); ok {
		request["SubType"] = v.(string)
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "dataworks-public"
	request["product"] = "dataworks-public"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_data_works_folder", action, AlibabacloudStackSdkGoERROR)
	}

	d.SetId(fmt.Sprint(response["Data"], ":", request["ProjectId"], ":", request["Name"]))

	return resourceAlibabacloudStackDataWorksConnectionRead(d, meta)
}
func resourceAlibabacloudStackDataWorksConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksConnection(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	d.Set("connection_id", parts[0])
	d.Set("project_id", parts[1])
	d.Set("name", parts[2])
	d.Set("connection_type", object["ConnectionType"].(string))

	// 由于密码返回为 *** 与原来不符，注释掉下面代码
	//var tempMap map[string]interface{}
	//err = json.Unmarshal([]byte(object["Content"].(string)), &tempMap)
	//d.Set("content", tempMap)
	return nil
}
func resourceAlibabacloudStackDataWorksConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ConnectionId": parts[0],
		"ProjectId":    parts[1],
	}
	if v, ok := d.GetOk("connection_type"); ok {
		request["ConnectionType"] = v.(string)
	}
	if d.HasChange("content") {
		v := d.Get("content")
		content, _ := json.Marshal(v)
		request["Content"] = string(content)

	}
	if d.HasChange("env_type") {
		request["EnvType"] = d.Get("env_type").(int)
	}
	if d.HasChange("description") {
		request["Description"] = d.Get("description").(string)
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "dataworks-public"
	request["product"] = "dataworks-public"
	request["OrganizationId"] = client.Department
	action := "UpdateConnection"
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("PUT"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return resourceAlibabacloudStackDataWorksConnectionRead(d, meta)
}
func resourceAlibabacloudStackDataWorksConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return WrapError(err)
	}
	action := "DeleteConnection"
	var response map[string]interface{}
	conn, err := client.NewDataworkspublicClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"ConnectionId": parts[0],
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "dataworks-public"
	request["product"] = "dataworks-public"
	request["OrganizationId"] = client.Department
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2020-05-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
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
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
