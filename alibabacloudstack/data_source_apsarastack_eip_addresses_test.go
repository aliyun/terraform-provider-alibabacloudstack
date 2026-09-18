package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEipsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_eips.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccCheckAlibabacloudstackEipsDataSourceConfig%d", rand),
		dataSourceEipAddressesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_eip.default.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_eip.default.0.id}_fake"},
		}),
	}

	ipsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ip_addresses": []string{"${alibabacloudstack_eip.default.0.ip_address}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ip_addresses": []string{"${alibabacloudstack_eip.default.0.ip_address}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_eip.default.0.id}"},
			"ip_addresses": []string{"${alibabacloudstack_eip.default.0.ip_address}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":          []string{"${alibabacloudstack_eip.default.0.id}"},
			"ip_addresses": []string{"${alibabacloudstack_eip.default.0.ip_address}_fake"},
		}),
	}

	dnsEipsCheckInfo.dataSourceTestCheck(t, rand, idsConf, ipsConf, allConf)
}

func dataSourceEipAddressesConfigDependence(name string) string {
	return fmt.Sprintf(`
resource "alibabacloudstack_eip" "default" {
  name = "%s"
  count = 2
  bandwidth = 5
}
`, name)
}

var existEipsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                "1",
		"names.#":              "1",
		"eips.#":               "1",
		"eips.0.id":            CHECKSET,
		"eips.0.status":        string(Available),
		"eips.0.ip_address":    CHECKSET,
		"eips.0.bandwidth":     "5",
		"eips.0.instance_id":   "",
		"eips.0.instance_type": "",
	}
}

var fakeEipsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"eips.#":  "0",
	}
}

var dnsEipsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_eips.default",
	existMapFunc: existEipsMapFunc,
	fakeMapFunc:  fakeEipsMapFunc,
}
