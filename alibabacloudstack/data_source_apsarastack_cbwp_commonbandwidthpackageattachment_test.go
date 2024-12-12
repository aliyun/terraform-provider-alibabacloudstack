package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_package_attachments.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_package_attachments.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_package_attachments.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cbwp_common_bandwidth_package_attachments.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackCbwpCommonBandwidthPackageAttachmentsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#":    "1",
		"attachments.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"attachments.#": "0",
	}
}

var AlibabacloudstackCbwpCommonBandwidthPackageAttachmentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cbwp_common_bandwidth_package_attachments.default",
	existMapFunc: existAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsMapFunc,
}

func testAccCheckAlibabacloudstackCbwpCommonBandwidthPackageAttachmentsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCbwpCommonBandwidthPackageAttachments%d"
}






data "alibabacloudstack_cbwp_common_bandwidth_package_attachments" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
