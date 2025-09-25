package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCenTransitRouterConnectPeersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cen_transit_router_connect_peers.default"

	attr := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"peers.0.name":                  fmt.Sprintf("tf-testAccTransitRouterConnectPeer%d", rand),
				"peers.0.local_ip":              "172.16.0.1",
				"peers.0.peer_ip":               "172.16.0.11",
				"peers.0.status":                CHECKSET,
				"peers.0.region_id":             CHECKSET,
				"peers.0.connect_attachment_id": CHECKSET,
				"peers.0.cen_id":                CHECKSET,
				"peers.0.peer_id":               CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"peers.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
		}),
		fakeConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"^fake-name$"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_cen_transit_router_connect_peer.default.id}"]`,
		}),
		fakeConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	transitRouterConnectPeerNameConf := dataSourceTestAccConfig{
		existConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"peer_name": `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
		}),
		fakeConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"peer_name": `"fake-peer-name"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
			"ids":        `["${alibabacloudstack_cen_transit_router_connect_peer.default.id}"]`,
			"peer_name":  `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
		}),
		fakeConfig: TransitRouterConnectPeerCommonTestCaseNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
			"ids":        `["fake-id"]`,
			"peer_name":  `"${alibabacloudstack_cen_transit_router_connect_peer.default.name}"`,
		}),
	}

	attr.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, transitRouterConnectPeerNameConf, allConf)
}

func TransitRouterConnectPeerCommonTestCaseNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccTransitRouterConnectPeer%d"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
  transit_router_cidrs {
    cidr = "172.16.0.0/16"
  }
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  physical_connection_id = "%s"
  vlan_id =                    99
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
  cen_id = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
  transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
  transit_router_attachment_name = "${var.name}"
  depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

resource "alibabacloudstack_cen_transit_router_connect_peer" "default" {
  name                  = "${var.name}"
  local_ip              = "172.16.0.1"
  peer_ip               = "172.16.0.11"
  cen_id                = "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_id}"
}

data "alibabacloudstack_cen_transit_router_connect_peers" "default" {
  cen_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.cen_id}"
  connect_attachment_id = "${alibabacloudstack_cen_transit_router_connect_peer.default.connect_attachment_id}"
  %s
}
`, rand, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), strings.Join(pairs, "\n   "))
}
