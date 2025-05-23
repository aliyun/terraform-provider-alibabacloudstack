package alibabacloudstack

import (
	"fmt"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackEcsDeploymentSets() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackEcsDeploymentSetsRead,
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
			"deployment_set_name": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"strategy": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"Availability"}, false),
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sets": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment_set_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"deployment_set_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"granularity": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_amount": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"strategy": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackEcsDeploymentSetsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := make(map[string]interface{})
	request["PageSize"] = PageSizeLarge
	request["PageNumber"] = 1

	if v, ok := d.GetOk("deployment_set_name"); ok {
		request["DeploymentSetName"] = v
	}

	if v, ok := d.GetOk("strategy"); ok {
		request["Strategy"] = v
	}

	var objects []map[string]interface{}
	var deploymentSetNameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		r, err := regexp.Compile(v.(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		deploymentSetNameRegex = r
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

	for {
		response, err := client.DoTeaRequest("POST", "Ecs", "2014-05-26", "DescribeDeploymentSets", "", nil, nil, request)
		if err != nil {
			return err
		}
		resp, err := jsonpath.Get("$.DeploymentSets.DeploymentSet", response)
		if err != nil {
			return errmsgs.WrapErrorf(err, errmsgs.FailedGetAttributeMsg, "DescribeDeploymentSets", "$.DeploymentSets.DeploymentSet", response)
		}
		result, _ := resp.([]interface{})
		for _, v := range result {
			item := v.(map[string]interface{})
			if deploymentSetNameRegex != nil && !deploymentSetNameRegex.MatchString(fmt.Sprint(item["DeploymentSetName"])) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[fmt.Sprint(item["DeploymentSetId"])]; !ok {
					continue
				}
			}
			objects = append(objects, item)
		}
		if len(result) < PageSizeLarge {
			break
		}
		request["PageNumber"] = request["PageNumber"].(int) + 1
	}
	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)
	for _, object := range objects {
		mapping := map[string]interface{}{
			"create_time":         object["CreationTime"],
			"deployment_set_id":   object["DeploymentSetId"],
			"deployment_set_name": object["DeploymentSetName"],
			"description":         object["DeploymentSetDescription"],
			"domain":              convertEcsDeploymentSetDomainResponse(object["Domain"]),
			"granularity":         convertEcsDeploymentSetGranularityResponse(object["Granularity"]),
			"instance_amount":     formatInt(object["InstanceAmount"]),
			"strategy":            object["DeploymentStrategy"],
		}
		if v, ok := object["InstanceIds"].(map[string]interface{}); ok {
			mapping["instance_ids"] = v["InstanceId"]
		}
		ids = append(ids, fmt.Sprint(mapping["DeploymentSetId"]))
		names = append(names, object["DeploymentSetName"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}

	if err := d.Set("sets", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}

	return nil
}
