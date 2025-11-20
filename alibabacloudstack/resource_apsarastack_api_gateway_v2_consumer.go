package alibabacloudstack

import (
	"fmt"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAPIGatewayV2Consumer() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"auth_type": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"payload": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"oauth2_payload": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"authorization_code": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"client_credentials": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"implicit_grant": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"password_grant": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"token_expiration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"refresh_token_expiration": {
							Type:     schema.TypeInt,
							Optional: true,
						},
						"pkce": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"scopes": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"client_id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"client_secret": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"redirect_uris": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"expire_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"gw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"app_secret": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"app_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cascade_link_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"app_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auth_type_name": {
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
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAPIGatewayV2ConsumerCreate, resourceAlibabacloudStackAPIGatewayV2ConsumerRead, resourceAlibabacloudStackAPIGatewayV2ConsumerUpdate, resourceAlibabacloudStackAPIGatewayV2ConsumerDelete)
	return resource
}

func resourceAlibabacloudStackAPIGatewayV2ConsumerCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["authType"] = d.Get("auth_type")
	if v, ok := d.GetOk("payload"); ok {
		payload := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			payload[k] = v.(string)
		}
		request["payload"] = payload
	}
	if v, ok := d.GetOk("groups"); ok {
		groups := make([]string, 0)
		for _, item := range v.([]interface{}) {
			groups = append(groups, item.(string))
		}
		request["groups"] = groups
	}
	if v, ok := d.GetOk("oauth2_payload"); ok && len(v.([]interface{})) > 0 {
		oauth2Payload := make(map[string]interface{})
		payload := v.([]interface{})[0].(map[string]interface{})
		if val, ok := payload["authorization_code"]; ok {
			oauth2Payload["authorizationCode"] = val.(bool)
		}
		if val, ok := payload["client_credentials"]; ok {
			oauth2Payload["clientCredentials"] = val.(bool)
		}
		if val, ok := payload["implicit_grant"]; ok {
			oauth2Payload["implicitGrant"] = val.(bool)
		}
		if val, ok := payload["password_grant"]; ok {
			oauth2Payload["passwordGrant"] = val.(bool)
		}
		if val, ok := payload["token_expiration"]; ok {
			oauth2Payload["tokenExpiration"] = val.(int)
		}
		if val, ok := payload["refresh_token_expiration"]; ok {
			oauth2Payload["refreshTokenExpiration"] = val.(int)
		}
		if val, ok := payload["pkce"]; ok {
			oauth2Payload["pkce"] = val.(bool)
		}
		if val, ok := payload["scopes"]; ok {
			oauth2Payload["scopes"] = val.(string)
		}
		if val, ok := payload["client_id"]; ok {
			oauth2Payload["clientId"] = val.(string)
		}
		if val, ok := payload["client_secret"]; ok {
			oauth2Payload["clientSecret"] = val.(string)
		}
		if val, ok := payload["redirect_uris"]; ok {
			oauth2Payload["redirectUris"] = val.(string)
		}
		request["oauth2Payload"] = oauth2Payload
	}
	if v, ok := d.GetOk("key"); ok {
		request["key"] = v.(string)
	}
	request["appName"] = d.Get("app_name").(string)
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v.(string)
	}
	if v, ok := d.GetOk("expire_time"); ok {
		request["expireTime"] = v.(int)
	}
	if v, ok := d.GetOk("password"); ok {
		request["password"] = v.(string)
	}
	request["gwInstanceId"] = d.Get("gw_instance_id").(string)
	if v, ok := d.GetOk("app_secret"); ok {
		request["appSecret"] = v.(string)
	}
	if v, ok := d.GetOk("app_code"); ok {
		request["appCode"] = v.(string)
	}
	action := "CreateApp"
	pattern := "/application/createApp"
	idpre := "app"
	if v, ok := d.GetOk("cascade_link_ids"); ok && v.(*schema.Set).Len() > 0 {
		action = "CreateSourceApplication"
		pattern = "/sourceApplication/createSourceApplication"
		idpre = "sourceApp"
		request["cascadeLinkIds"] = v.(*schema.Set).List()
	}
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, request)
	if err != nil {
		return err
	}

	appId := resp["data"].(string)
	gwInstanceId := d.Get("gw_instance_id").(string)
	d.SetId(fmt.Sprintf("%s:%s:%s", idpre, gwInstanceId, appId))

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ConsumerRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	service := ApiGateWayV2Service{client}

	object, err := service.DescribeAPIGatewayV2Consumer(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("app_name", object["appName"])
	d.Set("description", object["description"])
	d.Set("groups", object["groups"])
	d.Set("key", object["key"])
	d.Set("auth_type", object["authType"])
	d.Set("auth_type_name", object["authTypeName"])
	d.Set("app_id", object["appId"])
	d.Set("gw_instance_id", strings.Split(d.Id(), ":")[1])
	d.Set("token", object["token"])
	d.Set("use_white_list", object["useWhiteList"])
	d.Set("enable", object["enable"])

	if v, ok := object["expireTime"]; ok && v != nil {
		d.Set("expire_time", v)
	}

	if v, ok := object["appSecret"]; ok {
		d.Set("app_secret", v)
	}

	if v, ok := object["password"]; ok {
		d.Set("password", v)
	}

	if v, ok := object["appCode"]; ok {
		d.Set("app_code", v)
	}

	if v, ok := object["requestHeader"]; ok {
		d.Set("request_header", v)
	}

	if v, ok := object["secretKey"]; ok {
		d.Set("secret_key", v)
	}

	if v, ok := object["accessKey"]; ok {
		d.Set("access_key", v)
	}

	if oauth2Payload, ok := object["oauth2Payload"]; ok && oauth2Payload != nil {
		payload := oauth2Payload.(map[string]interface{})
		oauth2Data := make(map[string]interface{})

		if v, ok := payload["authorizationCode"]; ok {
			oauth2Data["authorization_code"] = v
		}
		if v, ok := payload["clientCredentials"]; ok {
			oauth2Data["client_credentials"] = v
		}
		if v, ok := payload["implicitGrant"]; ok {
			oauth2Data["implicit_grant"] = v
		}
		if v, ok := payload["passwordGrant"]; ok {
			oauth2Data["password_grant"] = v
		}
		if v, ok := payload["tokenExpiration"]; ok {
			oauth2Data["token_expiration"] = v
		}
		if v, ok := payload["refreshTokenExpiration"]; ok {
			oauth2Data["refresh_token_expiration"] = v
		}
		if v, ok := payload["pkce"]; ok {
			oauth2Data["pkce"] = v
		}
		if v, ok := payload["scopes"]; ok {
			oauth2Data["scopes"] = v
		}
		if v, ok := payload["clientId"]; ok {
			oauth2Data["client_id"] = v
		}
		if v, ok := payload["clientSecret"]; ok {
			oauth2Data["client_secret"] = v
		}
		if v, ok := payload["redirectUris"]; ok {
			oauth2Data["redirect_uris"] = v
		}

		d.Set("oauth2_payload", []map[string]interface{}{oauth2Data})
	}
	linkEntityRelations, ok := object["linkEntityRelations"]
	if ok && linkEntityRelations != nil {
		link_ids := make([]interface{}, 0)
		for _, item := range linkEntityRelations.([]interface{}) {
			data := item.(map[string]interface{})
			if linkId, ok := data["linkId"]; ok {
				link_ids = append(link_ids, linkId)
			}
		}
		d.Set("cascade_link_ids", link_ids)
	}

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ConsumerUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.IsNewResource() {
		return nil
	}
	params, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idpre := params[0]
	gwInstanceId := params[1]
	appId := params[2]
	request := make(map[string]interface{})
	request["gwInstanceId"] = gwInstanceId
	request["appId"] = appId
	request["appName"] = d.Get("app_name")
	if v, ok := d.GetOk("description"); ok {
		request["description"] = v
	}
	if v, ok := d.GetOk("groups"); ok {
		request["groups"] = v
	}
	if v, ok := d.GetOk("key"); ok {
		request["key"] = v
	}
	if v, ok := d.GetOk("expire_time"); ok {
		request["expireTime"] = v
	}
	if v, ok := d.GetOk("app_secret"); ok {
		request["appSecret"] = v
	}
	if v, ok := d.GetOk("password"); ok {
		request["password"] = v
	}
	if v, ok := d.GetOk("app_code"); ok {
		request["appCode"] = v
	}
	if v, ok := d.GetOk("auth_type"); ok {
		request["authType"] = v
	}
	if v, ok := d.GetOk("payload"); ok {
		payload := make(map[string]string)
		for k, v := range v.(map[string]interface{}) {
			payload[k] = v.(string)
		}
		request["payload"] = payload
	}
	if v, ok := d.GetOk("oauth2_payload"); ok && len(v.([]interface{})) > 0 {
		oauth2Payload := make(map[string]interface{})
		oauth2List := v.([]interface{})[0].(map[string]interface{})
		if val, ok := oauth2List["authorization_code"]; ok {
			oauth2Payload["authorizationCode"] = val
		}
		if val, ok := oauth2List["client_credentials"]; ok {
			oauth2Payload["clientCredentials"] = val
		}
		if val, ok := oauth2List["implicit_grant"]; ok {
			oauth2Payload["implicitGrant"] = val
		}
		if val, ok := oauth2List["password_grant"]; ok {
			oauth2Payload["passwordGrant"] = val
		}
		if val, ok := oauth2List["token_expiration"]; ok {
			oauth2Payload["tokenExpiration"] = val
		}
		if val, ok := oauth2List["refresh_token_expiration"]; ok {
			oauth2Payload["refreshTokenExpiration"] = val
		}
		if val, ok := oauth2List["pkce"]; ok {
			oauth2Payload["pkce"] = val
		}
		if val, ok := oauth2List["scopes"]; ok {
			oauth2Payload["scopes"] = val
		}
		if val, ok := oauth2List["client_id"]; ok {
			oauth2Payload["clientId"] = val
		}
		if val, ok := oauth2List["client_secret"]; ok {
			oauth2Payload["clientSecret"] = val
		}
		if val, ok := oauth2List["redirect_uris"]; ok {
			oauth2Payload["redirectUris"] = val
		}
		request["oauth2Payload"] = oauth2Payload
	}
	action := "ModifyApp"
	pattern := "/application/modifyApp"
	if idpre == "sourceApp" {
		cascadeLinkIds := make([]string, 0)
		cascade_link_ids := d.Get("cascade_link_ids").(*schema.Set).List()
		for _, item := range cascade_link_ids {
			cascadeLinkIds = append(cascadeLinkIds, item.(string))
		}
		request["cascadeLinkIds"] = cascadeLinkIds
		action = "ModifySourceApplication"
		pattern = "/sourceApplication/modifySourceApplication"
	}

	_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_api_gateway_v2_consumer", "ModifyApp", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}

func resourceAlibabacloudStackAPIGatewayV2ConsumerDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	params, err := ParseResourceId(d.Id(), 3)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	idpre := params[0]
	gwInstanceId := params[1]
	appId := params[2]

	reqQuery := map[string]interface{}{
		"appId":        appId,
		"gwInstanceId": gwInstanceId,
	}
	action := "DeleteApp"
	pattern := "/application/deleteApp"
	if idpre == "sourceApp" {
		action = "DeleteSourceApplication"
		pattern = "/sourceApplication/deleteSourceApplication"
	}

	_, err = client.DoTeaRequest("POST", "csb2", "2023-02-06", action, pattern, nil, nil, reqQuery)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
