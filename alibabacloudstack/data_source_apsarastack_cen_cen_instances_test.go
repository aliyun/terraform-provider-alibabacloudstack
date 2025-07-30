package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenInstancesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_instance.default.cen_instance_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_instance.default.cen_instance_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_cen_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_cen_instance.default.id}_fake" ]`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_instance.default.description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"description_regex": `"${alibabacloudstack_cen_instance.default.description}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_instance.default.cen_instance_name}"`,
			"description_regex": `"${alibabacloudstack_cen_instance.default.description}"`,
			"ids":               `[ "${alibabacloudstack_cen_instance.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand, map[string]string{
			"name_regex":        `"${alibabacloudstack_cen_instance.default.cen_instance_name}"`,
			"ids":               `[ "${alibabacloudstack_cen_instance.default.id}" ]`,
			"description_regex": `"${alibabacloudstack_cen_instance.default.description}_fake"`,
		}),
	}

	CenInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, allConf)
}

func testAccCheckAlibabacloudStackCenInstancesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testAccCenInstancesDatasource%d"
}

resource "alibabacloudstack_cen_instance" "default" {
    cen_instance_name = "${var.name}"
	description = "${var.name}"
}

data "alibabacloudstack_cen_instances" "default" {
	%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
}

var existCenInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#":                   "1",
		"ids.#":                    "1",
		"cens.0.cen_instance_name": CHECKSET,
		"cens.0.description":       CHECKSET,
		"cens.0.status":            CHECKSET,
		"cens.0.id":                CHECKSET,
		"cens.0.create_time":       CHECKSET,
		"cens.0.protection_level":  CHECKSET,
	}
}

var fakeCenInstancesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"cens.#": "0",
		"ids.#":  "0",
	}
}

var CenInstancesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_instances.default",
	existMapFunc: existCenInstancesMapFunc,
	fakeMapFunc:  fakeCenInstancesMapFunc,
}
