package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackVpcVpc0(t *testing.T) {

	var v map[string]interface{}

	// TODO Describe method，v 的类型

	resourceId := "alibabacloudstack_vpc_vpc.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccVpcVpcCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &VpcService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoVpcDescribevpcattributesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := acctest.RandIntRange(10000, 99999)
	name := fmt.Sprintf("tf-testacc%svpcvpc%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccVpcVpcBasicdependence)
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

					"description": "RDK更新",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDK更新",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"description": "RDK更新1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"description": "RDK更新1",
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

			{
				Config: testAccConfig(map[string]interface{}{

					"cidr_block": "172.16.0.0/12",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"cidr_block": "172.16.0.0/12",
					}),
				),
			},
			// {
			// 	Config: testAccConfig(map[string]interface{}{

			// 		"cidr_block": "192.168.0.0/16",
			// 	}),
			// 	Check: resource.ComposeTestCheckFunc(
			// 		testAccCheck(map[string]string{

			// 			"cidr_block": "192.168.0.0/16",
			// 		}),
			// 	),
			// },

			{
				Config: testAccConfig(map[string]interface{}{

					"name": name + "update1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name": name + "update1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"name": name + "update2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"name": name + "update2",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"enable_ipv6": "true",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"enable_ipv6": "true",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"enable_ipv6": "false",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"enable_ipv6": "false",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_name": name + "update_vpc_name1",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_name": name + "update_vpc_name1",
					}),
				),
			},
			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_name": name + "update_vpc_name2",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_name": name + "update_vpc_name2",
					}),
				),
			},
		},
	})
}

var AlibabacloudTestAccVpcVpcCheckmap = map[string]string{
	//  TODO  checkmap 和 case的资源对齐

	// "name": CHECKSET,

	// "vpc_name": CHECKSET,

	// "route_table_id":    CHECKSET,
	// "resource_group_id": CHECKSET,
	// // "router_table_id":   CHECKSET,
	// "router_id":       CHECKSET,
	// "ipv6_cidr_block": CHECKSET,
	// "status":          CHECKSET,
}

func AlibabacloudTestAccVpcVpcBasicdependence(name string) string {

	//  TODO  检查依赖变量

	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}

`, name)
}
