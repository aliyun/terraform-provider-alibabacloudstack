package alibabacloudstack

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserGroupCreate,
		Read:   resourceAlibabacloudStackAscmUserGroupRead,
		Update: resourceAlibabacloudStackAscmUserGroupUpdate,
		Delete: resourceAlibabacloudStackAscmUserGroupDelete,
		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_in_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Deprecated: "Attribute role_in_ids has been deprecated and replaced with role_ids.",
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackAscmUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client

	groupName := d.Get("group_name").(string)
	organizationId := d.Get("organization_id").(string)

	var loginNamesList []string

	if v, ok := d.GetOk("role_in_ids"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}
	
	if v, ok := d.GetOk("role_ids"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())

		for _, loginName := range loginNames {
			loginNamesList = append(loginNamesList, loginName)
		}
	}


	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	request.Headers["x-ascm-product-name"] = "ascm"
	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"groupName":        groupName,
		"organizationId":   organizationId,
		"roleIdList":       loginNamesList,
		"SecurityToken":    client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	requeststring, err := json.Marshal(QueryParams)

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["Content-Type"] = requests.Json
	request.SetContent(requeststring)
	request.PathPattern = "/roa/ascm/auth/user/createUserGroup"
	request.ApiName = "CreateUserGroup"
	request.RegionId = client.RegionId
	request.Headers["RegionId"] = client.RegionId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	log.Printf("response of raw CreateUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", raw)
	}

	addDebug("CreateUserGroup", raw, requestInfo, request)

	bresponse, _ := raw.(*responses.CommonResponse)

	if bresponse.GetHttpStatus() != 200 {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_user_group", "CreateUserGroup", AlibabacloudStackSdkGoERROR)
	}
	addDebug("CreateUserGroup", raw, requestInfo, bresponse.GetHttpContentString())

	d.SetId(groupName)

	return resourceAlibabacloudStackAscmUserGroupUpdate(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {

	return resourceAlibabacloudStackAscmUserGroupRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
	}
	if len(object.Data) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("user_group_id", strconv.Itoa(object.Data[0].Id))
	d.Set("group_name", object.Data[0].GroupName)
	d.Set("organization_id", strconv.Itoa(object.Data[0].OrganizationId))

	var roleIds []string
	if len(object.Data[0].Roles) >= 0 {
		for _, role := range object.Data[0].Roles {
			roleIds = append(roleIds, strconv.Itoa(role.Id))
		}
	}
	d.Set("role_ids", roleIds)

	return nil
}

func resourceAlibabacloudStackAscmUserGroupDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client

	check, err := ascmService.DescribeAscmUserGroup(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsUserGroupExist", AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsUserGroupExist", check, requestInfo, map[string]string{"groupName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {

		request := requests.NewCommonRequest()
		if client.Config.Insecure {
			request.SetHTTPSInsecure(client.Config.Insecure)
		}
		request.QueryParams = map[string]string{
			"RegionId":         client.RegionId,
			
			"Product":          "ascm",
			"Action":           "DeleteUserGroup",
			"Version":          "2019-05-10",
			"ProductName":      "ascm",
			"userGroupId":      strconv.Itoa(check.Data[0].Id),
			"SecurityToken":    client.Config.SecurityToken,
			"SignatureVersion": "1.0",
			"SignatureMethod":  "HMAC-SHA1",
		}

		request.Method = "POST"
		request.Product = "ascm"
		request.Version = "2019-05-10"
		request.ServiceCode = "ascm"
		request.Domain = client.Domain
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.ApiName = "DeleteUserGroup"
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.RegionId = client.RegionId

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})
		if err != nil {
			return resource.RetryableError(err)
		}
		check, err = ascmService.DescribeAscmUserGroup(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		addDebug("DeleteUserGroup", raw, request)
		return nil
	})
	return nil
}
