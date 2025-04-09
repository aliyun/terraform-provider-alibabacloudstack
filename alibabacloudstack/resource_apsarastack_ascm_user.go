package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUser() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cellphone_number": {
				Type:     schema.TypeString,
				Required: true,
			},
			"telephone_number": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mobile_nation_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"organization_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Computed:   true,
				Deprecated: "Field 'organization_id' has been deprecated from provider version 1.0.32. Use the organization to which the current user belongs",
			},
			"user_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"login_policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"init_password": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"role_ids": {
				Type:     schema.TypeSet,
				Computed: true,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackAscmUserCreate, resourceAlibabacloudStackAscmUserRead, resourceAlibabacloudStackAscmUserUpdate, resourceAlibabacloudStackAscmUserDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	lname := d.Get("login_name").(string)
	dname := d.Get("display_name").(string)
	email := d.Get("email").(string)
	cellnum := d.Get("cellphone_number").(string)
	mobnationcode := d.Get("mobile_nation_code").(string)
	loginpolicyid := d.Get("login_policy_id").(int)

	check, err := ascmService.DescribeAscmDeletedUser(lname)
	if check.Data != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "\"Login Name already exist in Historical Users, try with a different name.\"", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if check.Data == nil {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddUser", "/ascm/auth/user/addUser")
		mergeMaps(request.QueryParams, map[string]string{
			"loginName":        lname,
			"displayName":      dname,
			"cellphoneNum":     cellnum,
			"mobileNationCode": mobnationcode,
			"email":            email,
			"organizationId":   client.Department,
			"loginPolicyId":    fmt.Sprint(loginpolicyid),
		})
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("AddUser", bresponse, request, request.QueryParams)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			} else {
				return err
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "AddUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			} else {
				return err
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "AddUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	d.SetId(lname)
	init_password, err := ascmService.ExportInitPasswordByLoginName(lname)
	if err != nil {
		d.Set("init_password", init_password)
	}
	log.Printf("response of bresponse ExportInitPasswordByLoginName is : %s", init_password)

	return nil
}

func resourceAlibabacloudStackAscmUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)
	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ModifyUserInformation", "/ascm/auth/user/modifyUserInformation")
	update := false
	if d.HasChange("display_name") {
		update = true
		request.QueryParams["displayName"] = d.Get("display_name").(string)
	}
	if d.HasChange("cellphone_number") {
		update = true
		request.QueryParams["cellphoneNum"] = d.Get("cellphone_number").(string)
	}
	if d.HasChange("mobile_nation_code") {
		update = true
		request.QueryParams["mobileNationCode"] = d.Get("mobile_nation_code").(string)
	}

	if d.HasChange("email") {
		update = true
		request.QueryParams["email"] = d.Get("email").(string)
	}
	if d.HasChange("login_policy_id") {
		update = true
		request.QueryParams["loginPolicyId"] = fmt.Sprint(d.Get("login_policy_id").(int))
		request.QueryParams["policyId"] = fmt.Sprint(d.Get("login_policy_id").(int))
	}
	if update {
		request.QueryParams["loginName"] = lname
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"
		bresponse, err := client.ProcessCommonRequest(request)
		addDebug("ModifyUserInformation", bresponse, request, request.QueryParams)

		if err != nil || !bresponse.IsSuccess() {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			} else {
				return err
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ModifyUserInformationRequestFailed", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
		if err != nil {
			return errmsgs.WrapError(err)
		}
	}
	if len(d.Get("role_ids").(*schema.Set).List()) > 0 {
		role_ids := d.Get("role_ids").(*schema.Set).List()
		roleIdList := make([]string, 0)
		for _, role_id := range role_ids {
			roleIdList = append(roleIdList, role_id.(string))
		}
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "ResetRolesForUserByLoginName", "/ascm/auth/role/resetRolesForUserByLoginName")
		requeststring, err := json.Marshal(roleIdList)
		request.QueryParams["loginName"] = lname
		request.QueryParams["roleIdList"] = fmt.Sprint(requeststring)
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"
		bresponse, err := client.ProcessCommonRequest(request)

		log.Printf("response of bresponse ResetRolesForUserByLoginName is : %s", bresponse)
		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			} else {
				return err
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("ResetRolesForUserByLoginName", bresponse, request)
	}
	return nil
}

func resourceAlibabacloudStackAscmUserRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUser(d.Id())
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

	d.Set("user_id", object.Data[0].ID)
	d.Set("login_name", object.Data[0].LoginName)
	d.Set("display_name", object.Data[0].DisplayName)
	d.Set("email", object.Data[0].Email)
	d.Set("mobile_nation_code", object.Data[0].MobileNationCode)
	d.Set("cellphone_number", object.Data[0].CellphoneNum)
	d.Set("organization_id", client.Department)
	d.Set("login_policy_id", object.Data[0].LoginPolicy.ID)
	var user_roles []string
	for _, role := range object.Data[0].UserRoles {
		user_roles = append(user_roles, strconv.Itoa(role.ID))
	}
	d.Set("role_ids", user_roles)
	init_password, _ := ascmService.ExportInitPasswordByLoginName(object.Data[0].LoginName)
	if init_password != "" {
		d.Set("init_password", init_password)
	}
	log.Printf("Ascm User: %s init_password  : %s", object.Data[0].LoginName, init_password)

	return nil
}

func resourceAlibabacloudStackAscmUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	var requestInfo *sdk.Client
	check, err := ascmService.DescribeAscmUser(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsUserExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsUserExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveUserByLoginName", "/ascm/auth/user/removeUserByLoginName")
		request.QueryParams["loginName"] = d.Id()

		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"
		bresponse, err := client.ProcessCommonRequest(request)

		if err != nil {
			errmsg := ""
			if bresponse != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "RemoveUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		check, err = ascmService.DescribeAscmUser(d.Id())

		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return nil
}
