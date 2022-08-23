package apsarastack

import (
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceApsaraStackQuickBiUserGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceApsaraStackQuickBiUserGroupCreate,
		Read:   resourceApsaraStackQuickBiUserGroupRead,
		Update: resourceApsaraStackQuickBiUserGroupUpdate,
		Delete: resourceApsaraStackQuickBiUserGroupDelete,
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

func resourceApsaraStackQuickBiUserGroupCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
	action := "CreateUserGroup"
	request := make(map[string]interface{})
	conn, err := client.NewQuickbiClient()
	if err != nil {
		return WrapError(err)
	}
	if v, ok := d.GetOk("user_group_id"); ok {
		request["UserGroupId"] = v
	}

	UserGroupName = d.Get("user_group_name").(string)
	UserGroupDescription = d.Get("user_group_description").(string)
	ParentUserGroupId = d.Get("parent_user_group_id").(string)

	request["UserGroupName"] = UserGroupName
	request["UserGroupDescription"] = UserGroupDescription
	request["ParentUserGroupId"] = ParentUserGroupId

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	addDebug(action, response, request)
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "apsarastack_quick_bi_user", action, ApsaraStackSdkGoERROR)
	}
	responseResult := response["Result"].(string)
	UserGroupId = responseResult
	d.SetId(responseResult)

	return resourceApsaraStackQuickBiUserGroupRead(d, meta)
}
func resourceApsaraStackQuickBiUserGroupRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("user_group_id", UserGroupId)
	d.Set("user_group_name", UserGroupName)
	d.Set("user_group_description", UserGroupDescription)
	d.Set("parent_user_group_id", ParentUserGroupId)

	return nil
}
func resourceApsaraStackQuickBiUserGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	var response map[string]interface{}
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
		conn, err := client.NewQuickbiClient()
		if err != nil {
			return WrapError(err)
		}
		wait := incrementalWait(3*time.Second, 3*time.Second)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
			if err != nil {
				if NeedRetry(err) {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		addDebug(action, response, request)
		if err != nil {
			return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
		}
	}

	return resourceApsaraStackQuickBiUserGroupRead(d, meta)
}
func resourceApsaraStackQuickBiUserGroupDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	action := "DeleteUserGroup"
	var response map[string]interface{}
	conn, err := client.NewQuickbiClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"UserGroupId": d.Id(),
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{})
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		if IsExpectedErrors(err, []string{"User.Not.In.Organization"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, ApsaraStackSdkGoERROR)
	}
	return nil
}
