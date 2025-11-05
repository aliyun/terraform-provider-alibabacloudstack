package alibabacloudstack

import (
	"encoding/json"
	"sort"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2InstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2InstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
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
						"instance_class": {
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
						"spec_from": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uni_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2InstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
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
		"regionId":     client.RegionId,
		"status":       "Available",
		"resourceType": "COP",
	}

	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "SelectCommonSpec", "/ascm/manage/saleconf/commonSpec/select", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}
	v, ok := response["data"]
	if ok {
		for _, item := range v.([]interface{}) {
			data := item.(map[string]interface{})
			id := data["instanceClass"].(string)
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
			if v, ok := d.GetOk("cpu"); ok && int64(v.(int)) != cpu {
				continue
			}
			if v, ok := d.GetOk("memory"); ok && int64(v.(int)) != memory {
				continue
			}
			types = append(types, map[string]interface{}{
				"id":             data["instanceClass"],
				"cpu":            cpu,
				"memory":         memory,
				"instance_class": data["instanceClass"],
				"status":         data["status"],
				"product":        data["product"],
				"spec_from":      data["spec_from"],
				"uni_key":        data["uni_key"],
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
