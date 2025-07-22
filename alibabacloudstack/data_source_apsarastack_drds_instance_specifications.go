package alibabacloudstack

import (
	"encoding/json"
	"sort"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackDrdsInstanceSpecifications() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDrdsInstanceSpecificationsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"series": {
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
			// Computed values.
			"specifications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
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
						"series": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDrdsInstanceSpecificationsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	filterIds := map[string]string{}
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			filterIds[vv.(string)] = ""
		}
	}
	filterNames := map[string]string{}
	if v, ok := d.GetOk("names"); ok {
		for _, vv := range v.([]interface{}) {
			filterNames[vv.(string)] = ""
		}
	}

	existedId := map[string]string{}
	ids := []string{}
	names := []string{}
	specifications := []map[string]interface{}{}

	reqQuery := map[string]interface{}{
		"pageStart":    1,
		"pageSize":     500,
		"status":       "Available",
		"resourceType": "DRDS",
	}
	if v, ok := d.GetOk("series"); ok {
		reqQuery["seriesId"] = v
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
		id := data["spec"].(string)
		name := data["specLabel"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := filterNames[name]; len(filterNames) > 0 && !exists {
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
		specifications = append(specifications, map[string]interface{}{
			"id":     id,
			"name":   name,
			"cpu":    cpu,
			"memory": memory,
			"series": data["seriesId"],
		})
		existedId[id] = ""
		ids = append(ids, id)
		names = append(names, name)
	}

	sortedBy := d.Get("sorted_by").(string)
	if sortedBy != "" {
		sort.SliceStable(specifications, func(i, j int) bool {
			switch sortedBy {
			case "CPU":
				return specifications[i]["cpu"].(int64) < specifications[j]["cpu"].(int64)
			case "Memory":
				return specifications[i]["memory"].(int64) < specifications[j]["memory"].(int64)
			}
			return false
		})
	}
	
	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("specifications", specifications)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
