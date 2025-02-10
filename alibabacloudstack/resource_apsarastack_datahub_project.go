package alibabacloudstack

import (
	"strconv"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDatahubProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDatahubProjectCreate,
		Read:   resourceAlibabacloudStackDatahubProjectRead,
		//Update: resourceAlibabacloudStackDatahubProjectUpdate,
		Delete: resourceAlibabacloudStackDatahubProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringLenBetween(3, 32),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"comment": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true, // 当前专有云不支持修改comment
				Default:      "project added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
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
	return resourceAlibabacloudStackDatahubProjectRead(d, meta)
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
	return nil
}

func resourceAlibabacloudStackDatahubProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	/*
		client := meta.(*connectivity.AlibabacloudStackClient)

		if d.HasChange("comment") {

			projectName := d.Id()
			projectComment := d.Get("comment").(string)

			var requestInfo *datahub.DataHub

			raw, err := client.WithDataHubClient(func(dataHubClient datahub.DataHubApi) (interface{}, error) {
				requestInfo = dataHubClient.(*datahub.DataHub)
				return dataHubClient.UpdateProject(projectName, projectComment)
			})
			response, ok := raw.(*datahub.UpdateProjectResponse)
			if err != nil {
				errmsg := ""
				if ok {
					errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
				}
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateProject", errmsgs.AlibabacloudStackDatahubSdkGo, errmsg)
			}
			if debugOn() {
				requestMap := make(map[string]string)
				requestMap["ProjectName"] = projectName
				requestMap["ProjectComment"] = projectComment
				addDebug("UpdateProject", raw, requestInfo, requestMap)
			}
		}

	*/

	return resourceAlibabacloudStackDatahubProjectRead(d, meta)
}

func resourceAlibabacloudStackDatahubProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	projectName := d.Id()

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
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_datahub_project", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	return errmsgs.WrapError(datahubService.WaitForDatahubProject(d.Id(), Deleted, DefaultTimeout))
}
