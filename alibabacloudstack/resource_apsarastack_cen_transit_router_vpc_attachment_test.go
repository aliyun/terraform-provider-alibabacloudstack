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

					// "cen_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
					// "vpc_id": "${alibabacloudstack_cen_instance.default.transit_router_id}",
					// "zone_mappings": []map[string]string{
					// 	{"vswitch_id": "${alibabacloudstack_drds_database.default.0.drds_database_name}"},
					// },
					// "auto_create_vpc_route":           true,
					// "route_table_association_enabled": name,
					// "route_table_propagation_enabled": name,
					"cen_id": "cen-s1o6nsvvs9wc40p87v",
					"vpc_id": "vpc-j1p9canwb1cb7z9d6opol",
					"zone_mappings": []map[string]string{
						{"vswitch_id": "vsw-j1py5biwrpd81rkd6gl4b"},
					},
					"auto_create_vpc_route":           "true",
					"route_table_association_enabled": "true",
					"route_table_propagation_enabled": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": CHECKSET,

						"resource_type": CHECKSET,
						// "transit_router_route_table_type": CHECKSET,
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

			// 		"transit_router_route_table_name": modify_name,
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"transit_router_route_table_name": modify_name,
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

var AlibabacloudTestAccCenTransitRouterVpcAttachmentCheckmap = map[string]string{

	// "status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterVpcAttachmentBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
