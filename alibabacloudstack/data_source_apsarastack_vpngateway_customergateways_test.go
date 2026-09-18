package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpngatewayCustomergatewaiesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_vpn_customer_gateways.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccVpnCgwNameDataResource%d", rand),
		dataSourceVpnCustomerGatewaysConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpn_customer_gateway.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vpn_customer_gateway.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpn_customer_gateway.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vpn_customer_gateway.default.name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_vpn_customer_gateway.default.id}"},
			"name_regex": "${alibabacloudstack_vpn_customer_gateway.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_vpn_customer_gateway.default.id}"},
			"name_regex": "${alibabacloudstack_vpn_customer_gateway.default.name}_fake",
		}),
	}

	vpnCustomerGatewaysCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceVpnCustomerGatewaysConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_vpn_customer_gateway" "default" {
	name = "%s"
	ip_address = "40.104.22.228"
	description = "%s"
}
`, name, name)
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
