package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcVswitch0(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",

					"vswitch_name": name,

					"zone_id": "${data.alibabacloudstack_zones.default.zones.0.id}",

					"vpc_id": "${alibabacloudstack_vpc.default.id}",

					"cidr_block": "172.16.0.0/24",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",

						"vswitch_name": name,
					}),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF",
						"For":     "Test",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF",
						"tags.For":     "Test",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": map[string]string{
						"Created": "TF-update",
						"For":     "Test-update",
					},
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "2",
						"tags.Created": "TF-update",
						"tags.For":     "Test-update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"tags": REMOVEKEY,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"tags.%":       "0",
						"tags.Created": REMOVEKEY,
						"tags.For":     REMOVEKEY,
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcVswitchCheckmap = map[string]string{}

func AlibabacloudTestAccVpcVswitchBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s
%s
`, name, DataZoneCommonTestCase, VpcCommonTestCase)
}
func TestAccAlibabacloudStackVpcVswitch1(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse
	AlibabacloudTestAccVpcVswitchCheckmap["ipv6_cidr_block"] = CHECKSET
	resourceId := "alibabacloudstack_vpc_vswitch.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchIPV6dependence)
	resource.Test(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{
					"description": "modify_description",

					"vswitch_name": name,

					"zone_id": "${data.alibabacloudstack_zones.default.zones.0.id}",

					"vpc_id": "${alibabacloudstack_vpc.vpc.id}",

					"cidr_block": "172.16.0.0/24",

					"enable_ipv6": true,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name":            name,
						"description":     "modify_description",
						"ipv6_cidr_block": CHECKSET,
					}),
				),
			},
		},
	})
}

func AlibabacloudTestAccVpcVswitchIPV6dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

%s

resource "alibabacloudstack_vpc" "vpc" {
  vpc_name     = "${var.name}_vpc"
  enable_ipv6  = true
}

`, name, DataZoneCommonTestCase)
}
