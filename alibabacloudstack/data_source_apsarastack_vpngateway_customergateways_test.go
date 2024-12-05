package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackVpnCgwsDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vpn_customer_gateway.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vpn_customer_gateway.default.id}_fake" ]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_customer_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vpn_customer_gateway.default.name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids":        `[ "${alibabacloudstack_vpn_customer_gateway.default.id}" ]`,
			"name_regex": `"${alibabacloudstack_vpn_customer_gateway.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand, map[string]string{
			"ids":        `[ "${alibabacloudstack_vpn_customer_gateway.default.id}" ]`,
			"name_regex": `"${alibabacloudstack_vpn_customer_gateway.default.name}_fake"`,
		}),
	}

	vpnCustomerGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func testAccCheckAlibabacloudStackVpnCustomerGatewaysConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
resource "alibabacloudstack_vpn_customer_gateway" "default" {
	name = "tf-testAccVpnCgwNameDataResource%d"
	ip_address = "40.104.22.228"
	description = "tf-testAccVpnCgwNameDataResource%d"
}

data "alibabacloudstack_vpn_customer_gateways" "default" {
	%s
}
`, rand, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpnCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"gateways.#":             "1",
		"ids.#":                  "1",
		"names.#":                "1",
		"gateways.0.id":          CHECKSET,
		"gateways.0.name":        fmt.Sprintf("tf-testAccVpnCgwNameDataResource%d", rand),
		"gateways.0.ip_address":  "40.104.22.228",
		"gateways.0.description": fmt.Sprintf("tf-testAccVpnCgwNameDataResource%d", rand),
		"gateways.0.create_time": CHECKSET,
	}
}

var fakeVpnCustomerGatewaysMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"gateways.#": "0",
	}
}

var vpnCustomerGatewaysCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpn_customer_gateways.default",
	existMapFunc: existVpnCustomerGatewaysMapFunc,
	fakeMapFunc:  fakeVpnCustomerGatewaysMapFunc,
}
