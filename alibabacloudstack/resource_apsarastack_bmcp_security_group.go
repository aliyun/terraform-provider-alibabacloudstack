package alibabacloudstack

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
)

func resourceAlibabacloudStackBmcpSecurityGroup() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(2, 128),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 256),
			},
			"sg_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource_group": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"resource_group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"department": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"department_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ascm_create_user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackBmcpSecurityGroupCreate, resourceAlibabacloudStackBmcpSecurityGroupRead, resourceAlibabacloudStackBmcpSecurityGroupUpdate, resourceAlibabacloudStackBmcpSecurityGroupDelete)
	return resource
}

func resourceAlibabacloudStackBmcpSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"DeployType":    "bmcp",
		"AscmShareType": 0,
		"VpcId":         d.Get("vpc_id").(string),
		"Name":          d.Get("name").(string),
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "CreateSecurityGroup", "", nil, request, nil)
		addDebug("CreateSecurityGroup", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "TaskConflict", "UnknownError", errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bmcp_security_group", "CreateSecurityGroup", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("CreateSecurityGroup", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	data := response["data"].(map[string]interface{})
	sgId := data["data"].(map[string]interface{})["sgId"].(string)
	d.SetId(sgId)

	return nil
}

func resourceAlibabacloudStackBmcpSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"SgId": d.Id(),
	}

	var response map[string]interface{}
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroup", "", nil, request, nil)
		addDebug("ListSecurityGroup", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bmcp_security_group", "ListSecurityGroup", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ListSecurityGroup", raw, nil, request)
		response = raw
		return nil
	})
	if err != nil {
		return err
	}

	securityGroupList := response["data"].([]interface{})
	if len(securityGroupList) == 0 {
		log.Printf("[DEBUG] Resource alibabacloudstack_bmcp_security_group ListSecurityGroup Failed!!! SecurityGroup not found")
		d.SetId("")
		return nil
	}

	securityGroup := securityGroupList[0].(map[string]interface{})
	if securityGroup["sgId"].(string) != d.Id() {
		log.Printf("[DEBUG] Resource alibabacloudstack_bmcp_security_group ListSecurityGroup Failed!!! SgId mismatch")
		d.SetId("")
		return nil
	}

	d.Set("sg_id", securityGroup["sgId"])
	d.Set("name", securityGroup["name"])
	d.Set("description", securityGroup["description"])
	d.Set("resource_group", securityGroup["ResourceGroup"])
	d.Set("resource_group_name", securityGroup["ResourceGroupName"])
	d.Set("department", securityGroup["Department"])
	d.Set("department_name", securityGroup["DepartmentName"])
	d.Set("region_id", securityGroup["RegionId"])
	d.Set("ascm_create_user", securityGroup["AscmCreateUser"])
	d.Set("create_time", securityGroup["createTime"])
	d.Set("update_time", securityGroup["updateTime"])
	d.Set("vpc_id", securityGroup["vpcId"])

	return nil
}

func resourceAlibabacloudStackBmcpSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"SgId": d.Id(),
	}
	if d.HasChange("name") {
		request["Name"] = d.Get("name").(string)
	}
	if d.HasChange("description") {
		request["Description"] = d.Get("description").(string)
	}

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ModifySecurityGroup", "", nil, request, nil)
		addDebug("ModifySecurityGroup", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, errmsgs.Throttling) {
				time.Sleep(5 * time.Second)
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ModifySecurityGroup", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("ModifySecurityGroup", raw, nil, request)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func resourceAlibabacloudStackBmcpSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"SgId": d.Id(),
	}

	err := resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "DeleteSecurityGroup", "", nil, request, nil)
		addDebug("DeleteSecurityGroup", raw, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, "InvalidSecurityGroupId.NotFound", "Forbidden.SecurityGroupNotFound") {
				return nil
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteSecurityGroup", errmsgs.AlibabacloudStackSdkGoERROR, ""))
		}
		addDebug("DeleteSecurityGroup", raw, nil, request)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
