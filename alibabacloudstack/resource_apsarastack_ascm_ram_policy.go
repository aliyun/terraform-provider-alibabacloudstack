package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmRamPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringLenBetween(3, 64),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringLenBetween(0, 128),
			},
			"policy_document": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ram_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmRamPolicyCreate, resourceAlibabacloudStackAscmRamPolicyRead, resourceAlibabacloudStackAscmRamPolicyUpdate, resourceAlibabacloudStackAscmRamPolicyDelete)
	return resource
}

func resourceAlibabacloudStackAscmRamPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	name := d.Get("name").(string)

	policyDoc := d.Get("policy_document").(string)
	description := d.Get("description").(string)
	//resp := RamPolicies{}
	check, err := ascmService.DescribeAscmRamPolicy(name)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_ram_policy", "policy alreadyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if len(check.Data) == 0 {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "CreateRAMPolicy", "/ascm/auth/role/createRAMPolicy")
		request.QueryParams["policyName"] = name
		request.QueryParams["description"] = description
		request.QueryParams["policyDocument"] = policyDoc

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	err = resource.Retry(1*time.Minute, func() *resource.RetryError {
		check, err = ascmService.DescribeAscmRamPolicy(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.RetryableError(err)
	})
	d.SetId(check.Data[0].PolicyName + COLON_SEPARATED + fmt.Sprint(check.Data[0].ID))

	return nil
}

func resourceAlibabacloudStackAscmRamPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}

	response, err := ascmService.DescribeAscmRamPolicy(d.Id())
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		// Handle exceptions
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("name", did[0])
	d.Set("ram_id", did[1])
	d.Set("description", response.Data[0].Description)
	d.Set("policy_document", response.Data[0].PolicyDocument)
	return nil
}

func resourceAlibabacloudStackAscmRamPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	
	if d.IsNewResource() {
		return nil
	}
	
	_, err := ascmService.DescribeAscmRamPolicy(d.Id())
	if err != nil {
		return err
	}
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsInstanceExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	if d.HasChanges("name", "description", "policy_document") {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "UpdateRAMPolicy", "/ascm/auth/role/updateRAMPolicy")
		request.QueryParams["RamPolicyId"] = did[1]
		if d.HasChange("name") {
			request.QueryParams["NewPolicyName"] = d.Get("name").(string)
		}
		if d.HasChange("description") {
			request.QueryParams["NewDescription"] = d.Get("description").(string)
		}
		if d.HasChange("policy_document") {
			request.QueryParams["newPolicyDocument"] = d.Get("policy_document").(string)
		}

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}
	d.SetId(d.Get("name").(string) + COLON_SEPARATED + did[1])

	return nil
}

func resourceAlibabacloudStackAscmRamPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	_, err := ascmService.DescribeAscmRamPolicy(d.Id())
	if errmsgs.NotFoundError(err) {
		return nil
	}
	did := strings.Split(d.Id(), COLON_SEPARATED)

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsPolicyExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRAMPolicy", "/ascm/auth/role/removeRAMPolicy")
		request.QueryParams["ramPolicyId"] = did[1]

		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return resource.RetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_ram_policy", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})

	return nil
}
