package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackExpressconnectRouterInterfacesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectRouterInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_router_interfaces.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectRouterInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_router_interfaces.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackExpressconnectRouterInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_router_interfaces.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackExpressconnectRouterInterfacesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_expressconnect_router_interfaces.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackExpressconnectRouterInterfacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackExpressconnectRouterInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#":    "1",
		"interfaces.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackExpressconnectRouterInterfacesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"interfaces.#": "0",
	}
}

var AlibabacloudstackExpressconnectRouterInterfacesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_expressconnect_router_interfaces.default",
	existMapFunc: existAlibabacloudstackExpressconnectRouterInterfacesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackExpressconnectRouterInterfacesMapFunc,
}

func testAccCheckAlibabacloudstackExpressconnectRouterInterfacesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackExpressconnectRouterInterfaces%d"
}






data "alibabacloudstack_expressconnect_router_interfaces" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
