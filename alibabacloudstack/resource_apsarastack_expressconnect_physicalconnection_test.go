package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackExpressconnectPhysicalconnection0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_expressconnect_physicalconnection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccExpressconnectPhysicalconnectionCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribephysicalconnectionsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpress_connectphysical_connection%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccExpressconnectPhysicalconnectionBasicdependence)
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

					"description": "abcabc",

					"line_operator": "CO",

					"type": "VPC",

					"peer_location": "XX街道",

					"access_point_id": "ap-cn-hangzhou-jg-B",

					"port_type": "1000Base-LX",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "abcabc",

						"line_operator": "CO",

						"type": "VPC",

						"peer_location": "XX街道",

						"access_point_id": "ap-cn-hangzhou-jg-B",

						"port_type": "1000Base-LX",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"peer_location": "sssssss",

					"circuit_code": "longtel002",

					"description": "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"peer_location": "sssssss",

						"circuit_code": "longtel002",

						"description": "eeeeee",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"description": "dddd",

					"line_operator": "CT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "dddd",

						"line_operator": "CT",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"status": "Confirmed",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": "Confirmed",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"status": "Enabled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": "Enabled",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"status": "Terminated",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"status": "Terminated",
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

var AlibabacloudTestAccExpressconnectPhysicalconnectionCheckmap = map[string]string{

	"peer_location": CHECKSET,

	"redundant_physical_connection_id": CHECKSET,

	"business_status": CHECKSET,

	"ad_location": CHECKSET,

	"reservation_active_time": CHECKSET,

	"tags": CHECKSET,

	"vlan_id": CHECKSET,

	"status": CHECKSET,

	"circuit_code": CHECKSET,

	"physical_connection_name": CHECKSET,

	"reservation_internet_charge_type": CHECKSET,

	"access_point_id": CHECKSET,

	"port_number": CHECKSET,

	"port_type": CHECKSET,

	"description": CHECKSET,

	"end_time": CHECKSET,

	"line_operator": CHECKSET,

	"physical_connection_id": CHECKSET,

	"loa_status": CHECKSET,

	"bandwidth": CHECKSET,

	"payment_type": CHECKSET,

	"create_time": CHECKSET,

	"has_reservation_data": CHECKSET,

	"type": CHECKSET,

	"enabled_time": CHECKSET,

	"spec": CHECKSET,
}

func AlibabacloudTestAccExpressconnectPhysicalconnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
