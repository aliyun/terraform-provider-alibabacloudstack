package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackMqttInstance() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"max_sub": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  1000000,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "mqtt4Private",
				ForceNew: true,
			},
			"max_conn": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 100000),
			},
			"independent_naming": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
			"max_up_tps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(500, 20000),
			},
			"max_down_tps": {
				Type:         schema.TypeInt,
				Optional:     true,
				ValidateFunc: validation.IntBetween(1000, 100000),
			},
			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"store_type": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntInSlice([]int{1}),
			},
			"store_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoints": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackMqttInstanceCreate, resourceAlibabacloudStackMqttInstanceRead, resourceAlibabacloudStackMqttInstanceUpdate, resourceAlibabacloudStackMqttInstanceDelete)
	return resource
}

func resourceAlibabacloudStackMqttInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Prepare the request parameters for creating MQTT instance
	reqBody := map[string]interface{}{
		"OnsRegionId":       client.RegionId,
		"MaxSub":            d.Get("max_sub"),
		"InstanceName":      d.Get("instance_name"),
		"ClusterName":       d.Get("cluster_name"),
		"MaxConn":           d.Get("max_conn"),
		"IndependentNaming": d.Get("independent_naming"),
		"PreventCache":      time.Now().UnixNano() / 1e6,
		"MaxUpTps":          d.Get("max_up_tps"),
		"MaxDownTps":        d.Get("max_down_tps"),
		"Remark":            d.Get("remark"),
	}
	resp, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceCreate", "", nil, reqBody, nil)
	if err != nil {
		return err
	}

	instanceId, ok := resp["Data"]
	if !ok {
		return fmt.Errorf("create mqtt instance failed! response: %v", resp)
	}

	d.SetId(instanceId.(string))
	return nil
}

func resourceAlibabacloudStackMqttInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	mqttService := OnsService{client}
	object, err := mqttService.DescribeMqttInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	d.Set("instance_id", object["instanceId"])
	d.Set("instance_name", object["instanceName"])
	d.Set("max_conn", object["maxConn"])
	d.Set("max_sub", object["maxSub"])
	d.Set("max_up_tps", object["maxUpTps"])
	d.Set("max_down_tps", object["maxDownTps"])
	d.Set("independent_naming", object["independentNaming"])
	d.Set("remark", object["remark"])
	d.Set("namespace_id", object["namespaceId"])
	d.Set("store_instance_id", object["storeInstanceId"])
	d.Set("store_type", object["storeType"])

	if endpoints, ok := object["endpoints"].(map[string]interface{}); ok {
		d.Set("endpoints", endpoints)
	}

	return nil
}

func resourceAlibabacloudStackMqttInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	if d.HasChanges("store_instance_id", "store_type") && d.Get("store_instance_id") != "" {
		reqBody := map[string]interface{}{
			"Platform":       "onsConsole",
			"OnsRegionId":    client.RegionId,
			"Dauth_url_hash": "mqtt%2Fconsole%2Finstances",
			"MqttInstanceId": d.Id(),
			"StoreType":      d.Get("store_type"),
			"InstanceId":     d.Get("store_instance_id"),
			"PreventCache":   time.Now().UnixNano() / 1e6,
		}
		response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceBind", "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_instance", "ConsoleMqttInstanceBind", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprintf("%v", response["Code"]) != "200" {
			return errmsgs.WrapError(fmt.Errorf("ConsoleMqttInstanceBind: update mqtt store_instance_id failed! response: %v", response))
		}
	}

	if d.IsNewResource() {
		return nil
	}
	if d.HasChanges("instance_name", "remark") {
		reqBody := map[string]interface{}{
			"MqttInstanceId":    d.Id(),
			"MqttInstanceName":  d.Get("instance_name"),
			"Remark":            d.Get("remark"),
			"independentNaming": d.Get("independent_naming"),
			"MaxConn":           d.Get("max_conn"),
			"MaxSub":            d.Get("max_sub"),
			"MaxUpTps":          d.Get("max_up_tps"),
			"MaxDownTps":        d.Get("max_down_tps"),
			"OnsRegionId":       client.RegionId,
			"PreventCache":      time.Now().UnixNano() / 1e6,
		}

		response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceUpdate", "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_instance", "ConsoleMqttInstanceUpdate", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprintf("%v", response["Code"]) != "200" {
			return errmsgs.WrapError(fmt.Errorf("ConsoleMqttInstanceUpdate: update mqtt failed! response: %v", response))
		}
	}
	// ConsoleMqttInstanceChangeSpec
	if d.HasChanges("max_conn", "max_sub", "max_up_tps", "max_down_tps") {
		reqBody := map[string]interface{}{
			"MaxConn":        d.Get("max_conn"),
			"MaxSub":         d.Get("max_sub"),
			"MaxUpTps":       d.Get("max_up_tps"),
			"MaxDownTps":     d.Get("max_down_tps"),
			"MqttInstanceId": d.Id(),
			"OnsRegionId":    client.RegionId,
			"PreventCache":   time.Now().UnixNano() / 1e6,
		}

		response, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceChangeSpec", "", nil, reqBody, nil)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "alibabacloudstack_mqtt_instance", "ConsoleMqttInstanceChangeSpec", errmsgs.AlibabacloudStackSdkGoERROR)
		}
		if fmt.Sprintf("%v", response["Code"]) != "200" {
			return errmsgs.WrapError(fmt.Errorf("ConsoleMqttInstanceChangeSpec: update mqtt failed! response: %v", response))
		}
	}

	return nil
}

func resourceAlibabacloudStackMqttInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	requestQuery := map[string]interface{}{
		"MqttInstanceId": d.Id(),
		"regionId":       client.RegionId,
		"OnsRegionId":    client.RegionId,
		"PreventCache":   time.Now().UnixNano() / 1e6,
	}

	_, err := client.DoTeaRequest("POST", "Ons-inner", "2018-02-05", "ConsoleMqttInstanceDelete", "", nil, requestQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, d.Id(), "ConsoleMqttInstanceDelete", errmsgs.AlibabacloudStackSdkGoERROR)
	}

	return nil
}
