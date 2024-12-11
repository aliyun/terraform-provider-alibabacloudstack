package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackVpngatewayVpnGatewaysDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,
		}),
	}

	business_statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,
			"business_status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.BusinessStatus}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,
			"business_status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.BusinessStatus}_fake"`,
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,
			"status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.Status}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,
			"status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.Status}_fake"`,
		}),
	}

	vpc_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,
			"vpc_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpcId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":    `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,
			"vpc_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpcId}_fake"`,
		}),
	}

	vpn_instance_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,
			"vpn_instance_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpnInstanceId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids":             `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,
			"vpn_instance_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpnInstanceId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}"]`,

			"business_status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.BusinessStatus}"`,
			"status":          `"${alibabacloudstack_vpngateway_vpn_gateways.default.Status}"`,
			"vpc_id":          `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpcId}"`,
			"vpn_instance_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpnInstanceId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_vpn_gateways.default.id}_fake"]`,

			"business_status": `"${alibabacloudstack_vpngateway_vpn_gateways.default.BusinessStatus}_fake"`,
			"status":          `"${alibabacloudstack_vpngateway_vpn_gateways.default.Status}_fake"`,
			"vpc_id":          `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpcId}_fake"`,
			"vpn_instance_id": `"${alibabacloudstack_vpngateway_vpn_gateways.default.VpnInstanceId}_fake"`}),
	}

	AlibabacloudstackVpngatewayVpnGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, business_statusConf, statusConf, vpc_idConf, vpn_instance_idConf, allConf)
}

var existAlibabacloudstackVpngatewayVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":    "1",
		"gateways.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpngatewayVpnGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
	}
}

var AlibabacloudstackVpngatewayVpnGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpngateway_vpn_gateways.default",
	existMapFunc: existAlibabacloudstackVpngatewayVpnGatewaysMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpngatewayVpnGatewaysMapFunc,
}

func testAccCheckAlibabacloudstackVpngatewayVpnGatewaysSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpngatewayVpnGateways%d"
}






data "alibabacloudstack_vpngateway_vpn_gateways" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
