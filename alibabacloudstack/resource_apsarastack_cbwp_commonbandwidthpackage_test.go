package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackCbwpCommonbandwidthpackage0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_cbwp_commonbandwidthpackage.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccCbwpCommonbandwidthpackageCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribecommonbandwidthpackagesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%scbwpcommon_bandwidth_package%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccCbwpCommonbandwidthpackageBasicdependence)
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

					"description": "test",

					"resource_group_id": "rg-aek2xl5qajpkquq",

					"isp": "BGP",

					"bandwidth": "1000",

					"region_id": "cn-hangzhou",

					"internet_charge_type": "PayByBandwidth",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"resource_group_id": CHECKSET,

						"isp": "BGP",

						"bandwidth": "1000",

						"region_id": "cn-hangzhou",

						"internet_charge_type": "PayByBandwidth",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"region_id": "cn-hangzhou",

					"description": "test-update",

					"bandwidth": "1100",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"region_id": "cn-hangzhou",

						"description": "test-update",

						"bandwidth": "1100",
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

var AlibabacloudTestAccCbwpCommonbandwidthpackageCheckmap = map[string]string{

	"bandwidth_package_name": CHECKSET,

	"description": CHECKSET,

	"resource_group_id": CHECKSET,

	"business_status": CHECKSET,

	"reservation_order_type": CHECKSET,

	"bandwidth": CHECKSET,

	"expired_time": CHECKSET,

	"payment_type": CHECKSET,

	"public_ip_addresses": CHECKSET,

	"ratio": CHECKSET,

	"reservation_active_time": CHECKSET,

	"reservation_bandwidth": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"create_time": CHECKSET,

	"isp": CHECKSET,

	"has_reservation_data": CHECKSET,

	"deletion_protection": CHECKSET,

	"internet_charge_type": CHECKSET,

	"reservation_internet_charge_type": CHECKSET,

	"security_protection_types": CHECKSET,

	"region_id": CHECKSET,

	"common_bandwidth_package_id": CHECKSET,
}

func AlibabacloudTestAccCbwpCommonbandwidthpackageBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
