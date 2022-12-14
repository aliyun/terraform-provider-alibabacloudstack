package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlicloudNasZonesDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100, 999)
	regionIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudNasZonesDataSourceName(rand, map[string]string{}),
		fakeConfig:  "",
	}

	var existAlicloudNasZoneDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"zones.#": CHECKSET,
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

	alicloudNasZonesAccountBusesCheckInfo.dataSourceTestCheck(t, rand, regionIdConf)
}

func testAccCheckAlicloudNasZonesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
data "alibabacloudstack_nas_zones" "default" {  
   %s
}
`, strings.Join(pairs, " \n "))
	return config
}
