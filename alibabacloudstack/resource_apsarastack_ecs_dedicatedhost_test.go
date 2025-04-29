package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEcsDedicatedhost0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_ecs_dedicatedhost.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccEcsDedicatedhostCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EcsService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEcsDescribededicatedhostautorenewRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%secsdedicated_host%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccEcsDedicatedhostBasicdependence)
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

					"dedicated_host_type": "ddh.c5",

					"payment_type": "PostPaid",

					"region_id": "cn-hangzhou",

					"dedicated_host_name": "test-name",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"dedicated_host_type": "ddh.c5",

						"payment_type": "PostPaid",

						"region_id": "cn-hangzhou",

						"dedicated_host_name": "test-name",
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
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

var AlibabacloudTestAccEcsDedicatedhostCheckmap = map[string]string{

	"physical_gpus": CHECKSET,

	"machine_id": CHECKSET,

	"dedicated_host_id": CHECKSET,

	"description": CHECKSET,

	"resource_group_id": CHECKSET,

	"auto_renew": CHECKSET,

	"supported_custom_instance_type_families": CHECKSET,

	"network_attributes": CHECKSET,

	"capacity": CHECKSET,

	"dedicated_host_name": CHECKSET,

	"cpu_over_commit_ratio": CHECKSET,

	"expired_time": CHECKSET,

	"payment_type": CHECKSET,

	"sale_cycle": CHECKSET,

	"tags": CHECKSET,

	"status": CHECKSET,

	"zone_id": CHECKSET,

	"create_time": CHECKSET,

	"auto_placement": CHECKSET,

	"auto_renew_with_ecs": CHECKSET,

	"renewal_status": CHECKSET,

	"duration": CHECKSET,

	"dedicated_host_type": CHECKSET,

	"operation_locks": CHECKSET,

	"cores": CHECKSET,

	"sockets": CHECKSET,

	"gpu_spec": CHECKSET,

	"supported_instance_type_families": CHECKSET,

	"action_on_maintenance": CHECKSET,

	"region_id": CHECKSET,

	"supported_instance_types_list": CHECKSET,

	"dedicated_host_cluster_id": CHECKSET,

	"auto_release_time": CHECKSET,

	"period_unit": CHECKSET,
}

func AlibabacloudTestAccEcsDedicatedhostBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
