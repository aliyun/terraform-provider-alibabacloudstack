package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const resourceAlibabacloudStackCSAutoscalingConfigName = "resource_apsarastack_cs_autoscaling_config"

func resourceAlibabacloudStackCSAutoscalingConfig() *schema.Resource {
	resource := &schema.Resource{

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cool_down_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"unneeded_duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "10m",
			},
			"utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"gpu_utilization_threshold": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "0.5",
			},
			"ensure_required_role": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackCSAutoscalingConfigCreate,
		resourceAlibabacloudStackCSAutoscalingConfigRead, resourceAlibabacloudStackCSAutoscalingConfigUpdate,
		nil)
	return resource
}

func resourceAlibabacloudStackCSAutoscalingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	clusterId := d.Get("cluster_id").(string)

	reqBody := make(map[string]interface{})
	if v, ok := d.GetOk("cool_down_duration"); ok {
		reqBody["cool_down_duration"] = v.(string)
	}
	if v, ok := d.GetOk("unneeded_duration"); ok {
		reqBody["unneeded_duration"] = v.(string)
	}
	if v, ok := d.GetOk("utilization_threshold"); ok {
		reqBody["utilization_threshold"] = v.(string)
	}
	if v, ok := d.GetOk("gpu_utilization_threshold"); ok {
		reqBody["gpu_utilization_threshold"] = v.(string)
	}
	if v, ok := d.GetOk("ensure_required_role"); ok {
		reqBody["ensure_required_role"] = v.(bool)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "CS", "2015-12-15", "CreateAutoscalingConfig", fmt.Sprintf("/cluster/%s/autoscale/config/", clusterId), nil, nil, reqBody)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, resourceAlibabacloudStackCSAutoscalingConfigName, "CreateAutoscalingConfig", err)
	}

	d.SetId(clusterId)
	return resourceAlibabacloudStackCSAutoscalingConfigRead(d, meta)
}

func resourceAlibabacloudStackCSAutoscalingConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	csService := CsService{client}

	addonConfigRaw, err := csService.GetAutoscalingConfig(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("cluster_id", d.Id())
	if v, ok := addonConfigRaw["ScaleDownUnneededTime"]; ok && v != nil {
		d.Set("cool_down_duration", fmt.Sprint(v))
	}
	if v, ok := addonConfigRaw["ScaleDownDelayAfterAdd"]; ok && v != nil {
		d.Set("unneeded_duration", fmt.Sprint(v))
	}
	if v, ok := addonConfigRaw["ScaleDownUtilizationThreshold"]; ok && v != nil {
		d.Set("utilization_threshold", fmt.Sprint(v))
	}
	if v, ok := addonConfigRaw["ScaleDownGpuUtilizationThreshold"]; ok && v != nil {
		d.Set("gpu_utilization_threshold", fmt.Sprint(v))
	}
	// ensure_required_role is not returned by the API, use schema default
	d.Set("ensure_required_role", true)
	return nil
}

func resourceAlibabacloudStackCSAutoscalingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	// auto scaling config
	reqBody := make(map[string]interface{})
	if v, ok := d.GetOk("cool_down_duration"); ok {
		reqBody["cool_down_duration"] = v.(string)
	}
	if v, ok := d.GetOk("unneeded_duration"); ok {
		reqBody["unneeded_duration"] = v.(string)
	}
	if v, ok := d.GetOk("utilization_threshold"); ok {
		reqBody["utilization_threshold"] = v.(string)
	}
	if v, ok := d.GetOk("gpu_utilization_threshold"); ok {
		reqBody["gpu_utilization_threshold"] = v.(string)
	}
	if v, ok := d.GetOk("ensure_required_role"); ok {
		reqBody["ensure_required_role"] = v.(bool)
	}

	wait := incrementalWait(3*time.Second, 3*time.Second)
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "CS", "2015-12-15", "CreateAutoscalingConfig", fmt.Sprintf("/cluster/%s/autoscale/config/", d.Id()), nil, nil, reqBody)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, resourceAlibabacloudStackCSAutoscalingConfigName, "ModifyClusterAddon", err)
	}

	return nil
}
