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
	engineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": CHECKSET,
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccGpdbInstanceTypesDataSource_%d", rand),
		dataSourceGpdbInstanceTypesPresetDependence)

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

	GpdbInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, engineConf)
}

func dataSourceGpdbInstanceTypesConfigDependence(name string) string {
	return ""
}

func dataSourceGpdbInstanceTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_gpdb_instance_types" "preset" {
	}
`
}
