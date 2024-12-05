package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackCbwpCommonBandwidthPackagesDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}_fake"]`,
		}),
	}

	resource_group_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}"]`,
			"resource_group_id": `"${alibabacloudstack_cbwp_common_bandwidth_packages.default.ResourceGroupId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}_fake"]`,
			"resource_group_id": `"${alibabacloudstack_cbwp_common_bandwidth_packages.default.ResourceGroupId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}"]`,

			"resource_group_id": `"${alibabacloudstack_cbwp_common_bandwidth_packages.default.ResourceGroupId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_packages.default.id}_fake"]`,

			"resource_group_id": `"${alibabacloudstack_cbwp_common_bandwidth_packages.default.ResourceGroupId}_fake"`}),
	}

	AlibabacloudstackCbwpCommonBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, idsConf, resource_group_idConf, allConf)
}

var existAlibabacloudstackCbwpCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#":    "1",
		"packages.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCbwpCommonBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"packages.#": "0",
	}
}

var AlibabacloudstackCbwpCommonBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cbwp_common_bandwidth_packages.default",
	existMapFunc: existAlibabacloudstackCbwpCommonBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCbwpCommonBandwidthPackagesMapFunc,
}

func testAccCheckAlibabacloudstackCbwpCommonBandwidthPackagesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCbwpCommonBandwidthPackages%d"
}






data "alibabacloudstack_cbwp_common_bandwidth_packages" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
