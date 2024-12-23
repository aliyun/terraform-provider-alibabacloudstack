package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackVpcIpv6GatewaysDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,
		}),
	}

	ipv6_gateway_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,
			"ipv6_gateway_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,
			"ipv6_gateway_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayId}_fake"`,
		}),
	}

	ipv6_gateway_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,
			"ipv6_gateway_name": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,
			"ipv6_gateway_name": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayName}_fake"`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.ResourceGroupId}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.VpcId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_gateways.default.id}"]`,

			"ipv6_gateway_id":   `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayId}"`,
			"ipv6_gateway_name": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayName}"`,
			"resource_group_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.ResourceGroupId}"`,
			"vpc_id":            `"${alibabacloudstack_vpc_ipv6_gateways.default.VpcId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_gateways.default.id}_fake"]`,

			"ipv6_gateway_id":   `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayId}_fake"`,
			"ipv6_gateway_name": `"${alibabacloudstack_vpc_ipv6_gateways.default.Ipv6GatewayName}_fake"`,
			"resource_group_id": `"${alibabacloudstack_vpc_ipv6_gateways.default.ResourceGroupId}_fake"`,
			"vpc_id":            `"${alibabacloudstack_vpc_ipv6_gateways.default.VpcId}_fake"`}),
	}

	AlibabacloudstackVpcIpv6GatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, ipv6_gateway_idConf, ipv6_gateway_nameConf, resource_group_idConf, vpc_idConf, allConf)
}

var existAlibabacloudstackVpcIpv6GatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":    "1",
		"gateways.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcIpv6GatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
	}
}

var AlibabacloudstackVpcIpv6GatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_ipv6_gateways.default",
	existMapFunc: existAlibabacloudstackVpcIpv6GatewaysMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcIpv6GatewaysMapFunc,
}

func testAccCheckAlibabacloudstackVpcIpv6GatewaysSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcIpv6Gateways%d"
}






data "alibabacloudstack_vpc_ipv6_gateways" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
