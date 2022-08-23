package apsarastack

import (
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceApsaraStackElaticsearchZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackElaticsearchZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"zones": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"multi_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func dataSourceApsaraStackElaticsearchZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := elasticsearch.CreateGetRegionConfigurationRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": string(client.RegionId)}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "elasticsearch", "Department": client.Department, "ResourceGroup": client.ResourceGroup}

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.GetRegionConfiguration(request)
	})

	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_elasticsearch_zones", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.GetActionName(), request)
	zones, _ := raw.(*elasticsearch.GetRegionConfigurationResponse)
	for _, zoneID := range zones.Result.Zones {
		if multi && strings.Contains(zoneID, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zoneID)
			continue
		}
		if !multi && !strings.Contains(zoneID, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zoneID)
			continue
		}
	}
	if len(zoneIds) > 0 {
		sort.Strings(zoneIds)
	}

	var s []map[string]interface{}
	if !multi {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{"id": zoneId}
			s = append(s, mapping)
		}
	} else {
		for _, zoneId := range zoneIds {
			mapping := map[string]interface{}{
				"id":             zoneId,
				"multi_zone_ids": splitMultiZoneId(zoneId),
			}
			s = append(s, mapping)
		}
	}
	d.SetId(dataResourceIdHash(zoneIds))
	if err := d.Set("zones", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", zoneIds); err != nil {
		return WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
