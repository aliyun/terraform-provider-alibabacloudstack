package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackCloudmonitorserviceAlarmContactsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}_fake"]`,
		}),
	}

	alarm_contact_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}"]`,
			"alarm_contact_name": `"${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.AlarmContactName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids":                `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}_fake"]`,
			"alarm_contact_name": `"${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.AlarmContactName}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}"]`,

			"alarm_contact_name": `"${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.AlarmContactName}"`}),
		fakeConfig: testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.id}_fake"]`,

			"alarm_contact_name": `"${alibabacloudstack_cloudmonitorservice_alarm_contacts.default.AlarmContactName}_fake"`}),
	}

	AlibabacloudstackCloudmonitorserviceAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, idsConf, alarm_contact_nameConf, allConf)
}

var existAlibabacloudstackCloudmonitorserviceAlarmContactsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#":    "1",
		"contacts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackCloudmonitorserviceAlarmContactsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#": "0",
	}
}

var AlibabacloudstackCloudmonitorserviceAlarmContactsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_cloudmonitorservice_alarm_contacts.default",
	existMapFunc: existAlibabacloudstackCloudmonitorserviceAlarmContactsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackCloudmonitorserviceAlarmContactsMapFunc,
}

func testAccCheckAlibabacloudstackCloudmonitorserviceAlarmContactsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackCloudmonitorserviceAlarmContacts%d"
}






data "alibabacloudstack_cloudmonitorservice_alarm_contacts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
