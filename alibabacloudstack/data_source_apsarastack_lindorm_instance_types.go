package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"sort"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackLindormInstanceTypes() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackLindormInstanceTypesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"engine_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"lindorm", "tsdb", "solr", "lts"}, false),
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
						"rate": {
							Type:     schema.TypeInt,
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

func dataSourceAlibabacloudStackLindormInstanceTypesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	request := client.NewCommonRequest("GET", "hitsdb", "2020-06-15", "DescribeLindormSpecInfo", "")
	response := make(map[string]interface{})
	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_lindorm_instance_types", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_lindorm_instance_types", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR)
	}
	lindorm_specs, err := jsonpath.Get("$.Data.LindormSpecs", response)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	ids := []string{}
	types := []map[string]interface{}{}
	for _, lindorm_spec := range lindorm_specs.([]interface{}) {
		data := lindorm_spec.(map[string]interface{})
		engineType := data["EngineType"].(string)
		specs := data["Specs"].([]interface{})
		for _, s := range specs {
			spec := s.(map[string]interface{})
			id := fmt.Sprintf("%s:%s", engineType, spec["Name"].(string))
			if _, exists := idsMap[id]; len(idsMap) > 0 && !exists {
				continue
			}
			if cpu, ok := d.GetOk("cpu"); ok {
				if cpu.(int) != int(spec["CpuCount"].(float64)) {
					continue
				}
			}
			if memory, ok := d.GetOk("memory"); ok {
				if memory.(int) != int(spec["MemorySize"].(float64)) {
					continue
				}
			}
			types = append(types, map[string]interface{}{
				"id":     id,
				"cpu":    int(spec["CpuCount"].(float64)),
				"memory": int(spec["MemorySize"].(float64)),
				"rate":   spec["Rate"],
				"name":   spec["Name"],
			})
			ids = append(ids, id)
		}
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
