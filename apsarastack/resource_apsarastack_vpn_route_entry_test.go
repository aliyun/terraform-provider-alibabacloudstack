package apsarastack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/aliyun/terraform-provider-alibabaCloudStack/apsarastack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccApsaraStackVpnRouteEntry_basic(t *testing.T) {
	var v vpc.VpnRouteEntry

	resourceId := "apsarastack_vpn_route_entry.default"
	ra := resourceAttrInit(resourceId, vpnRouteEntrybasicMap)

	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%svpnRouteEntrybasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnRouteEntryConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${apsarastack_vpn_gateway.default.id}",
					"route_dest":     "10.0.0.0/24",
					"next_hop":       "${apsarastack_vpn_connection.default.id}",
					"weight":         "100",
					"publish_vpc":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_dest":     "10.0.0.0/24",
						"weight":         "100",
						"publish_vpc":    "false",
						"next_hop":       CHECKSET,
						"vpn_gateway_id": CHECKSET,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"publish_vpc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"publish_vpc": "true"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{"weight": "0"}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${apsarastack_vpn_gateway.default.id}",
					"route_dest":     "10.0.0.0/24",
					"next_hop":       "${apsarastack_vpn_connection.default.id}",
					"weight":         "100",
					"publish_vpc":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(vpnRouteEntrybasicMap),
				),
			},
		},
	})
}

func TestAccApsaraStackVpnRouteEntry_multi(t *testing.T) {
	var v vpc.VpnRouteEntry

	resourceId := "apsarastack_vpn_route_entry.default.1"
	ra := resourceAttrInit(resourceId, vpnRouteEntrybasicMap)
	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.ApsaraStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandInt()
	name := fmt.Sprintf("tf-testacc%svpnRouteEntrybasic%v", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnRouteEntryConfigDependence)

	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithAccountSiteType(t, IntlSite)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${apsarastack_vpn_gateway.default.id}",
					"route_dest":     "${element(var.route_dests, count.index)}",
					"next_hop":       "${apsarastack_vpn_connection.default.id}",
					"weight":         "0",
					"publish_vpc":    "false",
					"count":          "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vpn_gateway_id": CHECKSET,
						"route_dest":     "10.1.0.0/32",
						"next_hop":       CHECKSET,
						"weight":         "0",
						"publish_vpc":    "false",
					}),
				),
			},
		},
	})
}

func resourceVpnRouteEntryConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "route_dests" {
 default = ["10.1.0.0/24", "10.1.0.0/32"]
}
data "apsarastack_zones" "default"{
}
resource "apsarastack_vpc" "default" {
 name  = "%s"
 cidr_block = "10.1.0.0/21"
}
resource "apsarastack_vswitch" "default" {
 name			   = "${apsarastack_vpc.default.name}"
 vpc_id            = "${apsarastack_vpc.default.id}"
 cidr_block        = "10.1.1.0/24"
 availability_zone = "${data.apsarastack_zones.default.ids.0}"
}
resource "apsarastack_vpn_gateway" "default" {
 name                 = "${apsarastack_vpc.default.name}"
 vpc_id               = "${apsarastack_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = false
 vswitch_id			  = "${apsarastack_vswitch.default.id}"
}
resource "apsarastack_vpn_connection" "default" {
 name                = "${apsarastack_vpc.default.name}"
 customer_gateway_id = "${apsarastack_vpn_customer_gateway.default.id}"
 vpn_gateway_id      = "${apsarastack_vpn_gateway.default.id}"
 local_subnet        = ["192.168.2.0/24"]
 remote_subnet       = ["192.168.3.0/24"]
}
resource "apsarastack_vpn_customer_gateway" "default" {
 name       = "${apsarastack_vpc.default.name}"
 ip_address = "192.168.1.1"
}
`, name)
}

var vpnRouteEntrybasicMap = map[string]string{
	"vpn_gateway_id": CHECKSET,
	"route_dest":     "10.0.0.0/24",
	"next_hop":       CHECKSET,
	"weight":         "100",
	"publish_vpc":    "false",
}
