package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackMongoDBInstanceTypesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_mongodb_instance_types.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccMongodbInstanceTypesDataSource_%d", rand),
		dataSourceMongodbInstanceTypesConfigDependence)

	replicateConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
			"sorted_by":        "CPU",
		}),
	}

	mongosConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version":   "4.0",
			"db_instnace_type": "sharding",
			"node_type":        "mongos",
			"sorted_by":        "CPU",
		}),
	}
	shardConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version":   "4.0",
			"db_instnace_type": "sharding",
			"node_type":        "shard",
			"sorted_by":        "CPU",
		}),
	}
	configserverConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version":   "4.0",
			"db_instnace_type": "sharding",
			"node_type":        "configserver",
			"sorted_by":        "CPU",
		}),
	}

	testAccConfig = dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf_testAccMongodbInstanceTypesDataSource_%d", rand),
		dataSourceMongodbInstanceTypesPresetDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"${data.alibabacloudstack_mongodb_instance_types.preset.instance_types.0.id}"},
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
			"sorted_by":        "Memory",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":              []string{"xxxxx"},
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
			"sorted_by":        "Memory",
		}),
	}
	cpuConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":              "1",
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"cpu":              "800",
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
		}),
	}
	memoryConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory":           "2",
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory":           "3200",
			"engine_version":   "4.0",
			"db_instnace_type": "replicate",
		}),
	}

	var existMongodbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                           CHECKSET,
			"ids.0":                           CHECKSET,
			"instance_types.#":                CHECKSET,
			"instance_types.0.id":             CHECKSET,
			"instance_types.0.cpu":            CHECKSET,
			"instance_types.0.memory":         CHECKSET,
			"instance_types.0.engine_version": CHECKSET,
			"instance_types.0.cpu_type":       CHECKSET,
			"instance_types.0.series":         CHECKSET,
			"instance_types.0.connections":    CHECKSET,
			"instance_types.0.storage_min":    CHECKSET,
			"instance_types.0.storage_max":    CHECKSET,
		}
	}

	var fakeMongodbInstanceTypesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":            "0",
			"instance_types.#": "0",
		}
	}

	var MongodbInstanceTypesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMongodbInstanceTypesMapFunc,
		fakeMapFunc:  fakeMongodbInstanceTypesMapFunc,
	}

	MongodbInstanceTypesCheckInfo.dataSourceTestCheck(t, rand, replicateConf, mongosConf, shardConf, configserverConf, idsConf, cpuConf, memoryConf)
}

func dataSourceMongodbInstanceTypesConfigDependence(name string) string {
	return ""
}

func dataSourceMongodbInstanceTypesPresetDependence(name string) string {
	return `
	data "alibabacloudstack_mongodb_instance_types" "preset" {
		db_instnace_type = "replicate"
		engine_version = "4.0"
	}
`
}
