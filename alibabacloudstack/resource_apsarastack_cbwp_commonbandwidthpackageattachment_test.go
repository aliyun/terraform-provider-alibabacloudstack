package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCbwpCommonbandwidthpackageattachment0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cbwp_commonbandwidthpackageattachment.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCbwpCommonbandwidthpackageattachmentCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeCommonBandwidthPackageAttachment")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scbwpcommon_bandwidth_package_attachment%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCbwpCommonbandwidthpackageattachmentBasicdependence)
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {

			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,

		CheckDestroy: rac.checkResourceDestroy(),

		Steps: []resource.TestStep{

			{
				Config: testAccConfig(map[string]interface{}{

					"instance_id": "eip-bp1y3ih6jsour25vkjs1v",

					"bandwidth_package_id": "cbwp-bp1pz25h14sqddzke74ba",

					"region_id": "cn-hangzhou",

					"bandwidth_package_bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_id": "eip-bp1y3ih6jsour25vkjs1v",

						"bandwidth_package_id": "cbwp-bp1pz25h14sqddzke74ba",

						"region_id": "cn-hangzhou",

						"bandwidth_package_bandwidth": "5",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"bandwidth_package_bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"bandwidth_package_bandwidth": "5",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"instance_id": "eip-bp1y3ih6jsour25vkjs1v",

					"bandwidth_package_id": "cbwp-bp1pz25h14sqddzke74ba",

					"region_id": "cn-hangzhou",

					"bandwidth_package_bandwidth": "5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"instance_id": "eip-bp1y3ih6jsour25vkjs1v",

						"bandwidth_package_id": "cbwp-bp1pz25h14sqddzke74ba",

						"region_id": "cn-hangzhou",

						"bandwidth_package_bandwidth": "5",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccCbwpCommonbandwidthpackageattachmentCheckmap = map[string]string{

	"status": CHECKSET,

	"instance_id": CHECKSET,

	"bandwidth_package_id": CHECKSET,

	"region_id": CHECKSET,

	"bandwidth_package_bandwidth": CHECKSET,
}

func AlibabacloudTestAccCbwpCommonbandwidthpackageattachmentBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
