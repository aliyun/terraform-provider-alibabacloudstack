package alibabacloudstack

import (
	"testing"

	
)

func TestAccAlibabacloudStackSlbZonesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(10000,20000)
	resourceId := "data.alibabacloudstack_slb_zones.default"

	var existSlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                        CHECKSET,
			"zones.#":                      CHECKSET,
			"zones.0.slb_slave_zone_ids.#": CHECKSET,
		}
	}

	var fakeSlbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var slbZonesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existSlbZonesMapFunc,
		fakeMapFunc:  fakeSlbZonesMapFunc,
	}

	slbZonesCheckInfo.dataSourceTestCheck(t, rand)
}

func dataSourceslbZonesConfigDependence(name string) string {
	return ""
}
