package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
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
	var requestInfo *ecs.Client
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

			raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
				return ecsClient.ProcessCommonRequest(request)
			})
			log.Printf("response of raw AddRoleToUser Role(%s) is : %s", roleids[i], raw)

			bresponse, ok := raw.(*responses.CommonResponse)
			if err != nil || bresponse.GetHttpStatus() != 200 {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user_role_binding", "AddRoleToUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			}

			addDebug("AddRoleToUser", raw, requestInfo, bresponse.GetHttpContentString())
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
	var requestInfo *ecs.Client
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetRolesForUserByLoginName", "/roa/ascm/auth/user/ResetRolesForUserByLoginName")

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

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})

	log.Printf("response of raw ResetRolesForUserByLoginName is : %s", raw)

	if err != nil {
		errmsg := ""
		if bresponse, ok := raw.(*responses.CommonResponse); ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	bresponse, ok := raw.(*responses.CommonResponse)
	if !ok {
		return fmt.Errorf("Failed to cast response to CommonResponse")
	}
	addDebug("ResetRolesForUserByLoginName", raw, requestInfo, bresponse.GetHttpContentString())
	return resourceAlibabacloudStackAscmUserRoleBindingRead(d, meta)
}

func resourceAlibabacloudStackAscmUserRoleBindingDelete(d *schema.ResourceData, meta interface{}) error {
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
	check, err := ascmService.DescribeAscmUserRoleBinding(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsBindingExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsBindingExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		if flag {
			request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveRoleFromUser", "/ascm/auth/role/removeRoleFromUser")
			request.QueryParams["loginName"] = d.Id()
			request.QueryParams["roleId"] = fmt.Sprint(roleid)

			raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
				return csClient.ProcessCommonRequest(request)
			})
			if err != nil {
				errmsg := ""
				if bresponse, ok := raw.(*responses.CommonResponse); ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
				}
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "RemoveRoleFromUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			check, err = ascmService.DescribeAscmUserRoleBinding(d.Id())

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
