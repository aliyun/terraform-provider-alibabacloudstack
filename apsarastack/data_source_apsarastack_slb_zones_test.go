package apsarastack

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackSlbZonesDataSource_basic(t *testing.T) {
	rand := acctest.RandInt()
	resourceId := "data.apsarastack_slb_zones.default"

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
