package alibabacloudstack

import (
	"encoding/json"
	"fmt"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackPolardbxReadWriteSplittingConfig() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"db_instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"attend_htap_list": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				Optional: true,
			},
			"auto_attend_htap": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"delay_execution_strategy": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"enable_consistent_replica_read": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"enable_htap": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"master_read_weight": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
			"storage_delay_threshold": {
				Type:     schema.TypeInt,
				Computed: true,
				Optional: true,
			},
		},
	}
	setResourceFunc(
		resource,
		resourceAlibabacloudStackPolardbxReadWriteSplittingConfigCreate,
		resourceAlibabacloudStackPolardbxReadWriteSplittingConfigRead,
		resourceAlibabacloudStackPolardbxReadWriteSplittingConfigUpdate,
		resourceAlibabacloudStackPolardbxReadWriteSplittingConfigDelete,
	)
	return resource
}

func resourceAlibabacloudStackPolardbxReadWriteSplittingConfigCreate(d *schema.ResourceData, meta interface{}) error {
	// As per the documentation, this resource does not support creation.
	// We directly use the DBInstanceName as the resource ID.
	dbInstanceName := d.Get("db_instance_id").(string)
	d.SetId(dbInstanceName)

	return nil
}

func resourceAlibabacloudStackPolardbxReadWriteSplittingConfigRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	polardbxService := PolardbXService{client}

	object, err := polardbxService.DescribePolardbxReadWriteSplittingConfig(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	// Set basic fields from the response
	d.Set("db_instance_id", object["DbInstanceName"])

	// Parse ConfigValue JSON string to extract detailed fields
	configValueStr, ok := object["ConfigValue"].(string)
	if !ok {
		return fmt.Errorf("invalid type for ConfigValue")
	}

	var configValue map[string]interface{}
	if err := json.Unmarshal([]byte(configValueStr), &configValue); err != nil {
		return fmt.Errorf("failed to parse ConfigValue: %v", err)
	}

	// Set individual fields from parsed ConfigValue
	if v, ok := configValue["attendHtapList"].([]interface{}); ok {
		attendHtapList := make([]string, 0, len(v))
		for _, item := range v {
			if s, ok := item.(string); ok {
				attendHtapList = append(attendHtapList, s)
			}
		}
		d.Set("attend_htap_list", attendHtapList)
	}

	if v, ok := configValue["autoAttendHtap"].(string); ok {
		d.Set("auto_attend_htap", v == "true")
	}

	if v, ok := configValue["delayExecutionStrategy"].(float64); ok {
		d.Set("delay_execution_strategy", int(v))
	}

	if v, ok := configValue["enableConsistentReplicaRead"].(bool); ok {
		d.Set("enable_consistent_replica_read", v)
	}

	if v, ok := configValue["enableHtap"].(string); ok {
		d.Set("enable_htap", v == "true")
	}

	if v, ok := configValue["masterReadWeight"].(float64); ok {
		d.Set("master_read_weight", int(v))
	}

	if v, ok := configValue["storageDelayThreshold"].(float64); ok {
		d.Set("storage_delay_threshold", int(v))
	}

	return nil
}

func resourceAlibabacloudStackPolardbxReadWriteSplittingConfigUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Check if any of the config value related fields have changed
	if d.HasChanges("attend_htap_list", "auto_attend_htap", "delay_execution_strategy", "enable_consistent_replica_read", "enable_htap", "master_read_weight", "storage_delay_threshold") {
		reqQuery := map[string]interface{}{
			"DBInstanceName": d.Get("db_instance_id"),
			"ConfigName":     "htap",
		}

		// Construct ConfigValue from individual fields
		configValue := make(map[string]interface{})

		if v, ok := d.GetOk("attend_htap_list"); ok {
			attendHtapList := make([]string, 0)
			for _, item := range v.([]interface{}) {
				attendHtapList = append(attendHtapList, item.(string))
			}
			configValue["attendHtapList"] = attendHtapList
		}

		configValue["autoAttendHtap"] = d.Get("auto_attend_htap").(bool)
		configValue["enableHtap"] = d.Get("enable_htap").(bool)
		configValue["enableConsistentReplicaRead"] = d.Get("enable_consistent_replica_read").(bool)

		if v, ok := d.GetOk("delay_execution_strategy"); ok {
			configValue["delayExecutionStrategy"] = v.(int)
		}

		if v, ok := d.GetOk("master_read_weight"); ok {
			configValue["masterReadWeight"] = v.(int)
		}

		if v, ok := d.GetOk("storage_delay_threshold"); ok {
			configValue["storageDelayThreshold"] = v.(int)
		}

		configValueBytes, err := json.Marshal(configValue)
		if err != nil {
			return fmt.Errorf("failed to marshal config value: %v", err)
		}
		reqQuery["ConfigValue"] = string(configValueBytes)

		// Call ModifyDBInstanceConfig API
		if _, err := client.DoTeaRequest("POST", "polardbx", "2020-02-02", "ModifyDBInstanceConfig", "", nil, reqQuery, nil); err != nil {
			return fmt.Errorf("failed to modify DB instance config: %v", err)
		}
	}

	return nil
}

func resourceAlibabacloudStackPolardbxReadWriteSplittingConfigDelete(d *schema.ResourceData, meta interface{}) error {
	// As per the operation manual, this resource does not support deletion.
	// Hence, we simply return nil to indicate that no action is required.
	return nil
}
