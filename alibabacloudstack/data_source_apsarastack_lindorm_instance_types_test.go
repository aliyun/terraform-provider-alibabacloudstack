package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackLindormInstanceTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_lindorm_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceLindormInstanceTypesConfigDependence)

	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":    4,
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
			"memory": 8,
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

	// store 引擎的规格即 alibabacloudstack_lindorm_instance 中 core_spec 可取的值
	storeEngineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_type": "store",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":         32,
			"memory":      64,
			"engine_type": "lindorm",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":         999,
			"memory":      999,
		}),
	}

	var existLindormInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
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
			"instance_types.#": "0",
		}
	}

	var lindormInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existLindormInstanceTypesMapFunc,
		fakeMapFunc:  fakeLindormInstanceTypesMapFunc,
	}
	lindormInstanceTypesCheckInfo.dataSourceTestCheck(t, 0, cpuConf, MemoryConf, engineTypeConf, storeEngineConf, allConf)
}

func dataSourceLindormInstanceTypesConfigDependence(name string) string {
	return `
`
}
