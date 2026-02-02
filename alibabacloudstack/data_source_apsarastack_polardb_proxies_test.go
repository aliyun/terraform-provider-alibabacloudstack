package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbProxiesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_proxies.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sPolardbProxiesDataSource-%d", defaultRegionToTest, rand),
		dataSourcePolardbProxiesDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_polardb_proxy.default.db_instance_id}",
		}),
	}
	var existPolardbProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"db_proxies.#":                                 "1",
			"db_proxies.0.db_proxy_instance_type":          CHECKSET,
			"db_proxies.0.db_proxy_instance_num":           CHECKSET,
			"db_proxies.0.db_proxy_connect_string_items.#": "1",
			"db_proxies.0.db_proxy_endpoint_items.#":       "1",
		}
	}

	var fakePolardbProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"db_proxies.#": "0",
		}
	}

	var PolardbProxiesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbProxiesMapFunc,
		fakeMapFunc:  fakePolardbProxiesMapFunc,
	}

	PolardbProxiesCheckInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourcePolardbProxiesDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%v"
	}

%s

	resource "alibabacloudstack_polardb_proxy" "default" {
		db_instance_id = "${local.polardb_dbinstance_id}"
		db_proxy_instance_num = "1"
	}

 `, name, PolarDBCommonTestCase("MySQL",false))
}
