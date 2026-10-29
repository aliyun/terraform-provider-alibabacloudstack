package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackLindormLtsInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(90 * time.Minute),
			Update: schema.DefaultTimeout(30 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			"zone_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_alias": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cpu_brand": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"lts_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"topology_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "1azone",
				ValidateFunc: validation.StringInSlice([]string{"1azone", "multizone"}, false),
			},
			"standby_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arbiter_zone_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"standby_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"arbiter_vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource,
		resourceAlibabacloudStackLindormLtsInstanceCreate,
		resourceAlibabacloudStackLindormLtsInstanceRead,
		resourceAlibabacloudStackLindormLtsInstanceUpdate,
		resourceAlibabacloudStackLindormLtsInstanceDelete,
	)

	return resource
}

func resourceAlibabacloudStackLindormLtsInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}
	request := make(map[string]interface{})
	request["TopologyType"] = d.Get("topology_type")
	if d.Get("topology_type").(string) == "multizone" {
		request["ArchVersion"] = "2.0"
	}
	request["ZoneId"] = d.Get("zone_id")
	request["InstanceAlias"] = d.Get("instance_alias")
	request["CpuBrand"] = d.Get("cpu_brand")
	request["LtsSpec"] = d.Get("instance_type")
	request["LtsNum"] = d.Get("lts_num")
	request["ClusterType"] = "lindorm-for-lts"
	if v, ok := d.GetOk("vpc_id"); ok {
		request["VPCId"] = v
	}
	if v, ok := d.GetOk("vswitch_id"); ok {
		request["VSwitchId"] = v
	}
	if v, ok := d.GetOk("standby_zone_id"); ok {
		request["StandbyZoneId"] = v
	}
	if v, ok := d.GetOk("arbiter_zone_id"); ok {
		request["ArbiterZoneId"] = v
	}
	if v, ok := d.GetOk("standby_vswitch_id"); ok {
		request["StandbyVSwitchId"] = v
	}
	if v, ok := d.GetOk("arbiter_vswitch_id"); ok {
		request["ArbiterVSwitchId"] = v
	}
	resp, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "CreateLindormInstance", "", nil, request, nil)
	if err != nil {
		return err
	}
	instanceId, ok := resp["InstanceId"]
	if !ok {
		return errmsgs.Error(fmt.Sprintf("Create lindorm instance failed %#v", resp))
	}
	d.SetId(instanceId.(string))
	// Wait for instance to be active
	stateConf := BuildStateConf([]string{"CREATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutCreate), 60*time.Second, lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"CREATE_FAILED"}))
	if _, err := stateConf.WaitForState(); err != nil {
		return err
	}
	return nil
}

func resourceAlibabacloudStackLindormLtsInstanceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}

	data, err := lindormService.DescribeLindormInstance(d.Id())
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapError(err)
	}

	d.Set("instance_id", data["InstanceId"])
	d.Set("instance_alias", data["InstanceAlias"])
	d.Set("zone_id", data["ZoneId"])
	d.Set("instance_status", data["InstanceStatus"])
	d.Set("create_time", data["CreateTime"])
	d.Set("deletion_protection", data["DeletionProtection"])
	d.Set("cpu_brand", data["CpuBrand"])
	d.Set("vpc_id", data["VpcId"])
	d.Set("vswitch_id", data["VswitchId"])
	d.Set("standby_zone_id", data["StandbyZoneId"])
	d.Set("arbiter_zone_id", data["ArbiterZoneId"])
	d.Set("standby_vswitch_id", data["StandbyVSwitchId"])
	d.Set("arbiter_vswitch_id", data["ArbiterVSwitchId"])
	d.Set("service_type", data["ServiceType"])
	d.Set("network_type", data["NetworkType"])
	if archVer, ok := data["ArchVersion"].(string); ok {
		d.Set("arch_version", archVer)
		if archVer == "2.0" {
			d.Set("topology_type", "multizone")
		} else {
			d.Set("topology_type", "1azone")
		}
	}
	engineList := data["EngineList"].([]interface{})
	if len(engineList) > 0 {
		engineData := engineList[0].(map[string]interface{})
		d.Set("lts_num", engineData["CoreCount"])
		d.Set("instance_type", engineData["Specification"])
	}
	return nil
}

func resourceAlibabacloudStackLindormLtsInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}
	if d.HasChanges("instance_alias", "deletion_protection") {
		reqQuery := map[string]interface{}{
			"InstanceId":         d.Id(),
			"InstanceAlias":      d.Get("instance_alias"),
			"DeletionProtection": d.Get("deletion_protection"),
		}

		if _, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpdateLindormInstanceAttribute", "", nil, reqQuery, nil); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_lindorm_lts_instance", "UpdateLindormInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		// Wait for the instance to be updated
		stateConf := BuildStateConf([]string{"UPDATING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 3*time.Second,
			lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}
	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("instance_type", "lts_num") {
		old_type, new_type := d.GetChange("instance_type")
		old_num, new_num := d.GetChange("lts_num")
		if old_type != new_type {
			reqQuery := map[string]interface{}{
				"InstanceId":          d.Id(),
				"ZoneId":              d.Get("zone_id"),
				"EngineType":          "lts",
				"EngineCoreNum":       old_num,
				"EngineSpec":          new_type,
				"SourceEngineCoreNum": old_num,
				"SourceEngineSpec":    old_type,
			}

			if _, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpgradeLindormInstanceEngine", "", nil, reqQuery, nil); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_lindorm_lts_instance", "UpgradeLindormInstanceEngine", errmsgs.AlibabacloudStackSdkGoERROR)
			}

			// Wait for the instance to be updated
			stateConf := BuildStateConf([]string{"RESIZING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second,
				lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"FAILED"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}

		}
		if old_num != new_num {
			reqQuery := map[string]interface{}{
				"InstanceId":          d.Id(),
				"ZoneId":              d.Get("zone_id"),
				"EngineType":          "lts",
				"EngineCoreNum":       new_num,
				"EngineSpec":          new_type,
				"SourceEngineCoreNum": old_num,
				"SourceEngineSpec":    new_type,
			}

			if _, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpgradeLindormInstanceEngine", "", nil, reqQuery, nil); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_lindorm_lts_instance", "UpgradeLindormInstanceEngine", errmsgs.AlibabacloudStackSdkGoERROR)
			}

			// Wait for the instance to be updated
			stateConf := BuildStateConf([]string{"RESIZING"}, []string{"ACTIVATION"}, d.Timeout(schema.TimeoutUpdate), 20*time.Second,
				lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"FAILED"}))
			if _, err := stateConf.WaitForState(); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
			}
		}
	}

	return nil
}

func resourceAlibabacloudStackLindormLtsInstanceDelete(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*connectivity.AlibabacloudStackClient)
	reqQuery := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	// Call the release API to delete the instance
	_, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "ReleaseLindormInstance", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ReleaseLindormInstance", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	return nil
}
