package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackInstanceTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceInstanceTypesConfigDependence)

	eniConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"eni_amount": 1,
			"sorted_by":  "Price",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu_core_count": 999,
			"sorted_by":      "Price",
		}),
	}

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu_core_count": "${data.alibabacloudstack_instance_types.anyone.instance_types.0.cpu_core_count}",
			"sorted_by":      "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu_core_count": 999,
			"sorted_by":      "Memory",
		}),
	}

	MemoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory_size": "${data.alibabacloudstack_instance_types.anyone.instance_types.0.memory_size}",
			"sorted_by":   "CPU",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory_size": 999999999,
			"sorted_by":   "CPU",
		}),
	}

	// Test with instance_type_family that doesn't exist
	familyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_type_family": "${data.alibabacloudstack_instance_types.anyone.instance_types.0.family}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_type_family": "ecs.fake",
		}),
	}

	var existInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET, // Should contain at least one instance type
			"instance_types.#":                      CHECKSET, // Should contain at least one instance type
			"instance_types.0.id":                   CHECKSET,
			"instance_types.0.family":               CHECKSET,
			"instance_types.0.eni_amount":           CHECKSET,
			"instance_types.0.availability_zones.#": CHECKSET,
		}
	}

	var fakeInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var instanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existInstanceTypesMapFunc,
		fakeMapFunc:  fakeInstanceTypesMapFunc,
	}
	instanceTypesCheckInfo.dataSourceTestCheck(t, 0, cpuConf, MemoryConf, eniConf, familyConf)
}

func dataSourceInstanceTypesConfigDependence(name string) string {
	return `
data "alibabacloudstack_instance_types" "anyone" {
}
`
}
