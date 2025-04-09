package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/sts"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackQuickBiWorkspace() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackQuickBiWorkspaceCreate,
		resourceAlibabacloudStackQuickBiWorkspaceRead, nil, 
		resourceAlibabacloudStackQuickBiWorkspaceDelete)
	return resource
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
	WorkspaceName = d.Get("workspace_name").(string)
	WorkspaceDesc = d.Get("workspace_desc").(string)
	UseComment = d.Get("use_comment").(bool)
	AllowShare = d.Get("allow_share").(bool)
	AllowPublish = d.Get("allow_publish").(bool)

	request := client.NewCommonRequest("POST", "quickbi-public", "2022-03-01", "CreateWorkSpace", "")
	request.QueryParams["WorkspaceName"] = WorkspaceName
	request.QueryParams["WorkspaceDesc"] = WorkspaceDesc
	request.QueryParams["UseComment"] = fmt.Sprintf("%t", UseComment)
	request.QueryParams["AllowShare"] = fmt.Sprintf("%t", AllowShare)
	request.QueryParams["AllowPublish"] = fmt.Sprintf("%t", AllowPublish)
	endpoint := fmt.Sprintf("quickbi.%s.aliyuncs.com", client.Config.RegionId)

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
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := ""
			if raw != nil {
				errmsg = errmsgs.GetBaseResponseErrorMessage(raw.BaseResponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		err = json.Unmarshal(raw.GetHttpContentBytes(), &response)
		return nil
	})

	addDebug(action, response, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	WorkspaceId = response["Result"].(string)
	d.SetId(fmt.Sprint(WorkspaceId))

	return nil
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

func resourceAlibabacloudStackQuickBiWorkspaceDelete(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["WorkspaceId"] = d.Id()

	_, err = client.DoTeaRequest("POST", "quickbi-public", "2022-03-01", "DeleteWorkSpace", "", nil, nil, request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"Workspace.Not.In.Organization"}) {
			return nil
		}
		return err
	}
	return nil
}
