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

func TestAccAlibabacloudStackEssNotification_basic(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	rand := getAccTestRandInt(1000, 999999)
	var v ess.NotificationConfigurationModel
	resourceId := "alibabacloudstack_ess_notification.default"

	basicMap := map[string]string{
		"notification_types.#": "2",
		"scaling_group_id":     CHECKSET,
		"notification_arn":     CHECKSET,
	}

	ra := resourceAttrInit(resourceId, basicMap)
	rc := resourceCheckInit(resourceId, &v, func() interface{} {
		return &EssService{testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)}
	})
	rac := resourceAttrCheckInit(rc, ra)

	testAccCheck := rac.resourceAttrMapUpdateSet()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		IDRefreshName: resourceId,
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckEssNotificationDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssNotification(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				ResourceName:      resourceId,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccEssNotification_update_notification_types(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"notification_types.#": "4",
					}),
				),
			},
			{
				Config: testAccEssNotification_update_scaling_group_id(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
			{
				Config: testAccEssNotification_update_notification_arn(ECSInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(nil),
				),
			},
		},
	})

}

func testAccCheckEssNotificationDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_notification" {
			continue
		}
		if _, err := essService.DescribeEssNotification(rs.Primary.ID); err != nil {
			if errmsgs.IsExpectedErrors(err, []string{"NotificationConfigurationNotExist", "InvalidScalingGroupId.NotFound"}) {
				return nil
			}
			if errmsgs.NotFoundError(err) {
				continue
			}
			return err
		}
		return fmt.Errorf("ess notification %s still exists.", rs.Primary.ID)
	}
	return nil
}

func testAccEssNotification(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	data "alibabacloudstack_regions" "default" {
		current = true
	}

	data "alibabacloudstack_account" "default" {
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	}

	resource "alibabacloudstack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "alibabacloudstack_ess_notification" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR"]
		notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default.name}"
	}
	`, common, rand)
}

func testAccEssNotification_update_notification_types(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	data "alibabacloudstack_regions" "default" {
		current = true
	}

	data "alibabacloudstack_account" "default" {
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.name}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	}

	resource "alibabacloudstack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "alibabacloudstack_ess_notification" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default.name}"
	}
	`, common, rand)
}

func testAccEssNotification_update_scaling_group_id(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	variable "newname" {
		default = "tf-testAccEssNotification_new-%d"
	}

	data "alibabacloudstack_regions" "default" {
		current = true
	}

	data "alibabacloudstack_account" "default" {
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.newname}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	}

	resource "alibabacloudstack_mns_queue" "default"{
		name="${var.name}"
	}
	
	resource "alibabacloudstack_ess_notification" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default1.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default.name}"
	}
	`, common, rand, rand)
}

func testAccEssNotification_update_notification_arn(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testAccEssNotification-%d"
	}

	variable "newname" {
		default = "tf-testAccEssNotification-new-%d"
	}

	data "alibabacloudstack_regions" "default" {
		current = true
	}

	data "alibabacloudstack_account" "default" {
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default1" {
		min_size = 1
		max_size = 1
		scaling_group_name = "${var.newname}"
		removal_policies = ["OldestInstance", "NewestInstance"]
		vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	}

	resource "alibabacloudstack_mns_queue" "default1"{
		name="${var.newname}"
	}
	
	resource "alibabacloudstack_ess_notification" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default1.id}"
		notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS","AUTOSCALING:SCALE_OUT_ERROR","AUTOSCALING:SCALE_IN_SUCCESS","AUTOSCALING:SCALE_IN_ERROR"]
		notification_arn = "acs:ess:${data.alibabacloudstack_regions.default.regions.0.id}:${data.alibabacloudstack_account.default.id}:queue/${alibabacloudstack_mns_queue.default1.name}"
	}
	`, common, rand, rand)
}
