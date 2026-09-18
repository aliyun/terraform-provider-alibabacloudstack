package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSchedulerx2WorkflowsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_schedulerx2_workflows.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		"tf-test-1",
		dataSourceSchedulerx2WorkflowsConfigDependence)

	checkInfo := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"workflows.#":                 "1",
				"workflows.0.workflow_id":     CHECKSET,
				"workflows.0.group_id":        CHECKSET,
				"workflows.0.name":            CHECKSET,
				"workflows.0.description":     "Initial description",
				"workflows.0.time_type":       "cron",
				"workflows.0.time_expression": "34 14 14 */1 * ?",
				"workflows.0.max_concurrency": "1",
				"workflows.0.creator":         CHECKSET,
				"workflows.0.updater":         CHECKSET,
				"workflows.0.app_group_id":    CHECKSET,
				"ids.#":                       "1",
				"ids.0":                       CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"workflows.#": "0",
				"ids.#":       "0",
			}
		},
	}

	checkInfo.dataSourceTestCheck(t, -1,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_schedulerx2_workflow.example.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"nonexistent-group-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"group_id":   "${alibabacloudstack_schedulerx2_workflow.example.group_id}",
				"name_regex": "${alibabacloudstack_schedulerx2_workflow.example.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"group_id":   "${alibabacloudstack_schedulerx2_app_group.example.group_id}",
				"name_regex": "nonexistent-name",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_schedulerx2_workflow.example.name}",
				"ids":        []string{"${alibabacloudstack_schedulerx2_workflow.example.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_schedulerx2_workflow.example.name}",
				"ids":        []string{"999999"},
			}),
		},
	)
}

func dataSourceSchedulerx2WorkflowsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_schedulerx2_app_group" "example" {
  group_id         = "${var.name}.terra"
  app_name         = "${var.name}"
  description      = "${var.name}"
  max_jobs         = 30
  max_concurrency  = 10
  monitor_config {
    send_channel = "mail,ding"
    alarm_type   = "CustomContacts"
  }
  contacts {
    username    = "test"
    user_email  = "123@123.com"
    dingding_ak = "testakkkkk"
  }
  metrics_threshold {
    load5       = 10
    heap5_usage = 100
    disk_usage  = 100
  }
}

resource "alibabacloudstack_schedulerx2_workflow" "example" {
  group_id         = "${alibabacloudstack_schedulerx2_app_group.example.group_id}"
  name             = "${var.name}"
  description      = "Initial description"
  time_type        = "cron"
  time_zone		   = "PRC"
  time_expression  = "34 14 14 */1 * ?"
  max_concurrency  = 1
}
`, name)
}
