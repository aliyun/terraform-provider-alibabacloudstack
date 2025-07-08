package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackExpressconnectBgpNetworksDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 2000)
	resourceId := "data.alibabacloudstack_expressconnect_bgp_networks.default"
	testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")

	name := fmt.Sprintf("EcBgpNetworksDataSource-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceExpressconnectBgpNetworksDependence(rand))

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":       []string{"${alibabacloudstack_expressconnect_bgp_network.default.id}"},
			"router_id": "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id": "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"ids":       []string{"${alibabacloudstack_expressconnect_bgp_network.default.id}" + "fake"},
		}),
	}

	dstcidrblockConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"router_id":      "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"dst_cidr_block": "1.1.1.1",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"router_id":      "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
			"dst_cidr_block": "1.1.1.1-fake",
		}),
	}

	var existExpressconnectBgpNetworksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"bgp_networks.#":                "1",
			"bgp_networks.0.status":         CHECKSET,
			"bgp_networks.0.dst_cidr_block": CHECKSET,
			"bgp_networks.0.router_id":      CHECKSET,
		}
	}

	var fakeExpressconnectBgpNetworksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":          "0",
			"bgp_networks.#": "0",
		}
	}

	var ExpressconnectBgpNetworksCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existExpressconnectBgpNetworksMapFunc,
		fakeMapFunc:  fakeExpressconnectBgpNetworksMapFunc,
	}

	ExpressconnectBgpNetworksCheckInfo.dataSourceTestCheck(t, rand, idsConf, dstcidrblockConf)
}

func dataSourceExpressconnectBgpNetworksDependence(vlanId int) func(name string) string {
	return func(name string) string {
		return fmt.Sprintf(` 
	variable "name" {
	  default = "%s"
	}

	resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
		physical_connection_id =     %s
		vlan_id =                    %d
		local_gateway_ip =           "10.0.0.1"
		peer_gateway_ip =            "10.0.0.2"
		peering_subnet_mask =        "255.255.255.252"
		virtual_border_router_name = "${var.name}"
		description =                "TestAccAlibabacloudStackExpressconnectBgpgroup_basic0"
	}
	
	resource "alibabacloudstack_expressconnect_bgp_network" "default" {
		
        dst_cidr_block = "1.1.1.1"
		router_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	}
	`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
	}
}
