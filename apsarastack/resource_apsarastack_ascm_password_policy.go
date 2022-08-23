package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"strconv"
	"strings"

	"time"
)

func resourceApsaraStackAscmPasswordPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackAscmPasswordPolicyCreate,
		Read:   resourceApsaraStackAscmPasswordPolicyRead,
		Update: resourceApsaraStackAscmPasswordPolicyUpdate,
		Delete: resourceApsaraStackAscmPasswordPolicyDelete,
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
}

func resourceApsaraStackAscmPasswordPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	value123 := strconv.Itoa(d.Get("minimum_password_length").(int))

	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "SetPasswordPolicy"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"RegionId":              client.RegionId,
		"AccessKeySecret":       client.SecretKey,
		"Department":            client.Department,
		"ResourceGroup":         client.ResourceGroup,
		"Product":               "ascm",
		"Action":                "SetPasswordPolicy",
		"Version":               "2019-05-10",
		"ProductName":           "ascm",
		"minimumPasswordLength": value123,
	}
	var response = PasswordPolicy{}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_password_policy", "", raw)
	}

	addDebug("SetPasswordPolicy", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)
	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_password_policy", "", ApsaraStackSdkGoERROR)
	}
	addDebug("SetPasswordPolicy", raw, requestInfo, bresponse.GetHttpContentString())
	_ = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

	d.SetId(fmt.Sprint(response.Data.ID))

	return resourceApsaraStackAscmPasswordPolicyUpdate(d, meta)
}

func resourceApsaraStackAscmPasswordPolicyRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmPasswordPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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

func resourceApsaraStackAscmPasswordPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceApsaraStackAscmPasswordPolicyRead(d, meta)
}

func resourceApsaraStackAscmPasswordPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmPasswordPolicy(d.Id())

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsResourceGroupExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsResourceGroupExist", check, requestInfo, map[string]string{"resourceGroupName": d.Id()})
	err = resource.Retry(1*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":        client.RegionId,
			"AccessKeySecret": client.SecretKey,
			"Product":         "ascm",
			"Action":          "ResetPasswordPolicy",
			"Version":         "2019-05-10",
			"ProductName":     "ascm",
			"id":              d.Id(),
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "ResetPasswordPolicy"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		_, err = ascmService.DescribeAscmPasswordPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "ResetPasswordPolicy", ApsaraStackSdkGoERROR)
	}
	return nil
}
