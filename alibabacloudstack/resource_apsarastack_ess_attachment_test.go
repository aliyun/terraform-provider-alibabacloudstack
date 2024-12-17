package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ess"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/errmsgs"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccalibabacloudstackdEssAttachment_update(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	var v ess.ScalingGroup
	resourceId := "alibabacloudstack_ess_attachment.default"
	basicMap := map[string]string{
		"instance_ids.#":   "1",
		"scaling_group_id": CHECKSET,
	}
	ra := resourceAttrInit(resourceId, basicMap)

	testAccCheck := ra.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},

		// module name
		IDRefreshName: resourceId,

		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccEssAttachmentConfigInstance(SecurityGroupCommonTestCase, rand),
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
				Config: testAccEssAttachmentConfig(SecurityGroupCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckEssAttachmentExists(
						"alibabacloudstack_ess_attachment.default", &v),
					testAccCheck(map[string]string{
						"instance_ids.#": "2",
					}),
				),
			},
			{
				Config: testAccEssAttachmentConfigInstance(SecurityGroupCommonTestCase, rand),
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

func testAccEssAttachmentConfig(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAttachmentConfig-%d"
	}
	data "alibabacloudstack_images" "default" {
		name_regex  = "^ubuntu_18.*64"
		//name_regex  = "arm_centos_7_6_20G_20211110.raw"
		//name_regex  = "^arm_centos_7"
		most_recent = true
		owners      = "system"
	}

	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 100
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
		force_delete = true
		active = true
		enable = true
	}
	resource "alibabacloudstack_ecs_instance" "default" {
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		count = 2
		security_groups = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
		
		internet_max_bandwidth_out = "10"
		
		system_disk_category = "cloud_ssd"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
		instance_name = "${var.name}"
	}
	resource "alibabacloudstack_ess_attachment" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		instance_ids = ["${alibabacloudstack_ecs_instance.default.0.id}", "${alibabacloudstack_ecs_instance.default.1.id}"]
		force = true
	}
	`, common, rand)
}

func testAccEssAttachmentConfigInstance(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssAttachmentConfig-%d"
	}

	data "alibabacloudstack_images" "default" {
		name_regex  = "^ubuntu_18.*64"
		//name_regex  = "arm_centos_7_6_20G_20211110.raw"
		//name_regex  = "^arm_centos_7"
		most_recent = true
		owners      = "system"
	}
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 0
		max_size = 100
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
	}
	resource "alibabacloudstack_ess_scaling_configuration" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
		force_delete = true
		active = true
		enable = true
	}
	resource "alibabacloudstack_ecs_instance" "default" {
		image_id = "${data.alibabacloudstack_images.default.images.0.id}"
		instance_type = "ecs.e4.small"
		count = 2
		security_groups = ["${alibabacloudstack_ecs_securitygroup.default.id}"]
	
		internet_max_bandwidth_out = "10"
		
		system_disk_category = "cloud_ssd"
		vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
		instance_name = "${var.name}"
	}
	resource "alibabacloudstack_ess_attachment" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		instance_ids = ["${alibabacloudstack_ecs_instance.default.0.id}"]
		force = true
	}
	`, common, rand)
}
