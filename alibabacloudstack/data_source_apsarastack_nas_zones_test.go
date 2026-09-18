package alibabacloudstack

import (
	"testing"
)

func TestAccAlicloudNasZonesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	resourceId := "data.alibabacloudstack_nas_zones.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceNasZonesConfigDependence)

	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"zone_id": "${data.alibabacloudstack_nas_zones.anyone.zones.0.zone_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"zone_id": "fake_zone_id",
		}),
	}
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"protocol": "${data.alibabacloudstack_nas_zones.anyone.zones.0.protocols.0}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"protocol": "fake_protocol",
		}),
	}

	var existAlicloudNasZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#":         CHECKSET,
			"zones.0.zone_id": CHECKSET,
		}
	}
	var fakeNasZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": "0",
		}
	}
	var alicloudNasZonesAccountBusesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAlicloudNasZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeNasZonesMapFunc,
	}

	alicloudNasZonesAccountBusesCheckInfo.dataSourceTestCheck(t, rand, zoneIdConf, protocolConf)
}

func dataSourceNasZonesConfigDependence(name string) string {
	return `
data "alibabacloudstack_nas_zones" "anyone" {  
}
`
}
