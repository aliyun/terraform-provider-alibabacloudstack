package apsarastack

import (
	"encoding/json"
	"fmt"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"log"
	"strings"
)

func dataSourceApsaraStackAscmPasswordPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackAscmPasswordPoliciesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
				Computed: true,
			},
			"hard_expiry": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"require_numbers": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"require_symbols": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"require_lowercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"require_uppercase_characters": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"max_login_attempts": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_password_age": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"minimum_password_length": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password_reuse_prevention": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hard_expiry": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"require_numbers": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"require_symbols": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"require_lowercase_characters": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"require_uppercase_characters": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"max_login_attempts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"max_password_age": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"minimum_password_length": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"password_reuse_prevention": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackAscmPasswordPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}
	request.Method = "POST"
	request.Product = "ascm"
	request.Version = "2019-05-10"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.ApiName = "GetPasswordPolicy"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeyId": client.AccessKey, "AccessKeySecret": client.SecretKey, "Product": "ascm", "RegionId": client.RegionId, "Action": "GetPasswordPolicy", "Department": client.Department, "ResourceGroup": client.ResourceGroup, "Version": "2019-05-10"}
	response := PasswordPolicy{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw GetPasswordPolicy : %s", raw)

		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_ascm_password_policies", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		bresponse, _ := raw.(*responses.CommonResponse)

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if err != nil {
			return WrapError(err)
		}
		if response.Code == "200" || bresponse == nil {
			break
		}
	}
	var ids []string
	var s []map[string]interface{}
	mapping := map[string]interface{}{
		"hard_expiry":                  response.Data.HardExpiry,
		"require_numbers":              response.Data.RequireNumbers,
		"require_symbols":              response.Data.RequireSymbols,
		"require_lowercase_characters": response.Data.RequireLowercaseCharacters,
		"require_uppercase_characters": response.Data.RequireUppercaseCharacters,
		"max_login_attempts":           response.Data.MaxLoginAttemps,
		"max_password_age":             response.Data.MaxPasswordAge,
		"minimum_password_length":      response.Data.MinimumPasswordLength,
		"password_reuse_prevention":    response.Data.PasswordReusePrevention,
	}
	s = append(s, mapping)
	ids = append(ids, fmt.Sprint(response.Data.ID))

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("policies", s); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
