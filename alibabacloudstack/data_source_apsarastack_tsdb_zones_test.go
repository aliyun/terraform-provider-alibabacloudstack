package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackTsdbZonesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackTsdbZonesSourceConfig(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":           CHECKSET,
			"zones.#":         CHECKSET,
			"zones.0.zone_id": CHECKSET,
			//"zones.0.local_name": CHECKSET,
		}
	}

	var fakeTsdbZonesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"zones.#": "0",
		}
	}

	var tsdbZonesRecordsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_tsdb_zones.default",
		existMapFunc: existTsdbZonesMapFunc,
		fakeMapFunc:  fakeTsdbZonesMapFunc,
	}

	tsdbZonesRecordsCheckInfo.dataSourceTestCheck(t, rand, allConf)

}

func testAccCheckAlibabacloudStackTsdbZonesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
data "alibabacloudstack_tsdb_zones" "default"{
%s
}

`, strings.Join(pairs, "\n   "))
	return config
}
