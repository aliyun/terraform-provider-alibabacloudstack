package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackLogonPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "ALLOW",
				ValidateFunc: validation.StringInSlice([]string{
					"ALLOW", "DENY"}, false),
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackLogonPolicyCreate, 
		resourceAlibabacloudStackLogonPolicyRead, 
		resourceAlibabacloudStackLogonPolicyUpdate, 
		resourceAlibabacloudStackLogonPolicyDelete)
	return resource
}

func resourceAlibabacloudStackLogonPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	descr := d.Get("description").(string)
	rule := d.Get("rule").(string)
	object, err := ascmService.DescribeAscmLogonPolicy(name)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if len(object.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddLoginPolicy", "/ascm/auth/loginPolicy/addLoginPolicy")
		mergeMaps(request.QueryParams, map[string]string{
			"AccountInfo":            "123456",
			"SignatureVersion":       "1.0",
			"ProductName":            "ascm",
			"name":                   name,
			"description":            descr,
			"rule":                   rule,
			"organizationVisibility": "organizationVisibility.organization",
		})
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw AddLoginPolicy : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_login_policy", "AddLoginPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("AddLoginPolicy", raw, requestInfo, request)

		if !bresponse.IsSuccess() {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_login_policy", "AddLoginPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("AddLoginPolicy", raw, requestInfo, bresponse.GetHttpContentString())
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		object, err = ascmService.DescribeAscmLogonPolicy(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if len(object.Data) != 0 {
			return nil
		}
		return resource.RetryableError(err)
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_login_policy", "Failed to add login Policy", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.SetId(object.Data[0].Name)

	return nil
}

func resourceAlibabacloudStackLogonPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var name, rule, desc string
	if d.HasChange("name") {
		name = d.Get("name").(string)
	}
	if d.HasChange("rule") {
		rule = d.Get("rule").(string)
	}
	if d.HasChange("description") {
		desc = d.Get("description").(string)
	}
	policyId := fmt.Sprint(d.Get("policy_id").(int))

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ModifyLoginPolicy", "/ascm/auth/loginPolicy/modifyLoginPolicy")
	mergeMaps(request.QueryParams, map[string]string{
		"id":          policyId,
		"name":        name,
		"rule":        rule,
		"description": desc,
	})
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_login_policy", "LoginPolicyUpdateRequestFailed", raw, errmsg)
	}

	if !bresponse.IsSuccess() {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_login_policy", "LoginPolicyUpdateFailed", raw, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	object, err := ascmService.DescribeAscmLogonPolicy(name)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	d.SetId(object.Data[0].Name)

	return nil
}

func resourceAlibabacloudStackLogonPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmLogonPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", object.Data[0].Name)
	d.Set("description", object.Data[0].Description)
	d.Set("policy_id", object.Data[0].ID)
	d.Set("rule", object.Data[0].Rule)
	return nil
}

func resourceAlibabacloudStackLogonPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)

	check, err := ascmService.DescribeAscmLogonPolicy(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsLoginPolicyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsLoginPolicyExist", check, requestInfo, map[string]string{"loginpolicyName": d.Id()})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveLoginPolicyByName", "/ascm/auth/loginPolicy/removeLoginPolicyByName")
		mergeMaps(request.QueryParams, map[string]string{
			"AccountInfo":      "123456",
			"SignatureVersion": "1.0",
			"ProductName":      "ascm",
			"name":             name,
		})
		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_login_policy", "RemoveLoginPolicyByName", raw, errmsg))
		}

		_, err = ascmService.DescribeAscmLogonPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return nil
}
