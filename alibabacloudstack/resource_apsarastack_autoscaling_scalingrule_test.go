package alibabacloudstack

import (
	"fmt"

	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEssScalingRule_basic(t *testing.T) {
	var v ess.ScalingRule
	rand := getAccTestRandInt(1000, 999999)
	resourceId := "alibabacloudstack_ess_scaling_rule.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
		"cooldown":         "0",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfig(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
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

func TestAccAlibabacloudStackEssScalingRuleMulti(t *testing.T) {
	var v ess.ScalingRule
	rand := getAccTestRandInt(1000, 999999)
	resourceId := "alibabacloudstack_ess_scaling_rule.default.9"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"adjustment_type":  "TotalCapacity",
		"adjustment_value": "1",
	}
	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: testAccCheckEssScalingRuleDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingRuleConfigMulti(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})
}

func testAccCheckEssScalingRuleDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_scaling_rule" {
			continue
		}
		_, err := essService.DescribeEssScalingRule(rs.Primary.ID)

		// Verify the error is what we want
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("Scaling rule %s still exists.", rs.Primary.ID)
	}

	return nil
}

func testAccEssScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = "1"
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentType(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 1
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateAdjustmentValue(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateScalingRuleName(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = "${var.name}"
		cooldown = 0
	}
	`, common, rand)
}

func testAccEssScalingRuleUpdateCooldown(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "PercentChangeInCapacity"
		adjustment_value = 2
		scaling_rule_name = "${var.name}"
		cooldown = 200
	}
	`, common, rand)
}

func testAccEssScalingRuleConfigMulti(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssScalingRule-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 2
		default_cooldown = 20
		removal_policies = ["OldestInstance", "NewestInstance"]
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	
	resource "alibabacloudstack_ecs_deployment_set" "default" {
		strategy            = "Availability"
		domain              = "Default"
		granularity         = "Host"
		deployment_set_name = "example_value"
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		count = 10
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1
	}
	`, common, rand)
}

func testAccEssTargetTrackingScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssTargetTrackingScalingRuleConfig-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type  = "TotalCapacity"
  		adjustment_value = 1
	}​​​​
	`, common, rand)
}

func testAccEssStepScalingRuleConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssStepScalingRuleConfig-%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	resource "alibabacloudstack_ess_scaling_rule" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		adjustment_type = "TotalCapacity"
		adjustment_value = 1​​​
	}
	`, common, rand)
}
