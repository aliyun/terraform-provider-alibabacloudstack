package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitMulticastDomainAssociation0(t *testing.T) {
	var v *CbnDescribeTransitRouterMulticastDomainAssociationsResponse

	resourceId := "alibabacloudstack_cen_transit_router_multicast_domain_association.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitMulticastDomainAssociationCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterMuliticastDomainAssociationsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smulticast_domain_association%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitMulticastDomainAssociationBasicdependence)
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
					"transit_router_attachment_id":       "${alibabacloudstack_cen_transit_router_vpc_attachment.default.transit_router_attachment_id}",
					"transit_router_multicast_domain_id": "${alibabacloudstack_cen_transit_router_multicast_domain.default.transit_router_multicast_domain_id}",
					"vswitch_id":                         "${alibabacloudstack_vpc_vswitch.default.id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_attachment_id": CHECKSET,

						"transit_router_multicast_domain_id": CHECKSET,
						"vswitch_id":                         CHECKSET,
					}),
				),
			},

			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
				// ImportStateVerifyIgnore: []string{"transit_router_id"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitMulticastDomainAssociationCheckmap = map[string]string{

	"status": CHECKSET,
}

// func AlibabacloudTestAccCenTransitMulticastDomainAssociationBasicdependence(name string) string {
// 	return fmt.Sprintf(
// 		`
// 		variable "name" {
// 			  default = "%s"
// 			}
// 			`, name)
// }

func AlibabacloudTestAccCenTransitMulticastDomainAssociationBasicdependence(name string) string {
	return fmt.Sprintf(
		`
		variable "name" {
			  default = "%s"
			}

			%s

			resource "alibabacloudstack_cen_instance" "default" {
				cen_instance_name = "${var.name}"
				description = "${var.name}"
				transit_router_name = "${var.name}"
				transit_router_description = "${var.name}"
			}

			resource "alibabacloudstack_cen_transit_router_vpc_attachment" "default" {
				transit_router_attachment_name = "${var.name}"
				transit_router_attachment_description = "${var.name}"
				cen_id = "${alibabacloudstack_cen_instance.default.id}"
				vpc_id = "${alibabacloudstack_vpc_vswitch.default.vpc_id}"
				transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
				zone_mappings {
						vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
						 zone_id = "${alibabacloudstack_vpc_vswitch.default.zone_id}"
					}
			}

			resource "alibabacloudstack_cen_transit_router_multicast_domain" "default" {
				transit_router_multicast_domain_description = "${var.name}"
				transit_router_multicast_domain_name = "${var.name}"
				transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
			}
			`, name, VSwitchCommonTestCase)
}
