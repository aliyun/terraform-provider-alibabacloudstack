package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackSchedulerx2Workflow() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"namespace": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "system_namespace",
			},
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"time_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"cron", "api"}, false),
			},
			"time_expression": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"time_zone": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{"PRC", "Hongkong", "Japan", "Singapore", "GTM", "GTM-0", "GTM-1", "GTM-2", "GTM-3", "GTM-4", "GTM-5", "GTM-6", "GTM-7", "GTM-8", "GTM-9",
					"GTM-10", "GTM-11", "GTM-12", "GTM+1", "GTM+2", "GTM+3", "GTM+4", "GTM+5", "GTM+6", "GTM+7", "GTM+8", "GTM+9", "GTM+10", "GTM+11", "GTM+12", "GTM+13", "GTM+14"}, false),
			},
			"max_concurrency": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"workflow_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackSchedulerx2WorkflowCreate, resourceAlibabacloudStackSchedulerx2WorkflowRead, resourceAlibabacloudStackSchedulerx2WorkflowUpdate, resourceAlibabacloudStackSchedulerx2WorkflowDelete)
	return resource
}

func resourceAlibabacloudStackSchedulerx2WorkflowCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	reqBody := map[string]interface{}{
		"Namespace":      d.Get("namespace"),
		"GroupId":        d.Get("group_id"),
		"Name":           d.Get("name"),
		"Description":    d.Get("description"),
		"AcceptLanguage": "en",
	}
	time_type := d.Get("time_type").(string)
	switch time_type {
	case "cron":
		reqBody["TimeType"] = "1"
	case "api":
		reqBody["TimeType"] = "100"
	}
	if v, ok := d.GetOk("time_expression"); ok {
		reqBody["TimeExpression"] = v.(string)
	}
	if v, ok := d.GetOk("time_zone"); ok {
		reqBody["TimeZone"] = v.(string)
	}
	if v, ok := d.GetOk("max_concurrency"); ok {
		reqBody["MaxConcurrency"] = v.(int)
	}

	resp, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "CreateWorkflow", "", nil, nil, reqBody)
	if err != nil {
		return err
	}

	if fmt.Sprint(resp["Code"]) != "200" {
		return errmsgs.WrapError(fmt.Errorf("create workflow failed: %v", resp))
	}

	id, err := jsonpath.Get("$.Data.WorkflowId", resp)
	if err != nil {
		return errmsgs.WrapError(fmt.Errorf("failed to get WorkflowId from response: %v", resp))
	}

	d.SetId(fmt.Sprint(id))

	return nil
}

func resourceAlibabacloudStackSchedulerx2WorkflowRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	schedulerx2Service := Schedulerx2Service{client}

	object, err := schedulerx2Service.DescribeSchedulerx2Workflow(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("workflow_id", object["WorkflowId"])
	d.Set("group_id", object["GroupId"])
	d.Set("name", object["Name"])
	d.Set("description", object["Description"])
	d.Set("time_type", object["TimeType"])
	timeExpression := ""
	if object["TimeExpression"].(string) != "/" {
		timeExpression = object["TimeExpression"].(string)
	}
	d.Set("time_expression", timeExpression)
	d.Set("max_concurrency", object["MaxConcurrency"])
	if object["Status"].(string) == "enable" {
		d.Set("enabled", true)
	} else {
		d.Set("enabled", false)
	}
	return nil
}

func resourceAlibabacloudStackSchedulerx2WorkflowUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// Handle enable/disable status change
	if d.HasChange("enabled") {
		enabled := d.Get("enabled").(bool)
		var apiName string

		if enabled {
			apiName = "EnableWorkflow"
		} else {
			apiName = "DisableWorkflow"
		}
		reqBody := map[string]interface{}{
			"Namespace":  d.Get("namespace"),
			"GroupId":    d.Get("group_id"),
			"WorkflowId": d.Id(),
		}

		_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", apiName, "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_schedulerx2_workflow", apiName, errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	if d.IsNewResource() {
		return nil
	}

	// Update workflow basic info
	if d.HasChanges("name", "description", "time_type", "time_expression", "time_zone", "max_concurrency") {
		reqBody := map[string]interface{}{
			"WorkflowId":     d.Id(),
			"Namespace":      d.Get("namespace"),
			"GroupId":        d.Get("group_id"),
			"Name":           d.Get("name"),
			"Description":    d.Get("description"),
			"AcceptLanguage": "en",
		}
		time_type := d.Get("time_type").(string)
		switch time_type {
		case "cron":
			reqBody["TimeType"] = "1"
		case "api":
			reqBody["TimeType"] = "100"
		}
		if v, ok := d.GetOk("time_expression"); ok {
			reqBody["TimeExpression"] = v.(string)
		}
		if v, ok := d.GetOk("time_zone"); ok {
			reqBody["TimeZone"] = v.(string)
		}
		if v, ok := d.GetOk("max_concurrency"); ok {
			reqBody["MaxConcurrency"] = v.(int)
		}
		_, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "UpdateWorkflow", "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_schedulerx2_workflow", "UpdateWorkflow", errmsgs.AlibabacloudStackSdkGoERROR)
		}
	}

	return nil
}

func resourceAlibabacloudStackSchedulerx2WorkflowDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	requestQuery := map[string]interface{}{
		"Namespace":  d.Get("namespace").(string),
		"GroupId":    d.Get("group_id").(string),
		"WorkflowId": d.Id(),
	}

	raw, err := client.DoTeaRequest("POST", "schedulerx2", "2019-04-30", "DeleteWorkflow", "", nil, requestQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "DeleteWorkflow", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	if fmt.Sprint(raw["Code"]) != "200" {
		return errmsgs.WrapError(errmsgs.Error("schedulerx2 workflow delete failed for: " + raw["Message"].(string)))
	}
	return nil
}
