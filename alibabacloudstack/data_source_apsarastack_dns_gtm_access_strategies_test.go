package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsGtmAccessStrategiesDataSource_basic(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_dns_gtm_access_strategies.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourceDnsGtmAccessStrategiesConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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
		existConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"name_regex":      "${alibabacloudstack_dns_gtm_access_strategy.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"name_regex":      "^fake.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"ids":             []string{"${alibabacloudstack_dns_gtm_access_strategy.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"ids":             []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"ids":             []string{"${alibabacloudstack_dns_gtm_access_strategy.default.id}"},
			"name_regex":      "${alibabacloudstack_dns_gtm_access_strategy.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"gtm_instance_id": "${alibabacloudstack_dns_gtm_instance.default.id}",
			"ids":             []string{"fake-id"},
			"name_regex":      "^fake.*",
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceDnsGtmAccessStrategiesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
`, name)
}
