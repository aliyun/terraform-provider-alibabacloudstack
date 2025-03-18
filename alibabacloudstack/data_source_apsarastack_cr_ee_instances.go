package alibabacloudstack

import (
	"regexp"
	"sort"
	"strings"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCrEeInstances() *schema.Resource {
	return &schema.Resource{
		Read:	dataSourceAlibabacloudStackCrEeInstancesRead,
		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:		schema.TypeString,
				Optional:	true,
				ValidateFunc:	validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:		schema.TypeString,
				Optional:	true,
			},

			// Computed values
			"ids": {
				Type:		schema.TypeList,
				Optional:	true,
				Computed:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"names": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem:		&schema.Schema{Type: schema.TypeString},
			},
			"instances": {
				Type:		schema.TypeList,
				Computed:	true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"name": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"region": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"specification": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"namespace_quota": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"namespace_usage": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"repo_quota": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"repo_usage": {
							Type:		schema.TypeString,
							Computed:	true,
						},
						"vpc_endpoints": {
							Type:		schema.TypeList,
							Computed:	true,
							Elem:		&schema.Schema{Type: schema.TypeString},
						},
						"public_endpoints": {
							Type:		schema.TypeList,
							Computed:	true,
							Elem:		&schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackCrEeInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	crService := &CrService{client}
	pageNo := 1
	pageSize := 50

	var instances []interface{}
	for {
		resp, err := crService.ListCrEeInstances(pageNo, pageSize)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		respInstanceList := resp["Instances"].([]interface{})
		instances = append(instances, respInstanceList ...)
		if len(respInstanceList) < pageSize {
			break
		}
		pageNo++
	}

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	var idsMap map[string]string
	if v, ok := d.GetOk("ids"); ok {
		idsMap = make(map[string]string)
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}

	var targetInstances []map[string]interface{}
	for _, respInstance := range instances {
		instance := respInstance.(map[string]interface{})
		if nameRegex != nil && !nameRegex.MatchString(instance["InstanceName"].(string)) {
			continue
		}

		if idsMap != nil && idsMap[instance["InstanceId"].(string)] == "" {
			continue
		}

		targetInstances = append(targetInstances, instance)
	}


	sort.SliceStable(instances, func(i, j int) bool {
		return targetInstances[i]["CreateTime"].(float64) < targetInstances[j]["CreateTime"].(float64)
	})

	var (
		ids		[]string
		names		[]string
		instanceMaps	[]map[string]interface{}
	)

	for _, instance := range targetInstances {
		usageResp, err := crService.GetCrEeInstanceUsage(instance["InstanceId"].(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		endpointResp, err := crService.ListCrEeInstanceEndpoint(instance["InstanceId"].(string))
		if err != nil {
			return errmsgs.WrapError(err)
		}

		var (
			publicDomains	[]string
			vpcDomains	[]string
		)
		
		endpoints := endpointResp["Endpoints"].([]interface{})
		for _, endpointItem := range endpoints {
			endpoint := endpointItem.(map[string]interface{})
			if !endpoint["Enable"].(bool) {
				continue
			}
			domains := endpoint["Domains"].([]interface{})
			if endpoint["EndpointType"].(string) == "internet" {
				for _, domainItem := range domains {
					domain := domainItem.(map[string]interface{})
					publicDomains = append(publicDomains, domain["Domain"].(string))
				}
			} else if endpoint["EndpointType"].(string) == "vpc" {
				for _, domainItem := range domains {
					domain := domainItem.(map[string]interface{})
					vpcDomains = append(publicDomains, domain["Domain"].(string))
				}
			}
		}

		mapping := make(map[string]interface{})
		mapping["id"] = instance["InstanceId"].(string)
		mapping["name"] = instance["InstanceName"].(string)
		mapping["region"] = instance["RegionId"].(string)
		mapping["specification"] = instance["InstanceSpecification"].(string)
		mapping["namespace_quota"] = strings.TrimRight(strconv.FormatFloat(usageResp["NamespaceQuota"].(float64), 'f', -1, 64), "0")
		mapping["namespace_usage"] = strings.TrimRight(strconv.FormatFloat(usageResp["NamespaceUsage"].(float64), 'f', -1, 64), "0")
		mapping["repo_quota"] = strings.TrimRight(strconv.FormatFloat(usageResp["RepoQuota"].(float64), 'f', -1, 64), "0")
		mapping["repo_usage"] = strings.TrimRight(strconv.FormatFloat(usageResp["RepoUsage"].(float64), 'f', -1, 64), "0")
		mapping["vpc_endpoints"] = vpcDomains
		mapping["public_endpoints"] = publicDomains

		ids = append(ids, instance["InstanceId"].(string))
		names = append(names, instance["InstanceName"].(string))
		instanceMaps = append(instanceMaps, mapping)
	}

	d.SetId(dataResourceIdHash(names))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("instances", instanceMaps); err != nil {
		return errmsgs.WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok {
		if err := writeToFile(output.(string), instanceMaps); err != nil {
			return err
		}
	}

	return nil
}
