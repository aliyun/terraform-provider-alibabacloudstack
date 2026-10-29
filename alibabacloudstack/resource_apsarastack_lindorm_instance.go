package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackLindormInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
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
			"disk_category": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"engine_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_type": {
				Type:     schema.TypeString,
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
			"lindorm_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"local_disk_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"local_disk_size": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"arch_version": {
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
			"instance_storage": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"deletion_protection": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"disk_usage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enable_fs": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"switch_l_proxy_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"switch_ssl_encryption_flag": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ali_uid": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}

	setResourceFunc(resource,
		resourceAlibabacloudStackLindormInstanceCreate,
		resourceAlibabacloudStackLindormInstanceRead,
		resourceAlibabacloudStackLindormInstanceUpdate,
		resourceAlibabacloudStackLindormInstanceDelete,
	)

	return resource
}

func resourceAlibabacloudStackLindormInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}
	request := make(map[string]interface{})
	request["TopologyType"] = d.Get("topology_type")
	request["ZoneId"] = d.Get("zone_id")
	request["InstanceAlias"] = d.Get("instance_alias")
	request["CpuBrand"] = d.Get("cpu_brand")
	request["DiskCategory"] = d.Get("disk_category")
	request["EngineType"] = d.Get("engine_type")
	request["LindormSpec"] = d.Get("instance_type")
	request["VPCId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["LindormNum"] = d.Get("lindorm_num")
	request["LocalDiskSize"] = d.Get("local_disk_size")
	request["LocalDiskNum"] = d.Get("local_disk_num")
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

func resourceAlibabacloudStackLindormInstanceRead(d *schema.ResourceData, meta interface{}) error {
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
	d.Set("vpc_id", data["VpcId"])
	d.Set("vswitch_id", data["VswitchId"])
	d.Set("instance_storage", data["InstanceStorage"])
	d.Set("deletion_protection", data["DeletionProtection"])
	d.Set("disk_usage", data["DiskUsage"])
	d.Set("enable_fs", data["EnableFS"])
	d.Set("switch_l_proxy_flag", data["SwitchLProxyFlag"])
	d.Set("switch_ssl_encryption_flag", data["SwitchSSLEncryptionFlag"])
	d.Set("ali_uid", data["AliUid"])
	d.Set("cpu_brand", data["CpuBrand"])
	// d.Set("disk_category", data["DiskCategory"])
	if serviceType, ok := data["ServiceType"].(string); ok {
		if serviceType == "lindorm_multizone" {
			d.Set("engine_type", "lindorm")
		} else {
			d.Set("engine_type", serviceType)
		}
	}
	d.Set("standby_zone_id", data["StandbyZoneId"])
	d.Set("arbiter_zone_id", data["ArbiterZoneId"])
	d.Set("standby_vswitch_id", data["StandbyVSwitchId"])
	d.Set("arbiter_vswitch_id", data["ArbiterVSwitchId"])
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
		d.Set("lindorm_num", engineData["CoreCount"])
		d.Set("instance_type", engineData["Specification"])
	}
	return nil
}

func resourceAlibabacloudStackLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
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
				"alibabacloudstack_lindorm_instance", "UpdateLindormInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
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

	if d.HasChanges("instance_type", "lindorm_num") {
		old_type, new_type := d.GetChange("instance_type")
		old_num, new_num := d.GetChange("lindorm_num")
		if old_type != new_type {
			reqQuery := map[string]interface{}{
				"InstanceId":          d.Id(),
				"ZoneId":              d.Get("zone_id"),
				"EngineType":          "lindorm",
				"EngineCoreNum":       old_num,
				"EngineSpec":          new_type,
				"SourceEngineCoreNum": old_num,
				"SourceEngineSpec":    old_type,
			}

			if _, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpgradeLindormInstanceEngine", "", nil, reqQuery, nil); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_lindorm_instance", "UpgradeLindormInstanceEngine", errmsgs.AlibabacloudStackSdkGoERROR)
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
				"EngineType":          "lindorm",
				"EngineCoreNum":       new_num,
				"EngineSpec":          new_type,
				"SourceEngineCoreNum": old_num,
				"SourceEngineSpec":    new_type,
			}

			if _, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpgradeLindormInstanceEngine", "", nil, reqQuery, nil); err != nil {
				return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
					"alibabacloudstack_lindorm_instance", "UpdateLindormInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
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

func resourceAlibabacloudStackLindormInstanceDelete(d *schema.ResourceData, meta interface{}) error {

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
