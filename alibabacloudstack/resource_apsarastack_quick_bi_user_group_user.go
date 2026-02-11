package alibabacloudstack

import (
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackQuickBiUserGroupUser() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"account_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackQuickBiUserGroupUserCreate,
		resourceAlibabacloudStackQuickBiUserGroupUserRead, nil, resourceAlibabacloudStackQuickBiUserGroupUserDelete)
	return resource
}
func resourceAlibabacloudStackQuickBiUserGroupUserCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "AddUserGroupMembers"
	request := map[string]interface{}{
		"UserId":       d.Get("account_id"),
		"UserGroupIds": d.Get("user_group_id"),
	}

	_, err = client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	d.SetId(fmt.Sprintf("%s:%s", d.Get("user_group_id"), d.Get("account_id")))

	return nil
}

func resourceAlibabacloudStackQuickBiUserGroupUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	quickbiPublicService := QuickbiPublicService{client}
	_, err:= quickbiPublicService.DescribeQuickBiUserGroupUser(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return  nil
		}
		return err
	}
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	d.Set("account_id", parts[1])
	d.Set("user_group_id", parts[0])

	return nil
}

func resourceAlibabacloudStackQuickBiUserGroupUserDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts, err := ParseResourceId(d.Id(), 2)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	action := "DeleteUserGroupMembers"
	request := map[string]interface{}{
		"UserId":       parts[1],
		"UserGroupIds": parts[0],
	}

	_, err = client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", action, "", nil, request, nil)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return nil
		}
		return err
	}
	return nil
}
