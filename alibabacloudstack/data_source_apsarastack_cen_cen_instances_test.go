package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCenInstancesDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_instances.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccCenInstancesDatasource%d", rand),
		dataSourceCenInstancesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cen_instance.default.cen_instance_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cen_instance.default.cen_instance_name}_fake",
		}),
	}
	transit_router_nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_name_regex": "${alibabacloudstack_cen_instance.default.transit_router_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_name_regex": "${alibabacloudstack_cen_instance.default.transit_router_name}_fake",
		}),
	}
	transit_router_descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"transit_router_description_regex": "${alibabacloudstack_cen_instance.default.transit_router_description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"transit_router_description_regex": "${alibabacloudstack_cen_instance.default.transit_router_description}_fake",
		}),
	}
	cirConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cidr": "10.10.10.2/24",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cidr": "10.10.10.2/24_fake",
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cen_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cen_instance.default.id}_fake"},
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_cen_instance.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description_regex": "${alibabacloudstack_cen_instance.default.description}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_cen_instance.default.cen_instance_name}",
			"description_regex": "${alibabacloudstack_cen_instance.default.description}",
			"ids":               []string{"${alibabacloudstack_cen_instance.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex":        "${alibabacloudstack_cen_instance.default.cen_instance_name}",
			"ids":               []string{"${alibabacloudstack_cen_instance.default.id}"},
			"description_regex": "${alibabacloudstack_cen_instance.default.description}_fake",
		}),
	}

	CenInstancesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, descriptionConf, transit_router_nameRegexConf, transit_router_descriptionRegexConf, cirConf, allConf)
}

func dataSourceCenInstancesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_cen_instance" "default" {
    cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name="${var.name}"
	transit_router_description="${var.name}"
	transit_router_cidrs {
		cidr="10.10.10.2/24"
	} 
	transit_router_cidrs {
		cidr="10.10.11.2/24"
	}
}
`, name)
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
