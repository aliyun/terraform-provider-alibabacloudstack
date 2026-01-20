package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackADBZonesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_adb_zones.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAdbZonesConfigDependence)

	multiConfig := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"multi": "false",
		}),
	}

	var existAdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      CHECKSET,
			"zones.#":    CHECKSET,
			"ids.0":      CHECKSET,
			"zones.0.id": CHECKSET,
		}
	}

	var fakeAdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var adbZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAdbZonesMapFunc,
		fakeMapFunc:  fakeAdbZonesMapFunc,
	}

	adbZonesCheckInfo.dataSourceTestCheck(t, rand, multiConfig)
}

func dataSourceAdbZonesConfigDependence(name string) string {
	return ""
}
