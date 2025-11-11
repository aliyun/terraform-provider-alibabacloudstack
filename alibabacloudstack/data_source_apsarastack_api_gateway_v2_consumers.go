package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackAPIGatewayV2Consumers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2ConsumersRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"appid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"consumers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"auth_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"auth_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"token": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"use_white_list": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"enable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"expire_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"app_secret": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"password": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"app_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"request_header": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"secret_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"access_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"oauth2_payload": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"authorization_code": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"client_credentials": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"implicit_grant": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"password_grant": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"token_expiration": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"refresh_token_expiration": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"pkce": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"scopes": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"client_secret": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"redirect_uris": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2ConsumersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Parse gwInstanceId from schema
	gwInstanceId := d.Get("gw_instance_id").(string)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	request := map[string]interface{}{
		"gwInstanceId": gwInstanceId,
		"current":      1,
		"size":         10,
	}
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListApps", "/application/listApps", nil, nil, request)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	data, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	var consumers []map[string]interface{}
	var ids []string

	for _, v := range data.([]interface{}) {
		appData := v.(map[string]interface{})
		if v, ok := d.GetOk("appid"); ok && appData["appId"] != v {
			continue
		}
		appId := fmt.Sprintf("%s:%s", gwInstanceId, appData["appId"])
		appName, _ := appData["appName"].(string)
		if len(idsMap) > 0 {
			if _, exist := idsMap[appId]; !exist {
				continue
			}
		}
		if nameRegex != nil && !nameRegex.MatchString(appName) {
			continue
		}

		consumer := map[string]interface{}{
			"id":             appId,
			"app_name":       appName,
			"description":    appData["description"],
			"groups":         appData["groups"],
			"key":            appData["key"],
			"auth_type":      appData["authType"],
			"auth_type_name": appData["authTypeName"],
			"app_id":         appData["appId"],
			"token":          appData["token"],
			"use_white_list": appData["useWhiteList"],
			"enable":         appData["enable"],
		}

		// Handle optional fields
		if v, ok := appData["expireTime"]; ok && v != nil {
			consumer["expire_time"] = v
		}

		if v, ok := appData["appSecret"]; ok {
			consumer["app_secret"] = v
		}

		if v, ok := appData["password"]; ok {
			consumer["password"] = v
		}

		if v, ok := appData["appCode"]; ok {
			consumer["app_code"] = v
		}

		if v, ok := appData["requestHeader"]; ok {
			consumer["request_header"] = v
		}

		if v, ok := appData["secretKey"]; ok {
			consumer["secret_key"] = v
		}

		if v, ok := appData["accessKey"]; ok {
			consumer["access_key"] = v
		}

		// Handle oauth2_payload
		if oauth2Payload, ok := appData["oauth2Payload"]; ok && oauth2Payload != nil {
			if payloadMap, ok := oauth2Payload.(map[string]interface{}); ok {
				oauth2Data := make(map[string]interface{})

				if v, ok := payloadMap["authorizationCode"]; ok {
					oauth2Data["authorization_code"] = v
				}
				if v, ok := payloadMap["clientCredentials"]; ok {
					oauth2Data["client_credentials"] = v
				}
				if v, ok := payloadMap["implicitGrant"]; ok {
					oauth2Data["implicit_grant"] = v
				}
				if v, ok := payloadMap["passwordGrant"]; ok {
					oauth2Data["password_grant"] = v
				}
				if v, ok := payloadMap["tokenExpiration"]; ok {
					oauth2Data["token_expiration"] = v
				}
				if v, ok := payloadMap["refreshTokenExpiration"]; ok {
					oauth2Data["refresh_token_expiration"] = v
				}
				if v, ok := payloadMap["pkce"]; ok {
					oauth2Data["pkce"] = v
				}
				if v, ok := payloadMap["scopes"]; ok {
					oauth2Data["scopes"] = v
				}
				if v, ok := payloadMap["clientId"]; ok {
					oauth2Data["client_id"] = v
				}
				if v, ok := payloadMap["clientSecret"]; ok {
					oauth2Data["client_secret"] = v
				}
				if v, ok := payloadMap["redirectUris"]; ok {
					oauth2Data["redirect_uris"] = v
				}

				consumer["oauth2_payload"] = []map[string]interface{}{oauth2Data}
			}
		}

		consumers = append(consumers, consumer)
		ids = append(ids, appId)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("consumers", consumers); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
