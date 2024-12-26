package alibabacloudstack

import (
	"encoding/json"
	"log"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroupUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackAscmUserGroupUserCreate,
		Read:   resourceAlibabacloudStackAscmUserGroupUserRead,
		Update: resourceAlibabacloudStackAscmUserGroupUserUpdate,
		Delete: resourceAlibabacloudStackAscmUserGroupUserDelete,
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"login_names": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceAlibabacloudStackAscmUserGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	userGroupId := d.Get("user_group_id").(string)
	var loginNamesList []string

	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())
		loginNamesList = append(loginNamesList, loginNames...)
	}

	queryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNamesList,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/roa/ascm/auth/user/addUsersToUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"
	request.Headers["Content-Type"] = requests.Json
	requeststring, err := json.Marshal(queryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", "")
	}
	request.SetContent(requeststring)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddUsersToUserGroup is : %s", raw)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
	}

	addDebug("AddUsersToUserGroup", raw, nil, request)

	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddUsersToUserGroup", raw, nil, bresponse.GetHttpContentString())

	d.SetId(userGroupId)

	return resourceAlibabacloudStackAscmUserGroupUserRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUsergroupUser(d.Id())
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

	var loginNames []string
	for _, data := range object.Data {
		loginNames = append(loginNames, data.LoginName)
	}

	d.Set("login_names", loginNames)

	return nil
}

func resourceAlibabacloudStackAscmUserGroupUserUpdate(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)

	client := meta.(*connectivity.AlibabacloudStackClient)
	ascmService := AscmService{client}
	object, err := ascmService.DescribeAscmUsergroupUser(d.Id())
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

	var loginNames []string
	for _, data := range object.Data {
		loginNames = append(loginNames, data.LoginName)
	}
	userGroupId := d.Get("user_group_id").(string)

	queryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNames,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/roa/ascm/auth/user/RemoveUsersFromUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"
	request.Headers["Content-Type"] = requests.Json
	requeststring, err := json.Marshal(queryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", "")
	}
	request.SetContent(requeststring)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw RemoveUsersFromUserGroup is : %s", raw)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
	}

	var loginNamesList []string
	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())
		loginNamesList = append(loginNamesList, loginNames...)
	}

	queryParams = map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNamesList,
	}

	request = client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/roa/ascm/auth/user/addUsersToUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"
	request.Headers["Content-Type"] = requests.Json
	requeststring, err = json.Marshal(queryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", "")
	}
	request.SetContent(requeststring)

	raw, err = client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw AddUsersToUserGroup is : %s", raw)

	bresponse, ok = raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
	}

	return resourceAlibabacloudStackAscmUserGroupUserRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var loginNames []string
	userGroupId := d.Get("user_group_id").(string)
	if v, ok := d.GetOk("login_names"); ok {
		loginNames = expandStringList(v.(*schema.Set).List())
	}

	queryParams := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNames,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/roa/ascm/auth/user/RemoveUsersFromUserGroup")

	request.Headers["x-ascm-product-version"] = "2019-05-10"
	request.Headers["Content-Type"] = requests.Json
	requeststring, err := json.Marshal(queryParams)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", "")
	}
	request.SetContent(requeststring)

	raw, err := client.WithEcsClient(func(ecsClient *ecs.Client) (interface{}, error) {
		return ecsClient.ProcessCommonRequest(request)
	})
	log.Printf("response of raw RemoveUsersFromUserGroup is : %s", raw)

	bresponse, ok := raw.(*responses.CommonResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
	}

	return nil
}
