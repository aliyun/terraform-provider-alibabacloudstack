package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpngatewayVpnPbrRouteEntry_basic(t *testing.T) {
	var v *VpnGatewayVpnPbrRouteEntry

	resourceId := "alibabacloudstack_vpngateway_vpn_pbr_route_entry.default"
	ra := resourceAttrInit(resourceId, VpngatewayVpnpbrrouteentrybasicMap)

	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpnpbrrouteentriesRequest")

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testacc_Vpn_pbrrouteentrybasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpngatewayVpnpbrrouteentryConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${alibabacloudstack_vpn_gateway.default.id}",
					"route_dest":     "10.0.0.0/24",
					"route_source":   "192.168.0.0/24",
					"next_hop":       "${alibabacloudstack_vpn_connection.default.id}",
					"weight":         "100",
					"publish_vpc":    "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"route_dest":     "10.0.0.0/24",
						"route_source":   "192.168.0.0/24",
						"weight":         "100",
						"publish_vpc":    "false",
						"next_hop":       CHECKSET,
						"vpn_gateway_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"publish_vpc": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"publish_vpc": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"weight": "0",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"weight": "0",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// overlay_mode 不支持回读
				ImportStateVerifyIgnore: []string{"overlay_mode"},
			},
		},
	})
}

func resourceVpngatewayVpnpbrrouteentryConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
 default = "%s"
}

%s

resource "alibabacloudstack_vpn_connection" "default" {
 name                = "${var.name}"
 customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
 vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
 local_subnet        = ["192.168.2.0/24"]
 remote_subnet       = ["192.168.3.0/24"]
}
resource "alibabacloudstack_vpn_customer_gateway" "default" {
 name       = "${var.name}"
 ip_address = "192.168.1.1"
}
`, name, VpnGatewayCommonTestCase)
}

var VpngatewayVpnpbrrouteentrybasicMap = map[string]string{
	"vpn_gateway_id": CHECKSET,
	"route_dest":     "10.0.0.0/24",
	"next_hop":       CHECKSET,
	"weight":         "100",
	"publish_vpc":    "false",
}
