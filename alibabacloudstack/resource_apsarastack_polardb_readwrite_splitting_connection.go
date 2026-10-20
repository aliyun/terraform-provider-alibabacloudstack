package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbReadWriteSplittingConnection() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"connection_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"distribution_type": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"Standard", "Custom"}, false),
			},
			"weight": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeInt,
				},
				Computed: true,
				DiffSuppressFunc: func(k, oldValue, newValue string, d *schema.ResourceData) bool {
					if d.Get("distribution_type") == "Standard" {
						return true
					}
					return oldValue == newValue
				},
				DiffSuppressOnRefresh: true,
			},
			"max_delay_time": {
				Type:     schema.TypeInt,
				Optional: true,
				Default:  30,
			},
			"connection_string": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackPolardbReadWriteSplittingConnectionCreate,
		resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead,
		resourceAlibabacloudStackPolardbReadWriteSplittingConnectionUpdate,
		resourceAlibabacloudStackPolardbReadWriteSplittingConnectionDelete)
	return resource
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	dbInstanceId := d.Get("instance_id").(string)

	d.SetId(dbInstanceId)

	return nil
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbService := PolardbService{client}
	response, err := polardbService.DescribeDBReadWriteSplittingConnection(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_id", response["DBProxyEndpointId"])
	d.Set("connection_string", response["DBProxyConnectString"])
	d.Set("distribution_type", response["ReadOnlyInstanceDistributionType"])
	{
		var instances []interface{}
		rawWeight := response["ReadOnlyInstanceWeight"]
		if rawWeight == nil {
			// ReadWriteSplitting not yet configured; clear state
			d.SetId("")
			return nil
		}
		switch v := rawWeight.(type) {
		case string:
			if v != "" {
				if err := json.Unmarshal([]byte(v), &instances); err != nil {
					return fmt.Errorf("failed to parse ReadOnlyInstanceWeight: %w", err)
				}
			}
		case []interface{}:
			instances = v
		default:
			return fmt.Errorf("unsupported ReadOnlyInstanceWeight type: %T", rawWeight)
		}

		result := make(map[string]int)
		for _, item := range instances {
			inst, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			if avail, ok := inst["Availability"].(string); ok && avail == "Available" {
				if id, ok := inst["DBInstanceId"].(string); ok {
					if weight, err := toInt(inst["Weight"]); err == nil {
						result[id] = weight
					}
				}
			}
		}
		d.Set("weight", result)
	}
	if v, err := toInt(response["ReadOnlyInstanceMaxDelayTime"]); err == nil {
		d.Set("max_delay_time", v)
	}
	if v, err := toInt(response["DBProxyConnectStringPort"]); err == nil {
		d.Set("port", v)
	}
	return nil
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"DBInstanceId":                     d.Id(),
		"ConfigDBProxyFeatures":            "ReadWriteSpliting:1;",
		"ReadOnlyInstanceDistributionType": d.Get("distribution_type"),
		"DBProxyEndpointId":                d.Get("connection_id"),
	}
	if v, ok := d.GetOk("max_delay_time"); ok && v.(int) != 0 {
		reqQuery["ReadOnlyInstanceMaxDelayTime"] = v
	}
	if d.Get("distribution_type").(string) == "Custom" {
		weigth, err := json.Marshal(d.Get("weight").(map[string]interface{}))
		if err != nil {
			return err
		}
		reqQuery["ReadOnlyInstanceWeight"] = string(weigth)
	}

	if _, err := client.DoTeaRequest("POST", "polardb", "2024-01-30", "ModifyDBProxyEndpoint", "", nil, reqQuery, nil); err != nil {
		return err
	}
	polardbServce := PolardbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbServce.PolardbProxyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackPolardbReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"DBInstanceId":          d.Id(),
		"ConfigDBProxyFeatures": "ReadWriteSpliting:0;",
		"DBProxyEndpointId":     d.Get("connection_id"),
	}

	if _, err := client.DoTeaRequest("POST", "polardb", "2024-01-30", "ModifyDBProxyEndpoint", "", nil, reqQuery, nil); err != nil {
		return err
	}
	polardbServce := PolardbService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, polardbServce.PolardbProxyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}