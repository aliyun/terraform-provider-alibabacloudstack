package alibabacloudstack

import (
	"sort"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackSlbZones() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbZonesRead,

		Schema: map[string]*schema.Schema{
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_details": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
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
						"slb_slave_zone_ids": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"local_name": {
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

func dataSourceAlibabacloudStackSlbZonesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slaveZones := make(map[string][]string)
	localName := make(map[string][]string)
	request := slb.CreateDescribeZonesRequest()
	client.InitRpcRequest(*request.RpcRequest)

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeZones(request)
	})
	response, ok := raw.(*slb.DescribeZonesResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_zones", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	for _, zone := range response.Zones.Zone {
		slaveIds := []string{}
		for _, slavezone := range zone.SlaveZones.SlaveZone {
			slaveIds = append(slaveIds, slavezone.ZoneId)
			if len(slaveIds) > 0 {
				sort.Strings(slaveIds)
			}
		}
		localName[zone.ZoneId] = []string{zone.LocalName}
		slaveZones[zone.ZoneId] = slaveIds
	}

	var ids []string
	for v, _ := range slaveZones {
		ids = append(ids, v)
	}
	if len(ids) > 0 {
		sort.Strings(ids)
	}

	var s []map[string]interface{}
	for _, zoneId := range ids {
		mapping := map[string]interface{}{"id": zoneId}
		if len(slaveZones) > 0 {
			mapping["slb_slave_zone_ids"] = slaveZones[zoneId]
		}
		if len(localName) > 0 {
			mapping["local_name"] = localName[zoneId]
		}
		if !d.Get("enable_details").(bool) {
			s = append(s, mapping)
			continue
		}
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("zones", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
