package alibabacloudstack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
	"time"
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
				Config: fmt.Sprintf(testAccCheckAlibabacloudStackEssScheduledTasksDataSource, time.Now().Add(oneDay).Format("2006-01-02T15:04Z")),
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

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name = "${var.name}"
  cidr_block = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id = "${alibabacloudstack_vpc.default.id}"
  cidr_block = "10.1.0.0/16"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name = "${var.name}"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
	min_size = 0
	max_size = 2
	default_cooldown = 20
	removal_policies = ["OldestInstance", "NewestInstance"]
	scaling_group_name = "${var.name}"
	vswitch_ids = ["${alibabacloudstack_vswitch.default.id}"]
}

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = 2
  cooldown         = 60
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
