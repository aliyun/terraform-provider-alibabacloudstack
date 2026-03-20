package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackKVStoreInstanceClasses(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_kvstore_instance_classes.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "KVStore", kvstoreConfigHeader)

	EngineVersionConfRedis := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine":         "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.engine}",
			"engine_version": "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.engine_version}",
		}),
	}

	// At present, there are some limitation for sorted
	//prePaidSortedByConfRedis := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alibabacloudstack_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PrePaid",
	//		"sorted_by":            "Price",
	//	}),
	//	existChangMap: map[string]string{
	//		"classes.#":                CHECKSET,
	//		"classes.0.instance_class": CHECKSET,
	//		"classes.0.price":          CHECKSET,
	//	},
	//}
	//
	//postPaidSortedByConfRedis := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":              "${data.alibabacloudstack_zones.resources.zones.0.id}",
	//		"engine":               "Redis",
	//		"engine_version":       "5.0",
	//		"instance_charge_type": "PostPaid",
	//		"sorted_by":            "Price",
	//	}),
	//	existChangMap: map[string]string{
	//		"classes.#":                CHECKSET,
	//		"classes.0.instance_class": CHECKSET,
	//		"classes.0.price":          CHECKSET,
	//	},
	//}
	NotExisted := dataSourceTestAccConfig{
		fakeConfig: testAccConfig(map[string]interface{}{
			"memory": "0.01",
		}),
	}

	editionTypeCommunity := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"edition_type": "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.edition_type}",
		}),
	}
	cpu2 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"cpu":       "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.cpu}",
			"sorted_by": "Memory",
		}),
	}
	memory2 := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"memory":    "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.memory}",
			"sorted_by": "CPU",
		}),
	}
	ArchitectureStandard := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"architecture": "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.architecture}",
		}),
	}
	// Not all of zone support rwsplit
	//ArchitectureRwsplit := dataSourceTestAccConfig{
	//	existConfig: testAccConfig(map[string]interface{}{
	//		"zone_id":      "${data.alibabacloudstack_zones.resources.zones.0.id}",
	//		"architecture": "rwsplit",
	//	}),
	//}
	NodeType := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"node_type": "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.node_type}",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine":         "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.engine}",
			"engine_version": "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.engine_version}",
			"architecture":   "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.architecture}",
			"edition_type":   "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.edition_type}",
			"node_type":      "${data.alibabacloudstack_kvstore_instance_classes.anyone.instance_classes.0.node_type}",
		}),
	}

	var existKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#":    CHECKSET,
			"instance_classes.0.id": CHECKSET,
		}
	}

	var fakeKVStoreInstanceMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"instance_classes.#": "0",
		}
	}

	var KVStoreInstanceCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKVStoreInstanceMapFunc,
		fakeMapFunc:  fakeKVStoreInstanceMapFunc,
	}

	// At present, the datasource does not support memcache
	//KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, EngineVersionConfRedis, EngineVersionConfMemcache,
	//	ChargeTypeConfPostpaid, PerformanceTypeStandardPerformanceType, PerformanceTypeEnhancePerformanceType,
	//	StorageTypeInmemory, PackageTypeStandard, PackageTypeCustomized, ArchitectureStandard, ArchitectureCluster,
	//	ArchitectureRwsplit, NodeTypeDouble, NodeTypeSingle, NodeTypeReadone, NodeTypeReadthree, NodeTypeReadfive,
	//	ArchitectureStandard, allConf)
	KVStoreInstanceCheckInfo.dataSourceTestCheck(t, rand, cpu2, memory2, EngineVersionConfRedis,
		//prePaidSortedByConfRedis, postPaidSortedByConfRedis
		editionTypeCommunity,
		ArchitectureStandard,
		NodeType, allConf, NotExisted)
}

func kvstoreConfigHeader(name string) string {
	return `
	data "alibabacloudstack_kvstore_instance_classes" "anyone" {
	}`
}
