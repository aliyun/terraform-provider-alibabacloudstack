package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackSlbBackendServers() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbBackendServersRead,

		Schema: map[string]*schema.Schema{
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"backend_servers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"weight": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbBackendServersRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateDescribeLoadBalancerAttributeRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeLoadBalancerAttribute(request)
	})
	response, ok := raw.(*slb.DescribeLoadBalancerAttributeResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_backend_servers", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)

	var filteredBackendServersTemp []slb.BackendServerInDescribeLoadBalancerAttribute
	if len(idsMap) > 0 {
		for _, backendServer := range response.BackendServers.BackendServer {
			if len(idsMap) > 0 {
				if _, ok := idsMap[backendServer.ServerId]; !ok {
					continue
				}
			}

			filteredBackendServersTemp = append(filteredBackendServersTemp, backendServer)
		}
	} else {
		filteredBackendServersTemp = response.BackendServers.BackendServer
	}

	return slbBackendServersDescriptionAttributes(d, filteredBackendServersTemp)
}

func slbBackendServersDescriptionAttributes(d *schema.ResourceData, backendServers []slb.BackendServerInDescribeLoadBalancerAttribute) error {
	var ids []string
	var s []map[string]interface{}

	for _, backendServer := range backendServers {
		mapping := map[string]interface{}{
			"id":     backendServer.ServerId,
			"weight": backendServer.Weight,
		}

		ids = append(ids, backendServer.ServerId)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("backend_servers", s); err != nil {
		return errmsgs.WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
