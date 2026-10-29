package alibabacloudstack

import (
	"context"
	"fmt"
	"log"
	"strconv"
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
		CustomizeDiff: func(_ context.Context, d *schema.ResourceDiff, _ interface{}) error {
			// mini 规格的本地磁盘最小为 100GB
			if spec, ok := d.GetOk("local_disk_spec"); ok && spec.(string) == "mini" {
				if size, ok := d.GetOk("local_disk_size"); ok && size.(int) < 100 {
					return fmt.Errorf("local_disk_size must be at least 100 when local_disk_spec is 'mini', got %d", size.(int))
				}
			}
			return nil
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
			"local_disk_spec": {
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
				Type:     schema.TypeInt,
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
	// topology_type 不由用户填写，根据 standby_zone_id 是否设置自动判定拓扑类型
	topologyType := "1azone"
	if v, ok := d.GetOk("standby_zone_id"); ok && v.(string) != "" {
		topologyType = "multizone"
		request["ClusterType"] = "lindorm_dual_zone"
		request["Idc2ZoneId"] = v
		if sv, ok := d.GetOk("standby_vswitch_id"); ok {
			request["Idc2VswitchId"] = sv
		}
	} else {
		request["ClusterType"] = "lindorm_single_zone"
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
	diskCategory := d.Get("disk_category").(string)
	localDiskNum := d.Get("local_disk_num").(int)
	request["LocalDiskNum"] = localDiskNum

	size := d.Get("local_disk_size").(int)
	storageSize := size
	request["CoreDiskSize"] = size
	request["CoreSingleStorage"] = size
	if diskCategory == "SSD" || diskCategory == "HDD" {
		request["InstanceStorage"] = strconv.Itoa(size)
		request["LocalDiskSize"] = size
	}

	// lindormStore.storage 参数
	request["lindormStore"] = map[string]interface{}{
		"storage": map[string]interface{}{
			"diskNum":      localDiskNum,
			"storageSize":  storageSize,
			"storageClass": diskCategory,
		},
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
	// 测试框架中 d.Timeout 可能返回 0，需兜底
	timeout := d.Timeout(schema.TimeoutCreate)
	if timeout <= 0 {
		timeout = 120 * time.Minute
	}
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
	if v, ok := data["ArchVersion"].(string); ok && v != "" {
		d.Set("arch_version", v)
	}
	if sv, ok := data["ServiceType"].(string); ok && sv == "lindorm_dual_zone" {
		d.Set("topology_type", "multizone")
	} else {
		d.Set("topology_type", "1azone")
	}
	engineList := data["EngineList"].([]interface{})
	if len(engineList) > 0 {
		engineData := engineList[0].(map[string]interface{})
		d.Set("lindorm_num", engineData["CoreCount"])
		d.Set("instance_type", engineData["Specification"])
	}
	// 多可用区场景下 API 可能返回 CoreSpec 和 LocalDiskSpecs
	if v, ok := data["CoreSpec"]; ok && v != nil && v != "" {
		d.Set("core_spec", v)
	}
	if v, ok := data["LocalDiskSpecs"]; ok && v != nil {
		if specs, ok := v.([]interface{}); ok && len(specs) > 0 {
			if specMap, ok := specs[0].(map[string]interface{}); ok {
				if specName, ok := specMap["SpecName"]; ok && specName != nil && specName != "" {
					d.Set("local_disk_spec", specName)
				}
			}
		}
	}
	return nil
}

func resourceAlibabacloudStackLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}

	// 新创建的资源不应触发 Update，避免实例后台初始化期间报 Instance.IsUpdating
	if d.IsNewResource() {
		return nil
	}

	if d.HasChanges("instance_alias", "deletion_protection") {
		reqQuery := map[string]interface{}{
			"InstanceId":         d.Id(),
			"InstanceAlias":      d.Get("instance_alias"),
			"DeletionProtection": d.Get("deletion_protection"),
		}

		// 实例激活后后台可能仍在变更，遇到 Instance.IsUpdating 时等待并重试
		updateTimeout := d.Timeout(schema.TimeoutUpdate)
		if updateTimeout <= 0 {
			updateTimeout = 20 * time.Minute
		}
		wait := connectivity.IncrementalWait(3*time.Second, 3*time.Second)
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
	if deleteTimeout <= 0 {
		deleteTimeout = 20 * time.Minute
	}
	stateConf := BuildStateConf([]string{"ACTIVATION", "DELETING", "CREATING"}, []string{}, deleteTimeout, 10*time.Second,
		lindormService.LindormInstanceStateRefreshFunc(d.Id(), []string{}))
	if _, err := stateConf.WaitForState(); err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.IdMsg, d.Id())
	}

	// 等待底层网络资源（ENI 等）完全释放，避免 VSwitch/VPC 删除时 DependencyViolation
	log.Printf("[INFO] Lindorm instance %s deleted, waiting 5 minutes for underlying resources to be released...", d.Id())
	time.Sleep(5 * time.Minute)

	return nil
}
