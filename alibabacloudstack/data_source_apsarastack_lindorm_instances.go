package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackLindormInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackLindormInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				Deprecated:    "Field 'name_regex' is deprecated and will be removed in a future release. Please use new field 'description_regex' instead.",
				ConflictsWith: []string{"description_regex"},
			},
			"description_regex": {
				Type:          schema.TypeString,
				Optional:      true,
				ValidateFunc:  validation.StringIsValidRegExp,
				ConflictsWith: []string{"name_regex"},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_brand": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_storage": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ascm_create_user": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ali_uid": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"topology_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arch_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standby_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arbiter_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"standby_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"arbiter_vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackLindormInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := map[string]interface{}{
		"PageSize":   10000,
		"PageNumber": 1,
		"Department": client.Department,
		"RegionId":   client.RegionId,
	}
	response, err := client.DoTeaRequest("POST", "hitsdb", "2020-06-15", "GetLindormInstanceList", "", nil, request, nil)
	if err != nil {
		if response == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		return errmsgs.WrapError(err)
	}
	idsMap := getIdsStringFilter(d)

	// InstanceList 在环境无实例时可能不存在；直接断言 nil 会 panic。
	var instanceList []interface{}
	if v, ok := response["InstanceList"]; ok && v != nil {
		list, ok := v.([]interface{})
		if !ok {
			return fmt.Errorf("lindorm_instances: unexpected type %T for InstanceList in GetLindormInstanceList response", v)
		}
		instanceList = list
	}

	var ids []string
	datas := make([]interface{}, 0)
	for _, v := range instanceList {
		data, ok := v.(map[string]interface{})
		if !ok {
			continue
		}

		// InstanceId 是过滤及输出的主键，缺失则跳过该条目。
		instanceId, ok := data["InstanceId"].(string)
		if !ok || instanceId == "" {
			continue
		}

		if descriptionRegex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(descriptionRegex.(string))
			alias, _ := data["InstanceAlias"].(string)
			if !r.MatchString(alias) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[instanceId]; !exist {
				continue
			}
		}

		// service_type: List 响应的原始值，例如 "lindorm_dual_zone" 或 "lindorm"。
		serviceType, _ := data["ServiceType"].(string)
		// engine_type: 与资源 Read 对齐——将 "lindorm_dual_zone" 归一化回基础引擎名 "lindorm"，
		// 其余情况直接取 ServiceType 字符串值。
		engineType := serviceType
		if engineType == "lindorm_dual_zone" {
			engineType = "lindorm"
		}
		// topology_type: 与资源 Read 同样从 ServiceType 派生；
		// List 响应中无 ArchVersion 字段，不能用它判断拓扑。
		topologyType := "1azone"
		if serviceType == "lindorm_dual_zone" {
			topologyType = "multizone"
		}

		i := map[string]interface{}{
			"id":               instanceId,
			"cpu_brand":        data["CpuBrand"],
			"instance_storage": data["InstanceStorage"],
			"zone_id":          data["ZoneId"],
			"instance_id":      instanceId,
			"create_time":      data["CreateTime"],
			"ascm_create_user": data["AscmCreateUser"],
			"instance_alias":   data["InstanceAlias"],
			"network_type":     data["NetworkType"],
			"service_type":     serviceType,
			"engine_type":      engineType,
			"ali_uid":          data["AliUid"],
			// List 响应使用 Idc2ZoneId/Idc2VSwitchId（非 StandbyZoneId/StandbyVSwitchId）；
			// 单可用区实例这两个字段不存在，结果为 nil → 空字符串（正常）。
			"standby_zone_id":    data["Idc2ZoneId"],
			"arbiter_zone_id":    data["ArbiterZoneId"],
			"standby_vswitch_id": data["Idc2VSwitchId"],
			"arbiter_vswitch_id": data["ArbiterVSwitchId"],
			"topology_type":      topologyType,
		}
		// arch_version 在 List 响应中不存在；仅当有值时写入，否则保持零值。
		if av, ok := data["ArchVersion"].(string); ok {
			i["arch_version"] = av
		}

		datas = append(datas, i)
		ids = append(ids, instanceId)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", datas); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	return nil
}
