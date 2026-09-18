package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2NotifyGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_prometheus_v2_notify_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc-notifygroup-%d", rand),
		dataSourcePrometheusV2NotifyGroupsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                       "1",
				"notify_groups.#":             "1",
				"notify_groups.0.name":        CHECKSET,
				"notify_groups.0.type":        "WEBHOOK",
				"notify_groups.0.description": CHECKSET,
				"notify_groups.0.webhook_url": CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           "0",
				"notify_groups.#": "0",
			}
		},
	}

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_notify_group.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name.*",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_prometheus_v2_notify_group.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_notify_group.default.name}",
				"ids":        []string{"${alibabacloudstack_prometheus_v2_notify_group.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name.*",
				"ids":        []string{"fake-id"},
			}),
		},
	)
}

func dataSourcePrometheusV2NotifyGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
	name        = "${var.name}"
	type        = "WEBHOOK"
	description = "${var.name}"
	webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=56b4xxxxxxxxaea8"
	webhook_header_params {
		key   = "aaaa"
		value = "1111"
	}
}
`, name)
}
