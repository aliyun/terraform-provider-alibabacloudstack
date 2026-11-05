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
	// topology_type is not filled in by the user; the topology type is determined automatically based on whether standby_zone_id is set
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

	// Disk capacity maps to different fields read by the server depending on disk_category:
	// LV -> storageClass=yoda-lvm, diskNum fixed at 1, capacity taken from CoreDiskSize (a number; the server appends "Gi")
	// HDD/SSD -> fast-disks(-ssd), diskNum=LocalDiskNum, capacity taken from LocalDiskSize (a string with a unit)
	// InstanceStorage is a required server-side parameter regardless of disk type; as long as DiskCategory is non-empty it will be
	// overridden downstream and never reach the CR, so it serves as a fallback initial value. Because the VO is a String passed through verbatim and DiskCategory can be empty,
	// we always append a "G" suffix here, to avoid a bare number being MustParse'd as bytes by K8s when DiskCategory is empty.
	diskCategory := d.Get("disk_category").(string)
	size := d.Get("local_disk_size").(int)
	request["CoreDiskSize"] = size
	request["InstanceStorage"] = fmt.Sprintf("%dG", size)
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
	// DeletionProtection may not be returned early after instance creation; do not overwrite the state when it is nil, to avoid triggering an unnecessary Update
	if v, ok := data["DeletionProtection"]; ok && v != nil {
		d.Set("deletion_protection", v)
	}
	d.Set("disk_usage", data["DiskUsage"])
	d.Set("enable_fs", data["EnableFS"])
	d.Set("switch_l_proxy_flag", data["SwitchLProxyFlag"])
	d.Set("switch_ssl_encryption_flag", data["SwitchSSLEncryptionFlag"])
	d.Set("ali_uid", data["AliUid"])
	d.Set("cpu_brand", data["CpuBrand"])
	// DiskCategory may be empty in the response (e.g. for the LV type); do not force-write it, to avoid a diff
	if v, ok := data["DiskCategory"]; ok && v != nil && v != "" {
		d.Set("disk_category", v)
	}
	// ServiceType returns "lindorm_dual_zone" in multi-zone scenarios; map it back to the user-configured "lindorm"
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
	// Multi-zone read back: per the cloud product's definition, idc1=primary zone and idc2=standby zone; in single-zone scenarios the response does not contain these keys, so the state is not overwritten
	if v, ok := data["Idc2ZoneId"]; ok && v != nil && v != "" {
		d.Set("standby_zone_id", v)
	}
	if v, ok := data["Idc2VSwitchId"]; ok && v != nil && v != "" {
		d.Set("standby_vswitch_id", v)
	}
	// EngineList may be missing in some responses; as with the LTS resource, check existence before doing the type assertion,
	// to avoid a panic from asserting directly on nil / a non-slice.
	if raw, ok := data["EngineList"]; ok && raw != nil {
		if engineList, ok := raw.([]interface{}); ok && len(engineList) > 0 {
			if engineData, ok := engineList[0].(map[string]interface{}); ok {
				d.Set("lindorm_num", engineData["CoreCount"])
				d.Set("instance_type", engineData["Specification"])
			}
		}
	}
	// The API may return CoreSpec in multi-zone scenarios
	if v, ok := data["CoreSpec"]; ok && v != nil && v != "" {
		d.Set("core_spec", v)
	}
	return nil
}

func resourceAlibabacloudStackLindormInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	lindormService := LindormService{client}

	// DeletionProtection is not among the input parameters of CreateLindormInstance (verified against the API XML registered in the environment);
	// it can only be set by UpdateLindormInstanceAttribute. The setResourceFunc wrapper automatically
	// calls this function after a successful create, so this branch must come before the new-resource check, otherwise the deletion_protection configured at creation would be dropped.
	// instance_alias is already sent directly by CreateLindormInstance, so in the create scenario the only practically meaningful use of this update block
	// is to set deletion_protection. On create, d.HasChanges is always true (the old state is empty); without a guard, it would issue a no-op write request
	// for deletion_protection=false that matches the server-side default, hitting exactly the Instance.IsUpdating window right after the instance
	// activates while background changes are still running. Therefore, on create this runs only when deletion_protection is explicitly enabled; for non-create (real configuration changes), the original change-triggered logic is kept.
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

		// Background changes may still be running after the instance activates; wait and retry on Instance.IsUpdating
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

	// Configuration-change APIs only apply to existing resources; the first automatic call after create must be skipped,
	// to avoid Instance.IsUpdating while the instance is still initializing in the background
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

	// Dependency violations on VSwitch/VPC caused by underlying ENI release are handled by the
	// resource.Retry(TimeoutDelete) in the VSwitch/VPC Delete itself, which retries DependencyViolation-class errors; no hardcoded wait here.
	return nil
}
