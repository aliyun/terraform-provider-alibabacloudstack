package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackSchedulerx2JobsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_schedulerx2_jobs.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"jobs.#":                        "1",
				"jobs.0.name":                   fmt.Sprintf("tf-test%d", rand),
				"jobs.0.group_id":               fmt.Sprintf("tf-test%d.terra", rand),
				"jobs.0.job_type":               "python",
				"jobs.0.execute_mode":           "standalone",
				"jobs.0.description":            fmt.Sprintf("tf-test%d", rand),
				"jobs.0.priority":               "5",
				"jobs.0.parameters":             "testargs=1",
				"jobs.0.max_attempt":            "3",
				"jobs.0.attempt_interval":       "30",
				"jobs.0.max_concurrency":        "2",
				"jobs.0.time_type":              "1",
				"jobs.0.time_expression":        "8 59 15 */1 * ?",
				"jobs.0.content":                "python test",
				"jobs.0.monitor_timeout_enable": "true",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"jobs.#": "0",
			}
		},
	}

	schedulerx2JobBasic := buildBasicSchedulerx2JobNew

	// Test case 1: Filter by name_regex
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: schedulerx2JobBasic(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_schedulerx2_job.example.name}"`,
		}),
		fakeConfig: schedulerx2JobBasic(rand, map[string]string{
			"name_regex": `"nonexistent-job"`,
		}),
	}

	// Test case 2: Filter by ids
	idConf := dataSourceTestAccConfig{
		existConfig: schedulerx2JobBasic(rand, map[string]string{
			"ids": `["${alibabacloudstack_schedulerx2_job.example.id}"]`,
		}),
		fakeConfig: schedulerx2JobBasic(rand, map[string]string{
			"ids": `["nonexistent-id"]`,
		}),
	}

	// Test case 3: Filter by group_id
	groupIdConf := dataSourceTestAccConfig{
		existConfig: schedulerx2JobBasic(rand, map[string]string{
			"group_id": `"${alibabacloudstack_schedulerx2_job.example.group_id}"`,
			"ids":      `["${alibabacloudstack_schedulerx2_job.example.id}"]`,
		}),
		fakeConfig: schedulerx2JobBasic(rand, map[string]string{
			"group_id": `"nonexistent-group-id"`,
			"ids":      `["nonexistent-id"]`,
		}),
	}

	// Run tests
	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idConf, groupIdConf)
}

// buildBasicSchedulerx2JobNew generates the Terraform configuration for creating a SchedulerX2 job and its dependencies.
// It uses an app group as dependency which is created using alibabacloudstack_schedulerx2_app_group resource.
func buildBasicSchedulerx2JobNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`
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

resource "alibabacloudstack_schedulerx2_job" "example" {
  namespace          = "system_namespace"
  group_id           = alibabacloudstack_schedulerx2_app_group.example.group_id
  name               = var.name
  job_type           = "python"
  execute_mode       = "standalone"
  description        = var.name
  priority           = 5
  parameters         = "testargs=1"
  max_attempt        = 3
  attempt_interval   = 30
  max_concurrency    = 2
  time_type          = 1
  time_expression    = "8 59 15 */1 * ?"
  content            = "python test"
  monitor_timeout_enable       = true
  monitor_timeout_kill_enable  = true
  monitor_fail_enable          = true
  monitor_miss_worker_enable   = true
  monitor_timeout              = 3600
}

data "alibabacloudstack_schedulerx2_jobs" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
