package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackAutoscalingScalinggroup0(t *testing.T) {
	var v map[string]interface{}

	resourceId := "alibabacloudstack_autoscaling_scalinggroup.default"
	ra := resourceAttrInit(resourceId, AlibabacloudTestAccAutoscalingScalinggroupCheckmap)
	rc := resourceCheckInitWithDescribeMethod(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	}, "DoEssDescribescalinggroupsRequest")
	rac := resourceAttrCheckInit(rc, ra)
	testAccCheck := rac.resourceAttrMapUpdateSet()

	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testacc%sauto_scalingscaling_group%d", defaultRegionToTest, rand)

	testAccConfig := resourceTestAccConfigFunc(resourceId, name, AlibabacloudTestAccAutoscalingScalinggroupBasicdependence)
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

					"multi_az_policy": "PRIORITY",

					"scaling_policy": "release",

					"vswitch_id": "vsw-bp1jtqeiaavrfq53sljdz",

					"launch_template_id": "lt-bp184s52s0vz4r9u3csy",

					"scaling_group_name": "ScalingGroupNameTest11",

					"group_type": "ECS",

					"launch_template_version": "Default",

					"region_id": "cn-hangzhou",

					"health_check_type": "ECS",
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{

						"multi_az_policy": "PRIORITY",

						"scaling_policy": "release",

						"vswitch_id": "vsw-bp1jtqeiaavrfq53sljdz",

						"launch_template_id": "lt-bp184s52s0vz4r9u3csy",

						"scaling_group_name": "ScalingGroupNameTest11",

						"group_type": "ECS",

						"launch_template_version": "Default",

						"region_id": "cn-hangzhou",

						"health_check_type": "ECS",
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

var AlibabacloudTestAccAutoscalingScalinggroupCheckmap = map[string]string{

	"spot_instance_remedy": CHECKSET,

	"server_group": CHECKSET,

	"resource_group_id": CHECKSET,

	"active_scaling_configuration_id": CHECKSET,

	"vserver_groups": CHECKSET,

	"desired_capacity": CHECKSET,

	"on_demand_base_capacity": CHECKSET,

	"removal_policies": CHECKSET,

	"launch_template_overrides": CHECKSET,

	"tags": CHECKSET,

	"multi_az_policy": CHECKSET,

	"status": CHECKSET,

	"suspended_processes": CHECKSET,

	"removing_capacity": CHECKSET,

	"vswitch_ids": CHECKSET,

	"pending_capacity": CHECKSET,

	"scaling_group_id": CHECKSET,

	"vswitch_id": CHECKSET,

	"load_balancer_ids": CHECKSET,

	"spot_instance_pools": CHECKSET,

	"launch_template_id": CHECKSET,

	"custom_policy_arn": CHECKSET,

	"scaling_group_name": CHECKSET,

	"default_cooldown": CHECKSET,

	"group_type": CHECKSET,

	"vpc_id": CHECKSET,

	"launch_template_version": CHECKSET,

	"stopped_capacity": CHECKSET,

	"health_check_type": CHECKSET,

	"compensate_with_on_demand": CHECKSET,

	"on_demand_percentage_above_base_capacity": CHECKSET,

	"modification_time": CHECKSET,

	"total_instance_count": CHECKSET,

	"allocation_strategy": CHECKSET,

	"init_capacity": CHECKSET,

	"pending_wait_capacity": CHECKSET,

	"total_capacity": CHECKSET,

	"removing_wait_capacity": CHECKSET,

	"spot_allocation_strategy": CHECKSET,

	"protected_capacity": CHECKSET,

	"standby_capacity": CHECKSET,

	"scaling_policy": CHECKSET,

	"create_time": CHECKSET,

	"group_deletion_protection": CHECKSET,

	"max_size": CHECKSET,

	"active_capacity": CHECKSET,

	"min_size": CHECKSET,

	"alb_server_group": CHECKSET,

	"az_balance": CHECKSET,

	"system_suspended": CHECKSET,

	"monitor_group_id": CHECKSET,

	"load_balancer_config": CHECKSET,

	"region_id": CHECKSET,

	"max_instance_lifetime": CHECKSET,
}

func AlibabacloudTestAccAutoscalingScalinggroupBasicdependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
}



`, name)
}
