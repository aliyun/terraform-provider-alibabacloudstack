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
					"cen_id":            "cen-y5hsaev9m0fyefocz5",
					"vbr_id":            "vbr-i2qmch4xlydzx670c739p",
					"transit_router_id": "tr-xxxvhy5cbe7oniom3o1nm",
					// "cen_id":            "${alibabacloudstack_cen_instance.default.id}",
					// "vbr_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
					// "transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
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
			// {
			// 	Config: testAccConfig(map[string]interface{}{

			// 		"zone_mappings": []map[string]string{
			// 			{
			// 				"vswitch_id": "${alibabacloudstack_vpc_vswitch.vswitchv2.id}",
			// 				"zone_id":    "${alibabacloudstack_vpc_vswitch.vswitchv2.zone_id}",
			// 			},
			// 		},
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"zone_mappings.0.vswitch_id": CHECKSET,
			// 			"zone_mappings.#":            "1",
			// 		}),
			// 	),
			// },
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

	// "charge_type": CHECKSET,

	"resource_type": "VBR",

	"auto_publish_route_enabled": CHECKSET,

	"vbr_owner_id": CHECKSET,

	"vbr_id":                       CHECKSET,
	"transit_router_attachment_id": CHECKSET,
	"transit_router_id":            CHECKSET,
}

// func AlibabacloudTestAccCenTransitRouterVbrAttachmentBasicdependence(name string) string {
// 	return fmt.Sprintf(
// 		`
// variable "name" {
// 	default = "%s"
// }

// %s

// resource "alibabacloudstack_vpc_vswitch" "vswitchv2" {
// 	name = "${var.name}v2"
// 	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
// 	cidr_block = "172.16.0.0/24"
// 	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
//   }

// resource "alibabacloudstack_cen_instance" "default" {
// 	cen_instance_name = "${var.name}"
// 	description = "${var.name}"
// 	transit_router_name = "${var.name}"
// 	transit_router_description = "${var.name}"
// }

// `, name, VSwitchCommonTestCase)
// }

func AlibabacloudTestAccCenTransitRouterVbrAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
