package alibabacloudstack

import (
	"regexp"
	"strconv"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackOtsInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackOtsInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of OTS instance IDs.",
			},
			"name_regex": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "A regex string to filter results by instance name.",
			},
			"specification": {
				Type:        schema.TypeString,
				Optional:    true,
				ValidateFunc:  validation.StringInSlice([]string{"SSD", "HYBRID"}, false),
				Description: "Filter instances by instance specification. Valid values: SSD, HYBRID.",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"alias_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"sp_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"specification": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"payment_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"storage_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vcu_quota": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"instance_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag_key": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"tag_value": {
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

func dataSourceAlibabacloudStackOtsInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})

	// Set MaxResults to get all instances in one call
	request["MaxResults"] = 100

	var allInstances []interface{}
	var nextToken string

	// Loop through all pages
	for {
		if nextToken != "" {
			request["NextToken"] = nextToken
		}

		response, err := client.DoTeaRequest("GET", "tablestore", "2020-12-09", "ListInstances", "/v2/openapi/listinstances", nil, nil, request)
		if err != nil {
			return err
		}

		data, ok := response["Instances"]
		if !ok {
			break
		}

		instances, ok := data.([]interface{})
		if !ok {
			break
		}

		allInstances = append(allInstances, instances...)

		// Check for NextToken to continue pagination
		if token, exists := response["NextToken"]; exists && token != nil && token.(string) != "" {
			nextToken = token.(string)
		} else {
			break
		}
	}

	// Create idsMap for filtering if ids are provided
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			if vv == nil {
				continue
			}
			idsMap[vv.(string)] = vv.(string)
		}
	}

	// Create nameRegex for filtering if name_regex is provided
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Filter instances based on ids and name_regex
	filteredInstances := make([]interface{}, 0)
	for _, inst := range allInstances {
		instance, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}

		instanceName := ""
		if name, exists := instance["InstanceName"]; exists && name != nil {
			instanceName = name.(string)
		}

		// Check if instance name matches name_regex if provided
		if nameRegex != nil && !nameRegex.MatchString(instanceName) {
			continue
		}

		// Check if instance id is in ids list if provided
		if len(idsMap) > 0 {
			if _, ok := idsMap[instanceName]; !ok {
				continue
			}
		}

		// Additional filters based on schema fields
		if v, ok := d.GetOk("specification"); ok {
			spec := ""
			if val, exists := instance["InstanceSpecification"]; exists && val != nil {
				spec = val.(string)
			}
			if spec != v.(string) {
				continue
			}
		}

		filteredInstances = append(filteredInstances, instance)
	}

	// Prepare the result data
	ids := make([]string, 0)
	names := make([]string, 0)
	instances := make([]map[string]interface{}, 0)

	for _, inst := range filteredInstances {
		instance, ok := inst.(map[string]interface{})
		if !ok {
			continue
		}

		mapping := make(map[string]interface{})

		// Map the response fields to schema fields
		if v, exists := instance["InstanceName"]; exists {
			mapping["id"] = v.(string)
			mapping["name"] = v.(string)
			ids = append(ids, v.(string))
			names = append(names, v.(string))
		}

		if v, exists := instance["AliasName"]; exists {
			mapping["alias_name"] = v
		}

		if v, exists := instance["RegionId"]; exists {
			mapping["region_id"] = v
		}

		if v, exists := instance["ResourceGroupId"]; exists {
			mapping["resource_group_id"] = v
		}

		if v, exists := instance["SPInstanceId"]; exists {
			mapping["sp_instance_id"] = v
		}

		if v, exists := instance["UserId"]; exists {
			mapping["user_id"] = v
		}

		if v, exists := instance["InstanceDescription"]; exists {
			mapping["description"] = v
		}

		if v, exists := instance["InstanceSpecification"]; exists {
			mapping["specification"] = v
		}

		if v, exists := instance["PaymentType"]; exists {
			mapping["payment_type"] = v
		}

		if v, exists := instance["StorageType"]; exists {
			mapping["storage_type"] = v
		}

		if v, exists := instance["VCUQuota"]; exists {
			if quota, ok := v.(float64); ok {
				mapping["vcu_quota"] = int(quota)
			} else if quota, ok := v.(int); ok {
				mapping["vcu_quota"] = quota
			} else if quota, ok := v.(string); ok {
				if vv, err := strconv.Atoi(quota); err == nil {
					mapping["vcu_quota"] = vv
				}
			}
		}

		if v, exists := instance["InstanceStatus"]; exists {
			mapping["instance_status"] = v
		}

		if v, exists := instance["CreateTime"]; exists {
			mapping["create_time"] = v
		}

		// For network field, we need to get it from detailed instance info using GetInstance API
		// Since this is a data source, we'll set it to empty for now
		mapping["network"] = ""

		// Handle tags - initialize as empty list
		mapping["tags"] = make([]interface{}, 0)

		instances = append(instances, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return err
	}
	if err := d.Set("names", names); err != nil {
		return err
	}
	if err := d.Set("instances", instances); err != nil {
		return err
	}

	return nil
}

