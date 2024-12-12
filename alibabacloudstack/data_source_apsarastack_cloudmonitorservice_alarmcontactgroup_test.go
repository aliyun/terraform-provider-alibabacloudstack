package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackCloudmonitorserviceAlarmContactGroupsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contact_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contact_groups.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contact_groups.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactGroupsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contact_groups.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackCloudmonitorserviceAlarmContactGroupsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackCloudmonitorserviceAlarmContactGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#":    "1",
		"groups.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCloudmonitorserviceAlarmContactGroupsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"groups.#": "0",
	}
}

var AlibabacloudstackCloudmonitorserviceAlarmContactGroupsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cloudmonitorservice_alarm_contact_groups.default",
	existMapFunc: existAlibabacloudstackCloudmonitorserviceAlarmContactGroupsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCloudmonitorserviceAlarmContactGroupsMapFunc,
}

func testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactGroupsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCloudmonitorserviceAlarmContactGroups%d"
}






data "alibabacloudstack_cloudmonitorservice_alarm_contact_groups" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
