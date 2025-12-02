package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDnsGtmAddressPoolsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_dns_gtm_addresspools.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                        "1",
				"address_pools.#":              "1",
				"address_pools.0.name":         fmt.Sprintf("tfacc%d", rand),
				"address_pools.0.type":         "A",
				"address_pools.0.lba_strategy": "RATIO",
				"address_pools.0.addrs.#":      "2",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           "0",
				"address_pools.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_dns_gtm_addresspool.default.name}"`,
		}),
		fakeConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"name_regex": `"^fake.*"`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"type": `"A"`,
		}),
		fakeConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"type": `"AAAA"`,
		}),
	}

	lbaStrategyConf := dataSourceTestAccConfig{
		existConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"lba_strategy": `"RATIO"`,
		}),
		fakeConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"lba_strategy": `"ALL_RR"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_dns_gtm_addresspool.default.id}"]`,
		}),
		fakeConfig: testAccDnsGtmAddressPoolsDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, typeConf, lbaStrategyConf, idsConf)
}

func testAccDnsGtmAddressPoolsDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tfacc%d"
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name         = var.name
  type         = "A"
  lba_strategy = "RATIO"

  addrs {
    value      = "192.168.1.1"
    mode       = "SMART"
    lba_weight = 20
  }

  addrs {
    value      = "127.0.0.1"
    mode       = "SMART"
    lba_weight = 80
  }
}

data "alibabacloudstack_dns_gtm_addresspools" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))

	return config
}
