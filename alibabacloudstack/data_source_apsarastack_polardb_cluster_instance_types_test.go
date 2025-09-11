package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPolardbClusterInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_polardb_cluster_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbClusterInstanceTypesDataSource_%d", rand),
		dataSourcePolardbClusterInstanceTypesConfigDependence)

	baseConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"db_type":      "MySQL",
			"db_version":   "5.7",
			"cpu_type":     "intel",
			"sub_category": "normal_exclusive",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccPolardbClusterInstanceTypesDataSource_%d", rand),
		dataSourcePolardbClusterInstanceTypesPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${data.alibabacloudstack_polardb_cluster_instance_types.preset.instance_types.0.id}"},
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

	var existPolardbClusterInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                          CHECKSET,
			"ids.0":                          CHECKSET,
			"instance_types.#":               CHECKSET,
			"instance_types.0.id":            CHECKSET,
			"instance_types.0.cpu":           CHECKSET,
			"instance_types.0.memory":        CHECKSET,
			"instance_types.0.db_type":       CHECKSET,
			"instance_types.0.db_version":    CHECKSET,
			"instance_types.0.cpu_type":      CHECKSET,
			"instance_types.0.db_node_class": CHECKSET,
		}
	}

	var fakePolardbClusterInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var PolardbClusterInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existPolardbClusterInstanceTypesMapFunc,
		fakeMapFunc:  fakePolardbClusterInstanceTypesMapFunc,
	}

	PolardbClusterInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, baseConf, idsConf, cpuConf, memoryConf)
}

func dataSourcePolardbClusterInstanceTypesConfigDependence(name string) string {
	return ""
}

func dataSourcePolardbClusterInstanceTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_polardb_cluster_instance_types" "preset" {
	}
`
}
