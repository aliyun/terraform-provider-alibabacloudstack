package alibabacloudstack

import (
	"testing"

	
)

func TestAccAlibabacloudStackMongoDBZonesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	resourceId := "data.alibabacloudstack_mongodb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceMongoDBZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "false",
		}),
	}

	var existDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                    CHECKSET,
			"ids.0":                    CHECKSET,
			"zones.#":                  CHECKSET,
			"zones.0.id":               CHECKSET,
			"zones.0.multi_zone_ids.#": CHECKSET,
		}
	}

	var fakeMongDBZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var MongDBZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existDBZonesMapFunc,
		fakeMapFunc:  fakeMongDBZonesMapFunc,
	}

	MongDBZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceMongoDBZonesConfigDependence(name string) string {
	return ""
}
