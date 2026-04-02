package alibabacloudstack

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func dataSourceAlibabacloudStackBmcpSecurityGroupRules() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAlibabacloudStackBmcpSecurityGroupRulesRead,

		Schema: map[string]*schema.Schema{
			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
			},
			"ip_protocol": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "gre", "all"}, false),
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},
			"rules": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port_range": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cidr_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"source_security_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"priority": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceAlibabacloudStackBmcpSecurityGroupRulesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	sgId := d.Get("security_group_id").(string)
	direction := d.Get("type").(string)

	request := map[string]interface{}{
		"SgId":       sgId,
		"Direction":  direction,
		"PageNumber": 1,
		"PageSize":   100,
	}

	raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroupRule", "", nil, request, nil)
	addDebug("ListSecurityGroupRule", raw, nil, request)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId(dataResourceIdHash([]string{""}))
			d.Set("rules", make([]map[string]interface{}, 0))
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rules", "ListSecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	response := raw
	data := response["data"].(map[string]interface{})
	ruleList := data["data"].([]interface{})

	var rules []map[string]interface{}
	var ids []string

	for _, item := range ruleList {
		ruleMap := item.(map[string]interface{})

		// Filter by type
		if v, ok := d.GetOk("type"); ok && strings.ToLower(ruleMap["direction"].(string)) != v.(string) {
			continue
		}

		// Filter by ip_protocol
		if v, ok := d.GetOk("ip_protocol"); ok && strings.ToLower(ruleMap["ipProtocol"].(string)) != v.(string) {
			continue
		}

		// Filter by policy
		if v, ok := d.GetOk("policy"); ok && strings.ToLower(ruleMap["accessPolicy"].(string)) != v.(string) {
			continue
		}

		mapping := map[string]interface{}{
			"type":                     strings.ToLower(ruleMap["direction"].(string)),
			"ip_protocol":              strings.ToLower(ruleMap["ipProtocol"].(string)),
			"port_range":               ruleMap["portRange"],
			"cidr_ip":                  ruleMap["cidrIp"],
			"source_security_group_id": ruleMap["relatedGroupId"],
			"policy":                   strings.ToLower(ruleMap["accessPolicy"].(string)),
			"description":              ruleMap["description"],
		}

		ids = append(ids, fmt.Sprintf("%s:%s:%s:%s:%s:%s:%s",
			ruleMap["relatedGroupId"], ruleMap["direction"], ruleMap["ipProtocol"], ruleMap["portRange"],
			"intranet", ruleMap["cidrIp"], ruleMap["accessPolicy"]))

		var pri int
		switch v := ruleMap["priority"].(type) {
		case string:
			pri, _ = strconv.Atoi(v)
		case json.Number:
			pri, _ = strconv.Atoi(v.String())
		case float64:
			pri = int(v)
		case int:
			pri = v
		}
		mapping["priority"] = pri
		rules = append(rules, mapping)
	}

	d.SetId(dataResourceIdHash(ids))

	if err := d.Set("rules", rules); err != nil {
		return errmsgs.WrapError(err)
	}

	return nil
}
