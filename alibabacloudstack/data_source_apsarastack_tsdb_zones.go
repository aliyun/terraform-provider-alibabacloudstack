package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackTsdbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackTsdbZonesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						// 3.16 不支持该参数
						//						"local_name": {
						//							Type:     schema.TypeString,
						//							Computed: true,
						//						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackTsdbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["Product"] = "hitsdb"
	response, err := client.DoTeaRequest("POST", "hitsdb", "2017-06-01", "DescribeZones", "", nil, request)
	if err != nil {
		return err
	}

	resp, err := jsonpath.Get("$.ZoneList.ZoneModel", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeZones", "$.ZoneList.ZoneModel", response)
	}
	result, _ := resp.([]interface{})
	objects := make([]map[string]interface{}, len(result))
	for i, v := range result {
		item := v.(map[string]interface{})
		objects[i] = item
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"id":      fmt.Sprint(object["ZoneId"]),
			"zone_id": object["ZoneId"],
			//"local_name": object["LocalName"],
		}
		ids = append(ids, fmt.Sprint(object["ZoneId"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("zones", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
