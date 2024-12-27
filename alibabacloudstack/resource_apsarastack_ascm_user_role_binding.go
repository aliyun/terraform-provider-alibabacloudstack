package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserRoleBinding() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserRoleBindingCreate,
		Read:   resourceAlibabacloudStackAscmUserRoleBindingRead,
		Update: resourceAlibabacloudStackAscmUserRoleBindingUpdate,
		Delete: resourceAlibabacloudStackAscmUserRoleBindingDelete,
		Schema: map[string]*schema.Schema{
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackAscmUserRoleBindingCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)
	flag := false
	var roleids []string
	if v, ok := d.GetOk("role_ids"); ok {
		roleids = expandStringList(v.(*schema.Set).List())
	}
	log.Printf("roleids is %v", roleids)
	flag = true
	if flag {
		for i := range roleids {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddRoleToUser", "/ascm/auth/role/addRoleToUser")
			request.QueryParams["loginName"] = lname
			request.QueryParams["roleId"] = fmt.Sprint(roleids[i])

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil || bresponse.GetHttpStatus() != 200 {
				errmsg := ""
				if bresponse != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_role_binding", "AddRoleToUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			addDebug("AddRoleToUser", bresponse, request, request.QueryParams)
			log.Printf("response of queryparams AddRoleToUser is : %s", request.QueryParams)
		}
	}

	d.SetId(lname)

	return resourceAlibabacloudStackAscmUserRoleBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserRoleBindingRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
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
	d.Set("login_name", object.Data[0].LoginName)

	return nil
}

func resourceAlibabacloudStackAscmUserRoleBindingUpdate(d *schema.ResourceData, meta interface{}) error {
	var roleIdList []string

	if v, ok := d.GetOk("role_ids"); ok {
		roleids := expandStringList(v.(*schema.Set).List())

		for _, roleid := range roleids {
			roleIdList = append(roleIdList, roleid)
		}
	}
	lname := d.Get("login_name").(string)
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetRolesForUserByLoginName", "/ascm/auth/user/ResetRolesForUserByLoginName")

	request.Headers["x-ascm-product-version"] = "2019-05-10"

	QueryParams := map[string]interface{}{
		"loginName":        lname,
		"roleIdList":       roleIdList,
		"SecurityToken":    client.Config.SecurityToken,
		"SignatureVersion": "1.0",
		"SignatureMethod":  "HMAC-SHA1",
	}

	requeststring, err := json.Marshal(QueryParams)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	request.SetContent(requeststring)
	request.Headers["Content-Type"] = requests.Json

	bresponse, err := client.ProcessCommonRequest(request)

	log.Printf("response of raw ResetRolesForUserByLoginName is : %s", bresponse)

	if err != nil {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	addDebug("ResetRolesForUserByLoginName", bresponse, request, request.QueryParams)
	return resourceAlibabacloudStackAscmUserRoleBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
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
	_, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		if flag {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUser", "/ascm/auth/role/removeRoleFromUser")
			request.QueryParams["loginName"] = d.Id()
			request.QueryParams["roleId"] = fmt.Sprint(roleid)

			bresponse, err := client.ProcessCommonRequest(request)
			if err != nil {
				errmsg := ""
				if bresponse != nil {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveRoleFromUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			_, err = ascmService.DescribeAscmUserRoleBinding(d.Id())

			if err != nil {
				return resource.NonRetryableError(err)
			}
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "RemoveRoleFromUser", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
