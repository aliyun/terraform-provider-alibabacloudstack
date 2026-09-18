package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNatGatewaysBandwidthPackagesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_natgateway_bandwidth_packages.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccNatGatewaysBandwidthPackagesDatasource%d", rand),
		dataSourceNatgatewayBandwidthPackagesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.name}_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_natgateway_bandwidth_package.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_natgateway_bandwidth_package.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_natgateway_bandwidth_package.default.name}",
			"description_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.description}",
			"ids":               []string{"${alibabacloudstack_natgateway_bandwidth_package.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_natgateway_bandwidth_package.default.name}",
			"ids":               []string{"${alibabacloudstack_natgateway_bandwidth_package.default.id}"},
			"description_regex": "${alibabacloudstack_natgateway_bandwidth_package.default.description}_fake",
		}),
	}

	NatGatewaysBandwidthPackagesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func dataSourceNatgatewayBandwidthPackagesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
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
