package alibabacloudstack

import (
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackQuickBiWorkspace() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackQuickBiWorkspaceCreate,
		Read:   resourceAlibabacloudStackQuickBiWorkspaceRead,
		Update: resourceAlibabacloudStackQuickBiWorkspaceUpdate,
		Delete: resourceAlibabacloudStackQuickBiWorkspaceDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"workspace_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"use_comment": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_share": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"allow_publish": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

var WorkspaceId string
var WorkspaceName string
var WorkspaceDesc string
var UseComment bool
var AllowShare bool
var AllowPublish bool

func resourceAlibabacloudStackQuickBiWorkspaceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var response map[string]interface{}
	action := "CreateWorkSpace"
	request := make(map[string]interface{})
	conn, err := client.NewQuickbiClient()
	if err != nil {
		return WrapError(err)
	}

	WorkspaceName = d.Get("workspace_name").(string)
	WorkspaceDesc = d.Get("workspace_desc").(string)
	UseComment = d.Get("use_comment").(bool)
	AllowShare = d.Get("allow_share").(bool)
	AllowPublish = d.Get("allow_publish").(bool)

	request["WorkspaceName"] = WorkspaceName
	request["WorkspaceDesc"] = WorkspaceDesc
	request["UseComment"] = UseComment
	request["AllowShare"] = AllowShare
	request["AllowPublish"] = AllowPublish
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
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, AlibabacloudStackSdkGoERROR)
	}
	WorkspaceId = response["Result"].(string)
	d.SetId(fmt.Sprint(WorkspaceId))

	return resourceAlibabacloudStackQuickBiWorkspaceRead(d, meta)
}
func resourceAlibabacloudStackQuickBiWorkspaceRead(d *schema.ResourceData, meta interface{}) error {

	d.Set("workspace_id", WorkspaceId)
	d.Set("workspace_name", WorkspaceName)
	d.Set("workspace_desc", WorkspaceDesc)
	d.Set("use_comment", UseComment)
	d.Set("allow_share", AllowShare)
	d.Set("allow_publish", AllowPublish)

	return nil
}
func resourceAlibabacloudStackQuickBiWorkspaceUpdate(d *schema.ResourceData, meta interface{}) error {
	return resourceAlibabacloudStackQuickBiWorkspaceRead(d, meta)
}
func resourceAlibabacloudStackQuickBiWorkspaceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "DeleteWorkSpace"
	var response map[string]interface{}
	conn, err := client.NewQuickbiClient()
	if err != nil {
		return WrapError(err)
	}
	request := map[string]interface{}{
		"WorkspaceId": d.Id(),
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
		if IsExpectedErrors(err, []string{"Workspace.Not.In.Organization"}) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), action, AlibabacloudStackSdkGoERROR)
	}
	return nil
}
