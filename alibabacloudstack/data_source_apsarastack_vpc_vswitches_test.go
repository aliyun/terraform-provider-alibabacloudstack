package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpcVswitchesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_vswitches.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccVSwitchDatasource%d", rand),
		dataSourceVpcVswitchesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vswitch.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_vswitch.default.id}_fake"},
		}),
	}

	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"cidr_block": "172.16.0.0/24",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"cidr_block": "172.16.0.0/23",
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"is_default": "false",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"is_default": "true",
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}_fake",
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"zone_id":    "${data.alibabacloudstack_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"zone_id":    "${data.alibabacloudstack_zones.default.zones.0.id}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"ids":        []string{"${alibabacloudstack_vswitch.default.id}"},
			"cidr_block": "172.16.0.0/24",
			"is_default": "false",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
			"zone_id":    "${data.alibabacloudstack_zones.default.zones.0.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_vswitch.default.name}",
			"ids":        []string{"${alibabacloudstack_vswitch.default.id}"},
			"cidr_block": "172.16.0.0/24",
			"is_default": "false",
			"vpc_id":     "${alibabacloudstack_vpc.default.id}",
			"zone_id":    "${data.alibabacloudstack_zones.default.zones.0.id}_fake",
		}),
	}

	vswitchesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, cidrBlockConf, idDefaultConf, vpcIdConf, zoneIdConf, allConf)
}

func dataSourceVpcVswitchesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}
data "alibabacloudstack_zones" "default" {}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "alibabacloudstack_vswitch" "default" {
  vswitch_name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
`, name)
}

var existVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                                  "1",
		"names.#":                                "1",
		"vswitches.#":                            "1",
		"vswitches.0.id":                         CHECKSET,
		"vswitches.0.vpc_id":                     CHECKSET,
		"vswitches.0.zone_id":                    CHECKSET,
		"vswitches.0.name":                       fmt.Sprintf("tf-testAccVSwitchDatasource%d", rand),
		"vswitches.0.instance_ids.#":             "0",
		"vswitches.0.cidr_block":                 "172.16.0.0/24",
		"vswitches.0.description":                "",
		"vswitches.0.is_default":                 "false",
		"vswitches.0.creation_time":              CHECKSET,
		"vswitches.0.available_ip_address_count": CHECKSET,
	}
}

var fakeVSwitchesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":       "0",
		"names.#":     "0",
		"vswitches.#": "0",
	}
}

var vswitchesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_vswitches.default",
	existMapFunc: existVSwitchesMapFunc,
	fakeMapFunc:  fakeVSwitchesMapFunc,
}
