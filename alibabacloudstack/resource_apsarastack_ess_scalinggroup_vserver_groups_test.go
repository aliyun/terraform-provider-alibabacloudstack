package alibabacloudstack

import (
	"fmt"
	"testing"

	"github.com/aliyun/terraform-provider-alibabacloudstack/alibabacloudstack/connectivity"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAlibabacloudStackEssVserverGroups_basic(t *testing.T) {
	rand := acctest.RandIntRange(1000, 999999)
	resourceId := "alibabacloudstack_ess_scalinggroup_vserver_groups.default"
	basicMap := map[string]string{
		"scaling_group_id": CHECKSET,
		"vserver_groups.#": "2",
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
		CheckDestroy: testAccCheckEssVserverGroupsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccEssScalingGroupVserverGroup(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_groups.#": "2",
					}),
				),
			},
			{
				ResourceName:            resourceId,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"force"},
			},
			{
				Config: testAccEssScalingGroupVserverGroupUpdate(EcsInstanceCommonTestCase, rand),
				Check: resource.ComposeTestCheckFunc(
					testAccCheck(map[string]string{
						"vserver_groups.#": "1",
					}),
				),
			},
		},
	})
}

func testAccCheckEssVserverGroupsDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*connectivity.AlibabacloudStackClient)
	essService := EssService{client}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "alibabacloudstack_ess_scaling_group" {
			continue
		}

		scalingGroup, err := essService.DescribeEssScalingGroup(rs.Primary.ID)
		if err != nil {
			if NotFoundError(err) {
				continue
			}
			return WrapError(err)
		}

		if len(scalingGroup.VServerGroups.VServerGroup) > 0 {
			return WrapError(fmt.Errorf("There are still attached vserver groups."))
		}
	}
	return nil
}

func testAccEssScalingGroupVserverGroup(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testaccessscalinggroupupdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_ess_scalinggroup_vserver_groups" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		vserver_groups {
				loadbalancer_id = "${alibabacloudstack_slb.default.0.id}"
				vserver_attributes {
					vserver_group_id = "${alibabacloudstack_slb_server_group.vserver0.0.id}"
					port = "100"
					weight = "60"
				}
			}
      vserver_groups {
				loadbalancer_id = "${alibabacloudstack_slb.default.1.id}"
				vserver_attributes {
					vserver_group_id = "${alibabacloudstack_slb_server_group.vserver1.0.id}"
					port = "200"
					weight = "60"
				}
				vserver_attributes {
					vserver_group_id = "${alibabacloudstack_slb_server_group.vserver1.1.id}"
					port = "210"
					weight = "60"
				}
			}
	force = true
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_slb_server_group" "vserver0" {
 	  count = "2"
	  load_balancer_id = "${alibabacloudstack_slb.default.0.id}"
	  name = "test"
	}

	resource "alibabacloudstack_slb_server_group" "vserver1" {
 	  count = "2"
	  load_balancer_id = "${alibabacloudstack_slb.default.1.id}"
	  name = "test"
	}

	resource "alibabacloudstack_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alibabacloudstack_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	  health_check = "off"
	  sticky_session = "off"
	}
	`, common, rand)
}

func testAccEssScalingGroupVserverGroupUpdate(common string, rand int) string {
	return fmt.Sprintf(`
	%s
	variable "name" {
		default = "tf-testaccessscalinggroupupdate-%d"
	}
	
	resource "alibabacloudstack_ess_scaling_group" "default" {
	  min_size = "2"
	  max_size = "2"
      default_cooldown = 200
	  scaling_group_name = "${var.name}"
	  removal_policies = ["OldestInstance"]
	  vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
	  loadbalancer_ids = ["${alibabacloudstack_slb.default.0.id}","${alibabacloudstack_slb.default.1.id}"]
	  depends_on = ["alibabacloudstack_slb_listener.default"]
	}

	resource "alibabacloudstack_ess_scalinggroup_vserver_groups" "default" {
		scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
		vserver_groups {
				loadbalancer_id = "${alibabacloudstack_slb.default.0.id}"
				vserver_attributes {
					vserver_group_id = "${alibabacloudstack_slb_server_group.vserver0.1.id}"
					port = "110"
					weight = "60"
				}
			}
		force = false
	}

	resource "alibabacloudstack_slb" "default" {
	  count=2
	  name = "${var.name}"
	  vswitch_id = "${alibabacloudstack_vswitch.default.id}"
	}

	resource "alibabacloudstack_slb_server_group" "vserver0" {
 	  count = "2"
	  load_balancer_id = "${alibabacloudstack_slb.default.0.id}"
	  name = "test"
	}

	resource "alibabacloudstack_slb_server_group" "vserver1" {
 	  count = "2"
	  load_balancer_id = "${alibabacloudstack_slb.default.1.id}"
	  name = "test"
	}

	resource "alibabacloudstack_slb_listener" "default" {
	  count = 2
	  load_balancer_id = "${element(alibabacloudstack_slb.default.*.id, count.index)}"
	  backend_port = "22"
	  frontend_port = "22"
	  protocol = "tcp"
	  bandwidth = "10"
	  health_check_type = "tcp"
	  health_check = "off"
	  sticky_session = "off"
	}
	`, common, rand)
}
