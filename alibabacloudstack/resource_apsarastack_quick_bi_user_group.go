package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackQuickBiUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackQuickBiUserGroupCreate,
		Read:   resourceAlibabacloudStackQuickBiUserGroupRead,
		Update: resourceAlibabacloudStackQuickBiUserGroupUpdate,
		Delete: resourceAlibabacloudStackQuickBiUserGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_group_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"user_group_description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"parent_user_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "-1",
			},
		},
	}
}

var UserGroupName string
var UserGroupDescription string
var ParentUserGroupId string
var UserGroupId string

func resourceAlibabacloudStackQuickBiUserGroupCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateUserGroup"
	request := make(map[string]interface{})

	if v, ok := d.GetOk("user_group_id"); ok {
		request["UserGroupId"] = v
	}

	UserGroupName = d.Get("user_group_name").(string)
	UserGroupDescription = d.Get("user_group_description").(string)
	ParentUserGroupId = d.Get("parent_user_group_id").(string)

	request["UserGroupName"] = UserGroupName
	request["UserGroupDescription"] = UserGroupDescription
	request["ParentUserGroupId"] = ParentUserGroupId

	response, err = client.DoTeaRequest("POST", "QuickBI", "2022-03-01", action, "", nil, request)
	if err != nil {
		return err
	}
	responseResult := response["Result"].(string)
	UserGroupId = responseResult
	d.SetId(responseResult)

	return resourceAlibabacloudStackQuickBiUserGroupRead(d, meta)
}

func resourceAlibabacloudStackQuickBiUserGroupRead(d *schema.ResourceData, meta interface{}) error {
	d.Set("user_group_id", UserGroupId)
	d.Set("user_group_name", UserGroupName)
	d.Set("user_group_description", UserGroupDescription)
	d.Set("parent_user_group_id", ParentUserGroupId)

	return nil
}

func resourceAlibabacloudStackQuickBiUserGroupUpdate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	update := false
	request := map[string]interface{}{
		"UserGroupId": d.Id(),
	}
	if d.HasChange("user_group_name") || d.IsNewResource() {
		update = true
	}

	if d.HasChange("user_group_description") {
		update = true
	}

	if update {
		action := "UpdateUserGroup"
		UserGroupName = d.Get("user_group_name").(string)
		UserGroupDescription = d.Get("user_group_description").(string)
		ParentUserGroupId = d.Get("parent_user_group_id").(string)

		request["UserGroupName"] = UserGroupName
		request["UserGroupDescription"] = UserGroupDescription
		request["ParentUserGroupId"] = ParentUserGroupId

		_, err = client.DoTeaRequest("POST", "QuickBI", "2022-03-01", action, "", nil, request)
		if err != nil {
			return err
		}
	}

	return resourceAlibabacloudStackQuickBiUserGroupRead(d, meta)
}

func resourceAlibabacloudStackQuickBiUserGroupDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteUserGroup"
	request := map[string]interface{}{
		"UserGroupId": d.Id(),
	}

	_, err = client.DoTeaRequest("POST", "QuickBI", "2022-03-01", action, "", nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return nil
		}
		return err
	}
	return nil
}
