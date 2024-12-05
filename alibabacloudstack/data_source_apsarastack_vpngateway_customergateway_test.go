package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackVpngatewayCustomerGatewaysDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_customer_gateways.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_customer_gateways.default.id}_fake"]`,
		}),
	}

	customer_gateway_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpngateway_customer_gateways.default.id}"]`,
			"customer_gateway_id": `"${alibabacloudstack_vpngateway_customer_gateways.default.CustomerGatewayId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids":                 `["${alibabacloudstack_vpngateway_customer_gateways.default.id}_fake"]`,
			"customer_gateway_id": `"${alibabacloudstack_vpngateway_customer_gateways.default.CustomerGatewayId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_customer_gateways.default.id}"]`,

			"customer_gateway_id": `"${alibabacloudstack_vpngateway_customer_gateways.default.CustomerGatewayId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpngateway_customer_gateways.default.id}_fake"]`,

			"customer_gateway_id": `"${alibabacloudstack_vpngateway_customer_gateways.default.CustomerGatewayId}_fake"`}),
	}

	AlibabacloudstackVpngatewayCustomerGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, customer_gateway_idConf, allConf)
}

var existAlibabacloudstackVpngatewayCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":    "1",
		"gateways.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpngatewayCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#": "0",
	}
}

var AlibabacloudstackVpngatewayCustomerGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpngateway_customer_gateways.default",
	existMapFunc: existAlibabacloudstackVpngatewayCustomerGatewaysMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpngatewayCustomerGatewaysMapFunc,
}

func testAccCheckAlibabacloudstackVpngatewayCustomerGatewaysSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpngatewayCustomerGateways%d"
}






data "alibabacloudstack_vpngateway_customer_gateways" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
