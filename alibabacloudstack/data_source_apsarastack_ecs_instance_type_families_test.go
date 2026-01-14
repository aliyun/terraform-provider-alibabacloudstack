package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackInstanceTypeFamiliesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_instance_type_families.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccInstanceTypeFamiliesDataSource_%d", rand),
		dataSourceInstanceTypeFamiliesConfigDependence)

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alibabacloudstack_zones.default.zones.0.id}",
		}),
	}

	generationConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"generation":    "${data.alibabacloudstack_instance_type_families.anyone.families.0.generation}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"generation":    "fake-generation",
		}),
	}

	var existInstanceTypeFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                 CHECKSET,
			"ids.0":                 REGEXMATCH + "^ecs.*",
			"families.#":            CHECKSET,
			"families.0.id":         REGEXMATCH + "^ecs.*",
			"families.0.generation": REGEXMATCH + "^ecs-.*",
			"families.0.zone_ids.#": CHECKSET,
			"families.0.zone_ids.0": CHECKSET,
		}
	}

	var fakeInstanceTypeFamiliesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"families.#": CHECKSET,
		}
	}

	var instanceTypeFamiliesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstanceTypeFamiliesMapFunc,
		fakeMapFunc:  fakeInstanceTypeFamiliesMapFunc,
	}

	instanceTypeFamiliesCheckInfo.dataSourceTestCheck(t, rand, zoneIdConf, generationConf,)
}

func dataSourceInstanceTypeFamiliesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_zones" "default" {
	  available_resource_creation = "Instance"
	}
	data "alibabacloudstack_instance_type_families" "anyone" {
	}
`
}
