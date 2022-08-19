package apsarastack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackVpcIpv6EgressRulesDataSource(t *testing.T) {
	resourceId := "data.apsarastack_vpc_ipv6_egress_rules.default"
	rand := acctest.RandIntRange(1000000, 9999999)
	name := fmt.Sprintf("tf-testacc-vpcipv6egressrule-%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceVpcIpv6EgressRulesDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":      "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":      "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}-fake"},
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"status":          "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"status":          "Deleting",
		}),
	}
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"instance_id":     "${apsarastack_vpc_ipv6_egress_rule.default.instance_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":             []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"instance_id":     "${apsarastack_vpc_ipv6_egress_rule.default.instance_id}-fake",
		}),
	}
	ipv6EgressRuleNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":                   []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"ipv6_egress_rule_name": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"ids":                   []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"ipv6_egress_rule_name": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":            "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
			"ids":                   []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"status":                "Available",
			"ipv6_egress_rule_name": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ipv6_gateway_id":       "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_gateway_id}",
			"name_regex":            "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
			"ids":                   []string{"${apsarastack_vpc_ipv6_egress_rule.default.id}"},
			"status":                "Deleting",
			"ipv6_egress_rule_name": "${apsarastack_vpc_ipv6_egress_rule.default.ipv6_egress_rule_name}-fake",
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
	}

	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithEnvVariable(t, "ECS_WITH_IPV6_ADDRESS")
	}

	VpcIpv6EgressRuleCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, nameRegexConf, idsConf, statusConf, instanceIdConf, ipv6EgressRuleNameConf, allConf)
}

func dataSourceVpcIpv6EgressRulesDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

data "apsarastack_instances" "default" {
  name_regex = "no-deleteing-ipv6-address"
  status     = "Running"
}

data "apsarastack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.apsarastack_instances.default.instances.0.id
  status                 = "Available"
}

resource "apsarastack_vpc_ipv6_egress_rule" "default" {
  ipv6_egress_rule_name = var.name
  ipv6_gateway_id       = data.apsarastack_vpc_ipv6_addresses.default.addresses.0.ipv6_gateway_id
  instance_id           = data.apsarastack_vpc_ipv6_addresses.default.ids.0
  instance_type         = "Ipv6Address"
  description           = var.name
}`, name)
}
