package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudstackCmsAlarmContactGroupGroups_basic(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_cms_alarm_contact_groups.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testAccCmsAlarmContactGroupBisic-%d", rand),
		dataSourceCmsAlarmContactGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact_group.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact_group.default.id}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cms_alarm_contact_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cms_alarm_contact_group.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact_group.default.id}",
			"ids":        []string{"${alibabacloudstack_cms_alarm_contact_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact_group.default.id}_fake",
			"ids":        []string{"${alibabacloudstack_cms_alarm_contact_group.default.id}_fake"},
		}),
	}

	var existCmsAlarmContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"groups.#":                          "1",
			"groups.0.id":                       CHECKSET,
			"groups.0.alarm_contact_group_name": CHECKSET,
			"groups.0.describe":                 "For Test",
			"groups.0.enable_subscribed":        "true",
		}
	}

	var fakeCmsAlarmContactGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var cmsAlarmContactGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCmsAlarmContactGroupsMapFunc,
		fakeMapFunc:  fakeCmsAlarmContactGroupsMapFunc,
		PreCheck: func() {
			testAccPreCheck(t, resourceId)
		},
	}

	cmsAlarmContactGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceCmsAlarmContactGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}
resource "alibabacloudstack_cms_alarm_contact_group" "default" {
  alarm_contact_group_name = var.name
  describe                 = "For Test"
  enable_subscribed        = true
}
`, name)
}
