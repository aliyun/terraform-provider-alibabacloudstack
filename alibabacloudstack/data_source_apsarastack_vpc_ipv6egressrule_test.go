package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackVpcIpv6EgressRulesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
		}),
	}

	instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
			"instance_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
			"instance_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceId}_fake"`,
		}),
	}

	instance_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
			"instance_type": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":           `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
			"instance_type": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceType}_fake"`,
		}),
	}

	ipv6_egress_rule_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
			"ipv6_egress_rule_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
			"ipv6_egress_rule_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleId}_fake"`,
		}),
	}

	ipv6_egress_rule_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
			"ipv6_egress_rule_name": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":                   `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
			"ipv6_egress_rule_name": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleName}_fake"`,
		}),
	}

	ipv6_gateway_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,
			"ipv6_gateway_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6GatewayId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,
			"ipv6_gateway_id": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6GatewayId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}"]`,

			"instance_id":           `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceId}"`,
			"instance_type":         `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceType}"`,
			"ipv6_egress_rule_id":   `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleId}"`,
			"ipv6_egress_rule_name": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleName}"`,
			"ipv6_gateway_id":       `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6GatewayId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_egress_rules.default.id}_fake"]`,

			"instance_id":           `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceId}_fake"`,
			"instance_type":         `"${alibabacloudstack_vpc_ipv6_egress_rules.default.InstanceType}_fake"`,
			"ipv6_egress_rule_id":   `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleId}_fake"`,
			"ipv6_egress_rule_name": `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6EgressRuleName}_fake"`,
			"ipv6_gateway_id":       `"${alibabacloudstack_vpc_ipv6_egress_rules.default.Ipv6GatewayId}_fake"`}),
	}

	AlibabacloudstackVpcIpv6EgressRulesCheckInfo.dataSourceTestCheck(t, rand, idsConf, instance_idConf, instance_typeConf, ipv6_egress_rule_idConf, ipv6_egress_rule_nameConf, ipv6_gateway_idConf, allConf)
}

var existAlibabacloudstackVpcIpv6EgressRulesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":    "1",
		"rules.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcIpv6EgressRulesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var AlibabacloudstackVpcIpv6EgressRulesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_ipv6_egress_rules.default",
	existMapFunc: existAlibabacloudstackVpcIpv6EgressRulesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcIpv6EgressRulesMapFunc,
}

func testAccCheckAlibabacloudstackVpcIpv6EgressRulesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcIpv6EgressRules%d"
}






data "alibabacloudstack_vpc_ipv6_egress_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
