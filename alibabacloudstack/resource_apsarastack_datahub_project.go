package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDatahubProject() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(new, old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "project added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.EqualFold(new, old)
				},
			},
			"vpc_ids": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource,
		resourceAlibabacloudStackDatahubProjectCreate,
		resourceAlibabacloudStackDatahubProjectRead,
		resourceAlibabacloudStackDatahubProjectUpdate,
		resourceAlibabacloudStackDatahubProjectDelete)
	return resource
}

func resourceAlibabacloudStackDatahubProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	projectName := d.Get("name").(string)
	projectComment := d.Get("comment").(string)

	request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "CreateProject", "")
	request.QueryParams["ProjectName"] = projectName
	request.QueryParams["Comment"] = projectComment

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_project", request.GetActionName(), errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["ProjectComment"] = projectComment
		addDebug("CreateProject", response, requestMap)
	}

	d.SetId(strings.ToLower(projectName))
	return nil
}

func resourceAlibabacloudStackDatahubProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	if err := noUpdatesAllowedCheck(d, []string{"comment"}); err != nil {
		return err
	}

	client := meta.(*connectivity.AlibabacloudStackClient)

	projectName := d.Get("name").(string)

	if d.HasChange("vpc_ids") {
		o, n := d.GetChange("vpc_ids")
		oList := o.(*schema.Set).List()
		nList := n.(*schema.Set).List()
		oMap := map[string]interface{}{}
		for _, i := range oList {
			oMap[i.(string)] = nil
		}
		nMap := map[string]interface{}{}
		for _, i := range nList {
			nMap[i.(string)] = nil
		}
		for i := range oMap {
			if _, existed := nMap[i]; !existed {
				query := map[string]interface{}{
					"ProjectName": projectName,
					"VpcIds":      i,
				}
				if _, err := client.DoTeaRequest("POST", "datahub", "2019-11-20", "DeleteProjectVpcWhiteList", "", nil, query, nil); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteProjectVpcWhiteList", errmsgs.AlibabacloudStackSdkGoERROR)
				}
			}
		}
		for i := range nMap {
			if _, existed := oMap[i]; !existed {
				query := map[string]interface{}{
					"ProjectName": projectName,
					"VpcIds":      i,
				}
				if _, err := client.DoTeaRequest("POST", "datahub", "2019-11-20", "AddProjectVpcWhiteList", "", nil, query, nil); err != nil {
					return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "AddProjectVpcWhiteList", errmsgs.AlibabacloudStackSdkGoERROR)
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackDatahubProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubProject(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.SetId(strings.ToLower(d.Id()))
	d.Set("name", d.Id())
	d.Set("comment", object.Comment)
	d.Set("create_time", strconv.FormatInt(object.CreateTime, 10))
	d.Set("last_modify_time", strconv.FormatInt(object.LastModifyTime, 10))

	if resp, err := client.DoTeaRequest("GET", "datahub", "2019-11-20", "GetProjectVpcWhiteList", "", nil, map[string]interface{}{"ProjectName": d.Id()}, nil); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteProjectVpcWhiteList", errmsgs.AlibabacloudStackSdkGoERROR)
	} else {
		d.Set("vpc_ids", resp["VpcWhiteList"])
	}

	return nil
}

func resourceAlibabacloudStackDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	projectName := d.Id()

	// Step 1: Check if there are topics still under the project.
	// Terraform should NOT delete topics it doesn't manage.
	listTopicsReq := client.NewCommonRequest("GET", "datahub", "2019-11-20", "ListTopics", "")
	listTopicsReq.QueryParams["ProjectName"] = projectName
	listTopicsReq.QueryParams["PageNumber"] = "1"
	listTopicsReq.QueryParams["PageSize"] = "1000"

	listTopicsResp, err := client.ProcessCommonRequest(listTopicsReq)
	if err != nil {
		if !isDatahubNotExistError(err) {
			return errmsgs.WrapError(err)
		}
		// Project not found, nothing to delete
		return nil
	}

	listTopicResult := &ListTopicResult{}
	if err := json.Unmarshal(listTopicsResp.GetHttpContentBytes(), listTopicResult); err != nil {
		return errmsgs.WrapError(err)
	}

	if len(listTopicResult.List.Topic) > 0 {
		topicNames := make([]string, 0, len(listTopicResult.List.Topic))
		for _, t := range listTopicResult.List.Topic {
			topicNames = append(topicNames, t.TopicName)
		}
		return fmt.Errorf("cannot delete Datahub project %q: %d topic(s) still bound (%s). "+
			"Please remove topic resources managed by Terraform first", projectName, len(topicNames), strings.Join(topicNames, ", "))
	}

	// Step 2: No topics found, proceed to delete the project.
	// The backend may still have stale project-topic associations that haven't been
	// cleaned up yet, so DeleteProject can fail with a dependency error even though
	// ListTopics returned empty. Wait briefly before first attempt.
	time.Sleep(30 * time.Second)

	err = resource.Retry(30*time.Minute, func() *resource.RetryError {
		request := client.NewCommonRequest("GET", "datahub", "2019-11-20", "DeleteProject", "")
		request.QueryParams["ProjectName"] = projectName

		var requestInfo *datahub.DataHub
		response, err := client.ProcessCommonRequest(request)
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			addDebug("DeleteProject", response, requestInfo, requestMap)
		}
		if err != nil {
			if isDatahubNotExistError(err) {
				return nil
			}
			if response == nil {
				return resource.NonRetryableError(errmsgs.WrapErrorf(err, "Process Common Request Failed"))
			}
			// Only retry when the backend reports topic association not yet cleaned up:
			// ErrorCode: ApiProxyServiceError, Message: "Can not delete project ... which has: N topics."
			errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			if errmsgs.IsExpectedErrors(err, "ApiProxyServiceError") && regexp.MustCompile(`has: \d+ topics`).MatchString(errmsg) {
				return resource.RetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_project", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_project", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg))
		}
		return nil
	})

	if err != nil {
		return err
	}

	stateConf := BuildStateConf(
		[]string{"Deleting"},
		[]string{"Deleted"},
		d.Timeout(schema.TimeoutDelete),
		10*time.Second,
		datahubService.DatahubProjectStateRefreshFunc(d.Id()),
	)
	if _, err = stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}
	return nil
}
