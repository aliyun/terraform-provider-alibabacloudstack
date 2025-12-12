package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackApfsZonesDataSource_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_apfs_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", ApfsZonesDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existApfsZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#":                               CHECKSET,
			"zones.0.id":                            CHECKSET,
			"zones.0.zone_id":                       CHECKSET,
			"zones.0.clusters.#":                    CHECKSET,
			"zones.0.clusters.0.cluster_id":         CHECKSET,
			"zones.0.clusters.0.available_capacity": CHECKSET,
			"zones.0.clusters.0.storage_type":       CHECKSET,
		}
	}

	var fakeApfsZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}

	var apfsZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existApfsZonesMapFunc,
		fakeMapFunc:  fakeApfsZonesMapFunc,
	}

	apfsZonesCheckInfo.dataSourceTestCheck(t, 0, allConf)
}

func ApfsZonesDependence(name string) string {
	return ""
}
