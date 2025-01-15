package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAscmPasswordPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAscmPasswordPoliciesRead,
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

func dataSourceAlibabacloudStackAscmPasswordPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "GetPasswordPolicy", "/ascm/auth/user/getPasswordPolicy")
	response := PasswordPolicy{}

	for {
		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf(" response of raw GetPasswordPolicy : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_password_policies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)

		if err != nil {
			return errmsgs.WrapError(err)
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
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
