package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackTsdbZonesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_tsdb_zones.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceTsdbZonesConfigDependence)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
	}

	var existTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           CHECKSET,
			"zones.#":         CHECKSET,
			"zones.0.zone_id": CHECKSET,
		}
	}

	var fakeTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var tsdbZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existTsdbZonesMapFunc,
		fakeMapFunc:  fakeTsdbZonesMapFunc,
	}

	tsdbZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)
}

func dataSourceTsdbZonesConfigDependence(name string) string {
	return `
`
}
