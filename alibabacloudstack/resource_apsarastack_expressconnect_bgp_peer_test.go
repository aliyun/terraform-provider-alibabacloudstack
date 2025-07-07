package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectBgpPeer_basic0(t *testing.T) {
	var v *ExpressconnectBgpPeer

	resourceId := "alibabacloudstack_expressconnect_bgp_peer.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackExpressconnectBgpPeerCheckMap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &ExpressconnectService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribebgppeersRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testaccexpressconnect-bgp-peer%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackExpressconnectBgpPeerDependence0(rand))
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckWithEnvVariable(t, "ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID")
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"bgp_group_id":    "${alibabacloudstack_expressconnect_bgp_group.default.id}",
					"router_id":       "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
					"enable_bfd":      "true",
					"peer_ip_address": "192.168.0.1",
					"bfd_multi_hop":   "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_bfd":      "true",
						"peer_ip_address": "192.168.0.1",
						"bfd_multi_hop":   "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"enable_bfd":      "false",
					"peer_ip_address": "192.168.0.10",
					"bfd_multi_hop":   "12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"enable_bfd":      "false",
						"peer_ip_address": "192.168.0.10",
						"bfd_multi_hop":   "12",
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudStackExpressconnectBgpPeerCheckMap = map[string]string{
	"local_asn":   CHECKSET,
	"peer_asn":    CHECKSET,
	"route_limit": CHECKSET,
	"keepalive":   CHECKSET,
	"is_fake":     CHECKSET,
	"ip_version":  CHECKSET,
	"hold":        CHECKSET,
}

func AlibabacloudStackExpressconnectBgpPeerDependence0(vlanId int) func(string) string {
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
	description =                "TestAccAlibabacloudStackExpressconnectBgpPeer_basic0"
}

resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      65534
	peer_asn =       10
	router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	auth_key =       "%s"
}

`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId, getAccTestPassword(10))
	}
}
