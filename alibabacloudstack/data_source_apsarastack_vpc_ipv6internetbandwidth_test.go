package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackVpcIpv6InternetBandwidthsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}_fake"]`,
		}),
	}

	ipv6_internet_bandwidth_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids":                        `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}"]`,
			"ipv6_internet_bandwidth_id": `"${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.Ipv6InternetBandwidthId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids":                        `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}_fake"]`,
			"ipv6_internet_bandwidth_id": `"${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.Ipv6InternetBandwidthId}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}"]`,

			"ipv6_internet_bandwidth_id": `"${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.Ipv6InternetBandwidthId}"`}),
		fakeConfig: testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.id}_fake"]`,

			"ipv6_internet_bandwidth_id": `"${alibabacloudstack_vpc_ipv6_internet_bandwidths.default.Ipv6InternetBandwidthId}_fake"`}),
	}

	AlibabacloudstackVpcIpv6InternetBandwidthsCheckInfo.dataSourceTestCheck(t, rand, idsConf, ipv6_internet_bandwidth_idConf, allConf)
}

var existAlibabacloudstackVpcIpv6InternetBandwidthsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"bandwidths.#":    "1",
		"bandwidths.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackVpcIpv6InternetBandwidthsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"bandwidths.#": "0",
	}
}

var AlibabacloudstackVpcIpv6InternetBandwidthsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vpc_ipv6_internet_bandwidths.default",
	existMapFunc: existAlibabacloudstackVpcIpv6InternetBandwidthsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackVpcIpv6InternetBandwidthsMapFunc,
}

func testAccCheckAlibabacloudstackVpcIpv6InternetBandwidthsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackVpcIpv6InternetBandwidths%d"
}






data "alibabacloudstack_vpc_ipv6_internet_bandwidths" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
