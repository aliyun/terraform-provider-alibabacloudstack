package alibabacloudstack

import (
	"fmt"

	"testing"
)

func TestAccAlibabacloudStackEssNotificationsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ess_notifications.default"
	name := fmt.Sprintf("tf-essgroup%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackEssNotificationsDataSourceConfig)
	scalingGroupIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_notification.default.scaling_group_id}",
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_notification.default.scaling_group_id}",
			"ids":              []string{"${alibabacloudstack_ess_notification.default.notification_arn}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"scaling_group_id": "${alibabacloudstack_ess_notification.default.scaling_group_id}",
			"ids":              []string{"${alibabacloudstack_ess_notification.default.notification_arn}_fake"},
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
		resourceId:   resourceId,
		existMapFunc: existEssnotificationsMapFunc,
		fakeMapFunc:  fakeEssnotificationsMapFunc,
	}

	essNotificationsCheckInfo.dataSourceTestCheck(t, rand, scalingGroupIdConf, allConf)
}

func testAccCheckAlibabacloudStackEssNotificationsDataSourceConfig(name string) string {
	return fmt.Sprintf(`

variable "name" {
	default = "%s"
}

%s

resource "alibabacloudstack_ess_scaling_group" "default" {
    min_size = 1
    max_size = 1
    scaling_group_name = "${var.name}"
    removal_policies = ["OldestInstance", "NewestInstance"]
    vswitch_ids = [alibabacloudstack_vpc_vswitch.default.id,]
}


resource "alibabacloudstack_ess_notification" "default" {
    scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
    notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS"]
    notification_arn = "acs:ess"
}
`, name, ECSInstanceCommonTestCase)
}
