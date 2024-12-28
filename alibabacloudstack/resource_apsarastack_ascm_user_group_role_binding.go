package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
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
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRoleToUserGroup", "/ascm/auth/user/addRoleToUserGroup")
			mergeMaps(request.QueryParams, map[string]string{
				"ProductName":      "ascm",
				"userGroupId":      strconv.Itoa(userGroupId),
				"RoleId":           fmt.Sprint(roleids[i]),
				"SecurityToken":    client.Config.SecurityToken,
				"SignatureVersion": "1.0",
				"SignatureMethod":  "HMAC-SHA1",
			})
			bresponse, err := client.ProcessCommonRequest(request)
			log.Printf("response of raw AddRoleToUserGroup Role(%d) is : %s", roleids[i], bresponse)

			if err != nil {
				errmsg := ""
				if bresponse != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			addDebug("AddRoleToUserGroup", bresponse, request, request.QueryParams)
			if bresponse.GetHttpStatus() != 200 {
				errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_role_binding", "AddRoleToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}
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
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
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
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetRolesForUserGroup", "/ascm/auth/user/resetRolesForUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"userGroupId":      strconv.Itoa(user_group_id),
		"roleIdList":       roleIdList,
		"SecurityToken":    client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}

	requeststring, _ := json.Marshal(QueryParams)
	request.SetContent(requeststring)
	request.Headers["Content-Type"] = requests.Json

	bresponse, err := client.ProcessCommonRequest(request)
	log.Printf("response of raw ResetRolesForUserGroup is : %s", bresponse)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug("ResetRolesForUserGroup", bresponse, request)
	return resourceAlibabacloudStackAscmUserGroupRoleBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
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
	_, err := ascmService.DescribeAscmUserGroupRoleBinding(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		if flag {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUserGroup", "/ascm/auth/user/removeRoleFromUserGroup")
			mergeMaps(request.QueryParams, map[string]string{
				"ProductName": "ascm",
				"userGroupId": d.Id(),
				"roleId":      fmt.Sprint(roleid),
			})

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				errmsg := ""
				if bresponse != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveRoleFromUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			_, err = ascmService.DescribeAscmUserGroupRoleBinding(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
			addDebug("RemoveRoleFromUserGroup", bresponse, request)

		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRoleFromUserGroup", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
