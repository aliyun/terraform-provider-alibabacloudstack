package alibabacloudstack

import (
	"regexp"

	"github.com/PaesslerAG/jsonpath"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackApiGatewayV2CascadeLinks() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackApiGatewayV2CascadeLinksRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"source_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cascade_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_instance_address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cascade_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"link_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cascade_service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cascade_instance_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackApiGatewayV2CascadeLinksRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	// Build request parameters
	request := make(map[string]interface{})
	if v, ok := d.GetOk("source_instance_name"); ok {
		request["sourceInstanceName"] = v
	}
	if v, ok := d.GetOk("cascade_instance_name"); ok {
		request["cascadeInstanceName"] = v
	}
	request["current"] = 1
	request["size"] = 1000

	// Call List API
	response, err := client.DoTeaRequest("POST", "csb2", "2023-02-06", "ListCascadeLink", "/cascadeLink/listCascadeLink", nil, nil, request)
	if err != nil {
		return err
	}
	records, err := jsonpath.Get("$.data.records", response)
	if err != nil {
		return err
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

	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	links := make([]map[string]interface{}, 0)
	ids := make([]string, 0)

	for _, v := range records.([]interface{}) {
		item := v.(map[string]interface{})
		linkId := item["linkId"].(string)
		if nameRegex != nil {
			if linkName, ok := item["linkName"].(string); !ok || !nameRegex.MatchString(linkName) {
				continue
			}
		}

		if len(idsMap) > 0 {
			if linkId, ok := item["linkId"].(string); !ok || idsMap[linkId] == "" {
				continue
			}
		}

		mapping := map[string]interface{}{
			"id":                      linkId,
			"link_id":                 linkId,
			"source_instance_id":      item["sourceInstanceId"],
			"source_instance_address": item["sourceInstanceAddress"],
			"cascade_instance_id":     item["cascadeInstanceId"],
			"link_name":               item["linkName"],
			"cascade_service_id":      item["cascadeServiceId"],
			"source_instance_name":    item["sourceInstanceName"],
			"cascade_instance_name":   item["cascadeInstanceName"],
		}
		ids = append(ids, linkId)
		links = append(links, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("links", links); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
