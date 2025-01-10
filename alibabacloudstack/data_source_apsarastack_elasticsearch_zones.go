package alibabacloudstack

import (
	"sort"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/elasticsearch"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackElaticsearchZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackElaticsearchZonesRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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

func dataSourceAlibabacloudStackElaticsearchZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := elasticsearch.CreateGetRegionConfigurationRequest()
	client.InitRoaRequest(*request.RoaRequest)

	raw, err := client.WithElasticsearchClient(func(elasticsearchClient *elasticsearch.Client) (interface{}, error) {
		return elasticsearchClient.GetRegionConfiguration(request)
	})

	response, ok := raw.(*elasticsearch.GetRegionConfigurationResponse)
	addDebug(request.GetActionName(), raw, request.GetActionName(), request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_elasticsearch_zones", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	for _, zoneID := range response.Result.Zones {
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
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", zoneIds); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
