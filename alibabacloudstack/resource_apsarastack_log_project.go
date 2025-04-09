package alibabacloudstack

import (
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
	request := client.NewCommonRequest("POST", "SLS", "2020-03-31", "CreateProject", "")
	request.QueryParams["projectName"] = name
	request.QueryParams["Description"] = d.Get("description").(string)

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug("LogProject", bresponse)

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

		bresponse, err := client.ProcessCommonRequest(request)
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
	request.QueryParams["ProjectName"] = name

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if errmsgs.IsExpectedErrors(err, []string{"ProjectNotExist"}) {
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
