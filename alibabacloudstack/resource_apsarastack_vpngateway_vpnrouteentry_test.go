package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpnRouteEntry_basic(t *testing.T) {
	var v vpc.VpnRouteEntry

	resourceId := "alibabacloudstack_vpn_route_entry.default"
	ra := resourceAttrInit(resourceId, vpnRouteEntrybasicMap)

	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf-testaccvpnRouteEntrybasic%v", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnRouteEntryConfigDependence)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		// module name
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"vpn_gateway_id": "${alibabacloudstack_vpn_gateway.default.id}",
					"route_dest":     "10.0.0.0/24",
					"next_hop":       "${alibabacloudstack_vpn_connection.default.id}",
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
					"vpn_gateway_id": "${alibabacloudstack_vpn_gateway.default.id}",
					"route_dest":     "10.0.0.0/24",
					"next_hop":       "${alibabacloudstack_vpn_connection.default.id}",
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

func resourceVpnRouteEntryConfigDependence(name string) string {
	return fmt.Sprintf(`

variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_vpn_gateway" "default" {
 name                 = "${alibabacloudstack_vpc_vpc.default.name}"
 vpc_id               = "${alibabacloudstack_vpc_vpc.default.id}"
 bandwidth            = 10
 instance_charge_type = "PostPaid"
 enable_ssl           = false
 ipsec_vpn            = true
 vswitch_id			  = "${alibabacloudstack_vpc_vswitch.default.id}"
}
resource "alibabacloudstack_vpn_connection" "default" {
 name                = "${var.name}"
 customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
 vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
 local_subnet        = ["192.168.2.0/24"]
 remote_subnet       = ["192.168.3.0/24"]
}
resource "alibabacloudstack_vpn_customer_gateway" "default" {
 name       = "${var.name}"
 ip_address = "172.16.1.100"
}
`, name, VSwitchCommonTestCase)
}

var vpnRouteEntrybasicMap = map[string]string{
	"vpn_gateway_id": CHECKSET,
	"route_dest":     "10.0.0.0/24",
	"next_hop":       CHECKSET,
	"weight":         "100",
	"publish_vpc":    "false",
}
