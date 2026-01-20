package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackGpdbInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_gpdb_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccGpdbInstanceTypesDataSource_%d", rand),
		dataSourceGpdbInstanceTypesConfigDependence)

	// Configuration with engine version and sorting
	
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_gpdb_instance_types.anyone.instance_types.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake_id"},
		}),
	}

	engineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": "${data.alibabacloudstack_gpdb_instance_types.anyone.instance_types.0.engine_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"engine_version": "fake_engine_version",
		}),
	}

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":       "${data.alibabacloudstack_gpdb_instance_types.anyone.instance_types.0.cpu}",
			"sorted_by": "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":       "99999",
			"sorted_by": "Memory",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory":    "${data.alibabacloudstack_gpdb_instance_types.anyone.instance_types.0.memory}",
			"sorted_by": "CPU",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory":    "9999999",
			"sorted_by": "CPU",
		}),
	}

	var existGpdbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             CHECKSET,
			"ids.0":                             CHECKSET,
			"instance_types.#":                  CHECKSET,
			"instance_types.0.id":               CHECKSET,
			"instance_types.0.cpu":              CHECKSET,
			"instance_types.0.memory":           CHECKSET,
			"instance_types.0.engine_version":   CHECKSET,
			"instance_types.0.connections":      CHECKSET,
			"instance_types.0.storage_min":      CHECKSET,
			"instance_types.0.storage_max":      CHECKSET,
			"instance_types.0.specification":    CHECKSET,
			"instance_types.0.db_instance_mode": CHECKSET,
			"instance_types.0.node":             CHECKSET,
			"instance_types.0.region_id":        CHECKSET,
			"instance_types.0.status":           CHECKSET,
			"instance_types.0.product":          CHECKSET,
			"instance_types.0.storage":          CHECKSET,
			"instance_types.0.cpu_label":        CHECKSET,
			"instance_types.0.memory_label":     CHECKSET,
		}
	}

	var fakeGpdbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var GpdbInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existGpdbInstanceTypesMapFunc,
		fakeMapFunc:  fakeGpdbInstanceTypesMapFunc,
	}

	GpdbInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, idsConf, engineVersionConf, cpuConf, memoryConf)
}

func dataSourceGpdbInstanceTypesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_gpdb_instance_types" "anyone" {
	}
`
}
