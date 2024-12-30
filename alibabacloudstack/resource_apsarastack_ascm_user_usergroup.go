package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
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
				Type:     schema.TypeInt,
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
	userGroupId := d.Get("user_group_id").(int)
	var loginNamesList []string

	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())
		loginNamesList = append(loginNamesList, loginNames...)
	}

	body := map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": loginNamesList,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/ascm/auth/user/addUsersToUserGroup")
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	bresponse, err := client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
	}

	if bresponse.GetHttpStatus() != 200 {
		errmsg := ""
		if bresponse != nil {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("AddUsersToUserGroup", bresponse, nil, bresponse.GetHttpContentString())

	d.SetId(strconv.Itoa(userGroupId))

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
	userGroupId := d.Get("user_group_id").(int)
	body := map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": loginNames,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/ascm/auth/user/removeUsersFromUserGroup")
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	bresponse, err := client.ProcessCommonRequest(request)

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
	}

	var loginNamesList []string
	if v, ok := d.GetOk("login_names"); ok {
		loginNames := expandStringList(v.(*schema.Set).List())
		loginNamesList = append(loginNamesList, loginNames...)
	}

	body = map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": loginNamesList,
	}

	request = client.NewCommonRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/ascm/auth/user/addUsersToUserGroup")
	jsonData, err = json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	bresponse, err = client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
	}

	return resourceAlibabacloudStackAscmUserGroupUserRead(d, meta)
}

func resourceAlibabacloudStackAscmUserGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var loginNames []string
	userGroupId := d.Get("user_group_id").(int)
	if v, ok := d.GetOk("login_names"); ok {
		loginNames = expandStringList(v.(*schema.Set).List())
	}

	body := map[string]interface{}{
		"userGroupId":   userGroupId,
		"LoginNameList": loginNames,
	}

	request := client.NewCommonRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/ascm/auth/user/removeUsersFromUserGroup")
	jsonData, err := json.Marshal(body)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("Error marshaling to JSON: %v", err))
	}
	request.SetContentType(requests.Json)
	request.SetContent(jsonData)
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
	}

	return nil
}
