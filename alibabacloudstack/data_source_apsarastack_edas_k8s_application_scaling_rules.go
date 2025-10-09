package alibabacloudstack

import (
	"fmt"
	"regexp"
	"time"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEdasScalingRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEdasScalingRulesRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"names": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"scaling_rule_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"metric", "trigger"}, false),
			},
			"scaling_rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scaling_rule_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"scaling_rule_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"max_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"min_replicas": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"metrics": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"type": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"utilization": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"trigger_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_period": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"trigger_dryrun": {
							Type:     schema.TypeBool,
							Computed: true,
						},

						"trigger_timer_in_day": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"at_time": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"replicas": {
										Type:     schema.TypeInt,
										Computed: true,
									},
								},
							},
						},
						"trigger_timer_in_week": {
							Type:     schema.TypeSet,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEdasScalingRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "ListUserDefineRegion"
	request := make(map[string]interface{})
	var objects []map[string]interface{}
	var namespaceNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		namespaceNameRegex = r
	}

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
	wait := incrementalWait(3*time.Second, 3*time.Second)
	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		response, err = client.DoTeaRequest("POST", "Edas", "2017-08-01", action, "/pop/v5/user_region_defs", nil, request, nil)
		if err != nil {
			if errmsgs.NeedRetry(err) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})
	addDebug(action, response, request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DataDefaultErrorMsg, "alibabacloudstack_edas_namespaces", action, errmsgs.AlibabacloudStackSdkGoERROR)
	}
	resp, err := jsonpath.Get("$.UserDefineRegionList.UserDefineRegionEntity", response)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, action, "$.UserDefineRegionList.UserDefineRegionEntity", response)
	}
	result, _ := resp.([]interface{})
	for _, v := range result {
		item := v.(map[string]interface{})
		if namespaceNameRegex != nil && !namespaceNameRegex.MatchString(fmt.Sprint(item["RegionName"])) {
			continue
		}
		if len(idsMap) > 0 {
			if _, ok := idsMap[fmt.Sprint(item["Id"])]; !ok {
				continue
			}
		}
		objects = append(objects, item)
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			// 			"debug_enable":         object["DebugEnable"],
			"description":          object["Description"],
			"id":                   fmt.Sprint(object["Id"]),
			"namespace_id":         fmt.Sprint(object["Id"]),
			"namespace_logical_id": object["RegionId"],
			"namespace_name":       object["RegionName"],
			"user_id":              object["UserId"],
			"belong_region":        object["BelongRegion"],
		}
		ids = append(ids, fmt.Sprint(mapping["id"]))
		names = append(names, object["RegionName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("namespaces", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}

	return nil
}
