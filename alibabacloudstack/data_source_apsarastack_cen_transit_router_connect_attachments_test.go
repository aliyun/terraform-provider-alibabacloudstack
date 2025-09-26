package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterConnectAttachmentsDataSourceBasic(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_name}_fake"`,
		}),
	}
	IdsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_connect_attachment.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_connect_attachment.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_name}"`,
			"ids":        `["${alibabacloudstack_cen_transit_router_connect_attachment.default.id}" ]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_name}_fake"`,
			"ids":        `["${alibabacloudstack_cen_transit_router_connect_attachment.default.id}_fake"]`,
		}),
	}

	RouterConnectAttachmentsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, IdsConf, allConf)
}

var existRouterConnectAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transitrouterattachments.#": "1",
		"ids.#":                      "1",
		"transitrouterattachments.0.creation_time":                  CHECKSET,
		"transitrouterattachments.0.cen_id":                         CHECKSET,
		"transitrouterattachments.0.transport_type":                 CHECKSET,
		"transitrouterattachments.0.resource_type":                  CHECKSET,
		"transitrouterattachments.0.status":                         CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_id":   CHECKSET,
		"transitrouterattachments.0.resource_id":                    CHECKSET,
		"transitrouterattachments.0.resource_owner_id":              CHECKSET,
		"transitrouterattachments.0.transit_router_attachment_name": CHECKSET,
	}
}

var fakeRouterConnectAttachmentsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"transit_router_route_tables.#": "0",
		"ids.#":                         "0",
	}
}

var RouterConnectAttachmentsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cen_transit_router_connect_attachments.default",
	existMapFunc: existRouterConnectAttachmentsMapFunc,
	fakeMapFunc:  fakeRouterConnectAttachmentsMapFunc,
}

func testAccCheckAlibabacloudStacRouterConnectAttachmentsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-testacc_router_connect_attachment%v"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "${var.name}"
  cen_instance_name = "${var.name}"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id = "%s"
	vlan_id =                   %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" {
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	transit_router_attachment_name = "${var.name}"
	depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

data "alibabacloudstack_cen_transit_router_connect_attachments" "default" {
	cen_id = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
	transit_router_id = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_id}"
	%s
}`, rand, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), rand, strings.Join(pairs, "\n  "))
	return config
}
