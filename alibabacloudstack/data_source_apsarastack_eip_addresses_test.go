package alibabacloudstack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackEipsDataSourceBasic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_eip.default.0.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ids": `[ "${alibabacloudstack_eip.default.0.id}_fake" ]`,
		}),
	}

	ipsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ip_addresses": `[ "${alibabacloudstack_eip.default.0.ip_address}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ip_addresses": `[ "${alibabacloudstack_eip.default.0.ip_address}_fake" ]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ids":          `[ "${alibabacloudstack_eip.default.0.id}" ]`,
			"ip_addresses": `[ "${alibabacloudstack_eip.default.0.ip_address}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEipsDataSourceConfig(rand, map[string]string{
			"ids":          `[ "${alibabacloudstack_eip.default.0.id}" ]`,
			"ip_addresses": `[ "${alibabacloudstack_eip.default.0.ip_address}_fake" ]`,
		}),
	}

	dnsEipsCheckInfo.dataSourceTestCheck(t, rand, idsConf, ipsConf, allConf) //tagsConf, allConf)

}

func testAccCheckAlibabacloudStackEipsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`

resource "alibabacloudstack_eip" "default" {
  name = "tf-testAccCheckAlibabacloudstackEipsDataSourceConfig%d"
  count = 2
  bandwidth = 5
}
data "alibabacloudstack_eips" "default" {
  %s
}`, rand, strings.Join(pairs, "\n  "))
}

var existEipsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":                "1",
		"names.#":              "1",
		"eips.#":               "1",
		"eips.0.id":            CHECKSET,
		"eips.0.status":        string(Available),
		"eips.0.ip_address":    CHECKSET,
		"eips.0.bandwidth":     "5",
		"eips.0.instance_id":   "",
		"eips.0.instance_type": "",
		//"eips.0.internet_charge_type": string(PayByTraffic),
		//"eips.0.creation_time":        CHECKSET,
	}
}

var fakeEipsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"ids.#":   "0",
		"names.#": "0",
		"eips.#":  "0",
	}
}

var dnsEipsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_eips.default",
	existMapFunc: existEipsMapFunc,
	fakeMapFunc:  fakeEipsMapFunc,
}
