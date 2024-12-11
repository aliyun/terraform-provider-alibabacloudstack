package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackVSwitchesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vswitch.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_vswitch.default.id}_fake" ]`,
		}),
	}

	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/24"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/23"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"is_default": `"true"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}_fake"`,
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"zone_id":    `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"zone_id":    `"${data.alibabacloudstack_zones.default.zones.0.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"ids":        `[ "${alibabacloudstack_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
			"zone_id":    `"${data.alibabacloudstack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_vswitch.default.name}"`,
			"ids":        `[ "${alibabacloudstack_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${alibabacloudstack_vpc.default.id}"`,
			"zone_id":    `"${data.alibabacloudstack_zones.default.zones.0.id}_fake"`,
		}),
	}

	vswitchesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, cidrBlockConf, idDefaultConf, vpcIdConf, zoneIdConf, allConf)

}

func testAccCheckAlibabacloudStackVSwitchesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVSwitchDatasource%d"
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

data "alibabacloudstack_vswitches" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
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
