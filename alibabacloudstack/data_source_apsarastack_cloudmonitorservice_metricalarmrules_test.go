package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCms_Alarams_DataSource(t *testing.T) {
	// testAccPreCheckWithAPIIsNotSupport(t)
	rand := getAccTestRandInt(10000, 20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: dataSourceAlibabacloudStackcms_alarms(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm.default.rule_name}"`,
		}),
		fakeConfig: dataSourceAlibabacloudStackcms_alarms(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm.default.rule_name}_fake"`,
		}),
	}

	var existcmsAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"alarms.#":           "1",
			"alarms.0.rule_name": CHECKSET,
		}
	}

	var fakecmsAlarmsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"alarms.#": "0",
		}
	}

	var cmsAlarmContactsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_cms_alarms.default",
		existMapFunc: existcmsAlarmsMapFunc,
		fakeMapFunc:  fakecmsAlarmsMapFunc,
	}

	cmsAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceAlibabacloudStackcms_alarms(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	return fmt.Sprintf(`

variable "name" {
 default = "tf_testacc_cmsalarm%d"
}

resource "alibabacloudstack_slb" "basic" {
 name          = "${var.name}"
}
resource "alibabacloudstack_cms_alarm" "default" {
  name    = "${var.name}"
  project = "acs_slb_dashboard"
  metric  = "ActiveConnection"
  dimensions = {
    instanceId = alibabacloudstack_slb.basic.id
  }
  escalations_critical {
    statistics = "Average"
    comparison_operator = "<="
    threshold = 35
    times = 2
  }
  enabled =      true
  contact_groups     = ["test-group"]
  effective_interval = "0:00-2:00"
  
  lifecycle {
    ignore_changes = [
      dimensions,
      period,
    ]
  }
}

data "alibabacloudstack_cms_alarms" "default" {
%s
}
`, rand, strings.Join(pairs, "\n  "))
}
