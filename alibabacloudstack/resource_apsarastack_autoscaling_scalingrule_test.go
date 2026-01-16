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
	name := fmt.Sprintf("tf-testAccEssRule-%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccEssScalingRuleConfig)
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
				Config: testAccConfig(map[string]interface{}{
					"scaling_group_id":  "${alibabacloudstack_ess_scaling_group.default.id}",
					"adjustment_type":   "TotalCapacity",
					"adjustment_value":  "1",
					"cooldown":          10,
					"scaling_rule_name": name,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"scaling_group_id":  CHECKSET,
						"adjustment_type":   "TotalCapacity",
						"adjustment_value":  "1",
						"cooldown":          "10",
						"scaling_rule_name": name,
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
					"adjustment_type":   "PercentChangeInCapacity",
					"adjustment_value":  "2",
					"scaling_rule_name": name + "_update",
					"cooldown":          200,
				}),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"adjustment_type":   "PercentChangeInCapacity",
						"adjustment_value":  "2",
						"scaling_rule_name": name + "_update",
						"cooldown":          "200",
					}),
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

func testAccEssScalingRuleConfig(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	
	%s
	
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
		deployment_set_name = var.name
		description         = "example_value"
	}
	
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "${local.default_instance_type_id}"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
		system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
	}
	`, name, ECSInstanceCommonTestCase)
}
