package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2CascadeInstances() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2CascadeInstancesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
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
						"instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cascade_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2CascadeInstancesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build ids filter
	idsMap := getIdsStringFilter(d)

	// Build name regex filter
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	// Call List API to get all instances
	action := "ListCascadeInstances"
	request := make(map[string]interface{})
	request["current"] = 1
	request["size"] = 1000
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", action, "/cascadeInstance/listCascadeInstances", nil, nil, request)
	if err != nil {
		return err
	}

	// Parse response data
	records, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return err
	}
	resultInstances := make([]map[string]interface{}, 0)
	names := make([]string, 0)
	ids := make([]string, 0)

	for _, v := range records.([]interface{}) {
		instance := v.(map[string]interface{})
		id := instance["instanceId"].(string)
		name := instance["instanceName"].(string)

		// Apply id filter
		if len(idsMap) > 0 {
			if _, exists := idsMap[id]; !exists {
				continue
			}
		}

		// Apply name regex filter
		if nameRegex != nil && !nameRegex.MatchString(name) {
			continue
		}
		mapping := map[string]interface{}{
			"instance_type":       instance["instanceType"],
			"instance_name":       instance["instanceName"],
			"cascade_instance_id": instance["cascadeInstanceId"],
			"create_time":         instance["createTime"],
			"status":              instance["status"],
		}
		resultInstances = append(resultInstances, mapping)
		names = append(names, name)
		ids = append(ids, id)
	}

	// Set computed fields
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("instances", resultInstances); err != nil {
		return err
	}
	if err := d.Set("names", names); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}
