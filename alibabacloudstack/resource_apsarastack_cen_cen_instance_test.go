package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCenInstance0(t *testing.T) {
	var v *CenInstance

	resourceId := "alibabacloudstack_cen_instance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCenInstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &CenService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoCbnDescribecensRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testaccceninstance%d", rand)
	modify_name := fmt.Sprintf("tf-testaccceninstancemodify%d", rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCenInstanceBasicdependence)
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
					"cen_instance_name":          name,
					"description":                name,
					"transit_router_name":        name,
					"transit_router_description": name,
					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.0.0/16"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_id":            CHECKSET,
						"transit_router_id": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"cen_instance_name": modify_name,
					"description":       modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"cen_instance_name": modify_name,
						"description":       modify_name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_name":        modify_name,
					"transit_router_description": modify_name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_name":        modify_name,
						"transit_router_description": modify_name,
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.10.2/24"},
						{"cidr": "10.10.11.2/24"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_cidrs.#": "2",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"transit_router_cidrs": []map[string]interface{}{
						{"cidr": "10.10.10.2/24"},
						{"cidr": "10.10.12.2/24"},
						{"cidr": "10.10.13.2/24"},
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"transit_router_cidrs.#": "3",
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

var AlibabacloudTestAccCenInstanceCheckmap = map[string]string{

	"status": CHECKSET,

	"protection_level": CHECKSET,

	"create_time": CHECKSET,

	"cen_id": CHECKSET,
}

func AlibabacloudTestAccCenInstanceBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}
`, name)
}
