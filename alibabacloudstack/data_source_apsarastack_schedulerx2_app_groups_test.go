package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackSchedulerx2AppGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_schedulerx2_app_groups.default",
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                    "1",
				"groups.#":                 "1",
				"groups.0.app_name":        fmt.Sprintf("tf-test%d", rand),
				"groups.0.description":     fmt.Sprintf("tf-test%d", rand),
				"groups.0.group_id":        fmt.Sprintf("tf-test%d.terra", rand),
				"groups.0.max_jobs":        "30",
				"groups.0.max_concurrency": "10",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":    "0",
				"groups.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
		}),
		fakeConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name$"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_schedulerx2_app_group.example.id}"]`,
		}),
		fakeConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	appNameConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"app_name": `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
		}),
		fakeConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"app_name": `"fake-app-name"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
			"ids":        `["${alibabacloudstack_schedulerx2_app_group.example.id}"]`,
			"app_name":   `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
		}),
		fakeConfig: AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
			"ids":        `["fake-id"]`,
			"app_name":   `"${alibabacloudstack_schedulerx2_app_group.example.app_name}"`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, appNameConf, allConf)
}

func AlibabacloudStackSchedulerx2AppGroupDependenceNew(rand int, attrMap map[string]string) string {
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

data "alibabacloudstack_schedulerx2_app_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n"))
	return config
}
