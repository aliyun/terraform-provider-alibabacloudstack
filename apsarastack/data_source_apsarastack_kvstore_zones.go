package apsarastack

import (
	"fmt"
	"sort"
	"strings"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackKVStoreZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackKVStoreZoneRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"instance_charge_type": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      PostPaid,
				ValidateFunc: validation.StringInSlice([]string{"PrePaid", "PostPaid"}, false),
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

func dataSourceApsaraStackKVStoreZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string
	//instanceChargeType := d.Get("instance_charge_type").(string)

	request := r_kvstore.CreateDescribeRegionsRequest()
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "R-kvstore", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	//request.InstanceChargeType = instanceChargeType
	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeRegions(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_kvstore_zones", request.GetActionName(), ApsaraStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	zones, _ := raw.(*r_kvstore.DescribeRegionsResponse)
	if len(zones.RegionIds.KVStoreRegion) <= 0 {
		return WrapError(fmt.Errorf("[ERROR] There is no available zones for KVStore"))
	}
	for _, zone := range zones.RegionIds.KVStoreRegion {
		if multi && strings.Contains(zone.ZoneIds, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneIds)
			continue
		}
		if !multi && !strings.Contains(zone.ZoneIds, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneIds)
		}
	}

	if len(zoneIds) == 0 && len(zones.RegionIds.KVStoreRegion) == 1 {
		zoneIds = append(zoneIds, zones.RegionIds.KVStoreRegion[0].ZoneIds)
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
