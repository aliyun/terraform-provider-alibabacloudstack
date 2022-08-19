package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"os"
	"strings"
	"testing"
)

/*func TestAccApsaraStackVpcsDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	initVswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d"`, rand),
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d_fake"`, rand),
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vpc.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vpc.default.id}_fake" ]`,
		}),
	}
	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/12"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/0"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Pending"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"true"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}_fake"`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"ids":        `[ "${apsarastack_vpc.default.id}" ]`,
			"cidr_block": `"172.16.0.0/12"`,
			"status":     `"Available"`,
			"is_default": `"false"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"ids":        `[ "${apsarastack_vpc.default.id}" ]`,
			"cidr_block": `"172.16.0.0/16"`,
			"status":     `"Available"`,
			"is_default": `"false"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}_fake"`,
		}),
	}

	vpcsCheckInfo.dataSourceTestCheck(t, rand, initVswitchConf, nameRegexConf, idsConf, cidrBlockConf, statusConf, idDefaultConf, vswitchIdConf, resourceGroupIdConf, allConf)
}

func testAccCheckApsaraStackVpcsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpcsdatasource%d"
}

resource "apsarastack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "172.16.0.0/12"

}

data "apsarastack_zones" "default" {

}

resource "apsarastack_vswitch" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${apsarastack_vpc.default.id}"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}

data "apsarastack_vpcs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"names.#":               "1",
		"vpcs.#":                "1",
		"vpcs.0.id":             CHECKSET,
		"vpcs.0.region_id":      CHECKSET,
		"vpcs.0.status":         "Available",
		"vpcs.0.vpc_name":       fmt.Sprintf("tf-testAccVpcsdatasource%d", rand),
		"vpcs.0.vswitch_ids.#":  "1",
		"vpcs.0.cidr_block":     "172.16.0.0/12",
		"vpcs.0.vrouter_id":     CHECKSET,
		"vpcs.0.route_table_id": CHECKSET,
		"vpcs.0.description":    "",
		"vpcs.0.is_default":     "false",
		"vpcs.0.creation_time":  CHECKSET,
	}
}

var fakeVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpcs.#":  "0",
	}
}

var vpcsCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_vpcs.default",
	existMapFunc: existVpcsMapFunc,
	fakeMapFunc:  fakeVpcsMapFunc,
}*/

func TestAccApsaraStackVpcsDataSourceBasic(t *testing.T) {
	rand := acctest.RandInt()
	initVswitchConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d"`, rand),
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": fmt.Sprintf(`"tf-testAccVpcsdatasource%d_fake"`, rand),
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vpc.default.id}" ]`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${apsarastack_vpc.default.id}_fake" ]`,
		}),
	}
	cidrBlockConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/12"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"cidr_block": `"172.16.0.0/0"`,
		}),
	}
	statusConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Available"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"status":     `"Pending"`,
		}),
	}
	idDefaultConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"false"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"is_default": `"true"`,
		}),
	}
	vswitchIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}"`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"vswitch_id": `"${apsarastack_vswitch.default.id}_fake"`,
		}),
	}
	tagsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"tags": `{
							Created = "TF"
							For 	= "acceptance test"
					  }`,
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${var.name}"`,
			"tags": `{
							Created = "TF-fake"
							For 	= "acceptance test-fake"
					  }`,
		}),
	}
	resourceGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"ids":               `[ "${apsarastack_vpc.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/12"`,
			"status":            `"Available"`,
			"is_default":        `"false"`,
			"vswitch_id":        `"${apsarastack_vswitch.default.id}"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
		fakeConfig: testAccCheckApsaraStackVpcsDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${var.name}"`,
			"ids":               `[ "${apsarastack_vpc.default.id}" ]`,
			"cidr_block":        `"172.16.0.0/16"`,
			"status":            `"Available"`,
			"is_default":        `"false"`,
			"vswitch_id":        `"${apsarastack_vswitch.default.id}_fake"`,
			"resource_group_id": fmt.Sprintf(`"%s"`, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID")),
		}),
	}

	vpcsCheckInfo.dataSourceTestCheck(t, rand, initVswitchConf, nameRegexConf, idsConf, cidrBlockConf, statusConf, idDefaultConf, vswitchIdConf, tagsConf, resourceGroupIdConf, allConf)
}

func testAccCheckApsaraStackVpcsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccVpcsdatasource%d"
}

resource "apsarastack_vpc" "default" {
  vpc_name = "${var.name}"
  cidr_block = "172.16.0.0/12"
  tags 		= {
		Created = "TF"
		For 	= "acceptance test"
  }
  resource_group_id = "%s"
}

data "apsarastack_zones" "default" {

}

resource "apsarastack_vswitch" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/16"
	vpc_id = "${apsarastack_vpc.default.id}"
	availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
}

data "apsarastack_vpcs" "default" {
	enable_details = true
  %s
}
`, rand, os.Getenv("ALICLOUD_RESOURCE_GROUP_ID"), strings.Join(pairs, "\n  "))
	return config
}

var existVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                 "1",
		"names.#":               "1",
		"vpcs.#":                "1",
		"vpcs.0.id":             CHECKSET,
		"vpcs.0.region_id":      CHECKSET,
		"vpcs.0.status":         "Available",
		"vpcs.0.vpc_name":       fmt.Sprintf("tf-testAccVpcsdatasource%d", rand),
		"vpcs.0.vswitch_ids.#":  "1",
		"vpcs.0.cidr_block":     "172.16.0.0/12",
		"vpcs.0.vrouter_id":     CHECKSET,
		"vpcs.0.router_id":      CHECKSET,
		"vpcs.0.route_table_id": CHECKSET,
		"vpcs.0.description":    "",
		"vpcs.0.is_default":     "false",
		"vpcs.0.creation_time":  CHECKSET,
	}
}

var fakeVpcsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"vpcs.#":  "0",
	}
}

var vpcsCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_vpcs.default",
	existMapFunc: existVpcsMapFunc,
	fakeMapFunc:  fakeVpcsMapFunc,
}
