package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenTransitMulticastDomain0(t *testing.T) {
	var v *CbnDescribeTransitRouterMulticastDomainsResponse

	resourceId := "alibabacloudstack_cen_transit_router_multicast_domain.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenTransitMulticastDomainCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribeTransitRouterMuliticastDomainsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%smulticast_domain%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%smulticast_domain%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenTransitMulticastDomainBasicdependence)
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

					"transit_router_multicast_domain_description": name,
					"transit_router_multicast_domain_name":        name,
					"transit_router_id":                           "${alibabacloudstack_cen_instance.default.transit_router_id}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_multicast_domain_description": name,

						"transit_router_multicast_domain_name": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_multicast_domain_description": modify_name,

					"transit_router_multicast_domain_name": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_multicast_domain_description": modify_name,

						"transit_router_multicast_domain_name": modify_name,
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

var AlibabacloudTestAccCenTransitMulticastDomainCheckmap = map[string]string{

	"status": CHECKSET,

	"transit_router_id": CHECKSET,

	"transit_router_multicast_domain_id": CHECKSET,
}

func AlibabacloudTestAccCenTransitMulticastDomainBasicdependence(name string) string {
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
