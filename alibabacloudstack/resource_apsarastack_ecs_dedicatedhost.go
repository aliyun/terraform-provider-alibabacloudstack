package alibabacloudstack

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackEcsDedicatedHost() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(11 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
			Update: schema.DefaultTimeout(11 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"action_on_maintenance": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"Stop", "Migrate"}, false),
				Default:      "Stop",
			},
			"auto_placement": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"off", "on"}, false),
				Default:      "on",
			},
			"dedicated_host_cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dedicated_host_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"machine_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"core": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"physical_gpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"tags": tagsSchema(),
			"zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"supported_instance_type_families": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"supported_instance_types_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
		},
		CustomizeDiff: func(context context.Context, d *schema.ResourceDiff, i interface{}) error {
			old, _ := d.GetChange("dedicated_host_cluster_id")
			if d.HasChange("dedicated_host_cluster_id") && old.(string) != "" {
				return fmt.Errorf("can not change dedicated_host_cluster")
			}
			return nil
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackEcsDedicatedHostCreate, resourceAlibabacloudStackEcsDedicatedHostRead, resourceAlibabacloudStackEcsDedicatedHostUpdate, resourceAlibabacloudStackEcsDedicatedHostDelete)
	return resource
}

func resourceAlibabacloudStackEcsDedicatedHostCreate(d *schema.ResourceData, meta interface{}) (err error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	var response map[string]interface{}
	action := "AllocateDedicatedHosts"
	request := map[string]interface{}{
		"DedicatedHostType":   d.Get("dedicated_host_type"),
		"Quantity":            1,
		"ActionOnMaintenance": d.Get("action_on_maintenance"),
		"AutoPlacement":       d.Get("auto_placement"),
	}

	if v, ok := d.GetOk("dedicated_host_cluster_id"); ok {
		request["DedicatedHostClusterId"] = v
	}

	if v, ok := d.GetOk("dedicated_host_name"); ok {
		request["DedicatedHostName"] = v
	}

	if v, ok := d.GetOk("zone_id"); ok {
		request["ZoneId"] = v
	}

	response, err = client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
	if err != nil {
		return err
	}
	responseDedicatedHostIdSets := response["DedicatedHostIdSets"].(map[string]interface{})
	d.SetId(responseDedicatedHostIdSets["DedicatedHostId"].([]interface{})[0].(string))
	stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutCreate), 15*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{"PermanentFailure"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackEcsDedicatedHostRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}
	object, err := ecsService.DescribeEcsDedicatedHost(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			log.Printf("[DEBUG] Resource alibabacloudstack_ecs_dedicated_host ecsService.DescribeEcsDedicatedHost Failed!!! %s", err)
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}
	d.Set("action_on_maintenance", object.ActionOnMaintenance)
	d.Set("auto_placement", object.AutoPlacement)
	d.Set("dedicated_host_name", object.DedicatedHostName)
	d.Set("dedicated_host_type", object.DedicatedHostType)
	d.Set("dedicated_host_cluster_id", object.DedicatedHostClusterId)
	d.Set("status", object.Status)
	d.Set("machine_id", object.MachineId)
	d.Set("core", object.Cores)
	d.Set("physical_gpus", object.PhysicalGpus)
	d.Set("tags", ecsService.tagsToMap(object.Tags.Tag))
	d.Set("zone_id", object.ZoneId)
	d.Set("supported_instance_type_families", object.SupportedInstanceTypeFamilies.SupportedInstanceTypeFamily)
	d.Set("supported_instance_types_list", object.SupportedInstanceTypesList.SupportedInstanceTypesList)
	return nil
}

func resourceAlibabacloudStackEcsDedicatedHostUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	ecsService := EcsService{client}

	if d.HasChange("tags") {
		if err := ecsService.SetResourceTags(d, "ddh"); err != nil {
			return errmsgs.WrapError(err)
		}
	}
	update := false
	modifyDedicatedHostAttributeReq := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	if !d.IsNewResource() && d.HasChange("action_on_maintenance") {
		update = true
		modifyDedicatedHostAttributeReq["ActionOnMaintenance"] = d.Get("action_on_maintenance")
	}
	if !d.IsNewResource() && d.HasChange("auto_placement") {
		update = true
		modifyDedicatedHostAttributeReq["AutoPlacement"] = d.Get("auto_placement")
	}
	if !d.IsNewResource() && d.HasChange("dedicated_host_name") {
		update = true
		modifyDedicatedHostAttributeReq["DedicatedHostName"] = d.Get("dedicated_host_name")
	}
	if update {
		if _, ok := d.GetOk("dedicated_host_cluster_id"); ok {
			modifyDedicatedHostAttributeReq["DedicatedHostClusterId"] = d.Get("dedicated_host_cluster_id")
		}
		action := "ModifyDedicatedHostAttribute"
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, modifyDedicatedHostAttributeReq)
		if err != nil {
			return err
		}
		stateConf := BuildStateConf([]string{}, []string{"Available"}, d.Timeout(schema.TimeoutUpdate), 5*time.Second, ecsService.EcsDedicatedHostStateRefreshFunc(d.Id(), []string{}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	return nil
}

func resourceAlibabacloudStackEcsDedicatedHostDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	action := "ReleaseDedicatedHost"
	request := map[string]interface{}{
		"DedicatedHostId": d.Id(),
	}
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		_, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", action, "", nil, nil, request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"IncorrectHostStatus.Initializing"}) {
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	return err
}
