package alibabacloudstack

import (
	"encoding/json"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackVpcIpv6Isps() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackVpcIpv6IspsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				MinItems: 1,
			},
			"lock_status": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"service_provider": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ipv6_isps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_provider": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_block": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"available_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"in_use_count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"lock_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"need_declare": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"pool_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ula": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackVpcIpv6IspsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var VpcIpv6IspsResponseObj VpcIpv6IspsResponse
	request := client.NewCommonRequest("GET", "YaochiOps", "2022-05-10", "ListVpcIpv6Isps", "")

	bresponse, err := client.ProcessCommonRequest(request)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_vpc_ipv6_isps", "ListVpcIpv6Isps", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &VpcIpv6IspsResponseObj)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.DefaultErrorMsg,
			"alibabacloudstack_vpc_ipv6_isps", "ListVpcIpv6Isps", errmsgs.AlibabacloudStackSdkGoERROR)
	}
	service_provider := d.Get("service_provider").(string)
	lock_status := d.Get("lock_status").(string)

	idsMap := getIdsStringFilter(d)

	var ids []string
	datas := make([]interface{}, 0)
	for _, ipv6_isp := range VpcIpv6IspsResponseObj.VpcIpv6Isps {
		if lock_status != "" && ipv6_isp.LockStatus != lock_status {
			continue
		}

		if service_provider != "" && ipv6_isp.ServiceProvider != service_provider {
			continue
		}

		if len(idsMap) > 0 {
			if _, ok := idsMap[ipv6_isp.PoolId]; !ok {
				continue
			}
		}

		i := map[string]interface{}{

			"id": ipv6_isp.PoolId,

			"service_provider": ipv6_isp.ServiceProvider,

			"zone_id": ipv6_isp.AvailableZoneId,

			"type": ipv6_isp.Type,

			"cidr_block": ipv6_isp.CidrBlock,

			"available_count": ipv6_isp.AvailableCount,

			"in_use_count": ipv6_isp.InUseCount,

			"lock_status": ipv6_isp.LockStatus,

			"need_declare": ipv6_isp.NeedDeclare,

			"pool_id": ipv6_isp.PoolId,

			"ula": ipv6_isp.Ula,
		}

		datas = append(datas, i)

		pool_id := ipv6_isp.PoolId

		ids = append(ids, pool_id)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ipv6_isps", datas); err != nil {
		return err
	}
	if err := d.Set("ids", ids); err != nil {
		return err
	}

	return nil
}

type VpcIpv6IspsResponse struct {
	RequestId   string `json:"RequestId"`
	VpcIpv6Isps []struct {
		ServiceProvider string `json:"ServiceProvider"`
		AvailableZoneId string `json:"AvailableZoneId"`
		Type            string `json:"Type"`
		CidrBlock       string `json:"CidrBlock"`
		AvailableCount  int    `json:"AvailableCount"`
		InUseCount      int    `json:"InUseCount"`
		TotalCount      int    `json:"TotalCount"`
		LockStatus      string `json:"LockStatus"`
		NeedDeclare     int    `json:"NeedDeclare"`
		PoolId          string `json:"PoolId"`
		Ula             bool   `json:"Ula"`
	} `json:"VpcIpv6Isps"`
}
