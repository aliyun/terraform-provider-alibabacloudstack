package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpngatewaySslVpnServer0(t *testing.T) {
	var v *VpcDescribesslvpnserversResponse

	resourceId := "alibabacloudstack_vpngateway_ssl_vpnserver.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpngatewaySslVpnServerCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribesslvpnserversRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc-vpn_sslvpnserver%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpngatewaySslVpnServerBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"client_ip_pool": "10.8.0.0/24",

					"local_subnet": "192.168.1.0/24",

					"ssl_vpn_server_name": name,
					"vpn_gateway_id":      "${alibabacloudstack_vpngateway_vpngateway.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"client_ip_pool": "10.8.0.0/24",

						"local_subnet": "192.168.1.0/24",

						"ssl_vpn_server_name": name,

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

					"ssl_vpn_server_name": name + "_new",

					"client_ip_pool": "10.8.0.0/16",
					"local_subnet":   "192.168.10.0/24,192.168.1.0/16",
					"proto":          "TCP",
					"port":           "1193",
					"cipher":         "AES-128-CBC",
					"compress":       "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ssl_vpn_server_name": name + "_new",

						"client_ip_pool": "10.8.0.0/16",
						"local_subnet":   "192.168.10.0/24,192.168.1.0/16",
						"proto":          "TCP",
						"port":           "1193",
						"cipher":         "AES-128-CBC",
						"compress":       "true",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"ssl_vpn_server_name": name + "_new2",

					"client_ip_pool": "10.8.2.0/24",
					"local_subnet":   "192.168.2.0/24,192.168.3.0/16",
					"proto":          "UDP",
					"port":           "1194",
					"cipher":         "AES-192-CBC",
					"compress":       "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ssl_vpn_server_name": name + "_new2",

						"client_ip_pool": "10.8.2.0/24",
						"local_subnet":   "192.168.2.0/24,192.168.3.0/16",
						"proto":          "UDP",
						"port":           "1194",
						"cipher":         "AES-192-CBC",
						"compress":       "false",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpngatewaySslVpnServerCheckmap = map[string]string{}

func AlibabacloudTestAccVpngatewaySslVpnServerBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_vpngateway_vpngateway" "default" {
description= "${var.name}"

vpn_gateway_name = "${var.name}"

bandwidth = "5"

vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"

vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"

enable_ssl = "true"

instance_charge_type = "PostPaid"
}

`, name, VSwitchCommonTestCase)
}
