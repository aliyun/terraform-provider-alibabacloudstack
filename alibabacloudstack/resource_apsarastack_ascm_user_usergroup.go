package alibabacloudstack

import (
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackAscmUserGroupUser() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackAscmUserGroupUserCreate, resourceAlibabacloudStackAscmUserGroupUserRead, resourceAlibabacloudStackAscmUserGroupUserUpdate, resourceAlibabacloudStackAscmUserGroupUserDelete)
	return resource
}

func resourceAlibabacloudStackAscmUserGroupUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	userGroupId := d.Get("user_group_id").(int)
	var loginNames []string

	if v, ok := d.GetOk("login_names"); ok {
		loginNames = expandStringList(v.(*schema.Set).List())
	}

	body := map[string]interface{}{
		"userGroupId":   userGroupId,
		"loginNameList": loginNames,
	}
	addresponse, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/ascm/auth/user/addUsersToUserGroup", nil, nil, body)

	if err != nil {
		errmsg := errmsgs.GetAsapiErrorMessage(addresponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
	}

	d.SetId(strconv.Itoa(userGroupId))

	return nil
}

func resourceAlibabacloudStackAscmUserGroupUserRead(d *schema.ResourceData, meta interface{}) error {

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
	user_group_id, _ := strconv.Atoi(d.Id())
	d.Set("user_group_id", user_group_id)
	d.Set("login_names", loginNames)

	return nil
}

func resourceAlibabacloudStackAscmUserGroupUserUpdate(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("login_names") {
		o, n := d.GetChange("login_names")
		old, new := o.(*schema.Set).List(), n.(*schema.Set).List()
		removeLoginNames := expandStringList(old)
		body := map[string]interface{}{
			"userGroupId":   d.Id(),
			"loginNameList": removeLoginNames,
		}

		response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/ascm/auth/user/removeUsersFromUserGroup", nil, nil, body)

		if err != nil {
			errmsg := errmsgs.GetAsapiErrorMessage(response)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
		}
		addLoginNames := expandStringList(new)
		addbody := map[string]interface{}{
			"userGroupId":   d.Id(),
			"loginNameList": addLoginNames,
		}
		addresponse, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "AddUsersToUserGroup", "/ascm/auth/user/addUsersToUserGroup", nil, nil, addbody)

		if err != nil {
			errmsg := errmsgs.GetAsapiErrorMessage(addresponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "AddUsersToUserGroup", errmsg)
		}

	}

	return nil
}

func resourceAlibabacloudStackAscmUserGroupUserDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var loginNames []string
	if v, ok := d.GetOk("login_names"); ok {
		loginNames = expandStringList(v.(*schema.Set).List())
	}
	body := map[string]interface{}{
		"userGroupId":   d.Id(),
		"loginNameList": loginNames,
	}

	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "RemoveUsersFromUserGroup", "/ascm/auth/user/removeUsersFromUserGroup", nil, nil, body)

	if err != nil {
		errmsg := errmsgs.GetAsapiErrorMessage(response)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_ascm_usergroup_user", "RemoveUsersFromUserGroup", errmsg)
	}

	return nil
}
