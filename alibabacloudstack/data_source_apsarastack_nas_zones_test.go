package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlicloudNasZonesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	zoneIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasZonesDataSourceName(map[string]string{
			"zone_id": `"${data.alibabacloudstack_nas_zones.anyone.zones.0.zone_id}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasZonesDataSourceName(map[string]string{
			"zone_id": `"fake_zone_id"`,
		}),
	}
	protocolConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasZonesDataSourceName(map[string]string{
			"protocol": `"${data.alibabacloudstack_nas_zones.anyone.zones.0.protocols.0}"`,
		}),
		fakeConfig: testAccCheckAlicloudNasZonesDataSourceName(map[string]string{
			"protocol": `"fake_protocol"`,
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
		resourceId:   "data.alibabacloudstack_nas_zones.default",
		existMapFunc: existAlicloudNasZoneDataSourceNameMapFunc,
		fakeMapFunc:  fakeNasZonesMapFunc,
	}

	alicloudNasZonesAccountBusesCheckInfo.dataSourceTestCheck(t, rand, zoneIdConf, protocolConf)
}

func testAccCheckAlicloudNasZonesDataSourceName(attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s
data "alibabacloudstack_nas_zones" "anyone" {  
}
data "alibabacloudstack_nas_zones" "default" {  
   %s
}
`, DataZoneCommonTestCase, strings.Join(pairs, " \n "))
	return config
}
