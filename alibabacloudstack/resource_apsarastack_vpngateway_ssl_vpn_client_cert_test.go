package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpngatewaySslVpnClient0(t *testing.T) {
	var v *VpcDescribesslvpnclientcertResponse

	resourceId := "alibabacloudstack_vpngateway_sslvpnclientcert.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpngatewaySslVpnClientCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribesslvpnclientcertRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpn_gateway_sslVpnClient%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpngatewaySslVpnClientBasicdependence)
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

					"ssl_vpn_client_cert_name": name,

					"ssl_vpn_server_id": "${alibabacloudstack_vpngateway_ssl_vpnserver.default.id}",

					"vpn_gateway_id": "${alibabacloudstack_vpngateway_vpngateway.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ssl_vpn_client_cert_name": name,

						"ssl_vpn_server_id": CHECKSET,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"ssl_vpn_client_cert_name": name + "modify",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ssl_vpn_client_cert_name": name + "modify",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"vpn_gateway_id"},
			},
		},
	})
}

var AlibabacloudTestAccVpngatewaySslVpnClientCheckmap = map[string]string{
	"client_config":          CHECKSET,
	"client_cert":            CHECKSET,
	"client_key":             CHECKSET,
	"ca_cert":                CHECKSET,
	"create_time":            CHECKSET,
	"end_time":               CHECKSET,
	"status":                 CHECKSET,
	"ssl_vpn_client_cert_id": CHECKSET,
}

func AlibabacloudTestAccVpngatewaySslVpnClientBasicdependence(name string) string {
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

resource "alibabacloudstack_vpngateway_ssl_vpnserver" "default" {
	client_ip_pool =  "10.8.0.0/24"

	local_subnet =  "192.168.1.0/24"

	ssl_vpn_server_name =  "${var.name}"
	vpn_gateway_id =     "${alibabacloudstack_vpngateway_vpngateway.default.id}"
	}



`, name, VSwitchCommonTestCase)
}
