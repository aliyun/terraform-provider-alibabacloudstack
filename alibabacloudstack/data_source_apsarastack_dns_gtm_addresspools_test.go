package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsGtmAddressPoolsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_gtm_addresspools.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourceDnsGtmAddressPoolsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"name_regex": "${alibabacloudstack_dns_gtm_addresspool.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"name_regex": "^fake.*",
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"type": "A",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":  []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"type": "AAAA",
		}),
	}

	lbaStrategyConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"lba_strategy": "RATIO",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
			"lba_strategy": "ALL_RR",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_gtm_addresspool.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, typeConf, lbaStrategyConf)
}

func dataSourceDnsGtmAddressPoolsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
}
