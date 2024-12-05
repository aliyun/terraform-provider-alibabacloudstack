package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserCreate,
		Read:   resourceAlibabacloudStackAscmUserRead,
		Update: resourceAlibabacloudStackAscmUserUpdate,
		Delete: resourceAlibabacloudStackAscmUserDelete,
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
				Type:     schema.TypeString,
				Required: true,
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
}

func resourceAlibabacloudStackAscmUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	lname := d.Get("login_name").(string)
	dname := d.Get("display_name").(string)
	email := d.Get("email").(string)
	cellnum := d.Get("cellphone_number").(string)
	mobnationcode := d.Get("mobile_nation_code").(string)
	organizationid := d.Get("organization_id").(string)
	loginpolicyid := d.Get("login_policy_id").(int)

	check, err := ascmService.DescribeAscmDeletedUser(lname)
	if check.Data != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_resource_group", "\"Login Name already exist in Historical Users, try with a different name.\"", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if check.Data == nil {
		request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "AddUser", "")
		mergeMaps(request.QueryParams, map[string]string{
			"LoginName":        lname,
			"DisplayName":      dname,
			"CellphoneNum":     cellnum,
			"MobileNationCode": mobnationcode,
			"Email":            email,
			"OrganizationId":   organizationid,
			"LoginPolicyId":    fmt.Sprint(loginpolicyid),
		})
		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})
		log.Printf("response of raw AddUser is : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "AddUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("AddUser", raw, request)
		if bresponse.GetHttpStatus() != 200 {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "AddUser", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("AddUser", raw, request, bresponse.GetHttpContentString())
	}

	d.SetId(lname)
	init_password, err := ascmService.ExportInitPasswordByLoginName(lname)
	if err != nil {
		d.Set("init_password", init_password)
	}
	log.Printf("response of raw ExportInitPasswordByLoginName is : %s", init_password)

	return resourceAlibabacloudStackAscmUserUpdate(d, meta)
}

func resourceAlibabacloudStackAscmUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lname := d.Get("login_name").(string)
	organizationid := d.Get("organization_id").(string)
	var dname, cellnum, mobnationcode, email string
	var loginpolicyid int
	var roleIdList []string

	if v, ok := d.GetOk("role_ids"); ok {
		roleids := expandStringList(v.(*schema.Set).List())

		for _, roleid := range roleids {
			roleIdList = append(roleIdList, roleid)
		}
	}
	if d.HasChange("display_name") {
		dname = d.Get("display_name").(string)
	}
	if d.HasChange("cellphone_number") {
		cellnum = d.Get("cellphone_number").(string)
	}
	if d.HasChange("mobile_nation_code") {
		mobnationcode = d.Get("mobile_nation_code").(string)
	} else {
		mobnationcode = d.Get("mobile_nation_code").(string)
	}
	if d.HasChange("email") {
		email = d.Get("email").(string)
	}
	if d.HasChange("login_policy_id") {
		loginpolicyid = d.Get("login_policy_id").(int)
	}

	request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ModifyUserInformation", "")
	mergeMaps(request.QueryParams, map[string]string{
		"loginName":        lname,
		"display_name":     dname,
		"cellphone_num":    cellnum,
		"mobileNationCode": mobnationcode,
		"email":            email,
		"organization_id":  organizationid,
		"loginPolicyId":    fmt.Sprint(loginpolicyid),
		"policyId":         fmt.Sprint(loginpolicyid),
	})
	request.Headers["x-acs-content-type"] = "application/json"
	request.Headers["Content-Type"] = "application/json"

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	addDebug("ModifyUserInformation", raw, request, request.QueryParams)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ModifyUserInformationRequestFailed", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if !bresponse.IsSuccess() {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ModifyUserInformationFailed", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), bresponse)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	if len(d.Get("role_ids").(*schema.Set).List()) > 0 {
		request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "ResetRolesForUserByLoginName", "/roa/ascm/auth/user/ResetRolesForUserByLoginName")

		request.Headers["x-ascm-product-version"] = "2019-05-10"
		request.Headers["Content-Type"] = requests.Json

		queryParams := map[string]interface{}{
			"loginName":  lname,
			"roleIdList": roleIdList,
		}

		requeststring, err := json.Marshal(queryParams)
		request.SetContent(requeststring)

		raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
			return ecsClient.ProcessCommonRequest(request)
		})

		log.Printf("response of raw ResetRolesForUserByLoginName is : %s", raw)

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_user", "ResetRolesForUserByLoginName", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}

		addDebug("ResetRolesForUserByLoginName", raw, request)
	}
	return resourceAlibabacloudStackAscmUserRead(d, meta)
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
	d.Set("organization_id", strconv.Itoa(object.Data[0].Organization.ID))
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
	var requestInfo *ecs.Client
	check, err := ascmService.DescribeAscmUser(d.Id())
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "IsUserExist", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	addDebug("IsUserExist", check, requestInfo, map[string]string{"loginName": d.Id()})
	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("POST", "Ascm", "2019-05-10", "RemoveUserByLoginName", "")
		request.QueryParams["loginName"] = d.Id()

		request.Headers["x-acs-content-type"] = "application/json"
		request.Headers["Content-Type"] = "application/json"

		raw, err := client.WithEcsClient(func(csClient *ecs.Client) (interface{}, error) {
			return csClient.ProcessCommonRequest(request)
		})

		bresponse, ok := raw.(*responses.CommonResponse)
		if err != nil {
			errmsg := ""
			if ok {
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
