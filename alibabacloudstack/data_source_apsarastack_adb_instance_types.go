package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"sort"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAdbInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAdbInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"cpu_type": {
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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
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
						"cpu_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
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
						"mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"node_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"cluster_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"series": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAdbInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := getIdsStringFilter(d)
	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}
	filterCpu := d.Get("cpu").(int)
	filterMemroy := d.Get("memory").(int)
	filterStatus := d.Get("status").(string)
	filterCpuType := d.Get("cpu_type").(string)

	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"label":        "true",
		"resourceType": "ADB",
		"status":       "Available",
	}

	if v, ok := d.GetOk("cluster_type"); ok && v.(string) != "" {
		reqQuery["clusterType"] = v
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	for {
		response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
		if err != nil {
			return err
		}

		if len(response["data"].([]interface{})) == 0 {
			break
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

			if filterStatus != "" && data["status"].(string) != filterStatus {
				continue
			}

			if filterCpuType != "" && data["cpuType"].(string) != filterCpuType {
				continue
			}

			var cpu, memory int

			if v, existed := data["storageResource"]; existed {
				spec := v.(string)
				re := regexp.MustCompile(`^(\d+)Core(\d+)GB$`)
				matches := re.FindStringSubmatch(spec)
				if matches == nil {
					continue
				}

				cpu, err = strconv.Atoi(matches[1])
				if err != nil {
					continue
				}

				memory, err = strconv.Atoi(matches[2])
				if err != nil {
					continue
				}
			} else {
				continue
			}

			if filterCpu != 0 && filterCpu != cpu {
				continue
			}
			if filterMemroy != 0 && filterMemroy != memory {
				continue
			}

			var storage_min, storage_max, node_min, node_max int
			storage_min, err = strconv.Atoi(data["localStorageMinDiskSize"].(string))
			storage_max, err = strconv.Atoi(data["localStorageMaxDiskSize"].(string))
			if val, exists := data["minNode"]; exists && val != nil {
				if jsonNum, ok := val.(json.Number); ok {
					if i64, err := jsonNum.Int64(); err == nil {
						node_min = int(i64)
					}
				}
			}
			if val, exists := data["maxNode"]; exists && val != nil {
				if jsonNum, ok := val.(json.Number); ok {
					if i64, err := jsonNum.Int64(); err == nil {
						node_max = int(i64)
					}
				}
			}

			typeMap := map[string]interface{}{
				"id":               data["specification"],
				"cpu_type":         data["cpuType"],
				"cpu":              cpu,
				"memory":           memory,
				"storage_type":     data["storageType"],
				"storage_min":      storage_min,
				"storage_max":      storage_max,
				"mode":             data["mode"],
				"node_min":         node_min,
				"node_max":         node_max,
				"cluster_type":     data["clusterType"],
				"status":           data["status"],
				"series":           data["series"],
				"cluster_category": data["series"],
			}

			types = append(types, typeMap)
			existedId[id] = ""
			ids = append(ids, id)
		}
		reqQuery["pageStart"] = reqQuery["pageStart"].(int) + 1
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return types[i]["cpu"].(int) < types[j]["cpu"].(int)
			case "Memory":
				return types[i]["memory"].(int) < types[j]["memory"].(int)
			}
			return false
		})
	}

	d.Set("ids", ids)
	d.Set("instance_types", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
