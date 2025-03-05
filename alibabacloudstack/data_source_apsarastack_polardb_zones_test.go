package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackPolardbZonesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_polardb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourcePolardbZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "false",
		}),
	}

	var existDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
		}
	}

	var fakeDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var DBZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDBZonesMapFunc,
		fakeMapFunc:  fakeDBZonesMapFunc,
	}

	DBZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourcePolardbZonesConfigDependence(name string) string {
	return ""
}
