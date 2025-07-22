package alibabacloudstack

import (
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDrdsInstanceSeries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDrdsInstanceSeriesRead,

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
			// Computed values.
			"series": {
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
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDrdsInstanceSeriesRead(d *schema.ResourceData, meta interface{}) error {
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
	series := []map[string]string{}

	reqQuery := map[string]interface{}{
		"Label":        "true",
		"resourceType": "DRDS",
		"status":       "Available",
	}
	reqHeader := map[string]string{
		"x-acs-territory": "US",
		"x-acs-lang":      "EN",
	}
	response, err := client.DoTeaRequest("POST", "ascm", "2019-05-10", "DescribeSeriesIdFamilies", "/ascm/manage/saleconf/drdsSpec/describeSeriesIdFamilies", reqHeader, reqQuery, nil)
	if err != nil {
		return err
	}

	for _, d := range response["data"].([]interface{}) {
		data := d.(map[string]interface{})
		id := data["seriesId"].(string)
		name := data["seriesIdLabel"].(string)
		if _, exists := filterIds[id]; len(filterIds) > 0 && !exists {
			continue
		}
		if _, exists := filterNames[name]; len(filterNames) > 0 && !exists {
			continue
		}
		if _, exists := existedId[id]; exists {
			continue
		}

		series = append(series, map[string]string{
			"id":   id,
			"name": name,
		})
		existedId[id] = ""
		ids = append(ids, id)
		names = append(names, name)
	}

	d.Set("ids", ids)
	d.Set("names", names)
	d.Set("series", series)
	d.SetId(dataResourceIdHash(ids))

	return nil
}
