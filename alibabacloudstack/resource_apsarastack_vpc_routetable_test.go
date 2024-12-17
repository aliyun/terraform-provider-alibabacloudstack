package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackVpcRoutetable0(t *testing.T) {
	var v map[string]interface{}

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

					"vpc_id": "${{ref(variable, VpcId)}}",

					"route_table_name": "${{ref(variable, RouteTableName)}}",

					"description": "${{ref(variable, Description)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "${{ref(variable, VpcId)}}",

						"route_table_name": "${{ref(variable, RouteTableName)}}",

						"description": "${{ref(variable, Description)}}",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"vpc_id": "${{ref(variable, VpcId)}}",

					"route_table_name": "${{ref(variable, RouteTableName)}}",

					"description": "${{ref(variable, DescriptionUpdate)}}",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"vpc_id": "${{ref(variable, VpcId)}}",

						"route_table_name": "${{ref(variable, RouteTableName)}}",

						"description": "${{ref(variable, DescriptionUpdate)}}",
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

	"status": CHECKSET,

	"description": CHECKSET,

	"route_table_id": CHECKSET,

	"resource_group_id": CHECKSET,

	"vswitch_ids": CHECKSET,

	"create_time": CHECKSET,

	"router_id": CHECKSET,

	"route_table_type": CHECKSET,

	"vpc_id": CHECKSET,

	"router_type": CHECKSET,

	"route_table_name": CHECKSET,

	"tags": CHECKSET,
}

func AlibabacloudTestAccVpcRoutetableBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}


variable "vpc_id" {
    default = vpc-uf61ozax1zxo9y4shipw4
}

variable "region_id" {
    default = cn-shanghai
}

variable "description" {
    default = Description
}

variable "route_table_name" {
    default = RouteTableName
}

variable "description_update" {
    default = DescriptionUpdate
}




`, name)
}
