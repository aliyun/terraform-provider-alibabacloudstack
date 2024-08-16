package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
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
	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Product = "quickbi-public"
	request.Domain = client.Domain
	request.Version = "2022-03-01"
	request.ApiName = "CreateWorkSpace"
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"RegionId": client.RegionId,
		"Product":  "quickbi-public",
		"Version":  "2022-03-01",
		"Action":   action,
	}
	WorkspaceName = d.Get("workspace_name").(string)
	WorkspaceDesc = d.Get("workspace_desc").(string)
	UseComment = d.Get("use_comment").(bool)
	AllowShare = d.Get("allow_share").(bool)
	AllowPublish = d.Get("allow_publish").(bool)

	request.QueryParams["WorkspaceName"] = WorkspaceName
	request.QueryParams["WorkspaceDesc"] = WorkspaceDesc
	request.QueryParams["UseComment"] = fmt.Sprintf("%t", UseComment)
	request.QueryParams["AllowShare"] = fmt.Sprintf("%t", AllowShare)
	request.QueryParams["AllowPublish"] = fmt.Sprintf("%t", AllowPublish)
	endpoint := fmt.Sprintf("quickbi.%s.aliyuncs.com", client.Config.RegionId)
	request.Domain = endpoint
	var err error
	var popClient *sts.Client
	if client.Config.SecurityToken == "" {
		popClient, err = sts.NewClientWithAccessKey(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey)
	} else {
		popClient, err = sts.NewClientWithStsToken(client.Config.RegionId, client.Config.AccessKey, client.Config.SecretKey, client.Config.SecurityToken)
	}
	popClient.Domain = endpoint
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		raw, err := popClient.ProcessCommonRequest(request)
		addDebug(action, raw, request, request.QueryParams)
		if err != nil {
			if NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
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
		response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-03-01"), StringPointer("AK"), nil, request, &util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)})
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
