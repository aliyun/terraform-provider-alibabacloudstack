package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackAscmPasswordPolicy() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"hard_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"max_login_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"minimum_password_length": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(8, 32),
			},
			"password_reuse_prevention": {
				Type:     schema.TypeInt,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmPasswordPolicyCreate, resourceAlibabacloudStackAscmPasswordPolicyRead, resourceAlibabacloudStackAscmPasswordPolicyUpdate, resourceAlibabacloudStackAscmPasswordPolicyDelete)
	return resource
}

func resourceAlibabacloudStackAscmPasswordPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	value123 := strconv.Itoa(d.Get("minimum_password_length").(int))

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "SetPasswordPolicy", "/ascm/auth/user/setPasswordPolicy")
	request.QueryParams["minimumPasswordLength"] = value123

	var response = PasswordPolicy{}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_password_policy", "", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug("SetPasswordPolicy", raw, request)
	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_password_policy", "", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("SetPasswordPolicy", raw, request, bresponse.GetHttpContentString())
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	d.SetId(fmt.Sprint(response.Data.ID))

	return nil
}

func resourceAlibabacloudStackAscmPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmPasswordPolicy(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("hard_expiry", object.Data.HardExpiry)
	d.Set("require_numbers", object.Data.RequireNumbers)
	d.Set("require_symbols", object.Data.RequireSymbols)
	d.Set("require_lowercase_characters", object.Data.RequireLowercaseCharacters)
	d.Set("require_uppercase_characters", object.Data.RequireUppercaseCharacters)
	d.Set("max_login_attempts", object.Data.MaxLoginAttemps)
	d.Set("max_password_age", object.Data.MaxPasswordAge)
	d.Set("minimum_password_length", object.Data.MinimumPasswordLength)
	d.Set("password_reuse_prevention", object.Data.PasswordReusePrevention)

	return nil
}

func resourceAlibabacloudStackAscmPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	noUpdateAllowedFields := []string{"hard_expiry", "require_numbers", "require_lowercase_characters", "require_uppercase_characters", "max_login_attempts", "password_reuse_prevention", "require_symbols", "max_password_age", "minimum_password_length"}
	return noUpdatesAllowedCheck(d, noUpdateAllowedFields)
}

func resourceAlibabacloudStackAscmPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	check, err := ascmService.DescribeAscmPasswordPolicy(d.Id())

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsResourceGroupExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsResourceGroupExist", check, map[string]string{"resourceGroupName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetPasswordPolicy", "/ascm/auth/user/resetPasswordPolicy")
		request.QueryParams["id"] = d.Id()

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ResetPasswordPolicy", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		_, err = ascmService.DescribeAscmPasswordPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ResetPasswordPolicy", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
