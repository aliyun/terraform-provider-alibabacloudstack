package alibabacloudstack

import (
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogProject() *schema.Resource {
	return &schema.Resource{
		Create: resourceAlibabacloudStackLogProjectCreate,
		Read:   resourceAlibabacloudStackLogProjectRead,
		Update: resourceAlibabacloudStackLogProjectUpdate,
		Delete: resourceAlibabacloudStackLogProjectDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceAlibabacloudStackLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	name := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "SLS", "2020-03-31", "CreateProject", "")
	request.QueryParams["projectName"] = name
	request.QueryParams["Description"] = d.Get("description").(string)

	raw, err := client.WithEcsClient(func(alidnsClient *ecs.Client) (interface{}, error) {
		return alidnsClient.ProcessCommonRequest(request)
	})
	bresponse, ok := raw.(*responses.BaseResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_log_project", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("LogProject", raw)

	err = resource.Retry(2*time.Minute, func() *resource.RetryError {
		object, err := logService.DescribeLogProject(name)
		if err != nil {
			return resource.NonRetryableError(err)
		}
		if object.ProjectName != "" {
			return nil
		}
		return resource.RetryableError(errmsgs.Error("Failed to describe log project"))
	})
	d.SetId(name)
	return resourceAlibabacloudStackLogProjectRead(d, meta)
}

func resourceAlibabacloudStackLogProjectRead(d *schema.ResourceData, meta interface{}) error {
	waitSecondsIfWithTest(1)
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	object, err := logService.DescribeLogProject(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("name", object.ProjectName)
	d.Set("description", object.Description)

	return nil
}

func resourceAlibabacloudStackLogProjectUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *sls.Client

	name := d.Id()
	if d.HasChange("description") {
		request := client.NewCommonRequest("POST", "SLS", "2020-03-31", "UpdateProject", "")
		request.QueryParams["ProjectName"] = name
		request.QueryParams["description"] = d.Get("description").(string)

		raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.BaseResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "UpdateProject", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg)
		}
		addDebug("UpdateProject", raw, requestInfo, request)
	}

	return resourceAlibabacloudStackLogProjectRead(d, meta)
}

func resourceAlibabacloudStackLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *sls.Client
	name := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "SLS", "2020-03-31", "DeleteProject", "")
	request.QueryParams["ProjectName"] = name

	err := resource.Retry(3*time.Minute, func() *resource.RetryError {
		raw, err := client.WithEcsClient(func(slsClient *ecs.Client) (interface{}, error) {
			return slsClient.ProcessCommonRequest(request)
		})
		bresponse, ok := raw.(*responses.BaseResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.LogClientTimeout, "RequestTimeout"}) {
				return resource.RetryableError(err)
			}
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(bresponse)
			}
			return resource.NonRetryableError(errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "DeleteProject", errmsgs.AlibabacloudStackLogGoSdkERROR, errmsg))
		}
		addDebug("DeleteProject", raw, requestInfo, request)
		return nil
	})
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
			return nil
		}
		return err
	}

	return nil
}
