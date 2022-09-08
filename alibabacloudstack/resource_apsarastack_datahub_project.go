package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/aliyun-datahub-sdk-go/datahub"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDatahubProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackDatahubProjectCreate,
		Read:   resourceAlibabacloudStackDatahubProjectRead,
		Update: resourceAlibabacloudStackDatahubProjectUpdate,
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
				Default:      "project added by terraform",
				ValidateFunc: validation.StringLenBetween(0, 255),
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return strings.ToLower(new) == strings.ToLower(old)
				},
			},
			"create_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
			"last_modify_time": {
				Type:     schema.TypeString, //uint64 value from sdk
				Computed: true,
			},
		},
	}
}

func resourceAlibabacloudStackDatahubProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	//var dataService DatahubService = DatahubService{client: client}
	projectName := d.Get("name").(string)
	projectComment := d.Get("comment").(string)
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "datahub"
	request.Domain = client.Domain
	request.Version = "2019-11-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "CreateProject"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "datahub",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "CreateProject",
		"Version":         "2019-11-20",
		"ProjectName":     projectName,
		"Comment":         projectComment,
	}

	raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
		return dataHubClient.ProcessCommonRequest(request)
	})
	if err != nil {
		return WrapErrorf(err, DefaultErrorMsg, "alibabacloudstack_datahub_project", request.GetActionName(), AlibabacloudStackDatahubSdkGo)
	}
	if debugOn() {
		requestMap := make(map[string]string)
		requestMap["ProjectName"] = projectName
		requestMap["ProjectComment"] = projectComment
		addDebug("CreateProject", raw, requestMap)
	}

	d.SetId(strings.ToLower(projectName))
	return resourceAlibabacloudStackDatahubProjectRead(d, meta)
}

func resourceAlibabacloudStackDatahubProjectRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	datahubService := DatahubService{client}
	object, err := datahubService.DescribeDatahubProject(d.Id())
	if err != nil {
		if NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return WrapError(err)
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
			if err != nil {
				return WrapErrorf(err, DefaultErrorMsg, d.Id(), "UpdateProject", AlibabacloudStackDatahubSdkGo)
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
	request := requests.NewCommonRequest()
	request.Method = "GET"
	request.Product = "datahub"
	request.Domain = client.Domain
	request.Version = "2019-11-20"
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.ApiName = "DeleteProject"
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{
		"AccessKeySecret": client.SecretKey,
		"AccessKeyId":     client.AccessKey,
		"Product":         "datahub",
		"RegionId":        client.RegionId,
		"Department":      client.Department,
		"ResourceGroup":   client.ResourceGroup,
		"Action":          "DeleteProject",
		"Version":         "2019-11-20",
		"ProjectName":     projectName,
	}

	var requestInfo *datahub.DataHub
	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(dataHubClient *ecs.Client) (interface{}, error) {
			return dataHubClient.ProcessCommonRequest(request)
		})
		if err != nil {
			if isRetryableDatahubError(err) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		if debugOn() {
			requestMap := make(map[string]string)
			requestMap["ProjectName"] = projectName
			addDebug("DeleteProject", raw, requestInfo, requestMap)
		}
		return nil
	})
	if err != nil {
		if isDatahubNotExistError(err) {
			return nil
		}
		return WrapErrorf(err, DefaultErrorMsg, d.Id(), "DeleteProject", AlibabacloudStackDatahubSdkGo)
	}
	return WrapError(datahubService.WaitForDatahubProject(d.Id(), Deleted, DefaultTimeout))
}
