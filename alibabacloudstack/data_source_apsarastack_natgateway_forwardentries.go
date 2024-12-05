package alibabacloudstack

import (
	"regexp"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackForwardEntries() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackForwardEntriesRead,

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

func dataSourceAlibabacloudStackForwardEntriesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := vpc.CreateDescribeForwardTableEntriesRequest()
	client.InitRpcRequest(*request.RpcRequest)
	request.ForwardTableId = d.Get("forward_table_id").(string)
	request.PageSize = requests.NewInteger(PageSizeLarge)
	request.PageNumber = requests.NewInteger(1)

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
			return errmsgs.WrapError(err)
		}
	}

	var allForwardEntries []vpc.ForwardTableEntry
	invoker := NewInvoker()
	var raw interface{}
	for {
		invoker.Run(func() error {
			raw, err = client.WithVpcClient(func(vpcClient *vpc.Client) (interface{}, error) {
				return vpcClient.DescribeForwardTableEntries(request)
			})
			return err
		})
		response, ok := raw.(*vpc.DescribeForwardTableEntriesResponse)
		if err != nil {
			errmsg := ""
			if ok {
				errmsg = errmsgs.GetBaseResponseErrorMessage(response.BaseResponse)
			}
			return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_forward_entries", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		if !ok {
			return errmsgs.WrapErrorf(err, "Failed to parse response as *vpc.DescribeForwardTableEntriesResponse")
		}
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
			return errmsgs.WrapError(err)
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
