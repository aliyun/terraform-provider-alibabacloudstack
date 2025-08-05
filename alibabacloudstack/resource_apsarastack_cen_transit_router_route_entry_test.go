package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterEntry0(t *testing.T) {
	var v *CbnDescribeTransitRouterRouteEntriesResponse

	resourceId := "alibabacloudstack_cen_transit_router_route_entry.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterEtryCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterRouteEntriesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srouter_entry%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%srouter_entry_modify%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterEtryBasicdependence)
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

					"transit_router_route_entry_description": name,

					"transit_router_route_entry_destination_cidr_block": "10.10.10.1/32",

					"transit_router_route_entry_name":          name,
					"transit_router_route_entry_next_hop_type": "BlackHole",
					"transit_router_route_table_id":            "${alibabacloudstack_cen_transit_router_route_table.default.transit_router_route_table_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_route_entry_description": name,

						"transit_router_route_entry_destination_cidr_block": "10.10.10.1/32",
						"transit_router_route_entry_name":                   name,
						"transit_router_route_entry_next_hop_type":          "BlackHole",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_route_entry_description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_route_entry_description": modify_name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"transit_router_route_table_id"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitRouterEtryCheckmap = map[string]string{

	"status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterEtryBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_cen_instance" "default" {
	cen_instance_name = "${var.name}"
	description = "${var.name}"
	transit_router_name = "${var.name}"
	transit_router_description = "${var.name}"
}
resource "alibabacloudstack_cen_transit_router_route_table" "default" {
	transit_router_route_table_description = "${var.name}"
	transit_router_route_table_name = "${var.name}"
	transit_router_id = "${alibabacloudstack_cen_instance.default.transit_router_id}"
}
`, name)
}
