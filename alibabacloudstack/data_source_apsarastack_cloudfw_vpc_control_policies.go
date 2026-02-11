package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/responses"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackCloudfwVpcControlPolicies() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackCloudfwVpcControlPoliciesRead,

		Schema: map[string]*schema.Schema{
			"ids": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of policy IDs to filter results by. If not specified, all policies will be returned.",
			},
			"name_regex": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringIsValidRegExp,
				Description:  "A regex string to filter results by the policy description.",
			},
			"source": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The source address or address group.",
			},
			"destination": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The destination address or address group.",
			},
			"proto": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The protocol type, e.g., 'TCP', 'UDP'.",
			},
			"acl_action": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The action of the access control rule, e.g., 'accept', 'drop', 'log'.",
			},
			// Computed values
			"policies": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the policy.",
						},
						"acl_uuid": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The unique ID of the policy.",
						},
						"order": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The order of the policy in the rule list.",
						},
						"source": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The source address or address group.",
						},
						"source_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the source, e.g., 'net'.",
						},
						"destination": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination address or address group.",
						},
						"destination_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the destination, e.g., 'net'.",
						},
						"proto": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The protocol used, e.g., 'TCP', 'UDP'.",
						},
						"dest_port": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The destination port or port range.",
						},
						"dest_port_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of the destination port, e.g., 'port'.",
						},
						"acl_action": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The action taken on matching traffic: 'accept', 'drop', or 'log'.",
						},
						"application_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the application associated with the rule.",
						},
						"application_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the application associated with the rule.",
						},
						"direction": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The direction of the traffic flow, e.g., 'inout'.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Description of the policy.",
						},
						"hit_times": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "Number of times the policy has been matched.",
						},
						"release": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Whether the policy is released (enabled).",
						},
						"source_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of source CIDR blocks when source type is a group.",
						},
						"destination_group_cidrs": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of destination CIDR blocks when destination type is a group.",
						},
						"dest_port_group_ports": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
							Description: "List of ports when destination port type is a group.",
						},
						"dest_port_group": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Name of the destination port group.",
						},
					},
				},
			},
			"names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of policy descriptions (names) retrieved.",
			},
		},
	}
}

func dataSourceAlibabacloudStackCloudfwVpcControlPoliciesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	action := "DescribeVpcFirewallControlPolicy"

	idsMap := getIdsStringFilter(d)
	var nameRegex *regexp.Regexp
	if v, ok := d.GetOk("name_regex"); ok {
		nameRegex = regexp.MustCompile(v.(string))
	}

	request := client.NewCommonRequest("POST", "cloudfw", "2017-12-07", action, "")
	request.QueryParams["VpcFirewallId"] = ""
	request.QueryParams["CurrentPage"] = "1"
	request.QueryParams["PageSize"] = "100"
	request.QueryParams["yundun"] = "100"

	if v, ok := d.GetOk("destination"); ok {
		request.QueryParams["Destination"] = v.(string)
	}

	if v, ok := d.GetOk("source"); ok {
		request.QueryParams["Source"] = v.(string)
	}

	var bresponse *responses.CommonResponse
	var err error

	bresponse, err = client.ProcessCommonRequest(request)

	if err != nil {
		return errmsgs.WrapError(err)
	}

	var response map[string]interface{}
	err = json.Unmarshal(bresponse.GetHttpContentBytes(), &response)
	if err != nil {
		if bresponse == nil {
			return errmsgs.WrapErrorf(err, "Process Common Request Failed")
		}
		errmsg := errmsgs.GetBaseResponseErrorMessage(bresponse.BaseResponse)
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_cloudfw_vpc_control_policies", request.GetActionName(), errmsgs.AlibabacloudStackSdkGoERROR, errmsg)
	}

	// Check if Policys exists in response
	policiesData, ok := response["Policys"]
	if !ok {
		// If no policies found, return empty result
		d.SetId("")
		return nil
	}

	policiesList := policiesData.([]interface{})

	// Filter policies based on ids and name_regex
	var filteredPolicies []interface{}
	for _, policy := range policiesList {
		policyMap := policy.(map[string]interface{})

		id := fmt.Sprintf("%s:%s", policyMap["AclUuid"], policyMap["Direction"])
		if len(idsMap) > 0 {
			if _, ok := idsMap[id]; !ok {
				continue
			}
		}

		// Check if policy description matches name_regex
		if nameRegex != nil {
			description := policyMap["Description"].(string)
			if !nameRegex.MatchString(description) {
				continue
			}
		}
		if v, ok := d.GetOk("proto"); ok && v.(string) != policyMap["Proto"].(string) {
			continue
		}
		if v, ok := d.GetOk("acl_action"); ok && v.(string) != policyMap["AclAction"].(string) {
			continue
		}

		filteredPolicies = append(filteredPolicies, policy)
	}

	ids := make([]string, 0)
	names := make([]interface{}, 0)
	s := make([]map[string]interface{}, 0)

	for _, policy := range filteredPolicies {
		policyMap := policy.(map[string]interface{})

		mapping := map[string]interface{}{}
		id := fmt.Sprintf("%s:%s", policyMap["AclUuid"], policyMap["Direction"])
		mapping["id"] = id
		mapping["acl_uuid"] = policyMap["AclUuid"]
		mapping["order"] = policyMap["Order"]
		mapping["source"] = policyMap["Source"]
		mapping["source_type"] = policyMap["SourceType"]
		mapping["destination"] = policyMap["Destination"]
		mapping["destination_type"] = policyMap["DestinationType"]
		mapping["proto"] = policyMap["Proto"]
		mapping["dest_port"] = policyMap["DestPort"]
		mapping["dest_port_type"] = policyMap["DestPortType"]
		mapping["acl_action"] = policyMap["AclAction"]
		mapping["application_name"] = policyMap["ApplicationName"]
		mapping["application_id"] = policyMap["ApplicationId"]
		mapping["direction"] = policyMap["Direction"]
		mapping["description"] = policyMap["Description"]
		mapping["hit_times"], _ = strconv.Atoi(fmt.Sprintf("%v", policyMap["HitTimes"]))
		mapping["release"] = policyMap["Release"]

		// Handle array fields
		if sourceGroupCidrs, ok := policyMap["SourceGroupCidrs"].([]interface{}); ok {
			mapping["source_group_cidrs"] = sourceGroupCidrs
		} else {
			mapping["source_group_cidrs"] = []string{}
		}

		if destinationGroupCidrs, ok := policyMap["DestinationGroupCidrs"].([]interface{}); ok {
			mapping["destination_group_cidrs"] = destinationGroupCidrs
		} else {
			mapping["destination_group_cidrs"] = []string{}
		}

		if destPortGroupPorts, ok := policyMap["DestPortGroupPorts"].([]interface{}); ok {
			mapping["dest_port_group_ports"] = destPortGroupPorts
		} else {
			mapping["dest_port_group_ports"] = []string{}
		}

		mapping["dest_port_group"] = policyMap["DestPortGroup"]

		ids = append(ids, policyMap["AclUuid"].(string))
		names = append(names, policyMap["Description"])
		s = append(s, mapping)
	}

	d.SetId(dataResourceIdHash(ids))
	if err := d.Set("ids", ids); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("names", names); err != nil {
		return errmsgs.WrapError(err)
	}
	if err := d.Set("policies", s); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
