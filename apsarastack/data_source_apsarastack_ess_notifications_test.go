package apsarastack

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"strings"
	"testing"
)

func TestAccApsaraStackEssNotificationsDataSource(t *testing.T) {
	rand := acctest.RandIntRange(0, 500)
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_notification.default.scaling_group_id}"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckApsaraStackEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_notification.default.scaling_group_id}"`,
			"ids":              `["${apsarastack_ess_notification.default.notification_arn}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackEssNotificationsDataSourceConfig(rand, map[string]string{
			"scaling_group_id": `"${apsarastack_ess_notification.default.scaling_group_id}"`,
			"ids":              `["${apsarastack_ess_notification.default.notification_arn}_fake"]`,
		}),
	}

	var existEssnotificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                            "1",
			"notifications.#":                  "1",
			"notifications.0.notification_arn": CHECKSET,
			"notifications.0.scaling_group_id": CHECKSET,
		}
	}

	var fakeEssnotificationsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"notifications.#": "0",
			"ids.#":           "0",
		}
	}

	var essNotificationsCheckInfo = dataSourceAttr{
		resourceId:   "data.apsarastack_ess_notifications.default",
		existMapFunc: existEssnotificationsMapFunc,
		fakeMapFunc:  fakeEssnotificationsMapFunc,
	}

	essNotificationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, allConf)
}

func testAccCheckApsaraStackEssNotificationsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
%s

variable "name" {
	default = "tf-testAccDataSourceEssNs-%d"
}

resource "apsarastack_ess_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = ["${apsarastack_vswitch.default.id}"]
}


resource "apsarastack_ess_notification" "default" {
    scaling_group_id = "${apsarastack_ess_scaling_group.default.id}"
    notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS"]
    notification_arn = "acs:ess"
}

data "apsarastack_ess_notifications" "default"{
  %s
}
`, EcsInstanceCommonTestCase, rand, strings.Join(pairs, "\n  "))
	return config
}
