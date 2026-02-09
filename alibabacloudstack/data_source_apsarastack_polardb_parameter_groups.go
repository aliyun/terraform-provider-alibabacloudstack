package alibabacloudstack

import (
	"encoding/json"
	"regexp"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackPolardbParameterGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackPolardbParameterGroupsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: false,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"engine": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: false,
			},
			"parameter_group_type": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: false,
			},
			// Computed top-level attributes (only the list of instances)
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_desc": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameter_group_type": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"force_restart": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"param_counts": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"created": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"modified": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"parameters": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"param_name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"param_value": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAlibabacloudStackPolardbParameterGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeParameterGroups"

	request := client.NewCommonRequest("POST", "polardb", "2024-01-30", action, "")
	wait := incrementalWait(3*time.Second, 5*time.Second)
	var bresponse *responses.CommonResponse
	var err error

	err = resource.Retry(5*time.Minute, func() *resource.RetryError {
		bresponse, err = client.ProcessCommonRequest(request)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.ThrottlingUser}) {
				wait()
				return resource.RetryableError(err)
			}
			return resource.NonRetryableError(err)
		}
		return nil
	})

	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_polardb_parameter_groups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	var response DescribeParameterGroupsResponse
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Create filters based on schema
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Filter results
	var filteredGroups []ParameterGroupItem
	for _, group := range response.ParameterGroups.ParameterGroup {

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(group.ParameterGroupName) {
			continue
		}

		// Apply ids filter
		if len(idsMap) > 0 {
			if _, ok := idsMap[group.ParamGroupId]; !ok {
				continue
			}
		}

		filteredGroups = append(filteredGroups, group)
	}

	// Prepare result data
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	groups := make([]map[string]interface{}, 0)

	for _, group := range filteredGroups {
		mapping := map[string]interface{}{
			"id":                   group.ParamGroupId,
			"parameter_group_id":   group.ParamGroupId,
			"parameter_group_name": group.ParameterGroupName,
			"parameter_group_desc": group.ParameterGroupDesc,
			"engine":               group.Engine,
			"engine_version":       group.EngineVersion,
			"parameter_group_type": group.ParameterGroupType,
			"force_restart":        group.ForceRestart,
			"param_counts":         group.ParamCounts,
			"created":              group.Created,
			"modified":             group.Modified,
			"parameters":           []map[string]interface{}{}, // Initialize empty parameters list
		}

		// Add parameters if available in the response
		if group.ParamDetail != nil {
			parameters := make([]map[string]interface{}, 0)
			for _, param := range group.ParamDetail.ParameterDetail {
				paramMap := map[string]interface{}{
					"param_name":  param.ParamName,
					"param_value": param.ParamValue,
				}
				parameters = append(parameters, paramMap)
			}
			mapping["parameters"] = parameters
		}

		ids = append(ids, group.ParamGroupId)
		names = append(names, group.ParameterGroupName)
		groups = append(groups, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("groups", groups); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}

// Response structures for unmarshaling JSON
type DescribeParameterGroupsResponse struct {
	ParameterGroups ParameterGroups `json:"ParameterGroups"`
	RequestId       string          `json:"RequestId"`
	Success         bool            `json:"success"`
}

type ParameterGroups struct {
	ParameterGroup []ParameterGroupItem `json:"ParameterGroup"`
}

type ParameterGroupItem struct {
	ParameterGroupId   int          `json:"ParameterGroupId"`
	ParamGroupId       string       `json:"ParamGroupId"`
	ParameterGroupName string       `json:"ParameterGroupName"`
	ParameterGroupDesc string       `json:"ParameterGroupDesc"`
	Engine             string       `json:"Engine"`
	EngineVersion      string       `json:"EngineVersion"`
	ParameterGroupType int          `json:"ParameterGroupType"`
	ForceRestart       int          `json:"ForceRestart"`
	ParamCounts        int          `json:"ParamCounts"`
	Created            string       `json:"Created"`
	Modified           string       `json:"Modified"`
	ParamDetail        *ParamDetail `json:"ParamDetail,omitempty"`
}

type ParamDetail struct {
	ParameterDetail []ParameterDetailItem `json:"ParameterDetail"`
}

type ParameterDetailItem struct {
	ParamName  string `json:"ParamName"`
	ParamValue string `json:"ParamValue"`
}
