package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/apsara-stack/terraform-provider-apsarastack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"log"
	"strings"
	"time"
)

func resourceApsaraStackLogonPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackLogonPolicyCreate,
		Read:   resourceApsaraStackLogonPolicyRead,
		Update: resourceApsaraStackLogonPolicyUpdate,
		Delete: resourceApsaraStackLogonPolicyDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

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
			"time_range": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Default:  "flannel",
							Optional: true,
						},
						"config": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Default:  false,
						},
					},
				},
			},
		},
	}
}
func resourceApsaraStackLogonPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var requestInfo *ecs.Client
	ascmService := AscmService{client}
	name := d.Get("name").(string)
	descr := d.Get("description").(string)
	rule := d.Get("rule").(string)
	object, err := ascmService.DescribeAscmLogonPolicy(name)
	if err != nil {

		return WrapError(err)
	}
	if len(object.Data) == 0 {

		request := requests.NewCommonRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":               client.RegionId,
			"AccessKeySecret":        client.SecretKey,
			"Product":                "ascm",
			"Department":             client.Department,
			"ResourceGroup":          client.ResourceGroup,
			"Action":                 "AddLoginPolicy",
			"AccountInfo":            "123456",
			"Version":                "2019-05-10",
			"SignatureVersion":       "1.0",
			"ProductName":            "ascm",
			"name":                   name,
			"description":            descr,
			"rule":                   rule,
			"organizationVisibility": "organizationVisibility.organization",
		}
		request.Domain = client.Domain
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.ApiName = "AddLoginPolicy"
		request.RegionId = client.RegionId
		request.Headers = map[string]string{"RegionId": client.RegionId}
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw AddLoginPolicy : %s", raw)

		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "AddLoginPolicy", raw)
		}
		addDebug("AddLoginPolicy", raw, requestInfo, request)

		bresponse, _ := raw.(*responses.CommonResponse)
		if bresponse.GetHttpStatus() != 200 {
			return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "AddLoginPolicy", ApsaraStackSdkGoERROR)
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
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "Failed to add login Policy", ApsaraStackSdkGoERROR)
	}

	d.SetId(object.Data[0].Name)

	return resourceApsaraStackLogonPolicyUpdate(d, meta)
}

func resourceApsaraStackLogonPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	request := requests.NewCommonRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
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

	request.QueryParams = map[string]string{
		"RegionId":        client.RegionId,
		"AccessKeySecret": client.SecretKey,
		"Product":         "ascm",
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "ModifyLoginPolicy",
		"Version":         "2019-05-10",
		"ProductName":     "ascm",
		"id":              policyId,
		"Name":            name,
		"Rule":            rule,
		"Description":     desc,
	}
	request.Domain = client.Domain
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.ApiName = "ModifyLoginPolicy"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "LoginPolicyUpdateRequestFailed", raw)
	}
	bresponse, _ := raw.(*responses.CommonResponse)

	if !bresponse.IsSuccess() {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_ascm_login_policy", "LoginPolicyUpdateFailed", raw)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
	if err != nil {
		return WrapError(err)
	}
	object, err := ascmService.DescribeAscmLogonPolicy(name)
	if err != nil {
		return WrapError(err)
	}

	d.SetId(object.Data[0].Name)

	return resourceApsaraStackLogonPolicyRead(d, meta)
}
func resourceApsaraStackLogonPolicyRead(d *schema.ResourceData, meta interface{}) error {
	wiatSecondsIfWithTest(1)
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmLogonPolicy(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}

	d.Set("name", object.Data[0].Name)
	d.Set("description", object.Data[0].Description)
	d.Set("policy_id", object.Data[0].ID)
	d.Set("rule", object.Data[0].Rule)
	return nil
}
func resourceApsaraStackLogonPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	name := d.Get("name").(string)

	check, err := ascmService.DescribeAscmLogonPolicy(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsLoginPolicyExist", ApsaraStackSdkGoERROR)
	}
	addDebug("IsLoginPolicyExist", check, requestInfo, map[string]string{"loginpolicyName": d.Id()})
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		request := requests.NewCommonRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.RegionId = client.RegionId
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			"AccessKeySecret":  client.SecretKey,
			"Product":          "ascm",
			"Department":       client.Department,
			"ResourceGroup":    client.ResourceGroup,
			"Action":           "RemoveLoginPolicyByName",
			"AccountInfo":      "123456",
			"Version":          "2019-05-10",
			"SignatureVersion": "1.0",
			"ProductName":      "ascm",
			"Name":             name,
		}
		request.Domain = client.Domain
		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "RemoveLoginPolicyByName"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		_, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})

		if err != nil {
			return resource.RetryableError(err)
		}

		_, err = ascmService.DescribeAscmLogonPolicy(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}

		return nil
	})

	return nil
}
