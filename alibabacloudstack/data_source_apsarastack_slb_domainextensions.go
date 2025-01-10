package alibabacloudstack

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackSlbDomainExtensions() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbDomainExtensionsRead,
		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"load_balancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"frontend_port": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"extensions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_certificate_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbDomainExtensionsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateDescribeDomainExtensionsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.LoadBalancerId = d.Get("load_balancer_id").(string)
	request.ListenerPort = requests.NewInteger(d.Get("frontend_port").(int))

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeDomainExtensions(request)
	})
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, ok := raw.(*slb.DescribeDomainExtensionsResponse)
	if err != nil {
		errmsg := ""
		if ok {
			errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb_domain_extensions", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}
	var filteredDomainExtensionsTemp []slb.DomainExtension
	if len(idsMap) > 0 {
		for _, domainExtension := range response.DomainExtensions.DomainExtension {
			if _, ok := idsMap[domainExtension.DomainExtensionId]; !ok {
				continue
			}
			filteredDomainExtensionsTemp = append(filteredDomainExtensionsTemp, domainExtension)
		}
	} else {
		filteredDomainExtensionsTemp = response.DomainExtensions.DomainExtension
	}
	return slbDomainExtensionDescriptionAttributes(d, filteredDomainExtensionsTemp)
}

func slbDomainExtensionDescriptionAttributes(d *schema.ResourceData, domainExtensions []slb.DomainExtension) error {
	var ids []string
	var s []map[string]interface{}
	for _, domainExtension := range domainExtensions {
		mapping := map[string]interface{}{
			"id":                   domainExtension.DomainExtensionId,
			"domain":               domainExtension.Domain,
			"server_certificate_id": domainExtension.ServerCertificateId,
		}
		ids = append(ids, domainExtension.DomainExtensionId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("extensions", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		if err := writeToFile(output.(string), s); err != nil {
			return err
		}
	}
	return nil
}
