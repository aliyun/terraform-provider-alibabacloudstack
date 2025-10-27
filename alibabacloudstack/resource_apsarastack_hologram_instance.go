package alibabacloudstack

import (
	"fmt"
	"log"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackHologramInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"compute_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Follower"}, false),
			},
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"cpu": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "intel",
			},
			"node": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"leader_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_hive_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"endpoint": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackHologramInstanceCreate,
		resourceAlibabacloudStackHologramInstanceRead, resourceAlibabacloudStackHologramInstanceUpdate, resourceAlibabacloudStackHologramInstanceDelete)
	return resource
}

func resourceAlibabacloudStackHologramInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hologramService := HologramService{client}
	node := d.Get("node").(int)
	cluster := d.Get("cluster").(string)
	quota, err := hologramService.CalculateQuota(node, cluster)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	compute_type := d.Get("compute_type").(string)
	body := map[string]interface{}{
		"zoneId":         d.Get("zone_id"),
		"vpcId":          d.Get("vpc_id"),
		"vpcSwitchId":    d.Get("vswitch_id"),
		"instanceName":   d.Get("instance_name"),
		"computeType":    compute_type,
		"cpu":            d.Get("cpu"),
		"node":           node,
		"cluster":        cluster,
		"cu":             quota["Cu"],
		"memory":         quota["Memory"],
		"isNewFeature":   true,
		"ClusterEnvType": 0,
		"Department":     client.Department,
		"ResourceGroup":  client.ResourceGroup,
		"region":         client.Region,
	}
	if compute_type == "Follower" {
		body["ClusterEnvType"] = 2
		if v, ok := d.GetOk("leader_instance_id"); ok {
			stateConf := BuildStateConf([]string{"Allocating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, hologramService.HologramInstanceStateRefreshFunc(v.(string), []string{"Failed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("waiting for hologram leader instance %s to be Running failed: %v", v.(string), err)
			}
			body["leaderInstanceId"] = v
		} else {
			return errmsgs.Error("leader_instance_id is required when compute_type is Follower")
		}
	} else {
		if _, ok := d.GetOk("leader_instance_id"); ok {
			return errmsgs.Error("leader_instance_id can not be set when compute_type is not Follower")
		}
	}
	request := map[string]interface{}{
		"x-acs-body": body,
	}
	response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "CreateInstance", "/inner/v1/instances/createInstance", nil, nil, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "CreateInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	v, err := jsonpath.Get("$.Instance.InstanceId", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "alibabacloudstack_hologram_instance", "Instance.InstanceId", response)
	}
	instanceId := v.(string)
	stateConf := BuildStateConf([]string{"Allocating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 1*time.Minute, hologramService.HologramInstanceStateRefreshFunc(instanceId, []string{"Failed"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for hologram instance %s to be Running failed: %v", instanceId, err)
	}
	d.SetId(instanceId)
	return nil
}

func resourceAlibabacloudStackHologramInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.IsNewResource() {
		return nil
	}
	if d.HasChange("instance_name") {
		body := map[string]interface{}{
			"instanceName": d.Get("instance_name").(string),
		}
		request := map[string]interface{}{
			"x-acs-body": body,
		}
		pattern := fmt.Sprintf("/api/v1/instances/%s/instanceName", d.Id())
		response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "UpdateInstanceName", pattern, nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ScaleInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		log.Printf("[DEBUG] Hologres instance updated: %s", response)
	}
	if d.HasChange("node") {
		hologramService := HologramService{client}
		node := d.Get("node").(int)
		cluster := d.Get("cluster").(string)
		quota, err := hologramService.CalculateQuota(node, cluster)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		body := map[string]interface{}{
			"instanceId": d.Id(),
			"node":       node,
			"cluster":    cluster,
			"cu":         quota["Cu"],
			"memory":     quota["Memory"],
		}
		request := map[string]interface{}{
			"x-acs-body": body,
		}
		pattern := fmt.Sprintf("/inner/v1/instances/%s/scaleInstance", d.Id())
		response, err := client.DoTeaRequest("POST", "Hologram", "2022-06-01", "ScaleInstance", pattern, nil, nil, request)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "ScaleInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		log.Printf("[DEBUG] Hologres instance updated: %s", response)
		stateConf := BuildStateConf([]string{"Allocating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, hologramService.HologramInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for hologram instance %s to be Running failed: %v", d.Id(), err)
		}
	}
	if d.HasChange("leader_instance_id") {
		if d.Get("compute_type") != "Follower" {
			return errmsgs.Error("leader_instance_id can not be set when compute_type is not Follower")
		}
		// old, new := d.GetChange("leader_instance_id")
		hologramService := HologramService{client}
		var err error
		unBindPattern := fmt.Sprintf("/api/v1/instances/%s/unBindReadOnly", d.Id())
		_, err = client.DoTeaRequest("POST", "Hologram", "2022-06-01", "UnBindLeaderInstance", unBindPattern, nil, nil, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "UnBindLeaderInstance", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		stateConf := BuildStateConf([]string{"Allocating"}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, hologramService.HologramInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
		if _, err = stateConf.WaitForState(); err != nil {
			return fmt.Errorf("waiting for hologram instance %s to be Running failed: %v", d.Id(), err)
		}
		if v, ok := d.GetOk("leader_instance_id"); ok {
			request := map[string]interface{}{
				"x-acs-body": map[string]interface{}{
					"leaderInstanceId": v.(string),
					"RegionId":         client.RegionId,
				},
			}
			pattern := fmt.Sprintf("/api/v1/instances/%s/bindReadOnly", d.Id())
			_, err = client.DoTeaRequest("POST", "Hologram", "2022-06-01", "BindLeaderInstance", pattern, nil, nil, request)
			if err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "BindLeaderInstance", errmsgs.AlibabacloudStackSdkGoERROR)
			}
			stateConf := BuildStateConf([]string{"Allocating"}, []string{"Running"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, hologramService.HologramInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return fmt.Errorf("waiting for hologram instance %s to be Running failed: %v", d.Id(), err)
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackHologramInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hologramService := HologramService{client}
	instance, err := hologramService.DescribeHologramInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "DescribeHologramInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	d.Set("instance_name", instance["InstanceName"])
	d.Set("zone_id", instance["ZoneId"])
	d.Set("compute_type", instance["InstanceType"])
	d.Set("node", instance["ComputeNodeCount"])
	d.Set("instance_status", instance["InstanceStatus"])
	d.Set("cluster", instance["Cluster"])
	d.Set("creation_time", instance["CreationTime"])
	d.Set("cpu", instance["CpuBrand"])
	d.Set("version", instance["Version"])
	d.Set("enable_hive_access", instance["EnableHiveAccess"])
	if v, ok := instance["LeaderInstanceId"]; ok {
		d.Set("leader_instance_id", v)
	} else {
		d.Set("leader_instance_id", "")
	}

	if endpoints, ok := instance["Endpoints"].([]interface{}); ok {
		endpointList := make([]map[string]interface{}, 0)
		for _, endpointItem := range endpoints {
			if endpointMap, ok := endpointItem.(map[string]interface{}); ok {
				endpoint := make(map[string]interface{})
				if v, ok := endpointMap["Type"]; ok {
					endpoint["type"] = v
				}
				if v, ok := endpointMap["Endpoint"]; ok {
					endpoint["endpoint"] = v
				}
				if v, ok := endpointMap["Enabled"]; ok {
					endpoint["enabled"] = v
				}
				if v, ok := endpointMap["VpcId"]; ok {
					d.Set("vpc_id", v)
					endpoint["vpc_id"] = v
				}
				if v, ok := endpointMap["VSwitchId"]; ok {
					d.Set("vswitch_id", v)
					endpoint["vswitch_id"] = v
				}
				if v, ok := endpointMap["VpcInstanceId"]; ok {
					endpoint["vpc_instance_id"] = v
				}
				endpointList = append(endpointList, endpoint)
			}
		}
		d.Set("endpoints", endpointList)
	}

	return nil
}

func resourceAlibabacloudStackHologramInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	hologramService := HologramService{client}
	var err error
	request := map[string]interface{}{
		"instanceId": d.Id(),
	}
	stop_pattern := fmt.Sprintf("/api/v1/instances/%s/stop", d.Id())
	_, err = client.DoTeaRequest("POST", "Hologram", "2022-06-01", "StopInstance", stop_pattern, nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "StopInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	stateConf := BuildStateConf([]string{"Allocating"}, []string{"Suspended"}, d.Timeout(schema.TimeoutCreate), 10*time.Second, hologramService.HologramInstanceStateRefreshFunc(d.Id(), []string{"Failed"}))
	if _, err = stateConf.WaitForState(); err != nil {
		return fmt.Errorf("waiting for hologram instance %s to be Suspended failed: %v", d.Id(), err)
	}
	pattern := fmt.Sprintf("/api/v1/instances/%s/delete", d.Id())
	_, err = client.DoTeaRequest("POST", "Hologram", "2022-06-01", "DeleteInstance", pattern, nil, request, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_hologram_instance", "DeleteInstance", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	return nil
}
