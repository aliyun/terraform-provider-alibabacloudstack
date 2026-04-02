package alibabacloudstack

import (
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceAlibabacloudStackBmcpSecurityGroupRule() *schema.Resource {
	resource := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validation.StringInSlice([]string{"ingress", "egress"}, false),
				Description:  "Type of rule, ingress (inbound) or egress (outbound).",
			},

			"ip_protocol": {
				Type:         schema.TypeString,
				Required:     true,
				ValidateFunc: validation.StringInSlice([]string{"tcp", "udp", "icmp", "gre", "all"}, false),
			},
			"policy": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      GroupRulePolicyAccept,
				ValidateFunc: validation.StringInSlice([]string{"accept", "drop"}, false),
			},

			"port_range": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: bmcpSecurityGroupRulePortRangeDiffSuppressFunc,
			},

			"priority": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 100),
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"cidr_ip": {
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{"cidr_ip", "source_security_group_id"},
				ConflictsWith: []string{"source_security_group_id"},
			},

			"source_security_group_id": {
				Type:          schema.TypeString,
				Optional:      true,
				AtLeastOneOf:  []string{"cidr_ip", "source_security_group_id"},
				ConflictsWith: []string{"cidr_ip"},
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	setResourceFunc(resource, resourceAlibabacloudStackBmcpSecurityGroupRuleCreate, resourceAlibabacloudStackBmcpSecurityGroupRuleRead,
		resourceAlibabacloudStackBmcpSecurityGroupRuleUpdate, resourceAlibabacloudStackBmcpSecurityGroupRuleDelete)
	return resource
}

func resourceAlibabacloudStackBmcpSecurityGroupRuleCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	port := d.Get("port_range").(string)
	if port == "" {
		return errmsgs.WrapError(fmt.Errorf("'port_range': required field is not set or invalid."))
	}

	request, err := buildAlibabacloudStackBmcpSGRuleRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}

	// Call AddSecurityGroupRule API
	raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "AddSecurityGroupRule", "", nil, *request, nil)
	addDebug("AddSecurityGroupRule", raw, nil, *request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rule", "AddSecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	response := raw
	data := response["data"].(map[string]interface{})
	sgrIdList := data["data"].(map[string]interface{})["sgrId"].([]interface{})
	sgrId := sgrIdList[0].(string)

	d.SetId(fmt.Sprintf("%s:%s:%s", d.Get("security_group_id").(string), sgrId, d.Get("type").(string)))

	return resourceAlibabacloudStackBmcpSecurityGroupRuleRead(d, meta)
}

func resourceAlibabacloudStackBmcpSecurityGroupRuleRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid resource id: %s", d.Id())
	}
	direction := parts[2]
	sgId := parts[0]
	sgrId := parts[1]

	request := map[string]interface{}{
		"Direction":  direction,
		"SgId":       sgId,
		"PageNumber": 1,
		"PageSize":   100,
	}

	raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroupRule", "", nil, request, nil)
	addDebug("ListSecurityGroupRule", raw, nil, request)
	if err != nil {
		if errmsgs.NotFoundError(err) {
			d.SetId("")
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rule", "ListSecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	response := raw
	data := response["data"].(map[string]interface{})
	ruleList := data["data"].([]interface{})

	var found bool
	for _, rule := range ruleList {
		ruleMap := rule.(map[string]interface{})
		if ruleMap["sgrId"].(string) == sgrId {
			found = true
			d.Set("type", strings.ToLower(ruleMap["direction"].(string)))
			d.Set("ip_protocol", strings.ToLower(ruleMap["ipProtocol"].(string)))
			d.Set("policy", strings.ToLower(ruleMap["accessPolicy"].(string)))
			d.Set("port_range", ruleMap["portRange"])
			d.Set("description", ruleMap["description"])
			if priority, ok := ruleMap["priority"]; ok {
				d.Set("priority", priority)
			}
			d.Set("security_group_id", ruleMap["SgId"])
			// Handle cidr_ip vs source_security_group_id
			if cidrIp, ok := ruleMap["cidrIp"]; ok && cidrIp.(string) != "" {
				d.Set("cidr_ip", cidrIp)
				d.Set("source_security_group_id", "")
			} else if relatedGroupId, ok := ruleMap["relatedGroupId"]; ok && relatedGroupId.(string) != "" {
				d.Set("cidr_ip", "")
				d.Set("source_security_group_id", relatedGroupId)
			}
			return nil
		}
	}

	if !found {
		d.SetId("")
	}
	return nil
}

func resourceAlibabacloudStackBmcpSecurityGroupRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	request, err := buildAlibabacloudStackBmcpSGRuleRequest(d, meta)
	if err != nil {
		return errmsgs.WrapError(err)
	}
	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid resource id: %s", d.Id())
	}
	sgrId := parts[1]

	// Add SgrId to request for update
	(*request)["SgrId"] = sgrId

	raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ModifySecurityGroupRule", "", nil, *request, nil)
	addDebug("ModifySecurityGroupRule", raw, nil, *request)
	if err != nil {
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rule", "ModifySecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	return resourceAlibabacloudStackBmcpSecurityGroupRuleRead(d, meta)
}

func deleteBmcpSecurityGroupRule(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*connectivity.AlibabacloudStackClient)

	parts := strings.Split(d.Id(), ":")
	if len(parts) != 3 {
		return fmt.Errorf("invalid resource id: %s", d.Id())
	}
	sgrId := parts[1]
	sgId := parts[0]

	request := map[string]interface{}{
		"SgId":  sgId,
		"SgrId": sgrId,
	}

	raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "DeleteSecurityGroupRule", "", nil, request, nil)
	addDebug("DeleteSecurityGroupRule", raw, nil, request)
	if err != nil {
		if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, "InvalidSecurityGroupId.NotFound") {
			return nil
		}
		return errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rule", "DeleteSecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
	}

	return nil
}

func resourceAlibabacloudStackBmcpSecurityGroupRuleDelete(d *schema.ResourceData, meta interface{}) error {
	err := resource.Retry(5*time.Minute, func() *resource.RetryError {
		err := deleteBmcpSecurityGroupRule(d, meta)
		if err != nil {
			if errmsgs.NotFoundError(err) || errmsgs.IsExpectedErrors(err, "InvalidSecurityGroupId.NotFound") {
				return nil
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return errmsgs.WrapError(err)
	}
	return nil
}

func buildAlibabacloudStackBmcpSGRuleRequest(d *schema.ResourceData, meta interface{}) (*map[string]interface{}, error) {
	client := meta.(*connectivity.AlibabacloudStackClient)
	var request map[string]interface{}

	if !d.IsNewResource() {
		parts := strings.Split(d.Id(), ":")
		sgrId := parts[1]
		sgId := parts[0]
		direction := parts[2]

		// First, query the existing security group rule to get current values
		listRequest := map[string]interface{}{
			"Direction":  strings.ToUpper(direction),
			"SgId":       sgId,
			"PageNumber": 1,
			"PageSize":   100,
		}

		raw, err := client.DoTeaRequest("POST", "EasyAI", "2023-11-01", "ListSecurityGroupRule", "", nil, listRequest, nil)
		if err != nil {
			return nil, errmsgs.WrapErrorf(err, errmsgs.RequestV1ErrorMsg, "alibabacloudstack_bcmp_security_group_rule", "ListSecurityGroupRule", errmsgs.AlibabacloudStackSdkGoERROR, "")
		}

		response := raw
		data := response["data"].(map[string]interface{})
		ruleList := data["data"].([]interface{})

		// Find the rule and use its values as defaults
		var foundRule map[string]interface{}
		for _, rule := range ruleList {
			ruleMap := rule.(map[string]interface{})
			if ruleMap["sgrId"].(string) == sgrId {
				foundRule = ruleMap
				break
			}
		}
		if foundRule == nil {
			return nil, fmt.Errorf("security group rule %s not found", sgrId)
		}
		request = map[string]interface{}{
			"Direction":    foundRule["direction"],
			"AccessPolicy": foundRule["accessPolicy"],
			"IpProtocol":   foundRule["ipProtocol"],
			"PortRange":    foundRule["portRange"],
			"Priority":     foundRule["priority"],
			"SgId":         foundRule["SgId"],
			"CidrIp":       foundRule["cidrIp"],
			"Description":  foundRule["description"],
			"IpVersion":    foundRule["ipVersion"],
		}

	} else {
		request = map[string]interface{}{
			"Direction":    strings.ToUpper(d.Get("type").(string)),
			"AccessPolicy": strings.ToUpper(d.Get("policy").(string)),
			"IpProtocol":   strings.ToUpper(d.Get("ip_protocol").(string)),
			"PortRange":    d.Get("port_range").(string),
			"Priority":     d.Get("priority").(int),
			"SgId":         d.Get("security_group_id").(string),
			"Description":  d.Get("description").(string),
		}
		request["Direction"] = strings.ToUpper(d.Get("type").(string))
	}

	// Override with values from d.GetOk if they are set
	if v, ok := d.GetOk("policy"); ok {
		request["AccessPolicy"] = strings.ToUpper(v.(string))
	}
	if v, ok := d.GetOk("ip_protocol"); ok {
		request["IpProtocol"] = strings.ToUpper(v.(string))
		portRange := d.Get("port_range").(string)
		if v.(string) == "tcp" || v.(string) == "udp" {
			if portRange == "-1/-1" {
				return nil, fmt.Errorf("'tcp' and 'udp' can support port range: [1, 65535]. Please correct it and try again.")
			}
		} else if portRange != "-1/-1" {
			return nil, fmt.Errorf("'icmp', 'gre' and 'all' only support port range '-1/-1'. Please correct it and try again.")
		}
	}
	if v, ok := d.GetOk("port_range"); ok {
		request["PortRange"] = v.(string)
	}
	if v, ok := d.GetOk("priority"); ok {
		request["Priority"] = v.(int)
	}
	if v, ok := d.GetOk("description"); ok {
		request["Description"] = v.(string)
	}

	// Handle CIDR IP, or Source Security Group ID
	if v, ok := d.GetOk("cidr_ip"); ok {
		request["CidrIp"] = v.(string)
		request["IpVersion"] = "ipv4"
	} else if v, ok := d.GetOk("source_security_group_id"); ok {
		request["RelatedGroupId"] = v.(string)
		request["IpVersion"] = ""
	}
	return &request, nil
}

func bmcpSecurityGroupRulePortRangeDiffSuppressFunc(k, old, new string, d *schema.ResourceData) bool {
	return old == AllPortRange && new == ""
}
