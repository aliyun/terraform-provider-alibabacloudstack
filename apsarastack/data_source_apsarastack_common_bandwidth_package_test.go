package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackCommonBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${apsarastack_common_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"name_regex": `"${apsarastack_common_bandwidth_package.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${apsarastack_common_bandwidth_package.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids": `[ "${apsarastack_common_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${apsarastack_common_bandwidth_package.default.id}" ]`,
			"name_regex": `"${apsarastack_common_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand, map[string]string{
			"ids":        `[ "${apsarastack_common_bandwidth_package.default.id}_fake" ]`,
			"name_regex": `"${apsarastack_common_bandwidth_package.default.name}_fake"`,
		}),
	}
	commonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckApsaraStackCommonBandwidthPackagesDataSourceConfigBasic(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCommonBandwidthPackageDataSource%d"
}

resource "apsarastack_common_bandwidth_package" "default" {
  bandwidth = "2"
  name = "${var.name}"
  description = "${var.name}_description"

}

data "apsarastack_common_bandwidth_packages" "default"  {
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
	resourceId:   "data.apsarastack_common_bandwidth_packages.default",
	existMapFunc: existsCommonBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeCommonBandwidthPackagesMapFunc,
}
