package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectBgpPeerDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_expressconnect_bgp_peers.default"
	testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")

	name := fmt.Sprintf("tf-testAcc-ExpressconnectBgpPeerDataSource-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressconnectBgpPeerDependence(rand))

	descriptionRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id":    "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"bgp_group_id": "${alibabacloudstack_expressconnect_bgp_peer.default.bgp_group_id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id":    "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"bgp_group_id": "${alibabacloudstack_expressconnect_bgp_peer.default.bgp_group_id}-fakeTestAcccc",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id":    "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"bgp_group_id": "${alibabacloudstack_expressconnect_bgp_peer.default.bgp_group_id}",
			"ids":          []string{"${alibabacloudstack_expressconnect_bgp_peer.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id":    "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"bgp_group_id": "${alibabacloudstack_expressconnect_bgp_peer.default.bgp_group_id}-fakeTestAcccc",
			"ids":          []string{"${alibabacloudstack_expressconnect_bgp_peer.default.id}-fakeTestAcccc"},
		}),
	}

	var existExpressconnectBgpPeerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"ids.0":                   CHECKSET,
			"bgp_peers.#":             "1",
			"bgp_peers.0.hold":        CHECKSET,
			"bgp_peers.0.ip_version":  CHECKSET,
			"bgp_peers.0.is_fake":     CHECKSET,
			"bgp_peers.0.keepalive":   CHECKSET,
			"bgp_peers.0.local_asn":   CHECKSET,
			"bgp_peers.0.peer_asn":    CHECKSET,
			"bgp_peers.0.route_limit": CHECKSET,
			"bgp_peers.0.router_id":   CHECKSET,
		}
	}

	var fakeExpressconnectBgpPeerMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":       "0",
			"bgp_peers.#": "0",
		}
	}

	var ExpressconnectBgpPeerCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressconnectBgpPeerMapFunc,
		fakeMapFunc:  fakeExpressconnectBgpPeerMapFunc,
	}

	ExpressconnectBgpPeerCheckInfo.dataSourceTestCheck(t, rand, descriptionRegexConf, allConf)
}

func dataSourceExpressconnectBgpPeerDependence(vlanId int) func(name string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
		physical_connection_id =     "%s"
		vlan_id =                    %d
		local_gateway_ip =           "10.0.0.1"
		peer_gateway_ip =            "10.0.0.2"
		peering_subnet_mask =        "255.255.255.252"
		virtual_border_router_name = "${var.name}"
		description =                "TestAccAlibabacloudStackExpressconnectBgpgroup_basic0"
	}
	
	resource "alibabacloudstack_expressconnect_bgp_group" "default" {
		bgp_group_name = "${var.name}"
		description =    "${var.name}"
		local_asn =      "65534"
		peer_asn =       "10"
		router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	}
	
	resource "alibabacloudstack_expressconnect_bgp_peer" "default" {
		bgp_group_id =   "${alibabacloudstack_expressconnect_bgp_group.default.id}"
		router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
		enable_bfd =      "true"
		peer_ip_address = "192.168.0.1"
		bfd_multi_hop =   "10"
	}
	`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
