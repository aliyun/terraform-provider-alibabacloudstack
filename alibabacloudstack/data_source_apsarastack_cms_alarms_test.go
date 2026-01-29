package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCmsAlarmsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_cms_alarms.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf_testacc_cmsalarm%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, dataSourceAlibabacloudStackcms_alarms)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm.default.rule_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm.default.rule_name}_fake",
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
		resourceId:   resourceId,
		existMapFunc: existcmsAlarmsMapFunc,
		fakeMapFunc:  fakecmsAlarmsMapFunc,
	}

	cmsAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf)
}

func dataSourceAlibabacloudStackcms_alarms(name string) string {
	return fmt.Sprintf(`

variable "name" {
 default = "%s"
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
`, name)
}
