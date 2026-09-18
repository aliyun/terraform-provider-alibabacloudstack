package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackSchedulerx2AppGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_schedulerx2_app_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-test%d", rand),
		dataSourceSchedulerx2AppGroupsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
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

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_schedulerx2_app_group.example.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"app_name": "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"app_name": "fake-app-name",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
				"ids":        []string{"${alibabacloudstack_schedulerx2_app_group.example.id}"},
				"app_name":   "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
				"ids":        []string{"fake-id"},
				"app_name":   "${alibabacloudstack_schedulerx2_app_group.example.app_name}",
			}),
		},
	)
}

func dataSourceSchedulerx2AppGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_schedulerx2_app_group" "example" {
	group_id        = "${var.name}.terra"
	app_name        = "${var.name}"
	description     = "${var.name}"
	max_jobs        = 30
	max_concurrency = 10
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
`, name)
}
