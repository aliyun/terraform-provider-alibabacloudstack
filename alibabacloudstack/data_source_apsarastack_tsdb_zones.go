package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	util "github.com/alibabacloud-go/tea-utils/service"
	"github.com/alibabacloud-go/tea/tea"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
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

	action := "DescribeZones"
	request := make(map[string]interface{})
	request["RegionId"] = client.RegionId
	request["Product"] = "hitsdb"
	var objects []map[string]interface{}
	var response map[string]interface{}
	conn, err := client.NewHitsdbClient()
	if err != nil {
		return WrapError(err)
	}
	runtime := util.RuntimeOptions{IgnoreSSL: tea.Bool(client.Config.Insecure)}
	runtime.SetAutoretry(true)
	response, err = conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2017-06-01"), StringPointer("AK"), nil, request, &runtime)
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_tsdb_zones", action, AlibabacloudStackSdkGoERROR)
	}
	addDebug(action, response, request)

	resp, err := jsonpath.Get("$.ZoneList.ZoneModel", response)
	if err != nil {
		return WrapErrorf(err, FailedGetAttributeMsg, action, "$.ZoneList.ZoneModel", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
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
		return WrapError(err)
	}

	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
