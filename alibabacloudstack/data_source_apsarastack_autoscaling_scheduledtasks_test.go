package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccAlibabacloudStackEssScheduledtasksDataSource(t *testing.T) {
	oneDay, _ := time.ParseDuration("24h")
	ResourceTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCheckAlibabacloudStackEssScheduledTasksDataSource, ECSInstanceCommonTestCase, time.Now().Add(oneDay).Format("2006-01-02T15:04Z")),
				Check: resource.ComposeTestCheckFunc(

					testAccCheckAlibabacloudStackDataSourceID("data.alibabacloudstack_ess_scheduled_tasks.default"),
					resource.TestCheckResourceAttr("data.alibabacloudstack_ess_scheduled_tasks.default", "tasks.#", "1"),
					resource.TestCheckNoResourceAttr("data.alibabacloudstack_ess_scheduled_tasks.default", "tasks.1.id"),
					resource.TestCheckResourceAttrSet("data.alibabacloudstack_ess_scheduled_tasks.default", "ids.#"),
				),
			},
		},
	})
}

const testAccCheckAlibabacloudStackEssScheduledTasksDataSource = `

variable "name" {
	default = "tf-testAccDataSourceScheduledtas"
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

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action    = "${alibabacloudstack_ess_scaling_rule.default.ari}"
  launch_time         = "%s"
  scheduled_task_name = "${var.name}"
}

data "alibabacloudstack_ess_scheduled_tasks" "default"{
  ids = ["${alibabacloudstack_ess_scheduled_task.default.id}"]
}
`
