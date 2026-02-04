package alibabacloudstack

import (
	"encoding/json"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackDBReadWriteSplittingConnection() *schema.Resource {
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
	setResourceFunc(resource, resourceAlibabacloudStackDBReadWriteSplittingConnectionCreate,
		resourceAlibabacloudStackDBReadWriteSplittingConnectionRead,
		resourceAlibabacloudStackDBReadWriteSplittingConnectionUpdate,
		resourceAlibabacloudStackDBReadWriteSplittingConnectionDelete)
	return resource
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionCreate(d *schema.ResourceData, meta interface{}) error {

	dbInstanceId := d.Get("instance_id").(string)

	d.SetId(dbInstanceId)

	return nil
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionRead(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	rdsServce := RdsService{client}
	response, err := rdsServce.DescribeDBReadWriteSplittingConnection(d.Id())
	if err != nil {
		return err
	}

	d.Set("instance_id", d.Id())
	d.Set("connection_id", response["DBProxyEndpointId"])
	d.Set("connection_string", response["DBProxyConnectString"])
	d.Set("distribution_type", response["ReadOnlyInstanceDistributionType"])
	{
		var instances []map[string]interface{}
		if err := json.Unmarshal([]byte(response["ReadOnlyInstanceWeight"].(string)), &instances); err != nil {
			panic(err)
		}

		result := make(map[string]int)
		for _, inst := range instances {
			if avail, ok := inst["Availability"].(string); ok && avail == "Available" {
				if id, ok := inst["DBInstanceId"].(string); ok {
					if weight, ok := inst["Weight"].(float64); ok {
						result[id] = int(weight)
					} else if weight, ok := inst["Weight"].(int); ok {
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

func resourceAlibabacloudStackDBReadWriteSplittingConnectionUpdate(d *schema.ResourceData, meta interface{}) error {
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

	if _, err := client.DoTeaRequest("POST", "Rds", "2014-08-15", "ModifyDBProxyEndpoint", "", nil, reqQuery, nil); err != nil {
		return err
	}
	rdsService := RdsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, rdsService.RdsProxyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}

func resourceAlibabacloudStackDBReadWriteSplittingConnectionDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	reqQuery := map[string]interface{}{
		"DBInstanceId":          d.Id(),
		"ConfigDBProxyFeatures": "ReadWriteSpliting:0;",
		"DBProxyEndpointId":     d.Get("connection_id"),
	}

	if _, err := client.DoTeaRequest("POST", "Rds", "2014-08-15", "ModifyDBProxyEndpoint", "", nil, reqQuery, nil); err != nil {
		return err
	}
	rdsService := RdsService{client}
	stateConf := BuildStateConf([]string{}, []string{"Running"}, d.Timeout(schema.TimeoutUpdate), 10*time.Second, rdsService.RdsProxyStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
