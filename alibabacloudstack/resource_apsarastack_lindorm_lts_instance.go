package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackLindormLtsInstance() *schema.Resource {
	resource := &schema.Resource{
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(40 * time.Minute),
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
	// LTS 实例不支持多可用区（云产品口径），固定按单可用区传参
	request["TopologyType"] = "1azone"
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
	// DeletionProtection 在实例创建初期可能不返回，nil 时不要覆盖 state，避免触发不必要的 Update
	if v, ok := data["DeletionProtection"]; ok && v != nil {
		d.Set("deletion_protection", v)
	}
	d.Set("cpu_brand", data["CpuBrand"])
	d.Set("service_type", data["ServiceType"])
	d.Set("network_type", data["NetworkType"])
	// GetLindormInstance 对 LTS 实例不返回 VpcId/VswitchId 字段；
	// 这些属性是 Optional+ForceNew，若 d.Set(nil) 会清空 state 导致下次 plan 非空并触发 replace。
	// 仅在响应确实包含时写入，否则保持 Create 时已写入 state 的值。
	if v, ok := data["VpcId"]; ok && v != nil {
		d.Set("vpc_id", v)
	}
	if v, ok := data["VswitchId"]; ok && v != nil {
		d.Set("vswitch_id", v)
	}
	if raw, ok := data["EngineList"]; ok && raw != nil {
		if engineList, ok := raw.([]interface{}); ok && len(engineList) > 0 {
			if engineData, ok := engineList[0].(map[string]interface{}); ok {
				d.Set("lts_num", engineData["CoreCount"])
				d.Set("instance_type", engineData["Specification"])
			}
		}
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
	lindormService := LindormService{client}
	reqQuery := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	// Call the release API to delete the instance
	_, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "ReleaseLindormInstance", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ReleaseLindormInstance", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	// 与 lindorm_instance 对齐：等实例从 state 消失再返回，避免下游 VSwitch/VPC 立刻
	// 删除时撞上 ENI 未释放导致的 DependencyViolation。
	deleteTimeout := d.Timeout(schema.TimeoutDelete)
	stateConf := BuildStateConf([]string{"ACTIVATION", "DELETING", "CREATING"}, []string{}, deleteTimeout, 10*time.Second,
		lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	return nil
}
