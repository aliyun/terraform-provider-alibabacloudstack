package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testAccCommonBandwidthPackageDataSource%v", rand)
	testAccConfig := dataSourceTestAccConfigFunc(commonBandwidthPackagesCheckInfo.resourceId, name, testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_common_bandwidth_package.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_common_bandwidth_package.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_common_bandwidth_package.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_common_bandwidth_package.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_common_bandwidth_package.default.id}"},
			"name_regex": "${alibabacloudstack_common_bandwidth_package.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_common_bandwidth_package.default.id}_fake"},
			"name_regex": "${alibabacloudstack_common_bandwidth_package.default.name}_fake",
		}),
	}
	commonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}_description"

}
`, name)
}

var existsCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                      "1",
		"names.#":                    "1",
		"packages.#":                 "1",
		"packages.0.id":              CHECKSET,
		"packages.0.isp":             CHECKSET,
		"packages.0.creation_time":   CHECKSET,
		"packages.0.status":          CHECKSET,
		"packages.0.business_status": CHECKSET,
		"packages.0.bandwidth":       "2",
	}
}

var fakeCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":      "0",
		"names.#":    "0",
		"packages.#": "0",
	}
}

var commonBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_common_bandwidth_packages.default",
	existMapFunc: existsCommonBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeCommonBandwidthPackagesMapFunc,
}
