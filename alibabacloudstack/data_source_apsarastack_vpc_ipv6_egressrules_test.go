package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVPCIpv6EgressRulesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_vpc_ipv6_egress_rules.default"
	rand := getAccTestRandInt(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-vpcipv6egressrule-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6EgressRulesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":      "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":      "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"status":          "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"status":          "Deleting",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"instance_id":     "${alibabacloudstack_vpc_ipv6_egress_rule.default.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"instance_id":     "${alibabacloudstack_vpc_ipv6_egress_rule.default.instance_id}-fake",
		}),
	}
	ipv6EgressRuleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":                   []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"ipv6_egress_rule_name": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":                   []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"ipv6_egress_rule_name": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":            "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
			"ids":                   []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"status":                "Available",
			"ipv6_egress_rule_name": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":            "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
			"ids":                   []string{"${alibabacloudstack_vpc_ipv6_egress_rule.default.id}"},
			"status":                "Deleting",
			"ipv6_egress_rule_name": "${alibabacloudstack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
		}),
	}
	var existVpcIpv6EgressRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"ids.0":                         CHECKSET,
			"names.#":                       "1",
			"names.0":                       fmt.Sprintf("tf-testacc-vpcipv6egressrule-%d", rand),
			"rules.#":                       "1",
			"rules.0.id":                    CHECKSET,
			"rules.0.ipv6_egress_rule_name": fmt.Sprintf("tf-testacc-vpcipv6egressrule-%d", rand),
			"rules.0.description":           fmt.Sprintf("tf-testacc-vpcipv6egressrule-%d", rand),
			"rules.0.status":                "Available",
			"rules.0.ipv6_gateway_id":       CHECKSET,
			"rules.0.instance_type":         "Ipv6Address",
			"rules.0.instance_id":           CHECKSET,
			"rules.0.ipv6_egress_rule_id":   CHECKSET,
		}
	}

	var fakeVpcIpv6EgressRuleMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"rules.#": "0",
		}
	}

	var VpcIpv6EgressRuleCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpcIpv6EgressRuleMapFunc,
		fakeMapFunc:  fakeVpcIpv6EgressRuleMapFunc,
		PreCheck: func() {
			testAccPreCheck(t)
		},
	}

	VpcIpv6EgressRuleCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, statusConf, instanceIdConf, ipv6EgressRuleNameConf, allConf)
}

func dataSourceVpcIpv6EgressRulesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = alibabacloudstack_ecs_instance.default.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "default" {
  vpc_id            = alibabacloudstack_vpc_vpc.default.id
  ipv6_gateway_name = var.name
  description       = var.name
  spec              = "Medium"
}

resource "alibabacloudstack_vpc_ipv6_internet_bandwidth" "default" {
  ipv6_address_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.id
  ipv6_gateway_id      = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  internet_charge_type = "PayByBandwidth"
  bandwidth            = "20"
  depends_on = ["alibabacloudstack_vpc_ipv6_gateway.default"]
}

resource "alibabacloudstack_vpc_ipv6_egress_rule" "default" {
  ipv6_egress_rule_name = var.name
  ipv6_gateway_id       = data.alibabacloudstack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  instance_id           = data.alibabacloudstack_vpc_ipv6_addresses.default.ids.0
  instance_type         = "Ipv6Address"
  description           = var.name
  depends_on = ["alibabacloudstack_vpc_ipv6_internet_bandwidth.default"]
}`, name, ECSInstanceCommonTestCase)
}
