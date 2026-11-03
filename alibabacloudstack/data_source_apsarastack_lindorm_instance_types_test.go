package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackLindormInstanceTypesDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_lindorm_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceLindormInstanceTypesConfigDependence)

	// id 形如 "<engine_type>:<name>"，取自 anyone 数据源查出的真实规格
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

	// store 引擎的规格即 alibabacloudstack_lindorm_instance 中 core_spec 可取的值
	storeEngineConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_type": "store",
		}),
	}

	// 过滤器取值来自上方 anyone 数据源查出的真实规格，避免写死环境差异的 cpu/memory 组合
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
	// 先查出 lindorm 引擎族的全量规格供各过滤步骤取用；限定 engine_type，
	// 否则 instance_types.0 可能落在其它引擎族，叠加 cpu+memory+engine_type 后无法保证非空
	return `
data "alibabacloudstack_lindorm_instance_types" "anyone" {
  engine_type = "lindorm"
}
`
}
