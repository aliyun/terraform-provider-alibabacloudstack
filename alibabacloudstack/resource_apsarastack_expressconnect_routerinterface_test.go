package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectRouterinterface0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_expressconnect_routerinterface.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccExpressconnectRouterinterfaceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescriberouterinterfaceattributeRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpress_connectrouter_interface%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccExpressconnectRouterinterfaceBasicdependence)
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

					"opposite_region_id": "cn-hangzhou",

					"router_id": "vrt-bp17aybdbuvbl1m1yt7gb",

					"role": "InitiatingSide",

					"router_type": "VRouter",

					"payment_type": "PostPaid",

					"router_interface_name": "rdk-test",

					"spec": "Mini.2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test",

						"opposite_region_id": "cn-hangzhou",

						"router_id": "vrt-bp17aybdbuvbl1m1yt7gb",

						"role": "InitiatingSide",

						"router_type": "VRouter",

						"payment_type": "PostPaid",

						"router_interface_name": "rdk-test",

						"spec": "Mini.2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "test-update",

					"opposite_region_id": "cn-hangzhou",

					"router_id": "vrt-bp17aybdbuvbl1m1yt7gb",

					"role": "InitiatingSide",

					"router_type": "VRouter",

					"payment_type": "PostPaid",

					"router_interface_name": "rdk-test-name",

					"spec": "Mini.5",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "test-update",

						"opposite_region_id": "cn-hangzhou",

						"router_id": "vrt-bp17aybdbuvbl1m1yt7gb",

						"role": "InitiatingSide",

						"router_type": "VRouter",

						"payment_type": "PostPaid",

						"router_interface_name": "rdk-test-name",

						"spec": "Mini.5",
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

var AlibabacloudTestAccExpressconnectRouterinterfaceCheckmap = map[string]string{

	"opposite_interface_id": CHECKSET,

	"opposite_router_id": CHECKSET,

	"business_status": CHECKSET,

	"reservation_order_type": CHECKSET,

	"opposite_router_type": CHECKSET,

	"opposite_bandwidth": CHECKSET,

	"reservation_active_time": CHECKSET,

	"hc_threshold": CHECKSET,

	"reservation_bandwidth": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"opposite_interface_owner_id": CHECKSET,

	"opposite_region_id": CHECKSET,

	"health_check_source_ip": CHECKSET,

	"cross_border": CHECKSET,

	"reservation_internet_charge_type": CHECKSET,

	"role": CHECKSET,

	"opposite_vpc_instance_id": CHECKSET,

	"access_point_id": CHECKSET,

	"router_interface_name": CHECKSET,

	"health_check_target_ip": CHECKSET,

	"opposite_interface_status": CHECKSET,

	"description": CHECKSET,

	"end_time": CHECKSET,

	"router_id": CHECKSET,

	"bandwidth": CHECKSET,

	"connected_time": CHECKSET,

	"payment_type": CHECKSET,

	"create_time": CHECKSET,

	"has_reservation_data": CHECKSET,

	"hc_rate": CHECKSET,

	"opposite_interface_spec": CHECKSET,

	"router_type": CHECKSET,

	"opposite_interface_business_status": CHECKSET,

	"vpc_instance_id": CHECKSET,

	"opposite_access_point_id": CHECKSET,

	"spec": CHECKSET,

	"router_interface_id": CHECKSET,
}

func AlibabacloudTestAccExpressconnectRouterinterfaceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
