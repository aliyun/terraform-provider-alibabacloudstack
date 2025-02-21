package alibabacloudstack

import (
	"strconv"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackNasZones() *schema.Resource {
	return &schema.Resource{
		Read:    dataSourceAlibabacloudStackNasZonesRead,
		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"storage_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"protocol_type": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackNasZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	response, err := client.DoTeaRequest("POST", "Nas", "2017-06-26", "DescribeZones", "", nil, nil, request)
	if err != nil {
		return err
	}

	resp, err := jsonpath.Get("$.Zones.Zone", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeZones", "$.Zones.Zone", response)
	}
	result, _ := resp.([]interface{})
	objects := make([]map[string]interface{}, 0)
	for _, v := range result {
		item := v.(map[string]interface{})
		objects = append(objects, item)
	}
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"zone_id": object["ZoneId"],
		}
		InstanceTypes := make([]map[string]interface{}, 0)
		if InstanceTypesList, ok := object["InstanceTypes"]; ok {
			if data, ok := InstanceTypesList.(map[string]interface{}); ok {
				for _, v := range data {
					if m1, ok := v.([]interface{}); ok {
						for _, vv := range m1 {
							if res, ok := vv.(map[string]interface{}); ok {
								temp1 := map[string]interface{}{
									"storage_type": res["StorageType"],
									"protocol_type": res["ProtocolType"],
								}
								InstanceTypes = append(InstanceTypes, temp1)
							}
						}
					}
				}
			}
		}
		mapping["instance_types"] = InstanceTypes
		s = append(s, mapping)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 16))

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
