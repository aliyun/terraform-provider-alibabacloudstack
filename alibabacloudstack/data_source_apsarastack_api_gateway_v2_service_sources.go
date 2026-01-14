package alibabacloudstack

import (
	"fmt"
	"log"
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackAPIGatewayV2ServiceSources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackAPIGatewayV2ServiceSourcesRead,

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
			"sources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"update_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_type_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackAPIGatewayV2ServiceSourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Get instance_id from schema
	instanceId := d.Get("instance_id").(string)

	// Prepare request body for ListSources API
	requestBody := map[string]interface{}{
		"gwInstanceId": instanceId,
		"current":      1,
		"size":         100,
	}

	// Call ListSources API
	resp, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListSources", "/source/listSources", nil, nil, requestBody)
	if err != nil {
		return fmt.Errorf("failed to call ListSources API: %v", err)
	}

	// Check if response contains data
	records, err := jsonpath.Get("$.data.records", resp)
	if err != nil {
		d.SetId("")
		return errmsgs.WrapError(err)
	}

	// Build idsMap if ids are provided
	idsMap := getIdsStringFilter(d)

	var sourcesList []map[string]interface{}
	var ids []string

	for _, item := range records.([]interface{}) {
		source := item.(map[string]interface{})
		sourceid := fmt.Sprintf("%s:%s", instanceId, source["sourceId"])
		if description_regex, ok := connectivity.GetResourceDataOk(d, "description_regex", "name_regex"); ok {
			r := regexp.MustCompile(description_regex.(string))
			if !r.MatchString(source["sourceName"].(string)) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if _, exist := idsMap[sourceid]; !exist {
				continue
			}
		}
		mapping := map[string]interface{}{
			"id":               sourceid,
			"instance_id":      instanceId,
			"source_id":        source["sourceId"],
			"source_name":      source["sourceName"],
			"source_type":      source["sourceType"],
			"description":      source["description"],
			"create_time":      source["createTime"],
			"update_time":      source["updateTime"],
			"source_type_name": source["sourceTypeName"],
		}

		sourcesList = append(sourcesList, mapping)
		ids = append(ids, sourceid)
		log.Printf("[DEBUG] alibabacloudstack_api_gateway_v2_service_sources - sources: %s", mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("sources", sourcesList); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
