package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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

func resourceAlibabacloudStackDataWorksConnectionCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateConnection"
	request := make(map[string]interface{})

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

	response, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprint(response["Data"], ":", request["ProjectId"], ":", request["Name"]))

	return resourceAlibabacloudStackDataWorksConnectionRead(d, meta)
}

func resourceAlibabacloudStackDataWorksConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	dataworksPublicService := DataworksPublicService{client}
	object, err := dataworksPublicService.DescribeDataWorksConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_data_works_folder dataworksPublicService.DescribeDataWorksConnection Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
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
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
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
	action := "UpdateConnection"

	_, err = client.DoTeaRequest("PUT", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	return resourceAlibabacloudStackDataWorksConnectionRead(d, meta)
}

func resourceAlibabacloudStackDataWorksConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteConnection"
	request := map[string]interface{}{
		"ConnectionId": parts[0],
	}
	
	_, err = client.DoTeaRequest("POST", "dataworks-public", "2020-05-18", action, "", nil, request)
	if err != nil {
		return err
	}
	return nil
}
