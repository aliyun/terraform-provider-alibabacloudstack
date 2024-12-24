package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_common_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_common_bandwidth_package.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_common_bandwidth_package.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_common_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alibabacloudstack_common_bandwidth_package.default.id}" ]`,
			"name_regex": `"${alibabacloudstack_common_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${alibabacloudstack_common_bandwidth_package.default.id}_fake" ]`,
			"name_regex": `"${alibabacloudstack_common_bandwidth_package.default.name}_fake"`,
		}),
	}
	commonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackCommonBandwidthPackagesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCommonBandwidthPackageDataSource%d"
}

resource "alibabacloudstack_common_bandwidth_package" "default" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}_description"

}

data "alibabacloudstack_common_bandwidth_packages" "default"  {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
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
