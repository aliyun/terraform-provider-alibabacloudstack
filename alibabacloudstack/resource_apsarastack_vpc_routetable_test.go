package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcRoutetable0(t *testing.T) {
	var v vpc.RouterTableListType

	resourceId := "alibabacloudstack_vpc_routetable.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcRoutetableCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescriberoutetablelistRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcroute_table%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcRoutetableBasicdependence)
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

					"vpc_id": "${alibabacloudstack_vpc.default.id}",

					"name": name,

					"description": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": CHECKSET,

						"route_table_name": name,

						"description": name,
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

var AlibabacloudTestAccVpcRoutetableCheckmap = map[string]string{
	"test": NOSET,

	// "status": CHECKSET,

	// "description": CHECKSET,

	// "route_table_id": CHECKSET,

	// "resource_group_id": CHECKSET,

	// "vswitch_ids": CHECKSET,

	// "create_time": CHECKSET,

	// "router_id": CHECKSET,

	// "route_table_type": CHECKSET,

	// "vpc_id": CHECKSET,

	// "router_type": CHECKSET,

	// "route_table_name": CHECKSET,

	// "tags": CHECKSET,
}

func AlibabacloudTestAccVpcRoutetableBasicdependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	  }
	  
	  resource "alibabacloudstack_vpc" "default" {
		  cidr_block = "172.16.0.0/12"
		  name = "${var.name}"
	  }
`, name)
}
