package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackVpnGatewaysDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	preCheck := func() {
		testAccPreCheck(t)
		testAccPreCheckWithAccountSiteType(t, IntlSite)
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vpn_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vpn_gateway.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_gateway.default.name}_fake"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alibabacloudstack_vpn_gateway.default.vpc_id}"`,
		}),

		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"vpc_id": `"${alibabacloudstack_vpn_gateway.default.vpc_id}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"status":     `"Active"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"status":     `"Init"`,
		}),
	}

	businessStatusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"name_regex":      `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alibabacloudstack_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alibabacloudstack_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"Normal"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand, map[string]string{
			"ids":             `[ "${alibabacloudstack_vpn_gateway.default.id}" ]`,
			"name_regex":      `"${alibabacloudstack_vpn_gateway.default.name}"`,
			"vpc_id":          `"${alibabacloudstack_vpn_gateway.default.vpc_id}"`,
			"status":          `"Active"`,
			"business_status": `"FinancialLocked"`,
		}),
	}

	vpnGatewaysCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, vpcIdConf, statusConf, businessStatusConf, allConf)
}

func testAccCheckAlibabacloudStackVpnGatewaysDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAccVpcGatewayConfig%d"
}

resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	vpc_name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
	name = "${var.name}"
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	bandwidth = "10"
	enable_ssl = true
	instance_charge_type = "PostPaid"
	description = "${var.name}"
	vswitch_id = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_vpn_gateways" "default" {
	%s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":                      "1",
		"ids.#":                           "1",
		"names.#":                         "1",
		"gateways.0.id":                   CHECKSET,
		"gateways.0.vpc_id":               CHECKSET,
		"gateways.0.internet_ip":          CHECKSET,
		"gateways.0.create_time":          CHECKSET,
		"gateways.0.end_time":             CHECKSET,
		"gateways.0.name":                 fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.specification":        "10M",
		"gateways.0.description":          fmt.Sprintf("tf-testAccVpcGatewayConfig%d", rand),
		"gateways.0.enable_ssl":           "enable",
		"gateways.0.enable_ipsec":         "enable",
		"gateways.0.status":               "Active",
		"gateways.0.business_status":      "Normal",
		"gateways.0.instance_charge_type": string(PostPaid),
		"gateways.0.ssl_connections":      "5",
	}
}

var fakeVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"gateways.#": "0",
	}
}

var vpnGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpn_gateways.default",
	existMapFunc: existVpnGatewaysMapFunc,
	fakeMapFunc:  fakeVpnGatewaysMapFunc,
}
