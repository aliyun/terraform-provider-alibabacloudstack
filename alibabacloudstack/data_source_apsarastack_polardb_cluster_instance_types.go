package alibabacloudstack

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbClusterInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbClusterInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"db_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"cpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"intel", "arm64", "hygon"}, false),
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
			"sub_category": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"General", "Exclusive"}, false),
			},
			"db_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"MySQL", "PostgreSQL", "Oracle"}, false),
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
						"proxy_mem": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sub_category": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"memory": {
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
						"proxy_cpu": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"product": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_node_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"db_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpu": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"uni_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"proxy_class": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"gmt_create": {
							Type:     schema.TypeInt,
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
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbClusterInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
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
		"resourceType": "POLARDB",
		"status":       "Available",
	}
	if v, ok := d.GetOk("db_version"); ok {
		reqQuery["dbVersion"] = v
	}
	if v, ok := d.GetOk("cpu_type"); ok {
		reqQuery["cpuType"] = v
	}

	if v, ok := d.GetOk("db_type"); ok {
		reqQuery["dbType"] = v
	}
	if v, ok := d.GetOk("sub_category"); ok {
		sub_category := fmt.Sprintf("normal_%s", strings.ToLower(v.(string)))
		reqQuery["subCategory"] = sub_category
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	for _, v := range response["data"].([]interface{}) {
		data := v.(map[string]interface{})
		id := data["dbNodeClass"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := existedId[id]; exists {
			continue
		}
		cpu := data["cpu"].(string)
		cpu_count, err := strconv.Atoi(cpu)
		if err != nil {
			return errmsgs.WrapErrorf(err, "error data is :%#v", data)
		}
		memory := data["memory"].(string)
		memory_size, err := strconv.Atoi(memory)
		if err != nil {
			return errmsgs.WrapErrorf(err, "error data is :%#v", data)
		}
		if v, ok := d.GetOk("cpu"); ok && cpu_count != v.(int) {
			continue
		}
		if v, ok := d.GetOk("memory"); ok && memory_size != v.(int) {
			continue
		}

		types = append(types, map[string]interface{}{
			"id":            id,
			"cpu":           cpu_count,
			"memory":        memory_size,
			"proxy_mem":     data["proxyMem"],
			"sub_category":  data["subCategory"],
			"cpu_type":      data["cpuType"],
			"gmt_modify":    data["gmtModify"],
			"spec_from":     data["specFrom"],
			"proxy_cpu":     data["proxyCpu"],
			"product_type":  data["productType"],
			"db_node_class": data["dbNodeClass"],
			"product":       data["product"],
			"db_version":    data["dbVersion"],
			"proxy_type":    data["proxyType"],
			"db_node_num":   data["dbNodeNum"],
			"db_type":       data["dbType"],
			"uni_key":       data["uniKey"],
			"proxy_class":   data["proxyClass"],
			"gmt_create":    data["gmtCreate"],
			"region_id":     data["regionId"],
			"status":        data["status"],
		})
		existedId[id] = ""
		ids = append(ids, id)
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
