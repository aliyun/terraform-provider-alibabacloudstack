package alibabacloudstack

import (
	"fmt"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackWafInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackWafInstancesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_make_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_vswitch": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vswitch_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"detector_specs": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detector_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"detector_nodenum": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
}

func dataSourceAlibabacloudStackWafInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeWAFInstance"
	request := make(map[string]interface{})
	var objects []map[string]interface{}

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}
	var response map[string]interface{}
	var err error
	response, err = client.DoTeaRequest("GET", "waf-onecs", "2020-07-01", action, "", nil, nil, request)
	if err != nil {
		err = errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg, "", action, errmsgs.AlibabacloudStackSdkGoERROR)
		return err
	}
	addDebug(action, response, request)
	result, err := jsonpath.Get("$.Result", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "", "$.Result", response)
	}
	resp, err := jsonpath.Get("$.items", result)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "", "$.items", response)
	}
	for _, v := range resp.([]interface{}) {
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(v.(map[string]interface{})["instance_id"])]; !ok {
				continue
			}
		}
		objects = append(objects, v.(map[string]interface{}))
	}
	ids := make([]string, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"name":                 object["name"],
			"instance_status":      object["instance_status"],
			"vpc_vswitch":          object["vpc_vswitch"],
			"detector_version":     object["detector_version"],
			"instance_make_status": object["instance_make_status"],
			"detector_specs":       object["detector_specs"],
			"detector_nodenum":     object["detector_nodenum"],
		}
		ids = append(ids, fmt.Sprint(object["instance_id"]))
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	// if err := d.Set("instances", s); err != nil {
	// 	return errmsgs.WrapError(err)
	// }
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}

// d.Set("name", object["name"])
// d.Set("instance_status", object["instance_status"])
// d.Set("vpc_vswitch", object["vpc_vswitch"])
// d.Set("detector_version", object["detector_version"])
// d.Set("instance_make_status", object["instance_make_status"])
// d.Set("detector_specs", object["detector_specs"])
// d.Set("detector_nodenum", object["detector_node_num"])
