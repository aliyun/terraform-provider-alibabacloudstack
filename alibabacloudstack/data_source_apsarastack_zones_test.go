package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackZonesDataSource_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "basic"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"enable_details": true,
		}),
	}

	var existZonesBasicMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.slb_slave_zone_ids.#":          "0",
		}
	}

	var fakeZonesBasicMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesBasicCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesBasicMapFunc,
		fakeMapFunc:  fakeZonesBasicMapFunc,
	}
	zonesBasicCheckInfo.dataSourceTestCheck(t, 0, basicConf)
}

func TestAccAlibabacloudStackZonesDataSource_filter(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "filter"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	vswitchConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "VSwitch",
			"enable_details":              true,
		}),
	}
	kvStoreConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "KVStore",
			"enable_details":              true,
		}),
	}
	mongodbConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "MongoDB",
			"enable_details":              true,
		}),
	}
//	hbaseConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{
//			"available_resource_creation": "HBase",
//			"enable_details":              true,
//		}),
//	}
	adbConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "ADB",
			"enable_details":              true,
		}),
	}
	gpdbConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "Gpdb",
			"enable_details":              true,
		}),
	}
//	esConf := dataSourceTestAccConfig{
//		existConfig: testAccConfig(map[string]interface{}{
//			"available_resource_creation": "Elasticsearch",
//			"enable_details":              true,
//		}),
//	}

	var existZonesFilterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.slb_slave_zone_ids.#":          "0",
		}
	}

	var fakeZonesFilterMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesFilterCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesFilterMapFunc,
		fakeMapFunc:  fakeZonesFilterMapFunc,
		// FIXME: Need to increase deployment detection
		
	}
	zonesFilterCheckInfo.dataSourceTestCheck(t, 0, vswitchConf, kvStoreConf, mongodbConf, adbConf, gpdbConf)
}

func TestAccAlibabacloudStackZonesDataSource_filterIoOptimized(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "io_optimized"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	ioOptimizedConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "IoOptimized",
			"available_disk_category":     "${data.alibabacloudstack_zones.anyone.zones.0.available_disk_categories.0}",
			"enable_details":              true,
		}),
	}

	var existZonesIoOptimizedMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.slb_slave_zone_ids.#":          "0",
		}
	}

	var fakeZonesIoOptimizedMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesIoOptimizedCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesIoOptimizedMapFunc,
		fakeMapFunc:  fakeZonesIoOptimizedMapFunc,
	}
	zonesIoOptimizedCheckInfo.dataSourceTestCheck(t, 0, ioOptimizedConf)
}

func TestAccAlibabacloudStackZonesDataSource_unitRegion(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "unit_region"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	unitRegionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "VSwitch",
			"enable_details":              true,
		}),
	}

	var existZonesUnitRegionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.slb_slave_zone_ids.#":          "0",
		}
	}

	var fakeZonesUnitRegionMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesUnitRegionCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesUnitRegionMapFunc,
		fakeMapFunc:  fakeZonesUnitRegionMapFunc,
	}
	zonesUnitRegionCheckInfo.dataSourceTestCheck(t, 0, unitRegionConf)
}

func TestAccAlibabacloudStackZonesDataSource_multiZone(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "multi_zone"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	multiZoneConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation": "Rds",
			"multi":                       true,
			"enable_details":              true,
		}),
	}

	var existZonesMultiZoneMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.local_name":       CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
		}
	}

	var fakeZonesMultiZoneMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesMultiZoneCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesMultiZoneMapFunc,
		fakeMapFunc:  fakeZonesMultiZoneMapFunc,
	}
	zonesMultiZoneCheckInfo.dataSourceTestCheck(t, 0, multiZoneConf)
}

func TestAccAlibabacloudStackZonesDataSource_chargeType(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "charge_type"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	chargeTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_charge_type":        "PrePaid",
			"available_resource_creation": "Rds",
			"multi":                       true,
			"enable_details":              true,
		}),
	}

	var existZonesChargeTypeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.local_name":       CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
		}
	}

	var fakeZonesChargeTypeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesChargeTypeCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesChargeTypeMapFunc,
		fakeMapFunc:  fakeZonesChargeTypeMapFunc,
	}
	zonesChargeTypeCheckInfo.dataSourceTestCheck(t, 0, chargeTypeConf)
}

func TestAccAlibabacloudStackZonesDataSource_slb(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "slb"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	slbConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_resource_creation":      "Slb",
			"enable_details":                   true,
			"available_slb_address_ip_version": "ipv4",
			"available_slb_address_type":       "Vpc",
		}),
	}

	var existZonesSlbMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    CHECKSET,
			"zones.0.available_instance_types.#":    CHECKSET,
			"zones.0.available_resource_creation.#": CHECKSET,
			"zones.0.available_disk_categories.#":   CHECKSET,
			"zones.0.slb_slave_zone_ids.#":          CHECKSET,
		}
	}

	var fakeZonesSlbMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesSlbCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesSlbMapFunc,
		fakeMapFunc:  fakeZonesSlbMapFunc,
	}
	zonesSlbCheckInfo.dataSourceTestCheck(t, 0, slbConf)
}

func TestAccAlibabacloudStackZonesDataSource_enableDetails(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "enable_details"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	enableDetailsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existZonesEnableDetailsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                 CHECKSET,
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.local_name":                    "",
			"zones.0.available_instance_types.#":    "0",
			"zones.0.available_resource_creation.#": "0",
			"zones.0.available_disk_categories.#":   "0",
			"zones.0.slb_slave_zone_ids.#":          "0",
		}
	}

	var fakeZonesEnableDetailsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesEnableDetailsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesEnableDetailsMapFunc,
		fakeMapFunc:  fakeZonesEnableDetailsMapFunc,
	}
	zonesEnableDetailsCheckInfo.dataSourceTestCheck(t, 0, enableDetailsConf)
}

func TestAccAlibabacloudStackZonesDataSource_empty(t *testing.T) {
	resourceId := "data.alibabacloudstack_zones.default"
	name := "empty"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceZonesConfigDependence)

	emptyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"available_instance_type": "ecs.n1.fake",
		}),
	}

	var existZonesEmptyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var fakeZonesEmptyMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var zonesEmptyCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existZonesEmptyMapFunc,
		fakeMapFunc:  fakeZonesEmptyMapFunc,
	}
	zonesEmptyCheckInfo.dataSourceTestCheck(t, 0, emptyConf)
}

func dataSourceZonesConfigDependence(name string) string {
	return `
	data "alibabacloudstack_zones" "anyone" {
		enable_details = true
	}`
}
