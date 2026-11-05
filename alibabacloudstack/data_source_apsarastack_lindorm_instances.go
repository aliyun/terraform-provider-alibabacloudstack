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

	// InstanceList may be absent when the environment has no instances; asserting nil directly would panic.
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

		// InstanceId is the primary key for filtering and output; skip the entry if it is missing.
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

		// service_type: the raw value from the List response, e.g. "lindorm_dual_zone" or "lindorm".
		serviceType, _ := data["ServiceType"].(string)
		// engine_type: aligned with the resource Read; "lindorm_dual_zone" is normalized back to the base engine name "lindorm",
		// and in all other cases the ServiceType string value is taken directly.
		engineType := serviceType
		if engineType == "lindorm_dual_zone" {
			engineType = "lindorm"
		}
		// topology_type: derived from ServiceType, same as in the resource Read;
		// the List response has no ArchVersion field, so it cannot be used to determine the topology.
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
			// The List response uses Idc2ZoneId/Idc2VSwitchId (not StandbyZoneId/StandbyVSwitchId);
			// for single-zone instances these two fields are absent, so the result is nil -> empty string (expected).
			"standby_zone_id":    data["Idc2ZoneId"],
			"arbiter_zone_id":    data["ArbiterZoneId"],
			"standby_vswitch_id": data["Idc2VSwitchId"],
			"arbiter_vswitch_id": data["ArbiterVSwitchId"],
			"topology_type":      topologyType,
		}
		// arch_version is not present in the List response; write it only when it has a value, otherwise keep the zero value.
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
