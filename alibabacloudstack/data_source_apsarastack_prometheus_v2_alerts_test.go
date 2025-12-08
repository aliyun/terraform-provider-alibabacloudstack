package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2AlertsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_prometheus_v2_alerts.default",
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

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_alert.default.name}"`,
		}),
		fakeConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name$"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_prometheus_v2_alert.default.id}"]`,
		}),
		fakeConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_alert.default.name}"`,
			"ids":        `["${alibabacloudstack_prometheus_v2_alert.default.id}"]`,
		}),
		fakeConfig: testAccPrometheusV2AlertsConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_alert.default.name}"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccPrometheusV2AlertsConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tfacc%d"
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
  notify_recovered 	   = true
  is_check_all 		   = false
  tag_set 			   = ["aaa", "ccc"]
  trigger_clusters     = ["${alibabacloudstack_prometheus_v2_instance.default.id}"]
  trigger_promql 	   = "select testfield from testtable where testfield >= 0"
  trigger_severity     = "warning"
  trigger_cron 		   = "0 /5 * * * ?"
  trigger_period       = "5m"
  recover_notification = "Trigger condition\\\\$${alert_source} \\\\Hit record\\\\$${alert_time}"
  notification         = "Trigger condition: {condition}\nHit record :{alert_result}"
  notify_group_ids     = ["${alibabacloudstack_prometheus_v2_notify_group.default.id}"]
  notify_types         = ["EMAIL", "SMS"]
  notify_interval      = "10m"
}

data "alibabacloudstack_prometheus_v2_alerts" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
