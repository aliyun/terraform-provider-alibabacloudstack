package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

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
		resourceAlibabacloudStackQuickBiWorkspaceRead,
		resourceAlibabacloudStackQuickBiWorkspaceUpdate,
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

	var err error
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		bresponse, err := client.ProcessCommonRequest(request)
		log.Printf(" response of raw CreateWorkSpace : %v", bresponse)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if bresponse == nil {
			return resource.RetryableError(fmt.Errorf("received nil response from API"))
		}

		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		contentBytes := bresponse.GetHttpContentBytes()
		if contentBytes == nil || len(contentBytes) == 0 {
			return resource.RetryableError(fmt.Errorf("received empty response content from API"))
		}

		err = json.Unmarshal(contentBytes, &response)
		if err != nil {
			return resource.RetryableError(fmt.Errorf("error unmarshalling response: %v", err))
		}
		if result, exists := response["Result"]; !exists || result == nil {
			return resource.RetryableError(fmt.Errorf("API response does not contain Result field: %v", response))
		}

		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	resultValue, exists := response["Result"]
	if !exists || resultValue == nil {
		return errmsgs.WrapErrorf(fmt.Errorf("API response missing Result field"),
			errmsgs.DefaultErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	workspaceId, ok := resultValue.(string)
	if !ok {
		return errmsgs.WrapErrorf(fmt.Errorf("Result field is not a string: %T, value: %v", resultValue, resultValue),
			errmsgs.DefaultErrorMsg, "alibabacloudstack_quick_bi_Workspace", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}

	WorkspaceId = workspaceId
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
	noUpdateAllowedFields := []string{"workspace_name", "workspace_desc", "use_comment", "allow_share", "allow_publish"}
	return noUpdatesAllowedCheck(d, noUpdateAllowedFields)
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
