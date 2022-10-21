package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackClientExpressConnectPhysicalConnection_domesic(t *testing.T) {
	//t.Skipf("There is an api bug that its describe response does not return CircuitCode. If the bug fixed, reopen this case")
	checkoutSupportedRegions(t, true, connectivity.VbrSupportRegions)
	var v map[string]interface{}
	resourceId := "alibabacloudstack_express_connect_physical_connection.default"
	ra := resourceAttrInit(resourceId, AlibabacloudStackClientExpressConnectPhysicalConnectionMap0)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DescribeExpressConnectPhysicalConnection")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()
	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sexpressconnectphysicalconnection%d", defaultRegionToTest, rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudStackClientExpressConnectPhysicalConnectionBasicDependence)
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
					//"access_point_id":          "${data.alibabacloudstack_express_connect_access_points.default.ids.0}",
					"access_point_id": "ap-cn-qingdao-env17-d01-amtest17",
					//"redundant_physical_connection_id": "${data.alibabacloudstack_express_connect_physical_connections.nameRegex.connections.0.id}",
					"device_name":              "CSW-VM-VPC-G1.AMTEST17",
					"type":                     "VPC",
					"peer_location":            "testacc12345",
					"physical_connection_name": name,
					"description":              "${var.name}",
					"line_operator":            "CO",
					"port_type":                "10GBase-LR",
					"bandwidth":                "5",
					"circuit_code":             "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"access_point_id": CHECKSET,
						//"redundant_physical_connection_id": CHECKSET,
						"type":                     "VPC",
						"peer_location":            "testacc12345",
						"physical_connection_name": name,
						"description":              name,
						"line_operator":            "CO",
						"port_type":                "10GBase-LR",
						"bandwidth":                "5",
						"circuit_code":             "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"physical_connection_name": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"physical_connection_name": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"bandwidth": "10",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"bandwidth": "10",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"circuit_code": "longtel001",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"circuit_code": "longtel001",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"description": name + "_Update",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"description": name + "_Update",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{
					"peer_location": "浙江省---vfjdbg_21e",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"peer_location": "浙江省---vfjdbg_21e",
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

var AlibabacloudStackClientExpressConnectPhysicalConnectionMap0 = map[string]string{}

func AlibabacloudStackClientExpressConnectPhysicalConnectionBasicDependence(name string) string {
	return fmt.Sprintf(` 
variable "name" {
  default = "%s"
}

`, name)
}
