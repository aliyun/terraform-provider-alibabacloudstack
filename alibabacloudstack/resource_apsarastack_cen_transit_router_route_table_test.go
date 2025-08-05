package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitRouterTable0(t *testing.T) {
	var v *CbnDescribeTransitRouterRouteTablesResponse

	resourceId := "alibabacloudstack_cen_transit_router_route_table.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitRouterTableCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterRouteTablesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%srouter_table%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%srouter_table%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitRouterTableBasicdependence)
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

					"transit_router_route_table_description": name,
					"transit_router_route_table_name":        name,
					"transit_router_id":                      "${alibabacloudstack_cen_instance.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_route_table_description": name,

						"transit_router_route_table_name": name,
						"transit_router_route_table_type": "Custom",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_route_table_description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_route_table_description": modify_name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_route_table_name": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_route_table_name": modify_name,
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"transit_router_id"},
			},
		},
	})
}

var AlibabacloudTestAccCenTransitRouterTableCheckmap = map[string]string{

	// "status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccCenTransitRouterTableBasicdependence(name string) string {
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
`, name)
}
