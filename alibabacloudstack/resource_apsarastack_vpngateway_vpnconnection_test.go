package alibabacloudstack

import (
	"fmt"
	"testing"

	

	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func equalSubnet(astr string, bstr string) bool {
	aArray := strings.Split(astr, ",")
	bArray := strings.Split(bstr, ",")
	if len(aArray) != len(bArray) {
		return false
	}

	for _, item := range aArray {
		if !strings.Contains(bstr, item) {
			return false
		}
	}
	return true
}

func testAccCheckVpnConnectionAttr(vpnConn *vpc.DescribeVpnConnectionResponse, localSubnet, remoteSubnet string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if !equalSubnet(vpnConn.LocalSubnet, localSubnet) {
			return errmsgs.WrapError(errmsgs.Error("wrong local subnet, expect %s, get %s", localSubnet, vpnConn.LocalSubnet))
		}

		if !equalSubnet(vpnConn.RemoteSubnet, remoteSubnet) {
			return errmsgs.WrapError(errmsgs.Error("wrong remote subnet, expect %s, get %s", remoteSubnet, vpnConn.RemoteSubnet))
		}

		return nil
	}
}
func TestAccAlibabacloudstackVpnConnectionBasic(t *testing.T) {
	var v vpc.DescribeVpnConnectionResponse

	resourceId := "alibabacloudstack_vpn_connection.default"
	ra := resourceAttrInit(resourceId, testAccVpnConnectionCheckMap)

	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-testaccVpnConnectionBaisc%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnConnectionConfigDependence)

	ResourceTest(t, resource.TestCase{
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
					//"vpn_gateway_id": "${alibabacloudstack_vpn_gateway.default.id}",
					"vpn_gateway_id":      "vpn-bs1fy78mp728vf4i2y2ar",
					"customer_gateway_id": "${alibabacloudstack_vpn_customer_gateway.default.id}",
					"local_subnet":        []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"name":                "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.0.0/24,172.16.1.0/24", "10.0.0.0/24,10.0.1.0/24"),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"local_subnet": []string{"172.16.1.0/24", "172.16.2.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.0.0.0/24,10.0.1.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"remote_subnet": []string{"10.4.0.0/24", "10.0.3.0/24"},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.4.0.0/24,10.0.3.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name": "${var.name}_change",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_change",
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.4.0.0/24,10.0.3.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"effect_immediately": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"effect_immediately": "true",
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.4.0.0/24,10.0.3.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ike_config": []map[string]string{
						{
							"ike_auth_alg":  "md5",
							"ike_enc_alg":   "des",
							"ike_version":   "ikev2",
							"ike_mode":      "main",
							"ike_lifetime":  "86400",
							"psk":           "tf-testvpn2",
							"ike_pfs":       "group1",
							"ike_remote_id": "testbob2",
							"ike_local_id":  "testalice2",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ike_config.#":               "1",
						"ike_config.0.ike_auth_alg":  "md5",
						"ike_config.0.ike_enc_alg":   "des",
						"ike_config.0.ike_version":   "ikev2",
						"ike_config.0.ike_mode":      "main",
						"ike_config.0.ike_lifetime":  "86400",
						"ike_config.0.psk":           "tf-testvpn2",
						"ike_config.0.ike_pfs":       "group1",
						"ike_config.0.ike_remote_id": "testbob2",
						"ike_config.0.ike_local_id":  "testalice2",
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.4.0.0/24,10.0.3.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"ipsec_config": []map[string]string{
						{
							"ipsec_pfs":      "group5",
							"ipsec_enc_alg":  "des",
							"ipsec_auth_alg": "md5",
							"ipsec_lifetime": "8640",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"ipsec_config.#":                "1",
						"ipsec_config.0.ipsec_pfs":      "group5",
						"ipsec_config.0.ipsec_enc_alg":  "des",
						"ipsec_config.0.ipsec_auth_alg": "md5",
						"ipsec_config.0.ipsec_lifetime": "8640",
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.1.0/24,172.16.2.0/24", "10.4.0.0/24,10.0.3.0/24"),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"name":               "${var.name}",
					"local_subnet":       []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":      []string{"10.0.0.0/24", "10.0.1.0/24"},
					"effect_immediately": REMOVEKEY,
					"ike_config":         REMOVEKEY,
					"ipsec_config":       REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":               name,
						"effect_immediately": "false",
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.0.0/24,172.16.1.0/24", "10.0.0.0/24,10.0.1.0/24"),
				),
			},
		},
	})

}

func TestAccAlibabacloudstackVpnConnectionMulti(t *testing.T) {
	var v vpc.DescribeVpnConnectionResponse

	resourceId := "alibabacloudstack_vpn_connection.default.1"
	ra := resourceAttrInit(resourceId, testAccVpnConnectionCheckMap)

	serviceFunc := func() interface{} {
		return &VpnGatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}
	rc := resourceCheckInit(resourceId, &v, serviceFunc)

	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(10000,20000)
	name := fmt.Sprintf("tf-testaccVpnConnectionMulti%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, resourceVpnConnectionConfigDependence)

	ResourceTest(t, resource.TestCase{
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
					"count": "2",
					//"vpn_gateway_id": "${alibabacloudstack_vpn_gateway.default.id}",
					"vpn_gateway_id":      "vpn-bs1fy78mp728vf4i2y2ar",
					"customer_gateway_id": "${alibabacloudstack_vpn_customer_gateway.default.id}",
					"local_subnet":        []string{"172.16.0.0/24", "172.16.1.0/24"},
					"remote_subnet":       []string{"10.0.0.0/24", "10.0.1.0/24"},
					"name":                "${var.name}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name,
					}),
					testAccCheckVpnConnectionAttr(&v,
						"172.16.0.0/24,172.16.1.0/24", "10.0.0.0/24,10.0.1.0/24"),
				),
			},
		},
	})

}

var testAccVpnConnectionCheckMap = map[string]string{
	"vpn_gateway_id":                CHECKSET,
	"customer_gateway_id":           CHECKSET,
	"local_subnet.#":                "2",
	"remote_subnet.#":               "2",
	"effect_immediately":            "false",
	"ike_config.#":                  "1",
	"ike_config.0.ike_auth_alg":     "md5",
	"ike_config.0.ike_enc_alg":      "aes",
	"ike_config.0.ike_version":      "ikev2",
	"ike_config.0.ike_mode":         "main",
	"ike_config.0.ike_lifetime":     "86400",
	"ike_config.0.psk":              CHECKSET,
	"ike_config.0.ike_pfs":          "group2",
	"ike_config.0.ike_remote_id":    CHECKSET,
	"ike_config.0.ike_local_id":     CHECKSET,
	"ipsec_config.#":                "1",
	"ipsec_config.0.ipsec_pfs":      "group2",
	"ipsec_config.0.ipsec_enc_alg":  "aes",
	"ipsec_config.0.ipsec_auth_alg": "md5",
	"ipsec_config.0.ipsec_lifetime": "86400",
}

var resourceVpnConnectionConfigDependence = func(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_vpc" "default" {
	cidr_block = "172.16.0.0/12"
	name = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation= "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}



resource "alibabacloudstack_vpn_customer_gateway" "default" {
	name = "${var.name}"
	ip_address = "42.104.22.210"
	description = "testAccVpnConnectionDesc"
}

`, name)
}
