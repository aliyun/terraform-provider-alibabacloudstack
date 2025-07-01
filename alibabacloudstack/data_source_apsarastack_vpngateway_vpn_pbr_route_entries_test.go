package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackVpngatewayVpnPbrRouteEntriesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_vpngateway_vpn_pbr_route_entrys.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAcc%sVpngatewayVpnPbrRouteEntriesDataSource-%d", defaultRegionToTest, rand),
		dataSourceVpngatewayVpnPbrRouteEntriesDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"vpn_gateway_id": []string{"${alibabacloudstack_vpngateway_vpn_pbr_route_entry.default.vpn_gateway_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"vpn_gateway_id": []string{"${alibabacloudstack_vpngateway_vpn_pbr_route_entry.default.vpn_gateway_id}-fakeTestAcccc"},
		}),
	}

	var existVpngatewayVpnPbrRouteEntriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                  "1",
			"ids.0":                                  CHECKSET,
			"vpn_pbr_route_entries.#":                "1",
			"vpn_pbr_route_entries.0.description":    fmt.Sprintf("tf-testAcc%sVpngatewayVpnPbrRouteEntriesDataSource-%d", defaultRegionToTest, rand),
			"vpn_pbr_route_entries.0.next_hop":       CHECKSET,
			"vpn_pbr_route_entries.0.route_dest":     CHECKSET,
			"vpn_pbr_route_entries.0.status":         CHECKSET,
			"vpn_pbr_route_entries.0.vpn_gateway_id": CHECKSET,
		}
	}

	var fakeVpngatewayVpnPbrRouteEntriesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "0",
			"vpn_pbr_route_entries.#": "0",
		}
	}

	var VpngatewayVpnPbrRouteEntriesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existVpngatewayVpnPbrRouteEntriesMapFunc,
		fakeMapFunc:  fakeVpngatewayVpnPbrRouteEntriesMapFunc,
	}

	VpngatewayVpnPbrRouteEntriesCheckInfo.dataSourceTestCheck(t, rand, idsConf)
}

func dataSourceVpngatewayVpnPbrRouteEntriesDependence(name string) string {
	return fmt.Sprintf(`
data "alibabacloudstack_zones" "default"{
}

resource "alibabacloudstack_vpc" "default" {
 name  = "%s"
 cidr_block = "10.1.0.0/21"
}
resource "alibabacloudstack_vswitch" "default" {
 name			   = "${alibabacloudstack_vpc.default.name}"
 vpc_id            = "${alibabacloudstack_vpc.default.id}"
 cidr_block        = "10.1.1.0/24"
 availability_zone = "${data.alibabacloudstack_zones.default.ids.0}"
}
resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${alibabacloudstack_vpc.default.name}"
 vpc_id               = "${alibabacloudstack_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = true
 enable_ipsec		  = true
 vswitch_id			  = "${alibabacloudstack_vswitch.default.id}"
}
resource "alibabacloudstack_vpn_connection" "default" {
 name                = "${alibabacloudstack_vpc.default.name}"
 customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
 vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
 local_subnet        = ["192.168.2.0/24"]
 remote_subnet       = ["192.168.3.0/24"]
}
resource "alibabacloudstack_vpn_customer_gateway" "default" {
 name       = "${alibabacloudstack_vpc.default.name}"
 ip_address = "192.168.1.1"
}
resource "alibabacloudstack_vpngateway_vpn_pbr_route_entry" "default" {
	vpn_gateway_id = "${alibabacloudstack_vpn_gateway.default.id}"
	route_dest =     "10.0.0.0/24"
	route_source =   "192.168.0.0/24"
	next_hop =      "${alibabacloudstack_vpn_connection.default.id}"
	weight =       "100"
	publish_vpc =    "false"
}
 `, name)
}
