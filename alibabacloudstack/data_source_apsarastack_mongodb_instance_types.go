package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackMongoDBInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackMongodbInstanceTypesRead,

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
			"db_instnace_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"sharding", "replicate"}, false),
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"intel", "arm64"}, false),
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
			"node_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"configserver", "shard", "mongos"}, false),
			},
			// Computed values.
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
							Type:     schema.TypeInt,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"spec": {
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
						"storage_min": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"storage_max": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackMongodbInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	if d.Get("db_instnace_type").(string) == "sharding" {
		if _, ok := d.GetOk("node_type"); !ok {
			return fmt.Errorf("node_type is necessary while db_instnace_type is sharding.")
		}
	}

	filterIds := getIdsStringFilter(d)

	existedId := map[string]string{}
	ids := []string{}
	types := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"pageStart":      1,
		"pageSize":       500,
		"label":          "true",
		"resourceType":   "dds",
		"status":         "Available",
		"dbInstanceType": d.Get("db_instnace_type"),
		"engine":         "MongoDB",
	}
	if v, ok := d.GetOk("engine_version"); ok {
		reqQuery["engineVersion"] = v
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		reqQuery["cpuType"] = v
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
	if v, ok := d.GetOk("node_type"); ok {
		reqQuery["nodeType"] = v
	}
	territory := "US"
	lang := "EN"
	reqHeader := make(map[string]string)
	reqHeader["x-acs-territory"] = territory
	reqHeader["x-acs-lang"] = lang
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	if v, ok := response["data"].([]interface{}); ok {
		for _, d := range v {
			data := d.(map[string]interface{})
			id := data["spec"].(string)
			if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
				continue
			}
			if _, exists := existedId[id]; exists {
				continue
			}
			cpu, err := data["cpu"].(json.Number).Int64()
			if err != nil {
				return err
			}
			memory, err := data["memory"].(json.Number).Int64()
			if err != nil {
				return err
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
			storageMin, err := data["storageMin"].(json.Number).Int64()
			if err != nil {
				return err
			}
			storageMax, err := data["storageMax"].(json.Number).Int64()
			if err != nil {
				return err
			}
			types = append(types, map[string]interface{}{
				"id":             id,
				"cpu":            cpu,
				"memory":         memory,
				"series":         data["seriesId"],
				"spec":           data["spec"],
				"engine_version": data["engineVersionLabel"],
				"cpu_type":       data["cpuType"],
				"connections":    connections,
				"storage_min":    storageMin,
				"storage_max":    storageMax,
			})
			existedId[id] = ""
			ids = append(ids, id)
		}
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(types, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return types[i]["cpu"].(int64) < types[j]["cpu"].(int64)
			case "Memory":
				return types[i]["memory"].(int64) < types[j]["memory"].(int64)
			}
			return false
		})
	}

	d.Set("ids", ids)
	d.Set("instance_types", types)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
