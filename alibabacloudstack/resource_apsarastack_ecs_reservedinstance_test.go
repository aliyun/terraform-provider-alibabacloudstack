package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsReservedinstance0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_reservedinstance.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsReservedinstanceCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribereservedinstancesRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsreserved_instance%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsReservedinstanceBasicdependence)
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

					"zone_id": "cn-hangzhou-i",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::ResourceGroup.ResourceGroupId)}}",

					"reserved_instance_name": "fvt-ecs-reserved-fb71e8d3",

					"region_id": "cn-hangzhou",

					"instance_type": "ecs.t6-c4m1.large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "cn-hangzhou-i",

						"resource_group_id": CHECKSET,

						"reserved_instance_name": "fvt-ecs-reserved-fb71e8d3",

						"region_id": "cn-hangzhou",

						"instance_type": "ecs.t6-c4m1.large",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::ResourceGroup.ResourceGroupId)}}",

					"region_id": "cn-hangzhou",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"resource_group_id": CHECKSET,

						"region_id": "cn-hangzhou",
					}),
				),
			},

			{
				Config: testAccConfig(map[string]interface{}{

					"zone_id": "cn-hangzhou-i",

					"resource_group_id": "${{ref(resource, ResourceManager::ResourceGroup::3.0.0::ResourceGroup.ResourceGroupId)}}",

					"reserved_instance_name": "fvt-ecs-reserved-fb71e8d3",

					"region_id": "cn-hangzhou",

					"instance_type": "ecs.t6-c4m1.large",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"zone_id": "cn-hangzhou-i",

						"resource_group_id": CHECKSET,

						"reserved_instance_name": "fvt-ecs-reserved-fb71e8d3",

						"region_id": "cn-hangzhou",

						"instance_type": "ecs.t6-c4m1.large",
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

var AlibabacloudTestAccEcsReservedinstanceCheckmap = map[string]string{

	"description": CHECKSET,

	"platform": CHECKSET,

	"resource_group_id": CHECKSET,

	"instance_amount": CHECKSET,

	"expired_time": CHECKSET,

	"reserved_instance_name": CHECKSET,

	"instance_type": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"allocation_status": CHECKSET,

	"zone_id": CHECKSET,

	"create_time": CHECKSET,

	"start_time": CHECKSET,

	"operation_locks": CHECKSET,

	"offering_type": CHECKSET,

	"scope": CHECKSET,

	"reserved_instance_id": CHECKSET,

	"region_id": CHECKSET,
}

func AlibabacloudTestAccEcsReservedinstanceBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
