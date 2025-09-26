package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenVbrHealthChecksDataSource(t *testing.T) {
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_cen_vbr_health_checks.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                                      "1",
				"vbr_health_checks.#":                        "1",
				"vbr_health_checks.0.cen_id":                 CHECKSET,
				"vbr_health_checks.0.vbr_instance_id":        CHECKSET,
				"vbr_health_checks.0.vbr_instance_region_id": CHECKSET,
				"vbr_health_checks.0.health_check_source_ip": "172.16.0.1",
				"vbr_health_checks.0.health_check_target_ip": "10.0.0.1",
				"vbr_health_checks.0.health_check_interval":  "2",
				"vbr_health_checks.0.healthy_threshold":      "3",
				"vbr_health_checks.0.health_check_only":      "true",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":               "0",
				"vbr_health_checks.#": "0",
			}
		},
	}

	CenVbrHealthCheckCommonTestCaseNew := func(rand int, attrMap map[string]string) string {
		var pairs []string
		for k, v := range attrMap {
			pairs = append(pairs, k+" = "+v)
		}
		return fmt.Sprintf(`
variable "name" {
  default = "tf-testAcc%d"
}

resource "alibabacloudstack_cen_instance" "default" {
  cen_instance_name = var.name
  description       = var.name
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "%s"
  virtual_border_router_name = var.name
  vlan_id                    = %d
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_vbr_health_check" "default" {
    cen_id                 = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
    vbr_instance_id        = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
    health_check_source_ip = "172.16.0.1"
    health_check_target_ip = "10.0.0.1"
    healthy_threshold      = 3
    health_check_interval  = 2
    health_check_only      = true
    depends_on             = [alibabacloudstack_cen_transit_router_vbr_attachment.default]
}

data "alibabacloudstack_cen_vbr_health_checks" "default" {
    %s
}
`, rand, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), rand, strings.Join(pairs, "\n   "))
	}

	rand := getAccTestRandInt(1000, 2000)
	idsConf := dataSourceTestAccConfig{
		existConfig: CenVbrHealthCheckCommonTestCaseNew(rand, map[string]string{
			"cen_id": `"${alibabacloudstack_cen_vbr_health_check.default.cen_id}"`,
			"ids":    `["${alibabacloudstack_cen_vbr_health_check.default.cen_id}:${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"]`,
		}),
		fakeConfig: CenVbrHealthCheckCommonTestCaseNew(rand, map[string]string{
			"cen_id": `"${alibabacloudstack_cen_vbr_health_check.default.cen_id}"`,
			"ids":    `["cen-fakeid:vbr-fakeid"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: CenVbrHealthCheckCommonTestCaseNew(rand, map[string]string{
			"cen_id":          `"${alibabacloudstack_cen_vbr_health_check.default.cen_id}"`,
			"vbr_instance_id": `"${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"`,
			"ids":             `["${alibabacloudstack_cen_vbr_health_check.default.cen_id}:${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"]`,
		}),
		fakeConfig: CenVbrHealthCheckCommonTestCaseNew(rand, map[string]string{
			"cen_id":          `"${alibabacloudstack_cen_vbr_health_check.default.cen_id}"`,
			"vbr_instance_id": `"${alibabacloudstack_cen_vbr_health_check.default.vbr_instance_id}"`,
			"ids":             `["cen-fakeid:vbr-fakeid"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, idsConf, allConf)
}
