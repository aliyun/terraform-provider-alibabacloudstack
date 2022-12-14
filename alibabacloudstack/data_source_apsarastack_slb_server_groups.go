package alibabacloudstack

import (
	"regexp"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackSlbServerGroups() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackSlbServerGroupsRead,

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
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				ForceNew:     true,
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
			"slb_server_groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"servers": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"instance_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"port": {
										Type:     schema.TypeInt,
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
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackSlbServerGroupsRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request := slb.CreateDescribeVServerGroupsRequest()
	if strings.ToLower(client.Config.Protocol) == "https" {
		request.Scheme = "https"
	} else {
		request.Scheme = "http"
	}
	request.RegionId = client.RegionId
	request.Headers = map[string]string{"RegionId": client.RegionId}
	request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
	request.LoadBalancerId = d.Get("load_balancer_id").(string)

	idsMap := make(map[string]string)
	if v, ok := d.GetOk("ids"); ok {
		for _, vv := range v.([]interface{}) {
			idsMap[Trim(vv.(string))] = Trim(vv.(string))
		}
	}

	raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
		return slbClient.DescribeVServerGroups(request)
	})
	if err != nil {
		return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_slb_server_groups", request.GetActionName(), AlibabacloudStackSdkGoERROR)
	}
	addDebug(request.GetActionName(), raw, request.RpcRequest, request)
	response, _ := raw.(*slb.DescribeVServerGroupsResponse)
	var filteredServerGroupsTemp []slb.VServerGroup
	nameRegex, ok := d.GetOk("name_regex")
	if (ok && nameRegex.(string) != "") || (len(idsMap) > 0) {
		var r *regexp.Regexp
		if nameRegex != "" {
			r = regexp.MustCompile(nameRegex.(string))
		}
		for _, serverGroup := range response.VServerGroups.VServerGroup {
			if r != nil && !r.MatchString(serverGroup.VServerGroupName) {
				continue
			}
			if len(idsMap) > 0 {
				if _, ok := idsMap[serverGroup.VServerGroupId]; !ok {
					continue
				}
			}

			filteredServerGroupsTemp = append(filteredServerGroupsTemp, serverGroup)
		}
	} else {
		filteredServerGroupsTemp = response.VServerGroups.VServerGroup
	}

	return slbServerGroupsDescriptionAttributes(d, filteredServerGroupsTemp, meta)
}

func slbServerGroupsDescriptionAttributes(d *schema.ResourceData, serverGroups []slb.VServerGroup, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	var ids []string
	var names []string
	var s []map[string]interface{}

	for _, serverGroup := range serverGroups {
		mapping := map[string]interface{}{
			"id":   serverGroup.VServerGroupId,
			"name": serverGroup.VServerGroupName,
		}

		request := slb.CreateDescribeVServerGroupAttributeRequest()
		if strings.ToLower(client.Config.Protocol) == "https" {
			request.Scheme = "https"
		} else {
			request.Scheme = "http"
		}
		request.Headers = map[string]string{"RegionId": client.RegionId}
		request.QueryParams = map[string]string{"AccessKeySecret": client.SecretKey, "Product": "slb", "Department": client.Department, "ResourceGroup": client.ResourceGroup}
		request.VServerGroupId = serverGroup.VServerGroupId
		raw, err := client.WithSlbClient(func(slbClient *slb.Client) (interface{}, error) {
			return slbClient.DescribeVServerGroupAttribute(request)
		})
		if err != nil {
			return WrapErrorf(err, DataDefaultErrorMsg, "alibabacloudstack_slb_server_groups", request.GetActionName(), AlibabacloudStackSdkGoERROR)
		}
		addDebug(request.GetActionName(), raw, request.RpcRequest, request)
		response, _ := raw.(*slb.DescribeVServerGroupAttributeResponse)
		if response != nil && len(response.BackendServers.BackendServer) > 0 {
			var backendServerMappings []map[string]interface{}
			for _, backendServer := range response.BackendServers.BackendServer {
				backendServerMapping := map[string]interface{}{
					"instance_id": backendServer.ServerId,
					"port":        backendServer.Port,
					"weight":      backendServer.Weight,
				}
				backendServerMappings = append(backendServerMappings, backendServerMapping)
			}
			mapping["servers"] = backendServerMappings
		}

		ids = append(ids, serverGroup.VServerGroupId)
		names = append(names, serverGroup.VServerGroupName)
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("slb_server_groups", s); err != nil {
		return WrapError(err)
	}
	if err := d.Set("ids", ids); err != nil {
		return WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return WrapError(err)
	}

	// create a json file in current directory and write data source to it.
	if output, ok := d.GetOk("output_file"); ok && output.(string) != "" {
		writeToFile(output.(string), s)
	}
	return nil
}
