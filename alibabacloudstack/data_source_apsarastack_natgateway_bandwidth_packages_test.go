package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_natgateway_bandwidth_package.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_natgateway_bandwidth_package.default.id}_fake" ]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_natgateway_bandwidth_package.default.name}"`,
			"description_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.description}"`,
			"ids":               `[ "${alibabacloudstack_natgateway_bandwidth_package.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_natgateway_bandwidth_package.default.name}"`,
			"ids":               `[ "${alibabacloudstack_natgateway_bandwidth_package.default.id}" ]`,
			"description_regex": `"${alibabacloudstack_natgateway_bandwidth_package.default.description}_fake"`,
		}),
	}

	NatGatewaysBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccNatGatewaysBandwidthPackagesDatasource%d"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	name = "${var.name}"
}

resource "alibabacloudstack_natgateway_bandwidth_package" "default" {
    name = "${var.name}"
	bandwidth = "5"
    natgateway_id = "${alibabacloudstack_nat_gateway.default.id}"
	description = "${var.name}"
	ip_count = "2"
}

data "alibabacloudstack_natgateway_bandwidth_packages" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existNatGatewaysBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"bandwidth_packages.#":                                     "1",
		"ids.#":                                                    "1",
		"bandwidth_packages.0.name":                                CHECKSET,
		"bandwidth_packages.0.bandwidth":                           CHECKSET,
		"bandwidth_packages.0.bandwidth_package_id":                CHECKSET,
		"bandwidth_packages.0.business_status":                     CHECKSET,
		"bandwidth_packages.0.creation_time":                       CHECKSET,
		"bandwidth_packages.0.description":                         CHECKSET,
		"bandwidth_packages.0.natgateway_id":                       CHECKSET,
		"bandwidth_packages.0.ip_count":                            CHECKSET,
		"bandwidth_packages.0.isp":                                 CHECKSET,
		"bandwidth_packages.0.public_ip_addresses.0.allocation_id": CHECKSET,
		"bandwidth_packages.0.public_ip_addresses.0.ip_address":    CHECKSET,
		"bandwidth_packages.0.status":                              CHECKSET,
		"bandwidth_packages.0.instance_charge_type":                CHECKSET,
		"bandwidth_packages.0.internet_charge_type":                CHECKSET,
	}
}

var fakeNatGatewaysBandwidthPackagesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"bandwidth_packages.#": "0",
		"ids.#":                "0",
	}
}

var NatGatewaysBandwidthPackagesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_natgateway_bandwidth_packages.default",
	existMapFunc: existNatGatewaysBandwidthPackagesMapFunc,
	fakeMapFunc:  fakeNatGatewaysBandwidthPackagesMapFunc,
}
