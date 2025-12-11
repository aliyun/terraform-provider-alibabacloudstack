package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackSchedulerx2WorkflowsDataSource(t *testing.T) {
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_schedulerx2_workflows.default",
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

	testAcc.dataSourceTestCheck(t, -1,
		dataSourceTestAccConfig{
			existConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"ids": `["${alibabacloudstack_schedulerx2_workflow.example.id}"]`,
			}),
			fakeConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"ids": `["nonexistent-group-id"]`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"group_id":   `"${alibabacloudstack_schedulerx2_workflow.example.group_id}"`,
				"name_regex": `"${alibabacloudstack_schedulerx2_workflow.example.name}"`,
			}),
			fakeConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"group_id":   `"${alibabacloudstack_schedulerx2_app_group.example.group_id}"`,
				"name_regex": `"nonexistent-name"`,
			}),
		},
		dataSourceTestAccConfig{
			existConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"name_regex": `"${alibabacloudstack_schedulerx2_workflow.example.name}"`,
				"ids":        `["${alibabacloudstack_schedulerx2_workflow.example.id}"]`,
			}),
			fakeConfig: AlibabacloudStackSchedulerx2WorkflowDependenceNew(-1, map[string]string{
				"name_regex": `"${alibabacloudstack_schedulerx2_workflow.example.name}"`,
				"ids":        `["999999"]`,
			}),
		},
	)
}

func AlibabacloudStackSchedulerx2WorkflowDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, "  "+k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
  default = "tf-test%d"
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

data "alibabacloudstack_schedulerx2_workflows" "default" {
%s
}
`, rand, strings.Join(pairs, "\n"))
	return config
}
