package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2AlertsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_prometheus_v2_alerts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourcePrometheusV2AlertsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"alerts.#":                      "1",
				"alerts.0.id":                   CHECKSET,
				"alerts.0.name":                 CHECKSET,
				"alerts.0.notify_recovered":     CHECKSET,
				"alerts.0.trigger_promql":       CHECKSET,
				"alerts.0.trigger_period":       CHECKSET,
				"alerts.0.trigger_severity":     CHECKSET,
				"alerts.0.trigger_cron":         CHECKSET,
				"alerts.0.notification":         CHECKSET,
				"alerts.0.recover_notification": CHECKSET,
				"alerts.0.tag_set.#":            CHECKSET,
				"alerts.0.notify_types.#":       CHECKSET,
				"alerts.0.notify_group_ids.#":   CHECKSET,
				"alerts.0.notify_interval":      CHECKSET,
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"alerts.#": "0",
			}
		},
	}

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_alert.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_prometheus_v2_alert.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_alert.default.name}",
				"ids":        []string{"${alibabacloudstack_prometheus_v2_alert.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_alert.default.name}",
				"ids":        []string{"fake-id"},
			}),
		},
	)
}

func dataSourcePrometheusV2AlertsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_prometheus_v2_instance" "default" {
	cluster_name = "${var.name}"
	tags         = ["test1", "test2"]
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
	name        = "${var.name}_notify_group"
	type        = "WEBHOOK"
	description = "${var.name}_notify_group_description"
	webhook_url = "https://test.com"
	webhook_header_params {
		key   = "Content-Type"
		value = "application/json"
	}
}

resource "alibabacloudstack_prometheus_v2_alert" "default" {
	name                 = "${var.name}"
	notify_recovered     = true
	is_check_all         = false
	tag_set              = ["aaa", "ccc"]
	trigger_clusters     = ["${alibabacloudstack_prometheus_v2_instance.default.id}"]
	trigger_promql       = "select testfield from testtable where testfield >= 0"
	trigger_severity     = "warning"
	trigger_cron         = "0 /5 * * * ?"
	trigger_period       = "5m"
	recover_notification = "Trigger condition\\\\$${alert_source} \\\\Hit record\\\\$${alert_time}"
	notification         = "Trigger condition: {condition}\nHit record :{alert_result}"
	notify_group_ids     = ["${alibabacloudstack_prometheus_v2_notify_group.default.id}"]
	notify_types         = ["EMAIL", "SMS"]
	notify_interval      = "10m"
}
`, name)
}
