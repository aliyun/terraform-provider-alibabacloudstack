package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackRdsDbProxiesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_db_proxies.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sRdsDbProxiesDataSource-%d", defaultRegionToTest, rand),
		dataSourceRdsDbProxiesDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_instance_id": "${alibabacloudstack_db_proxy.default.db_instance_id}",
		}),
	}
	var existRdsDbProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                        "1",
			"db_proxies.#":                                 "1",
			"db_proxies.0.db_proxy_instance_type":          CHECKSET,
			"db_proxies.0.db_proxy_instance_num":           CHECKSET,
			"db_proxies.0.db_proxy_connect_string_items.#": "1",
			"db_proxies.0.db_proxy_endpoint_items.#":       "1",
		}
	}

	var fakeRdsDbProxiesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"db_proxies.#": "0",
		}
	}

	var RdsDbProxiesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existRdsDbProxiesMapFunc,
		fakeMapFunc:  fakeRdsDbProxiesMapFunc,
	}

	RdsDbProxiesCheckInfo.dataSourceTestCheck(t, rand, basicConf)
}

func dataSourceRdsDbProxiesDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
  default = "%s"
}

%s

%s

resource "alibabacloudstack_db_proxy" "default" {
	db_instance_id = "${alibabacloudstack_db_instance.default.id}"
	db_proxy_instance_num = "1"
}

 `, name, VSwitchCommonTestCase, RdsMysqlCommonTestCase())
}
