package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDnsGtmAccessStrategiesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_dns_gtm_access_strategies.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":             "1",
				"strategies.#":      "1",
				"strategies.0.name": fmt.Sprintf("tfacc%d", rand),
				"strategies.0.default_gtm_address_pool_id":  CHECKSET,
				"strategies.0.failover_gtm_address_pool_id": CHECKSET,
				"strategies.0.line_ids.#":                   "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":        "0",
				"strategies.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_dns_gtm_access_strategy.default.name}"`,
		}),
		fakeConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"name_regex": `"^fake.*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_dns_gtm_access_strategy.default.id}"]`,
		}),
		fakeConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_dns_gtm_access_strategy.default.id}"]`,
			"name_regex": `"${alibabacloudstack_dns_gtm_access_strategy.default.name}"`,
		}),
		fakeConfig: buildDnsGtmAccessStrategiesDataSourceConfig(rand, map[string]string{
			"ids":        `["fake-id"]`,
			"name_regex": `"^fake.*"`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func buildDnsGtmAccessStrategiesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tfacc%d"
}

resource "alibabacloudstack_dns_private_domain" "default" {
  name = "${var.name}.local."
}

resource "alibabacloudstack_dns_gtm_instance" "default" {
  name     = var.name
  prefix   = var.name
  zone_id  = alibabacloudstack_dns_private_domain.default.id
  ttl      = 300
}

resource "alibabacloudstack_dns_private_line" "default" {
  name         = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

resource "alibabacloudstack_dns_gtm_addresspool" "default" {
  name        = "${var.name}-default"
  type        = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "1.1.1.1"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_addresspool" "failover" {
  name        = "${var.name}-failover"
  type        = "A"
  lba_strategy = "ALL_RR"
  addrs {
    value = "2.2.2.2"
    mode  = "SMART"
  }
}

resource "alibabacloudstack_dns_gtm_access_strategy" "default" {
  name                           = var.name
  gtm_instance_id                = alibabacloudstack_dns_gtm_instance.default.id
  switch_mode                    = "BY_PROBE_RESULT"
  default_min_available_addr_num = 1
  failover_min_available_addr_num= 1
  line_ids                       = [alibabacloudstack_dns_private_line.default.id]
  default_gtm_address_pool_id    = alibabacloudstack_dns_gtm_addresspool.default.id
  default_gtm_address_pool_type  = "IPV4"
  failover_gtm_address_pool_id   = alibabacloudstack_dns_gtm_addresspool.failover.id
  failover_gtm_address_pool_type = "IPV4"
}

data "alibabacloudstack_dns_gtm_access_strategies" "default" {
	gtm_instance_id = alibabacloudstack_dns_gtm_instance.default.id
%s
}
`, rand, strings.Join(pairs, "\n"))
}
