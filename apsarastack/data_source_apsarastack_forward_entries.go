package apsarastack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceApsaraStackForwardEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceApsaraStackForwardEntriesRead,

		Schema: map[string]*schema.Schema{
			"forward_table_id": {
				Type:     schema.TypeString,
				ForceNew: true,
				Required: true,
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
			},
			"output_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"external_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"internal_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Computed: true,
				ForceNew: true,
			},

			"entries": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"external_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"internal_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
func dataSourceApsaraStackForwardEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.ApsaraStackClient)

	request := vpc.CreateDescribeForwardTableEntriesRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "vpc", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.PageSize = requests.NewInteger(PageSizeLarge)

	request.PageNumber = requests.NewInteger(1)
	request.ForwardTableId = d.Get("forward_table_id").(string)
	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}
	var r *regexp.Regexp
	var err error
	if nameRegex, ok := d.GetOk("name_regex"); ok {
		if r, err = regexp.Compile(nameRegex.(string)); err != nil {
			return WrapError(err)
		}
	}
	var allForwardEntries []vpc.ForwardTableEntry
	invoker := NewInvoker()
	var raw interface{}
	for {
		if err := invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeForwardTableEntries(request)
			})
			return err
		}); err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "apsarastack_forward_entries", request.GetActionName(), ApsaraStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*vpc.DescribeForwardTableEntriesResponse)
		if len(response.ForwardTableEntries.ForwardTableEntry) < 1 {
			break
		}

		for _, entries := range response.ForwardTableEntries.ForwardTableEntry {
			if r != nil && !r.MatchString(entries.ForwardEntryName) {
				continue
			}
			if external_ip, ok := d.GetOk("external_ip"); ok && entries.ExternalIp != external_ip.(string) {
				continue
			}
			if internal_ip, ok := d.GetOk("internal_ip"); ok && entries.InternalIp != internal_ip.(string) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[entries.ForwardEntryId]; !ok {
					continue
				}
			}
			allForwardEntries = append(allForwardEntries, entries)
		}

		if len(response.ForwardTableEntries.ForwardTableEntry) < PageSizeLarge {
			break
		}

		if page, err := getNextpageNumber(request.PageNumber); err != nil {
			return WrapError(err)
		} else {
			request.PageNumber = page
		}
	}

	return ForwardEntriesDecriptionAttributes(d, allForwardEntries, meta)
}

func ForwardEntriesDecriptionAttributes(d *schema.ResourceData, entries []vpc.ForwardTableEntry, meta interface{}) error {
	var ids []string
	var s []map[string]interface{}
	for _, entry := range entries {
		mapping := map[string]interface{}{
			"id":            entry.ForwardEntryId,
			"external_ip":   entry.ExternalIp,
			"internal_ip":   entry.InternalIp,
			"external_port": entry.ExternalPort,
			"internal_port": entry.InternalPort,
			"ip_protocol":   entry.IpProtocol,
			"status":        entry.Status,
		}
		ids = append(ids, entry.ForwardEntryId)
		s = append(s, mapping)
	}
	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("entries", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}

	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil

}
