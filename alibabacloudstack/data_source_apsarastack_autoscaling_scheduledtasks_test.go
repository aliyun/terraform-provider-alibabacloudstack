package alibabacloudstack

import (
	"fmt"
	"testing"
	"time"
)

func TestAccAlibabacloudStackEssScheduledTasksDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_scheduled_tasks.default"
	name := fmt.Sprintf("tf-testacc-essscheduledtask%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEssScheduledTasksConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ess_scheduled_task.default.scheduled_task_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scheduled_task.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ess_scheduled_task.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ess_scheduled_task.default.id}"},
			"name_regex": "${alibabacloudstack_ess_scheduled_task.default.scheduled_task_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_ess_scheduled_task.default.id}_fake"},
			"name_regex": "${alibabacloudstack_ess_scheduled_task.default.scheduled_task_name}_fake",
		}),
	}

	var existEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                                   "1",
			"tasks.#":                                 "1",
			"tasks.0.id":                              CHECKSET,
			"tasks.0.name":                            name,
			"tasks.0.scheduled_action":                CHECKSET,
			"tasks.0.description":                     "",
			"tasks.0.launch_expiration_time":          "600",
			"tasks.0.launch_time":                     CHECKSET,
			"tasks.0.max_value":                       "0",
			"tasks.0.min_value":                       "0",
			"tasks.0.recurrence_end_time":             "",
			"tasks.0.recurrence_value":                "",
			"tasks.0.recurrence_type":                 "",
			"tasks.0.task_enabled":                    "true",
			"tasks.0.scaling_group_id":                CHECKSET,
		}
	}

	var fakeEssScheduledTasksMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"tasks.#": "0",
		}
	}

	var essScheduledTasksCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEssScheduledTasksMapFunc,
		fakeMapFunc:  fakeEssScheduledTasksMapFunc,
	}
	essScheduledTasksCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceEssScheduledTasksConfigDependence(name string) string {
	oneDay, _ := time.ParseDuration("24h")
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  default_cooldown   = 20
  removal_policies   = ["OldestInstance", "NewestInstance"]
  scaling_group_name = var.name
  vswitch_ids        = [alibabacloudstack_vpc_vswitch.default.id]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id         = alibabacloudstack_ess_scaling_group.default.id
  image_id                 = data.alibabacloudstack_images.default.images.0.id
  instance_type            = alibabacloudstack_ecs_instance.default.instance_type
  system_disk_category     = alibabacloudstack_ecs_instance.default.system_disk_category
  security_group_ids       = [alibabacloudstack_ecs_securitygroup.default.id]
  force_delete             = true
  active                   = true
  enable                   = true
  deployment_set_id        = alibabacloudstack_ecs_deployment_set.default.id
}

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id   = alibabacloudstack_ess_scaling_group.default.id
  adjustment_type    = "TotalCapacity"
  adjustment_value   = "1"
  cooldown           = 0
}

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action       = alibabacloudstack_ess_scaling_rule.default.ari
  scheduled_task_name    = var.name
  launch_time            = "%s"
}
`, name, ECSInstanceCommonTestCase, time.Now().Add(oneDay).Format("2006-01-02T15:04Z"))
}
