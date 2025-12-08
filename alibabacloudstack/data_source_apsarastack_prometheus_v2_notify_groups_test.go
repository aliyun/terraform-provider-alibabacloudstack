package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2NotifyGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_prometheus_v2_notify_groups.default",
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

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_notify_group.default.name}"`,
		}),
		fakeConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name.*"`,
		}),
	}

	idConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_prometheus_v2_notify_group.default.id}"]`,
		}),
		fakeConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_notify_group.default.name}"`,
			"ids":        `["${alibabacloudstack_prometheus_v2_notify_group.default.id}"]`,
		}),
		fakeConfig: resourcePrometheusV2NotifyGroupConfigDependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name.*"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idConf, allConf)
}

// Dependency template generation method
func resourcePrometheusV2NotifyGroupConfigDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		if v == "" {
			pairs = append(pairs, k)
		} else {
			pairs = append(pairs, k+" = "+v)
		}
	}
	config := `
variable "name" {
  default = "tfacc-notifygroup-%d"
}

resource "alibabacloudstack_prometheus_v2_notify_group" "default" {
  name        = "${var.name}"
  type        = "WEBHOOK"
  description = "${var.name}"
  webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=56b42bc6e7cad53bab514a583847db73c68fa1804b0e72af7167954b66f7aea8"
  webhook_header_params {
    key   = "aaaa"
    value = "1111"
  }
}

data "alibabacloudstack_prometheus_v2_notify_groups" "default" {
%s
}
`
	return fmt.Sprintf(config, rand, strings.Join(pairs, "\n"))
}
