package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAlibabacloudStackDnsGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackDnsGroupsRead,

		Schema: map[string]*schema.Schema{
			"name_regex": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},

			// Computed values
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackDnsGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := alidns.CreateDescribeDomainGroupsRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)
	idsMap := make(map[string]string)
	var allGroups []alidns.DomainGroup
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[vv.(string)] = vv.(string)
		}
	}
	for {
		raw, err := client.WithDnsClient(func(dnsClient *alidns.Client) (interface{}, error) {
			return dnsClient.DescribeDomainGroups(request)
		})
		response, ok := raw.(*alidns.DescribeDomainGroupsResponse)
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_dns_groups", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		groups := response.DomainGroups.DomainGroup
		for _, domainGroup := range groups {
			if len(idsMap) > 0 {
				if _, ok := idsMap[domainGroup.GroupId]; !ok {
					continue
				}
			}
			allGroups = append(allGroups, domainGroup)
		}
		if len(groups) < PageSizeLarge {
			break
		}
		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return errmsgs.WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	var filteredGroups []alidns.DomainGroup
	if v, ok := d.GetOk("name_regex"); ok && v.(string) != "" {
		r := regexp.MustCompile(v.(string))

		for _, group := range allGroups {
			if r.MatchString(group.GroupName) {
				filteredGroups = append(filteredGroups, group)
			}
		}
	} else {
		filteredGroups = allGroups[:]
	}

	return groupsDecriptionAttributes(d, filteredGroups, meta)
}

func groupsDecriptionAttributes(d *schema.ResourceData, groupTypes []alidns.DomainGroup, meta interface{}) error {
	var ids []string
	var names []string
	var s []map[string]interface{}
	for _, group := range groupTypes {
		mapping := map[string]interface{}{
			"group_id":   group.GroupId,
			"group_name": group.GroupName,
		}
		ids = append(ids, group.GroupId)
		names = append(names, group.GroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("groups", s); err != nil {
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
		writeToFile(output.(string), s)
	}
	return nil
}
