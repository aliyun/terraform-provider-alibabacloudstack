package alibabacloudstack

import (
	"regexp"
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cr_ee"
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

	var instances []cr_ee.InstancesItem
	for {
		resp, err := crService.ListCrEeInstances(pageNo, pageSize)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		instances = append(instances, resp.Instances...)
		if len(resp.Instances) < pageSize {
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

	var targetInstances []cr_ee.InstancesItem
	for _, instance := range instances {
		if nameRegex != nil && !nameRegex.MatchString(instance.InstanceName) {
			continue
		}

		if idsMap != nil && idsMap[instance.InstanceId] == "" {
			continue
		}

		targetInstances = append(targetInstances, instance)
	}

	instances = targetInstances

	sort.SliceStable(instances, func(i, j int) bool {
		return instances[i].CreateTime < instances[j].CreateTime
	})

	var (
		ids		[]string
		names		[]string
		instanceMaps	[]map[string]interface{}
	)

	for _, instance := range instances {
		usageResp, err := crService.GetCrEeInstanceUsage(instance.InstanceId)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		endpointResp, err := crService.ListCrEeInstanceEndpoint(instance.InstanceId)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		var (
			publicDomains	[]string
			vpcDomains	[]string
		)
		for _, endpoint := range endpointResp.Endpoints {
			if !endpoint.Enable {
				continue
			}
			if endpoint.EndpointType == "internet" {
				for _, d := range endpoint.Domains {
					publicDomains = append(publicDomains, d.Domain)
				}
			} else if endpoint.EndpointType == "vpc" {
				for _, d := range endpoint.Domains {
					vpcDomains = append(vpcDomains, d.Domain)
				}
			}
		}

		mapping := make(map[string]interface{})
		mapping["id"] = instance.InstanceId
		mapping["name"] = instance.InstanceName
		mapping["region"] = instance.RegionId
		mapping["specification"] = instance.InstanceSpecification
		mapping["namespace_quota"] = usageResp.NamespaceQuota
		mapping["namespace_usage"] = usageResp.NamespaceUsage
		mapping["repo_quota"] = usageResp.RepoQuota
		mapping["repo_usage"] = usageResp.RepoUsage
		mapping["vpc_endpoints"] = vpcDomains
		mapping["public_endpoints"] = publicDomains

		ids = append(ids, instance.InstanceId)
		names = append(names, instance.InstanceName)
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
