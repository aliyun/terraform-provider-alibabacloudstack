package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAdbInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_adb_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccAdbInstanceTypesDataSource_%d", rand),
		dataSourceAdbInstanceTypesConfigDependence)

	// Configuration with engine version and sorting

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_adb_instance_types.anyone.instance_types.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake_id"},
		}),
	}

	statusConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"status": "Available",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"status": "FakeStatus",
		}),
	}

	clusterTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cluster_type": "${data.alibabacloudstack_adb_instance_types.anyone.instance_types.0.cluster_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cluster_type": "fake_cluster_type",
		}),
	}

	cpuTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu_type": "${data.alibabacloudstack_adb_instance_types.anyone.instance_types.0.cpu_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu_type": "fake_cpu_type",
		}),
	}

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":       "${data.alibabacloudstack_adb_instance_types.anyone.instance_types.0.cpu}",
			"sorted_by": "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":       "99999",
			"sorted_by": "Memory",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory":    "${data.alibabacloudstack_adb_instance_types.anyone.instance_types.0.memory}",
			"sorted_by": "CPU",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory":    "9999999",
			"sorted_by": "CPU",
		}),
	}

	var existGpdbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         CHECKSET,
			"ids.0":                         CHECKSET,
			"instance_types.#":              CHECKSET,
			"instance_types.0.id":           CHECKSET,
			"instance_types.0.cpu":          CHECKSET,
			"instance_types.0.memory":       CHECKSET,
			"instance_types.0.storage_min":  CHECKSET,
			"instance_types.0.storage_max":  CHECKSET,
			"instance_types.0.mode":         CHECKSET,
			"instance_types.0.node_min":     CHECKSET,
			"instance_types.0.node_max":     CHECKSET,
			"instance_types.0.status":       CHECKSET,
			"instance_types.0.storage_type": CHECKSET,
		}
	}

	var fakeGpdbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var AdbInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existGpdbInstanceTypesMapFunc,
		fakeMapFunc:  fakeGpdbInstanceTypesMapFunc,
	}

	AdbInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, idsConf, statusConf, clusterTypeConf, cpuTypeConf, cpuConf, memoryConf)
}

func dataSourceAdbInstanceTypesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_adb_instance_types" "anyone" {
		status = "Available"
	}
`
}
