package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackVSwitchesDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vswitch.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vswitch.default.id}_fake" ]`,
		}),
	}

	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/24"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"cidr_block": `"172.16.0.0/23"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"is_default": `"true"`,
		}),
	}

	vpcIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}_fake"`,
		}),
	}

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"zone_id":    `"${data.apsarastack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"zone_id":    `"${data.apsarastack_zones.default.zones.0.id}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"ids":        `[ "${apsarastack_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
			"zone_id":    `"${data.apsarastack_zones.default.zones.0.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVSwitchesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${apsarastack_vswitch.default.name}"`,
			"ids":        `[ "${apsarastack_vswitch.default.id}" ]`,
			"cidr_block": `"172.16.0.0/24"`,
			"is_default": `"false"`,
			"vpc_id":     `"${apsarastack_vpc.default.id}"`,
			"zone_id":    `"${data.apsarastack_zones.default.zones.0.id}_fake"`,
		}),
	}

	vswitchesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, cidrBlockConf, idDefaultConf, vpcIdConf, zoneIdConf, allConf)

}

func testAccCheckApsaraStackVSwitchesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVSwitchDatasource%d"
}
data "apsarastack_zones" "default" {}

resource "apsarastack_vpc" "default" {
  cidr_block = "172.16.0.0/16"
  name = "${var.name}"
}

resource "apsarastack_vswitch" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/24"
  vpc_id = "${apsarastack_vpc.default.id}"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  
}

data "apsarastack_vswitches" "default" {
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
	resourceId:   "data.apsarastack_vswitches.default",
	existMapFunc: existVSwitchesMapFunc,
	fakeMapFunc:  fakeVSwitchesMapFunc,
}
