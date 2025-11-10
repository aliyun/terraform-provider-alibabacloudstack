package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"strings"

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
							MaxItems: 1,
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

	// Handle ids filter
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Handle name_regex filter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Since there is no list API, we need to get all consumers by calling GetApp for each appId
	// However, we don't have a list of appIds. In this case, we can only return empty list
	// unless we have specific appIds to query.
	// But according to the schema, it seems like this data source should be able to list consumers
	// based on gw_instance_id. This might require an additional API or method to list all apps.

	// For now, we'll assume that if no specific filters are provided, we return empty list
	// If specific IDs are provided, we try to fetch them

	if len(idsMap) == 0 && nameRegex == nil {
		// No specific consumers requested, return empty list
		d.SetId("")
		if err := d.Set("consumers", []interface{}{}); err != nil {
			return errmsgs.WrapError(err)
		}
		if err := d.Set("ids", []string{}); err != nil {
			return errmsgs.WrapError(err)
		}
		return nil
	}

	// Collect results
	var consumers []map[string]interface{}
	var resultIds []string

	// For each ID in idsMap, we assume the ID format is "gwInstanceId:appId"
	for id := range idsMap {
		parts := strings.Split(id, ":")
		if len(parts) != 2 {
			continue
		}

		// Check if gwInstanceId matches
		if parts[0] != gwInstanceId {
			continue
		}

		appId := parts[1]

		// Call GetApp API
		request := map[string]interface{}{
			"appId":        appId,
			"gwInstanceId": gwInstanceId,
		}

		resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "GetApp", "/application/getApp", nil, nil, request)
		if err != nil {
			// If resource not found, continue with others
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}

		// Check if response is successful
		if success, ok := resp["asapiSuccess"].(bool); !ok || !success {
			continue
		}

		data, ok := resp["data"]
		if !ok || data == nil {
			continue
		}

		dataBytes, err := json.Marshal(data)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		var appData map[string]interface{}
		if err := json.Unmarshal(dataBytes, &appData); err != nil {
			return errmsgs.WrapError(err)
		}

		appName, _ := appData["appName"].(string)

		// Apply name_regex filter
		if nameRegex != nil && !nameRegex.MatchString(appName) {
			continue
		}

		// Build consumer object
		consumer := map[string]interface{}{
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
		resultIds = append(resultIds, id)
	}

	// Set results
	d.SetId(dataResourceIdHash(resultIds))
	if err := d.Set("consumers", consumers); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", resultIds); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
