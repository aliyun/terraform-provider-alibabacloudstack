package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackLindormInstanceTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_lindorm_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceLindormInstanceTypesConfigDependence)

	// id has the form "<engine_type>:<name>", taken from a real spec queried via the anyone data source
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${data.alibabacloudstack_lindorm_instance_types.anyone.instance_types.0.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"lindorm:fake.not.exist.spec"},
		}),
	}

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":       "${data.alibabacloudstack_lindorm_instance_types.anyone.instance_types.0.cpu}",
			"sorted_by": "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":    999,
			"sorted_by": "Memory",
		}),
	}
	
	MemoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"sorted_by": "CPU",
			"memory":    "${data.alibabacloudstack_lindorm_instance_types.anyone.instance_types.0.memory}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"sorted_by": "CPU",
			"memory": 999,
		}),
	}


	engineTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_type": "lindorm",
		}),
	}

	// The specs of the store engine are exactly the values that core_spec can take in alibabacloudstack_lindorm_instance
	storeEngineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_type": "store",
		}),
	}

	// Filter values come from the real specs queried via the anyone data source above, avoiding hardcoded cpu/memory combinations that differ across environments
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":         "${data.alibabacloudstack_lindorm_instance_types.anyone.instance_types.0.cpu}",
			"memory":      "${data.alibabacloudstack_lindorm_instance_types.anyone.instance_types.0.memory}",
			"engine_type": "lindorm",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":         999,
			"memory":      999,
		}),
	}

	var existLindormInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   CHECKSET,
			"ids.0":                   CHECKSET,
			"instance_types.#":        CHECKSET,
			"instance_types.0.id":     CHECKSET,
			"instance_types.0.cpu":    CHECKSET,
			"instance_types.0.memory": CHECKSET,
			"instance_types.0.rate":   CHECKSET,
			"instance_types.0.name":   CHECKSET,
		}
	}

	var fakeLindormInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var lindormInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existLindormInstanceTypesMapFunc,
		fakeMapFunc:  fakeLindormInstanceTypesMapFunc,
	}
	lindormInstanceTypesCheckInfo.dataSourceTestCheck(t, 0, idsConf, cpuConf, MemoryConf, engineTypeConf, storeEngineConf, allConf)
}

func dataSourceLindormInstanceTypesConfigDependence(name string) string {
	// First query the full spec list of the lindorm engine family for the filter steps to use; engine_type is restricted,
	// otherwise instance_types.0 could fall in another engine family, and combining cpu+memory+engine_type could not guarantee a non-empty result
	return `
data "alibabacloudstack_lindorm_instance_types" "anyone" {
  engine_type = "lindorm"
}
`
}
