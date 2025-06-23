package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackNatGatewaysBandwidthPackage0(t *testing.T) {
	var v *VpcDescribebandwidthpackagesResponse

	resourceId := "alibabacloudstack_natgateway_bandwidth_package.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccBandwidthPackageCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &NatgatewayService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribebandwidthpackagesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%snat_gatewaybandwi%d", defaultRegionToTest, rand)
	modify_name := fmt.Sprintf("tf-testacc%snat_gatewaybandmodify%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccBandwidthPackageBasicdependence)
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

					"name": name,

					"bandwidth": "5",

					"natgateway_id": "${alibabacloudstack_nat_gateway.default.id}",
					"description":   name,
					"ip_count":      "2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ip_count": "2",

						"description": name,

						"name": name,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"name": modify_name,

					"bandwidth":   "10",
					"description": modify_name,
					"ip_count":    "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ip_count": "5",

						"description": modify_name,

						"name":      modify_name,
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"ip_count":  "1",
					"bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"ip_count":  "1",
						"bandwidth": "5",
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

var AlibabacloudTestAccBandwidthPackageCheckmap = map[string]string{

	// "status": CHECKSET,

	// "source_cidr": CHECKSET,

	// "snat_ip": CHECKSET,

	// "snat_table_id": CHECKSET,

	// "source_vswitch_id": CHECKSET,

	// "snat_entry_name": CHECKSET,

	// "snat_entry_id": CHECKSET,
}

func AlibabacloudTestAccBandwidthPackageBasicdependence(name string) string {
	return fmt.Sprintf(
		`
variable "name" {
	default = "%s"
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
	name = "${var.name}"
	cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
	vpc_id = "${alibabacloudstack_vpc.default.id}"
	cidr_block = "172.16.0.0/21"
	availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
	name = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
	vpc_id = "${alibabacloudstack_vswitch.default.vpc_id}"
	name = "${var.name}"
}
`, name)
}
