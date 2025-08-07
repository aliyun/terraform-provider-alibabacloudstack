package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbxInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardbx_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbxInstanceTypesDataSource_%d", rand),
		dataSourcePolardbxInstanceTypesConfigDependence)

	baseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": "5.7",
			"cpu_type":       "intel",
			"series":         "enterprise",
			"sorted_by":      "CPU",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbxInstanceTypesDataSource_%d", rand),
		dataSourcePolardbxInstanceTypesPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${data.alibabacloudstack_polardbx_instance_types.preset.instance_types.0.id}"},
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

	var existPolardbxInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   CHECKSET,
			"ids.0":                   CHECKSET,
			"instance_types.#":        CHECKSET,
			"instance_types.0.id":     CHECKSET,
			"instance_types.0.cpu":    CHECKSET,
			"instance_types.0.memory": CHECKSET,
		}
	}

	var fakePolardbxInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var PolardbxInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbxInstanceTypesMapFunc,
		fakeMapFunc:  fakePolardbxInstanceTypesMapFunc,
	}

	PolardbxInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, baseConf, idsConf, cpuConf, memoryConf)
}

func dataSourcePolardbxInstanceTypesConfigDependence(name string) string {
	return ""
}

func dataSourcePolardbxInstanceTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_polardbx_instance_types" "preset" {
	}
`
}
