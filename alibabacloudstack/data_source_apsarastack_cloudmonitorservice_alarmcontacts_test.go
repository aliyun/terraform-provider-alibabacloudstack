package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudstackCmsAlarmContacts_basic(t *testing.T) {
	testAccPreCheckWithAPIIsNotSupport(t)
	rand := getAccTestRandInt(10000,20000)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm_contact.default.id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm_contact.default.id}_fake"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cms_alarm_contact.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm_contact.default.id}"`,
			"ids":        `["${alibabacloudstack_cms_alarm_contact.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_cms_alarm_contact.default.id}_fake"`,
			"ids":        `["${alibabacloudstack_cms_alarm_contact.default.id}_fake"]`,
		}),
	}

	var existcmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"contacts.#":                    "1",
			"contacts.0.id":                 CHECKSET,
			"contacts.0.alarm_contact_name": CHECKSET,
		}
	}

	var fakecmsAlarmContactsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}

	var cmsAlarmContactsCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_cms_alarm_contacts.default",
		existMapFunc: existcmsAlarmContactsMapFunc,
		fakeMapFunc:  fakecmsAlarmContactsMapFunc,
	}

	cmsAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
		variable "name" {
			default = "tf-testAccCmsAlarmContactBisic-%d"
		}
		resource "alibabacloudstack_cms_alarm_contact" "default" {
			alarm_contact_name = var.name
		    describe           = "For Test"
		    channels_mail      = "hello.uuuu@aaa.com"
			lifecycle {
				ignore_changes = [channels_mail]
  			}	
		}

		data "alibabacloudstack_cms_alarm_contacts" "default" {
		  %s
		}
`, rand, strings.Join(pairs, "\n  "))
	return config
}
