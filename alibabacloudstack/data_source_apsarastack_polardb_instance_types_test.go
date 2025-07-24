package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbInstanceTypesDataSource_%d", rand),
		dataSourcePolardbInstanceTypesConfigDependence)

	baseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine":         "MySQL",
			"engine_version": "5.7",
			"cpu_type":       "intel",
			"series":         "dual_ha",
			"sorted_by":      "CPU",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbInstanceTypesDataSource_%d", rand),
		dataSourcePolardbInstanceTypesPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${data.alibabacloudstack_polardb_instance_types.preset.instance_types.0.id}"},
			"sorted_by": "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"xxxxx"},
			"sorted_by": "Memory",
		}),
	}
	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu": "4",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu": "800",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory": "8",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory": "3200",
		}),
	}

	var existPolardbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"instance_types.#":                CHECKSET,
			"instance_types.0.id":             CHECKSET,
			"instance_types.0.cpu":            CHECKSET,
			"instance_types.0.memory":         CHECKSET,
			"instance_types.0.engine":         CHECKSET,
			"instance_types.0.engine_version": CHECKSET,
			"instance_types.0.cpu_type":       CHECKSET,
			"instance_types.0.series":         CHECKSET,
			"instance_types.0.connections":    CHECKSET,
			"instance_types.0.storage_type":   CHECKSET,
			"instance_types.0.storage_min":    CHECKSET,
			"instance_types.0.storage_max":    CHECKSET,
		}
	}

	var fakePolardbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var PolardbInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbInstanceTypesMapFunc,
		fakeMapFunc:  fakePolardbInstanceTypesMapFunc,
	}

	PolardbInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, baseConf, idsConf, cpuConf, memoryConf)
}

func dataSourcePolardbInstanceTypesConfigDependence(name string) string {
	return ""
}

func dataSourcePolardbInstanceTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_polardb_instance_types" "preset" {
	}
`
}
