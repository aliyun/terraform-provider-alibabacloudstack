package alibabacloudstack

import (
	"encoding/json"
	"fmt"
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

func resourceAlibabacloudStackAscmUserGroupRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserGroupRoleBindingCreate,
		Read:   resourceAlibabacloudStackAscmUserGroupRoleBindingRead,
		Update: resourceAlibabacloudStackAscmUserGroupRoleBindingUpdate,
		Delete: resourceAlibabacloudStackAscmUserGroupRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeInt},
			},
		},
	}
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	userGroupId := d.Get("user_group_id").(int)
	flag := false
	var roleids []int
	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandIntList(v.(*schema.Set).List())
	}
	log.Printf("roleids is %v", roleids)
	flag = true
	if flag {
		for i := range roleids {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "Ascm",
				"Action":           "AddRoleToUserGroup",
				"Version":          "2019-05-10",
				"ProductName":      "ascm",
				"userGroupId":      strconv.Itoa(userGroupId),
				"RoleId":           fmt.Sprint(roleids[i]),
				"SecurityToken":    client.Config.SecurityToken,
				"SignatureVersion": "1.0",
				"SignatureMethod":  "HMAC-SHA1",
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
			request.ApiName = "AddRoleToUserGroup"
			request.RegionId = client.RegionId
			request.Headers = map[string]string{"RegionId": client.RegionId}
			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			log.Printf("response of raw AddRoleToUserGroup Role(%d) is : %s", roleids[i], raw)

			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", raw)
			}

			addDebug("AddRoleToUserGroup", raw, requestInfo, request)

			bresponse, _ := raw.(*responses.CommonResponse)
			if bresponse.GetHttpStatus() != 200 {
				return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", AlibabacloudStackSdkGoERROR)
			}
			addDebug("AddRoleToUserGroup", raw, requestInfo, bresponse.GetHttpContentString())
			log.Printf("response of queryparams AddRoleToUserGroup is : %s", request.QueryParams)
		}

	}

	d.SetId(strconv.Itoa(userGroupId))

	return resourceAlibabacloudStackAscmUserGroupRoleBindingUpdate(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserGroupRoleBinding(d.Id())
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
	atoi, err := strconv.Atoi(d.Id())
	d.Set("user_group_id", atoi)

	return nil
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	var roleIdList []string

	if v, ok := d.GetOk("role_ids"); ok {
		roleids := expandIntList(v.(*schema.Set).List())

		for _, roleid := range roleids {
			roleIdList = append(roleIdList, strconv.Itoa(roleid))
		}
	}
	user_group_id := d.Get("user_group_id").(int)
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *ecs.Client
	request := requests.NewCommonRequest()
	if client.Config.Insecure {
		request.SetHTTPSInsecure(client.Config.Insecure)
	}

	request.Headers["x-ascm-product-name"] = "ascm"
	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"userGroupId":      strconv.Itoa(user_group_id),
		"roleIdList":       roleIdList,
		"SecurityToken":    client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}

	request.Method = "POST"
	request.Product = "Ascm"
	request.Version = "2019-05-10"
	request.ServiceCode = "ascm"
	request.Domain = client.Domain
	requeststring, _ := json.Marshal(QueryParams)

	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.Headers["Content-Type"] = requests.Json
	request.SetContent(requeststring)
	request.PathPattern = "/roa/ascm/auth/user/resetRolesForUserGroup"
	request.ApiName = "ResetRolesForUserGroup"
	request.RegionId = client.RegionId
	request.Headers["RegionId"] = client.RegionId

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	log.Printf("response of raw ResetRolesForUserGroup is : %s", raw)

	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserGroup", raw)
	}

	addDebug("ResetRolesForUserGroup", raw, requestInfo, request)
	return resourceAlibabacloudStackAscmUserGroupRoleBindingRead(d, meta)

}

func resourceAlibabacloudStackAscmUserGroupRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *ecs.Client
	var roleid int
	flag := false
	var roleids []int

	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandIntList(v.(*schema.Set).List())
		for i := range roleids {
			if len(roleids) > 1 {
				roleid = roleids[i]
				flag = true
			} else {
				roleid = roleids[0]
				flag = true
			}
		}
	}
	log.Printf("roleid is %v", roleid)
	log.Printf("roleids is %v", roleids)
	check, err := ascmService.DescribeAscmUserGroupRoleBinding(d.Id())
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "IsBindingExist", AlibabacloudStackSdkGoERROR)
	}

	addDebug("IsBindingExist", check, requestInfo, map[string]string{"userGroupId": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		if flag {
			request := requests.NewCommonRequest()
			if client.Config.Insecure {
				request.SetHTTPSInsecure(client.Config.Insecure)
			}
			request.QueryParams = map[string]string{
				"RegionId":         client.RegionId,
				"AccessKeySecret":  client.SecretKey,
				"Product":          "ascm",
				"Action":           "RemoveRoleFromUserGroup",
				"Version":          "2019-05-10",
				"ProductName":      "ascm",
				"userGroupId":      d.Id(),
				"RoleId":           fmt.Sprint(roleid),
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
			request.ApiName = "RemoveRoleFromUserGroup"
			request.Headers = map[string]string{"RegionId": client.RegionId}
			request.RegionId = client.RegionId

			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			if err != nil {
				return resource.RetryableError(err)
			}
			check, err = ascmService.DescribeAscmUserGroupRoleBinding(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
			addDebug("RemoveRoleFromUserGroup", raw, request)

		}
		return nil
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "RemoveRoleFromUserGroup", AlibabacloudStackSdkGoERROR)
	}
	return nil
}
