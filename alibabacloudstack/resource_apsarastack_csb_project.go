package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackCsbProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackCsbProjectCreate,
		Read:   resourceAlibabacloudStackCsbProjectRead,
		Update: resourceAlibabacloudStackCsbProjectUpdate,
		Delete: resourceAlibabacloudStackCsbProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Timeouts: &schema.ResourceTimeout{
			Update: schema.DefaultTimeout(5 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data2": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"csb_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"project_name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(1, 128),
			},
			"project_owner_email": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_owner_phone_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gmt_modified": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gmt_create": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"owner_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"api_num": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"src_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_flag": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cs_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"project_owner_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackCsbProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateProject"
	request := make(map[string]interface{})
	conn, err := client.NewCsbClient()
	if err != nil {
		return WrapError(err)
	}

	if v, ok := d.GetOk("csb_id"); ok {
		request["CsbId"] = v
	}

	if v, ok := d.GetOk("data"); ok {
		request["Data"] = v
	}

	request["RegionId"] = client.RegionId
	request["product"] = "CSB"
	request["OrganizationId"] = client.Department

	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_csb_project", action, AlibabacloudStackSdkGoERROR)
	}

	request1 := make(map[string]interface{})
	if v, ok := d.GetOk("project_name"); ok {
		request1["ProjectName"] = v
	}
	addDebug(action, response, request)
	d.SetId(fmt.Sprint(request["CsbId"], ":", request1["ProjectName"]))
	return resourceAlibabacloudStackCsbProjectRead(d, meta)
}
func resourceAlibabacloudStackCsbProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbProjectDetail(d.Id())
	if err != nil {
		if NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_csb_project csbService.DescribeCsbProjectDetail Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	fmt.Sprint(object["ProjectName"])
	d.Set("project_name", fmt.Sprint(object["ProjectName"]))
	d.Set("csb_id", fmt.Sprint(object["CsbId"]))
	//d.Set("project_owner_name", fmt.Sprint(object["projectOwnerName"]))
	//d.Set("gmt_modified", fmt.Sprint(object["GmtModified"]))
	//d.Set("gmt_create", fmt.Sprint(object["GmtCreate"]))
	//d.Set("owner_id", fmt.Sprint(object["OwnerId"]))
	//d.Set("api_num", fmt.Sprint(object["ApiNum"]))
	//d.Set("user_id", fmt.Sprint(object["UserId"]))
	//d.Set("delete_flag", fmt.Sprint(object["DeleteFlag"]))
	//d.Set("cs_id", fmt.Sprint(object["Id"]))
	//d.Set("status", fmt.Sprint(object["Status"]))
	//d.Set("src_type", "0")

	return nil
}
func resourceAlibabacloudStackCsbProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	var response map[string]interface{}
	request := map[string]interface{}{}
	if d.HasChange("project_name") {
		update = true
		if v, ok := d.GetOk("project_name"); ok {
			request["ProjectName"] = v
		}
	}
	if v, ok := d.GetOk("data2"); ok {
		request["Data"] = v
	}
	request["RegionId"] = client.RegionId
	request["Product"] = "CSB"
	request["OrganizationId"] = client.Department
	request["ResourceId"] = client.ResourceGroup
	if update {
		action := "UpdateProject"
		conn, err := client.NewCsbClient()
		if err != nil {
			return WrapError(err)
		}

		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
		}

	}
	return resourceAlibabacloudStackCsbProjectRead(d, meta)
}
func resourceAlibabacloudStackCsbProjectDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	csbService := CsbService{client}
	object, err := csbService.DescribeCsbProjectDetail(d.Id())
	action := "DeleteProject"
	var response map[string]interface{}
	conn, err := client.NewCsbClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{}

	request["ProjectId"] = object["Id"]
	request["CsbId"] = object["CsbId"]

	request["RegionId"] = client.RegionId
	request["Product"] = "CSB"
	request["OrganizationId"] = client.Department

	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-11-18"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})

	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
