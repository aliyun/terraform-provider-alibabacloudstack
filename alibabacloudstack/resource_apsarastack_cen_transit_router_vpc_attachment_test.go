package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterVpcAttachment0(t *testing.T) {
	var v *CbnDescribeTransitRouterVpcAttachmentsResponse

	resourceId := "alibabacloudstack_cen_transit_router_vpc_attachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterVpcAttachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterVpcAttachmentsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srouter_vpc_attachment%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%srouter_vpc_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterVpcAttachmentBasicdependence)
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
					"cen_id":            "${alibabacloudstack_cen_instance.default.id}",
					"vpc_id":            "${alibabacloudstack_vpc_vswitch.default.vpc_id}",
					"transit_router_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
					"zone_mappings": []map[string]string{
						{
							"vswitch_id": "${alibabacloudstack_vpc_vswitch.default.id}",
							"zone_id":    "${alibabacloudstack_vpc_vswitch.default.zone_id}",
						},
					},
					"auto_create_vpc_route":           "true",
					"route_table_association_enabled": "true",
					"route_table_propagation_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_mappings.0.vswitch_id": CHECKSET,
						"zone_mappings.#":            "1",
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
				Config: testAccConfig(map[string]interface{}{

					"zone_mappings": []map[string]string{
						{
							"vswitch_id": "${alibabacloudstack_vpc_vswitch.vswitchv2.id}",
							"zone_id":    "${alibabacloudstack_vpc_vswitch.vswitchv2.zone_id}",
						},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_mappings.0.vswitch_id": CHECKSET,
						"zone_mappings.#":            "1",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"cen_id", "auto_create_vpc_route", "route_table_association_enabled", "route_table_propagation_enabled"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitRouterVpcAttachmentCheckmap = map[string]string{

	"status": CHECKSET,

	"creation_time": CHECKSET,

	"charge_type": CHECKSET,

	"resource_type": "VPC",

	"auto_publish_route_enabled": CHECKSET,

	"vpc_owner_id": CHECKSET,

	"vpc_id":                       CHECKSET,
	"transit_router_attachment_id": CHECKSET,
	"transit_router_id":            CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterVpcAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

%s


resource "alibabacloudstack_vpc_vswitch" "vswitchv2" {
	name = "${var.name}v2"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	cidr_block = "172.16.0.0/24"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  }

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}

`, name, VSwitchCommonTestCase)
}
