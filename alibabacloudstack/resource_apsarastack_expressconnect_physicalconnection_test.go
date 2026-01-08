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
					"physical_connection_name": name,
					"description":              "abcabc",
					"line_operator":            "CO",
					"type":                     "VPC",
					"peer_location":            "XX Street",
					"access_point_id":          "${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.access_point_id}",
					"port_type":                "${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.port_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"description":              "abcabc",
						"line_operator":            "CO",
						"type":                     "VPC",
						"peer_location":            "XX Street",
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
					"peer_location": "sssssss",
					"circuit_code":  "longtel002",
					"description":   "eeeeee",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "sssssss",
						"circuit_code":  "longtel002",
						"description":   "eeeeee",
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
					"description":   "dddd",
					"line_operator": "CT",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description":   "dddd",
						"line_operator": "CT",
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
		},
	})
}

func TestAccAlibabacloudStackExpressconnectPhysicalconnection1(t *testing.T) {
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
					"physical_connection_name": name,
					"description":              "abcabc",
					"line_operator":            "CO",
					"type":                     "VPC",
					"peer_location":            "XX Street",
					"access_point_id":          "${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.access_point_id}",
					"port_type":                "${data.alibabacloudstack_expressconnect_physical_connections.anyone.connections.0.port_type}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name,
						"description":              "abcabc",
						"line_operator":            "CO",
						"type":                     "VPC",
						"peer_location":            "XX Street",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"status": "Canceled",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"status": "Canceled",
					}),
				),
			},
		},
	})
}


var AlibabacloudTestAccExpressconnectPhysicalconnectionCheckmap = map[string]string{
	"peer_location":            CHECKSET,
	"status":                   CHECKSET,
	"physical_connection_name": CHECKSET,
	"access_point_id":          CHECKSET,
	"port_type":                CHECKSET,
	"description":              CHECKSET,
	"line_operator":            CHECKSET,
	"bandwidth":                CHECKSET,
	"type":                     CHECKSET,
}

func AlibabacloudTestAccExpressconnectPhysicalconnectionBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

data "alibabacloudstack_expressconnect_physical_connections" "anyone" {
}


`, name)
}
