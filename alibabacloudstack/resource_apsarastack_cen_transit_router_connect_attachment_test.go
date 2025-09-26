package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterConnectAttachment0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cen_transit_router_connect_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterConnectAttachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCenTransitRouterConnectAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(1000, 2000)
	name := fmt.Sprintf("tf-testacc-router_connect_attachment%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterConnectAttachmentBasicdependence)
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
					"cen_id":                         "${alibabacloudstack_cen_instance.default.id}",
					"transit_router_id":              "${alibabacloudstack_cen_instance.default.transit_router_id}",
					"transit_router_attachment_name": "${var.name}",
					"depends_on":                     []string{`alibabacloudstack_cen_transit_router_vbr_attachment.default`},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_attachment_name": name,
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

var AlibabacloudTestAccCenTransitRouterConnectAttachmentCheckmap = map[string]string{

	"status": CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterConnectAttachmentBasicdependence(name string) string {
	rand := getAccTestRandInt(1000, 2000)
	return fmt.Sprintf(
		`
variable "name" {
  default = "%v"
}

resource "alibabacloudstack_cen_instance" "default" {
  description = "tf-testaccceninstance48958"
  cen_instance_name = "tf-testaccceninstance48958"
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

`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), rand)
}
