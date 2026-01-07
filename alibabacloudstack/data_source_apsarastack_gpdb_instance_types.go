package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackGpdbInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackGpdbInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"sorted_by": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"CPU", "Memory"}, false),
			},
			"series": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"dual_ha", "read_only"}, false),
			},
			"instance_types": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"series": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"connections": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_instance_mode_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version_label": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"gmt_modify": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"spec_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackGpdbInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := map[string]string{}
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			filterIds[vv.(string)] = ""
		}
	}

	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"label":        "true",
		"resourceType": "gpdb",
		"status":       "Available",
	}
	if v, ok := d.GetOk("engine_version"); ok {
		reqQuery["engineVersion"] = v
	}

	if v, ok := d.GetOk("series"); ok {
		reqQuery["series"] = v
	}
	if v, ok := d.GetOk("cpu"); ok {
		reqQuery["cpu"] = v
	}
	if v, ok := d.GetOk("memory"); ok {
		reqQuery["memory"] = v
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	for _, d := range response["data"].([]interface{}) {
		data := d.(map[string]interface{})
		id := data["specification"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := existedId[id]; exists {
			continue
		}

		var cpu int64
		var memory string
		var storage string

		if cpuData, ok := data["cpu"]; ok {
			if cpuNum, ok := cpuData.(json.Number); ok {
				cpu, err = cpuNum.Int64()
				if err != nil {
					if cpuStr, ok := cpuData.(string); ok {
						cpu = parseCpuFromSpec(cpuStr)
					}
				}
			} else if cpuStr, ok := cpuData.(string); ok {
				cpu = parseCpuFromSpec(cpuStr)
			}
		}

		if memoryData, ok := data["memory"]; ok {
			if memoryNum, ok := memoryData.(json.Number); ok {
				memoryVal, err := memoryNum.Int64()
				if err != nil {
					memory = memoryData.(string)
				} else {
					memory = fmt.Sprintf("%d GB", memoryVal)
				}
			} else {
				memory = memoryData.(string)
			}
		}

		if storageData, ok := data["storage"]; ok {
			storage = storageData.(string)
		}

		var connections int
		if v, ok := data["connections"].(json.Number); ok {
			vv, err := v.Int64()
			if err != nil {
				return err
			}
			connections = int(vv)
		} else if v, ok := data["connections"].(string); ok {
			if v == "Unlimited" {
				connections = -1
			} else {
				vv, err := strconv.Atoi(v)
				if err != nil {
					return err
				}
				connections = vv
			}
		}

		typeMap := map[string]interface{}{
			"id":                     data["specification"],
			"cpu":                    cpu,
			"memory":                 memory,
			"engine_version":         data["engineVersionLabel"],
			"connections":            connections,
			"storage_min":            data["storageMin"],
			"storage_max":            data["storageMax"],
			"specification":          data["specification"],
			"specification_label":    data["specificationLabel"],
			"db_instance_mode":       data["dbInstanceMode"],
			"db_instance_mode_label": data["dbInstanceModeLabel"],
			"node":                   data["node"],
			"region_id":              data["regionId"],
			"status":                 data["status"],
			"product":                data["product"],
			"storage":                storage,
			"cpu_label":              data["cpuLabel"],
			"memory_label":           data["memoryLabel"],
			"storage_label":          data["storageLabel"],
			"engine_version_label":   data["engineVersionLabel"],
			"gmt_create":             data["gmtCreate"],
			"gmt_modify":             data["gmtModify"],
			"spec_from":              data["specFrom"],
		}

		types = append(types, typeMap)
		existedId[id] = ""
		ids = append(ids, id)
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return types[i]["cpu"].(int64) < types[j]["cpu"].(int64)
			case "Memory":
				memI := parseMemoryValue(types[i]["memory"].(string))
				memJ := parseMemoryValue(types[j]["memory"].(string))
				return memI < memJ
			}
			return false
		})
	}

	d.Set("ids", ids)
	d.Set("instance_types", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}

func parseCpuFromSpec(cpuStr string) int64 {
	re := regexp.MustCompile(`(\d+)\s*Core`)
	matches := re.FindStringSubmatch(cpuStr)
	if len(matches) > 1 {
		if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			return val
		}
	}
	return 0
}

func parseMemoryValue(memoryStr string) int64 {
	re := regexp.MustCompile(`(\d+)\s*GB`)
	matches := re.FindStringSubmatch(memoryStr)
	if len(matches) > 1 {
		if val, err := strconv.ParseInt(matches[1], 10, 64); err == nil {
			return val
		}
	}
	return 0
}
