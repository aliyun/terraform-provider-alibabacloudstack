package alibabacloudstack

import (
	"fmt"
	"time"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceAlibabacloudStackLindormInstance() *schema.Resource {
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
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_zone_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"standby_vswitch_id"},
			},
			"standby_vswitch_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				RequiredWith: []string{"standby_zone_id"},
			},
			"lindorm_num": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"core_spec": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"local_disk_num": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 10),
			},
			"local_disk_size": {
				Type:         schema.TypeInt,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.IntBetween(100, 5000),
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
				Default:  false,
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
	// topology_type 不由用户填写，根据 standby_zone_id 是否设置自动判定拓扑类型
	topologyType := "1azone"
	if v, ok := d.GetOk("standby_zone_id"); ok && v.(string) != "" {
		topologyType = "multizone"
		request["ClusterType"] = "lindorm_dual_zone"
		request["Idc2ZoneId"] = v
		if sv, ok := d.GetOk("standby_vswitch_id"); ok {
			request["Idc2VswitchId"] = sv
		}
	}
	request["ZoneId"] = d.Get("zone_id")
	request["InstanceAlias"] = d.Get("instance_alias")
	request["CpuBrand"] = d.Get("cpu_brand")
	request["DiskCategory"] = d.Get("disk_category")
	request["EngineType"] = d.Get("engine_type")
	request["LindormSpec"] = d.Get("instance_type")
	request["VPCId"] = d.Get("vpc_id")
	request["VSwitchId"] = d.Get("vswitch_id")
	request["LindormNum"] = d.Get("lindorm_num")
	if v, ok := d.GetOk("core_spec"); ok {
		request["CoreSpec"] = v
	}
	request["LocalDiskNum"] = d.Get("local_disk_num")

	// 磁盘容量按 disk_category 映射到服务端读不同的字段：
	// LV -> storageClass=yoda-lvm，diskNum 固定 1，容量取 CoreDiskSize（数字，服务端补 "Gi"）
	// HDD/SSD -> fast-disks(-ssd)，diskNum=LocalDiskNum，容量取 LocalDiskSize（带单位字符串）
	// InstanceStorage 是服务端必传参数，不区分磁盘类型
	diskCategory := d.Get("disk_category").(string)
	size := d.Get("local_disk_size").(int)
	request["CoreDiskSize"] = size
	request["InstanceStorage"] = size
	if diskCategory == "SSD" || diskCategory == "HDD" {
		request["LocalDiskSize"] = fmt.Sprintf("%dG", size)
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
	d.Set("topology_type", topologyType)
	// Wait for instance to be active
	timeout := d.Timeout(schema.TimeoutCreate)
	stateConf := BuildStateConf([]string{"CREATING"}, []string{"ACTIVATION"}, timeout, 60*time.Second, lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"CREATE_FAILED"}))
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
	// DeletionProtection 在实例创建初期可能不返回，nil 时不要覆盖 state，避免触发不必要的 Update
	if v, ok := data["DeletionProtection"]; ok && v != nil {
		d.Set("deletion_protection", v)
	}
	d.Set("disk_usage", data["DiskUsage"])
	d.Set("enable_fs", data["EnableFS"])
	d.Set("switch_l_proxy_flag", data["SwitchLProxyFlag"])
	d.Set("switch_ssl_encryption_flag", data["SwitchSSLEncryptionFlag"])
	d.Set("ali_uid", data["AliUid"])
	d.Set("cpu_brand", data["CpuBrand"])
	// DiskCategory 在响应中可能为空（如 LV 类型），不强制写入避免 diff
	if v, ok := data["DiskCategory"]; ok && v != nil && v != "" {
		d.Set("disk_category", v)
	}
	// ServiceType 在多可用区场景下返回 "lindorm_dual_zone"，需映射回用户配置的 "lindorm"
	engineType := data["ServiceType"]
	if v, ok := engineType.(string); ok && v == "lindorm_dual_zone" {
		engineType = "lindorm"
	}
	d.Set("engine_type", engineType)
	if sv, ok := data["ServiceType"].(string); ok && sv == "lindorm_dual_zone" {
		d.Set("topology_type", "multizone")
	} else {
		d.Set("topology_type", "1azone")
	}
	// 多可用区回读：云产品口径 idc1=主可用区、idc2=备可用区；单可用区场景响应不含该 key，不覆盖 state
	if v, ok := data["Idc2ZoneId"]; ok && v != nil && v != "" {
		d.Set("standby_zone_id", v)
	}
	if v, ok := data["Idc2VSwitchId"]; ok && v != nil && v != "" {
		d.Set("standby_vswitch_id", v)
	}
	// EngineList 在部分响应中可能缺失；与 LTS 资源一致，先判存在再做类型断言，
	// 避免对 nil / 非 slice 直接断言导致 panic。
	if raw, ok := data["EngineList"]; ok && raw != nil {
		if engineList, ok := raw.([]interface{}); ok && len(engineList) > 0 {
			if engineData, ok := engineList[0].(map[string]interface{}); ok {
				d.Set("lindorm_num", engineData["CoreCount"])
				d.Set("instance_type", engineData["Specification"])
			}
		}
	}
	// 多可用区场景下 API 可能返回 CoreSpec
	if v, ok := data["CoreSpec"]; ok && v != nil && v != "" {
		d.Set("core_spec", v)
	}
	return nil
}

func resourceAlibabacloudStackLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}

	// DeletionProtection 不在 CreateLindormInstance 的入参中（已核对环境注册的 API XML），
	// 只能由 UpdateLindormInstanceAttribute 设置；setResourceFunc 的包装器在 create 成功后会自动
	// 调用本函数，因此该分支必须放在新建判定之前，否则用户新建时配置的 deletion_protection 会被丢弃。
	// instance_alias 已由 CreateLindormInstance 直接下发，故新建场景下本段 update 唯一有实际意义的用途
	// 就是补设 deletion_protection；而新建时 d.HasChanges 恒为 true（旧 state 为空），若不加限定，会对一个
	// 与服务端默认一致的 deletion_protection=false 发出空操作写请求，恰好撞上实例刚激活、后台仍在变更时的
	// Instance.IsUpdating 窗口。因此新建时仅当 deletion_protection 被显式开启才执行，非新建（真实变配）时保持原有的按变更触发逻辑。
	needUpdate := d.HasChanges("instance_alias", "deletion_protection")
	if d.IsNewResource() {
		needUpdate = d.Get("deletion_protection").(bool)
	}
	if needUpdate {
		reqQuery := map[string]interface{}{
			"InstanceId":         d.Id(),
			"InstanceAlias":      d.Get("instance_alias"),
			"DeletionProtection": d.Get("deletion_protection"),
		}

		// 实例激活后后台可能仍在变更，遇到 Instance.IsUpdating 时等待并重试
		updateTimeout := d.Timeout(schema.TimeoutUpdate)
		wait := connectivity.IncrementalWait(5*time.Second, 5*time.Second)
		err := resource.Retry(updateTimeout, func() *resource.RetryError {
			_, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "UpdateLindormInstanceAttribute", "", nil, reqQuery, nil)
			if err != nil {
				if errmsgs.IsExpectedErrors(err, "Instance.IsUpdating") {
					wait()
					return resource.RetryableError(err)
				}
				return resource.NonRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg,
				"alibabacloudstack_lindorm_instance", "UpdateLindormInstanceAttribute", errmsgs.AlibabacloudStackSdkGoERROR)
		}

		// Wait for the instance to be updated
		stateConf := BuildStateConf([]string{"UPDATING"}, []string{"ACTIVATION"}, updateTimeout, 3*time.Second,
			lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{"FAILED"}))
		if _, err := stateConf.WaitForState(); err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
		}
	}

	// 变配类接口只针对已存在的资源，新建后的首次自动调用必须跳过，
	// 避免实例后台仍在初始化时报 Instance.IsUpdating
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
	lindormService := LindormService{client}
	reqQuery := map[string]interface{}{
		"InstanceId": d.Id(),
	}
	// Call the release API to delete the instance
	_, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "ReleaseLindormInstance", "", nil, reqQuery, nil)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, d.Id(), "ReleaseLindormInstance", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	// Wait for the instance to be deleted
	deleteTimeout := d.Timeout(schema.TimeoutDelete)
	stateConf := BuildStateConf([]string{"ACTIVATION", "DELETING", "CREATING"}, []string{}, deleteTimeout, 10*time.Second,
		lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	// 底层 ENI 释放带来的 VSwitch/VPC 依赖冲突，由 VSwitch/VPC 自身的 Delete 里
	// resource.Retry(TimeoutDelete) 应对 DependencyViolation 类错误重试，此处不再硬编码等待。
	return nil
}
