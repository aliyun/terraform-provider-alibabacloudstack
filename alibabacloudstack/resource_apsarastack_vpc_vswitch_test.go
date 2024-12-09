package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcVswitch0(t *testing.T) {

	var v vpc.DescribeVSwitchAttributesResponse

	resourceId := "alibabacloudstack_vpc_vswitch.test_default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVswitchCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevswitchattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchBasicdependence)
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
					"vswitch_name": "${var.name}",

					"zone_id": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${alibabacloudstack_vpc.default.id}",

					"cidr_block": "${alibabacloudstack_vpc.default.cidr_block}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vswitch_name": name,

						"zone_id": "cn-wulan-env212-amtest212001-a",

						"cidr_block": "172.16.0.0/24",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",

					"vswitch_name": "modify_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",

						"vswitch_name": "modify_name",
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
`, name, VpcCommonTestCase)
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

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvswitch%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVswitchIpv6dependence)
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
					"name": "${var.name}_ipv6",

					"availability_zone": "cn-wulan-env212-amtest212001-a",

					"vpc_id": "${alibabacloudstack_vpc.default.id}",

					"cidr_block": "10.1.0.0/21",

					"ipv6_cidr_block": "${alibabacloudstack_vpc.vpc.ipv6_cidr_block}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"name": name + "_ipv6",

						"ipv6_cidr_block": CHECKSET,
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "modify_description",

					"vswitch_id": "alibabacloudstack_vswitch.default.id",

					"vswitch_name": "modify_name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "modify_description",

						"vswitch_id": "alibabacloudstack_vswitch.default.id",

						"vswitch_name": "modify_name",
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

func AlibabacloudTestAccVpcVswitchIpv6dependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

resource "alibabacloudstack_vpc" "vpc" {
  vpc_name       = "${var.name}"
  cidr_block = "${var.cidr_block}"
  enable_ipv6    = true
}

`, name)
}
