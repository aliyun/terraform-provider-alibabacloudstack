package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// This test case builds a whole chain in its dependency block: a CEN carrying a transit
// router and a CIDR block, a VBR, a VBR attachment, a connect attachment and the peer
// itself. Its leftovers are swept here with the prefixes this file owns, reusing the
// cleanup functions defined in the CEN and VBR test files: a failed run otherwise strands
// a transit router and a VBR that nothing else can reach.
func init() {
	resource.AddTestSweepers("alibabacloudstack_cen_transit_router_connect_peer", &resource.Sweeper{
		Name: "alibabacloudstack_cen_transit_router_connect_peer",
		F: func(region string) error {
			return SweeperFunc(
				sweeperEntry{fn: testSweepCenInstances, precise: []string{"tf-testaccceninstance48958"}, broad: broadPrefixes, fatal: true, waitSec: 10},
				sweeperEntry{fn: testSweepVirtualBorderRouters, precise: []string{"tf-testAccTransitRouterConnectPeer"}},
			)
		},
	})
}

func TestAccAlibabacloudStackCenTransitRouterConnectPeer_basic(t *testing.T) {
	var v map[string]interface{}
	resourceId := "alibabacloudstack_cen_transit_router_connect_peer.default"
	ra := resourceAttrInit(resourceId, nil)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCenTransitRouterConnectPeer")
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testAccTransitRouterConnectPeer%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, TransitRouterConnectPeerCommonTestCase)

	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfig(map[string]interface{}{
					"name":                  name,
					"local_ip":              "172.16.0.1",
					"peer_ip":               "172.16.0.11",
					"cen_id":                "${alibabacloudstack_cen_transit_router_connect_attachment.default.cen_id}",
					"connect_attachment_id": "${alibabacloudstack_cen_transit_router_connect_attachment.default.transit_router_attachment_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":                  name,
						"local_ip":              "172.16.0.1",
						"peer_ip":               "172.16.0.11",
						"connect_attachment_id": CHECKSET,
						"cen_id":                CHECKSET,
						"peer_id":               CHECKSET,
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

func TransitRouterConnectPeerCommonTestCase(randInt string) string {
	vlanId := getAccTestRandInt(1000, 2000)
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testAccTransitRouterConnectPeer%s"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
  transit_router_cidrs {
  	cidr = "172.16.0.0/16"
  }
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id = "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
}

resource "alibabacloudstack_cen_transit_router_vbr_attachment" "default" {
    vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
	cen_id = "${alibabacloudstack_cen_instance.default.id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}

resource "alibabacloudstack_cen_transit_router_connect_attachment" "default" { 
	cen_id = "${alibabacloudstack_cen_transit_router_vbr_attachment.default.cen_id}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
	transit_router_attachment_name = "${var.name}"
	depends_on = ["alibabacloudstack_cen_transit_router_vbr_attachment.default"]
}

`, randInt, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlanId)
}
