package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSlbs() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbsRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
				MinItems: 1,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
			},
			"master_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"slave_availability_zone": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"network_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"vswitch_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"address": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
			},
			"output_file": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.",
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			// Computed values
			"slbs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_availability_zone": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"network_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vpc_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"vswitch_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"address": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Computed: true,
						},
						"address_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_zone_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"delete_protection": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internet_charge_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"load_balancer_spec": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	slbService := &SlbService{client}
	request := client.NewCommonRequest("POST", "Slb", "2014-05-15", "DescribeLoadBalancers", "")

	if v, ok := d.GetOk("master_availability_zone"); ok && v.(string) != "" {
		request.QueryParams["MasterZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("slave_availability_zone"); ok && v.(string) != "" {
		request.QueryParams["SlaveZoneId"] = v.(string)
	}
	if v, ok := d.GetOk("network_type"); ok && v.(string) != "" {
		request.QueryParams["NetworkType"] = v.(string)
	}
	if v, ok := d.GetOk("vpc_id"); ok && v.(string) != "" {
		request.QueryParams["VpcId"] = v.(string)
	}
	if v, ok := d.GetOk("vswitch_id"); ok && v.(string) != "" {
		request.QueryParams["VSwitchId"] = v.(string)
	}
	if v, ok := d.GetOk("address"); ok && v.(string) != "" {
		request.QueryParams["Address"] = v.(string)
	}

	if v, ok := d.GetOk("tags"); ok {
		var tags []Tag

		for key, value := range v.(map[string]interface{}) {
			tags = append(tags, Tag{
				Key:   key,
				Value: value.(string),
			})
		}
		request.QueryParams["Tags"] = toSlbTagsString(tags)
	}

	idsMap := getIdsStringFilter(d)

	var allLoadBalancers []LoadBalancerNew
	request.QueryParams["PageSize"] = fmt.Sprintf("%d", PageSizeLarge)
	request.QueryParams["PageNumber"] = fmt.Sprintf("%d", 1)

	for {
		bresponse, err := client.ProcessCommonRequest(request)
		if err != nil {
			if bresponse == nil {
				return errmsgs.WrapErrorf(err, "Process Common Request Failed")
			}
			errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_slb", request.GetActionName(), request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		responseobj := DescribeLoadBalancersNewResponse{}

		err = json.Unmarshal(bresponse.GetHttpContentBytes(), &responseobj)
		if len(responseobj.LoadBalancers.LoadBalancer) < 1 {
			break
		}
		err = json.Unmarshal(bresponse.BaseResponse.GetHttpContentBytes(), bresponse)
		if err != nil {
			return errmsgs.WrapError(err)
		}
		allLoadBalancers = append(allLoadBalancers, responseobj.LoadBalancers.LoadBalancer...)

		if len(responseobj.LoadBalancers.LoadBalancer) < PageSizeLarge {
			break
		}

		page, err := getNextpageNumber(requests.Integer(request.QueryParams["PageNumber"]))
		if err != nil {
			return errmsgs.WrapError(err)
		}
		request.QueryParams["PageNumber"] = string(page)
	}

	var filteredLoadBalancersTemp []LoadBalancerNew

	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, balancer := range allLoadBalancers {
			if r != nil && !r.MatchString(balancer.LoadBalancerName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[balancer.LoadBalancerId]; !ok {
					continue
				}
			}

			filteredLoadBalancersTemp = append(filteredLoadBalancersTemp, balancer)
		}
	} else {
		filteredLoadBalancersTemp = allLoadBalancers
	}

	return slbsDescriptionAttributes(d, filteredLoadBalancersTemp, slbService)
}

func slbsDescriptionAttributes(d *schema.ResourceData, loadBalancers []LoadBalancerNew, slbService *SlbService) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, loadBalancer := range loadBalancers {
		tags, _ := slbService.DescribeTags(loadBalancer.LoadBalancerId, nil, TagResourceInstance)
		mapping := map[string]interface{}{
			"id":                       loadBalancer.LoadBalancerId,
			"region_id":                loadBalancer.RegionId,
			"master_availability_zone": loadBalancer.MasterZoneId,
			"slave_availability_zone":  loadBalancer.SlaveZoneId,
			"name":                     loadBalancer.LoadBalancerName,
			"network_type":             loadBalancer.NetworkType,
			"vpc_id":                   loadBalancer.VpcId,
			"vswitch_id":               loadBalancer.VSwitchId,
			"address":                  loadBalancer.Address,
			"tags":                     slbService.tagsToMap(tags),
			"address_type":             loadBalancer.AddressType,
			"delete_protection":        loadBalancer.DeleteProtection,
			"internet_charge_type":     loadBalancer.InternetChargeType,
			"create_time":              loadBalancer.CreateTime,
			"load_balancer_spec":       loadBalancer.LoadBalancerSpec,
		}

		ids = append(ids, loadBalancer.LoadBalancerId)
		names = append(names, loadBalancer.LoadBalancerName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slbs", s); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
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
