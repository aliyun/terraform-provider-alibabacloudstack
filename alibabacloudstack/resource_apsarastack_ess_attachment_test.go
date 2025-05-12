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

func TestAccalibabacloudstackdEssAttachment_update(t *testing.T) {
	rand := getAccTestRandInt(10000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_attachment.default"
	basicMap := map[string]string{
		"instance_ids.#":   "2",
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAttachmentConfigInstance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alibabacloudstack_ess_attachment.default", &v),
					testAccCheck(nil),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccEssAttachmentConfigInstance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alibabacloudstack_ess_attachment.default", &v),
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccEssAttachmentConfigRemoveInstance(rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alibabacloudstack_ess_attachment.default", &v),
					testAccCheck(map[string]string{
						"instance_ids.#": "1",
					}),
				),
			},
		},
	})
}

func testAccCheckEssAttachmentExists(n string, d *ess.ScalingGroup) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ESS attachment ID is set")
		}

		client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
		essService := EssService{client}
		group, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			return errmsgs.WrapError(err)
		}

		instances, err := essService.DescribeEssAttachment(rs.Primary.ID, make([]string, 0))

		if err != nil {
			return errmsgs.WrapError(err)
		}

		if len(instances) < 1 {
			return errmsgs.WrapError(errmsgs.Error("Scaling instances not found"))
		}

		*d = group
		return nil
	}
}

func testAccCheckEssAttachmentDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_scaling_configuration" {
			continue
		}

		_, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if errmsgs.NotFoundError(err) {
				continue
			}
			return errmsgs.WrapError(err)
		}

		instances, err := essService.DescribeEssAttachment(rs.Primary.ID, make([]string, 0))

		if err != nil && !errmsgs.IsExpectedErrors(err, []string{"InvalidScalingGroupId.NotFound"}) {
			return errmsgs.WrapError(err)
		}

		if len(instances) > 0 {
			return errmsgs.WrapError(fmt.Errorf("There are still ECS instances in the scaling group."))
		}
	}

	return nil
}

func testAccEssAttachmentConfig(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tf-testAccEssAttachmentConfig-%d"
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

	%s

	resource "alibabacloudstack_ecs_instance" "new" {
		image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type        = "${local.default_instance_type_id}"
		system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
		system_disk_size     = 20
		system_disk_name     = "test_sys_disk"
		security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
		instance_name        = "${var.name}_ecs"
		vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
		zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
		is_outdated          = false
		lifecycle {
			ignore_changes = [
				instance_type
			]
		}
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
	}

	resource "alibabacloudstack_ess_attachment" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		instance_ids = [alibabacloudstack_ecs_instance.default.id, alibabacloudstack_ecs_instance.new.id]
		force = true
	}
	`, rand, ECSInstanceCommonTestCase)
}

func testAccEssAttachmentConfigInstance(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tftestAcc%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 20
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
	
	resource "alibabacloudstack_ecs_instance" "new" {
		image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type        = "${local.default_instance_type_id}"
		system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
		system_disk_size     = 20
		system_disk_name     = "test_sys_disk"
		security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
		instance_name        = "${var.name}_ecs"
		vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
		zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
		is_outdated          = false
		lifecycle {
			ignore_changes = [
				instance_type
			]
		}
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
	}

	%s

	resource "alibabacloudstack_ess_attachment" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		instance_ids = [alibabacloudstack_ecs_instance.default.id, alibabacloudstack_ecs_instance.new.id]
		// 没有就绪的alibabacloudstack_ess_scaling_configuration会导致 ess_scaling_group 状态非活跃， 无法正常操作
		depends_on = ["alibabacloudstack_ess_scaling_configuration.default"]
		force = true
	}
	`, rand, ECSInstanceCommonTestCase)
}

func testAccEssAttachmentConfigRemoveInstance(rand int) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "tftestAcc%d"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 20
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
		instance_type = "${local.default_instance_type_id}"
		security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
		force_delete = true
		active = true
		enable = true
		deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
	}

	%s

	resource "alibabacloudstack_ess_attachment" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		instance_ids = [alibabacloudstack_ecs_instance.default.id]
		// 没有就绪的alibabacloudstack_ess_scaling_configuration会导致 ess_scaling_group 状态非活跃， 无法正常操作
		depends_on = ["alibabacloudstack_ess_scaling_configuration.default"]
		force = true
	}

	`, rand, ECSInstanceCommonTestCase)
}
