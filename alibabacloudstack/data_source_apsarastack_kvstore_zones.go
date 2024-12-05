package alibabacloudstack

import (
	"fmt"
	"sort"
	"strings"

	r_kvstore "github.com/aliyun/alibaba-cloud-sdk-go/services/r-kvstore"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackKVStoreZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackKVStoreZoneRead,

		Schema: map[string]*schema.Schema{
			"multi": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
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

func dataSourceAlibabacloudStackKVStoreZoneRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	multi := d.Get("multi").(bool)
	var zoneIds []string

	request := r_kvstore.CreateDescribeRegionsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	//request.InstanceChargeType = instanceChargeType
	raw, err := client.WithRkvClient(func(rkvClient *r_kvstore.Client) (interface{}, error) {
		return rkvClient.DescribeRegions(request)
	})
	response, ok := raw.(*r_kvstore.DescribeRegionsResponse)
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_kvstore_zones", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	if len(response.RegionIds.KVStoreRegion) <= 0 {
		return errmsgs.WrapError(fmt.Errorf("[ERROR] There is no available zones for KVStore"))
	}
	for _, zone := range response.RegionIds.KVStoreRegion {
		if multi && strings.Contains(zone.ZoneIds, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneIds)
			continue
		}
		if !multi && !strings.Contains(zone.ZoneIds, MULTI_IZ_SYMBOL) {
			zoneIds = append(zoneIds, zone.ZoneIds)
		}
	}

	if len(zoneIds) == 0 && len(response.RegionIds.KVStoreRegion) == 1 {
		zoneIds = append(zoneIds, response.RegionIds.KVStoreRegion[0].ZoneIds)
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
				"id":            zoneId,
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
		writeToFile(output.(string), s)
	}
	return nil
}
