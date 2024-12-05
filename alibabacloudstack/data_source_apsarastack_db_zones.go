package alibabacloudstack

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/rds"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDBZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDBZonesRead,

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

func dataSourceAlibabacloudStackDBZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	multi := d.Get("multi").(bool)
	var zoneIds []string
	request := rds.CreateDescribeRegionsRequest()
	client.InitRpcRequest(*request.RpcRequest)

	var response *rds.DescribeRegionsResponse
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		raw, err := client.WithRdsClient(func(rdsClient *rds.Client) (i interface{}, err error) {
			return rdsClient.DescribeRegions(request)
		})
		response, ok := raw.(*rds.DescribeRegionsResponse)
		if err != nil {
			if errmsgs.IsExpectedErrors(err, []string{errmsgs.Throttling}) {
				time.Sleep(time.Duration(3) * time.Second)
				return resource.RetryableError(err)
			}
			
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			err = errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "Process Common Request failed", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
			return resource.RetryableError(err)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		return nil
	})
	if err != nil {
		return err
	}
	if len(response.Regions.RDSRegion) <= 0 {
		return errmsgs.WrapError(fmt.Errorf("[ERROR] There is no available zone for RDS."))
	}
	for _, r := range response.Regions.RDSRegion {
		if multi && strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
			zoneIds = append(zoneIds, r.ZoneId)
			continue
		}
		if !multi && !strings.Contains(r.ZoneId, MULTI_IZ_SYMBOL) && r.RegionId == string(client.Region) {
			zoneIds = append(zoneIds, r.ZoneId)
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
		writeToFile(output.(string), s)
	}
	return nil
}

func splitMultiZoneId(id string) (ids []string) {
	if !(strings.Contains(id, MULTI_IZ_SYMBOL) || strings.Contains(id, "(")) {
		return
	}
	firstIndex := strings.Index(id, MULTI_IZ_SYMBOL)
	secondIndex := strings.Index(id, "(")
	for _, p := range strings.Split(id[secondIndex+1:len(id)-1], COMMA_SEPARATED) {
		ids = append(ids, id[:firstIndex]+string(p))
	}
	return
}
