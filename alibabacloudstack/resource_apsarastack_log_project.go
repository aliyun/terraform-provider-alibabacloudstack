package alibabacloudstack

import (
	"log"
	"os"
	"time"

	sls "github.com/aliyun/aliyun-log-go-sdk"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLogProject() *schema.Resource {
	resource := &schema.Resource{
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
	setResourceFunc(resource, resourceAlibabacloudStackLogProjectCreate,
		resourceAlibabacloudStackLogProjectRead, resourceAlibabacloudStackLogProjectUpdate, resourceAlibabacloudStackLogProjectDelete)
	return resource
}

func resourceAlibabacloudStackLogProjectCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	logService := LogService{client}
	name := d.Get("name").(string)
	description := d.Get("description").(string)

	// Try new API first (SLS 2019-10-23), fallback to old API (SLS 2020-03-31)
	var err error

	// Attempt 1: New API (2019-10-23)
	request := client.NewCommonRequest("POST", "Sls", "2019-10-23", "CreateProject", "")
	request.QueryParams["projectName"] = name
	request.QueryParams["description"] = description

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)

	// If new API fails, fallback to old API
	if err != nil && errmsgs.IsExpectedErrors(err, "InvalidVersion") {
		// TODO: Remove old API logic in version 3.20.0
		log.Printf("[WARN] SLS 2019-10-23 CreateProject failed: %v, fallback to 2020-03-31 API", err)

		// Attempt 2: Old API (2020-03-31) - will be removed in 3.20.0
		request = client.NewCommonRequest("POST", "SLS", "2020-03-31", "CreateProject", "")
		request.SetDomain(client.Config.Endpoints[connectivity.ASAPICode])
		request.QueryParams["projectName"] = name
		request.QueryParams["Description"] = description

		bresponse, err = client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
	}

	// Wait for project to be created
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
	return nil
}

func resourceAlibabacloudStackLogProjectRead(d *schema.ResourceData, meta interface{}) error {
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
		request.SetDomain(os.Getenv("ALIBABACLOUDSTACK_ASAPI_ENDPOINT"))
		request.QueryParams["ProjectName"] = name
		request.QueryParams["description"] = d.Get("description").(string)

		bresponse, err := client.ProcessCommonRequest(request)
		addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug("UpdateProject", bresponse, requestInfo, request)
	}

	return nil
}

func resourceAlibabacloudStackLogProjectDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var requestInfo *sls.Client
	name := d.Get("name").(string)
	request := client.NewCommonRequest("POST", "SLS", "2020-03-31", "DeleteProject", "")
	request.SetDomain(os.Getenv("ALIBABACLOUDSTACK_ASAPI_ENDPOINT"))
	request.QueryParams["ProjectName"] = name

	bresponse, err := client.ProcessCommonRequest(request)
	addDebug(request.GetActionName(), bresponse, request, request.QueryParams)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, "ProjectNotExist") {
			return nil
		}
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("DeleteProject", bresponse, requestInfo, request)
	return nil
}
