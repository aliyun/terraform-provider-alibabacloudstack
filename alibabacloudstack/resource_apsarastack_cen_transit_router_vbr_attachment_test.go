package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterVbrAttachment0(t *testing.T) {
	var v *CbnDescribeTransitRouterVbrAttachmentsResponse

	resourceId := "alibabacloudstack_cen_transit_router_vbr_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterVbrAttachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterVbrAttachmentsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srouter_vbr_attachment%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%srouter_vbr_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterVbrAttachmentBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		// CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				// d.SetId(fmt.Sprintf("%s:%s:%s:%s", "cen-xga10z1cg13l3vizpj", "tr-xxxvhy5cbe7oniom3o1nm", "tr-attach-qfikjqwq6fb0ld5ygi", "vbr-i2qmch4xlydzx670c739p"))
				Config: testAccConfig(map[string]interface{}{
					"vbr_id":            "${alibabacloudstack_express_connect_virtual_border_router.default.id}",
					"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
					"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_attachment_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_attachment_name":        modify_name,
					"transit_router_attachment_description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_attachment_name":        modify_name,
						"transit_router_attachment_description": modify_name,
					}),
				),
			},

			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitRouterVbrAttachmentCheckmap = map[string]string{

	"status": CHECKSET,

	"creation_time": CHECKSET,

	"resource_type": "VBR",

	"auto_publish_route_enabled": CHECKSET,

	"vbr_owner_id": CHECKSET,

	"vbr_id":                       CHECKSET,
	"transit_router_attachment_id": CHECKSET,
	"transit_router_id":            CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterVbrAttachmentBasicdependence(name string) string {
	vlan_id := getAccTestRandInt(1000, 2000)
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id = "%s"
	vlan_id =                    %d
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	enable_ipv6              = true
	local_ipv6_gateway_ip = "2408:4004:cc:400::1"
	peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
	peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
}


resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}
`, name, getAccTestOsEnv("ALIBABACLOUDSTACK_PHYSICAL_CONNECTION_ID"), vlan_id)
}
