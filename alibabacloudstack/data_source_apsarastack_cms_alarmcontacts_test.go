package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudstackCmsAlarmContacts_basic(t *testing.T) {
	resourceId := "data.alibabacloudstack_cms_alarm_contacts.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccCmsAlarmContactBisic%d", rand)
	testAccConfig := resourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig)
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact.default.id}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cms_alarm_contact.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_cms_alarm_contact.default.id}_fake"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact.default.id}",
			"ids":        []string{"${alibabacloudstack_cms_alarm_contact.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cms_alarm_contact.default.id}_fake",
			"ids":        []string{"${alibabacloudstack_cms_alarm_contact.default.id}_fake"},
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
		PreCheck: func(){
			testAccPreCheckWithAPIIsNotSupport(t)
		},
	}

	cmsAlarmContactsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func testAccCheckAlibabacloudstackCmsAlarmContactsDataSourceConfig(name string) string {
	return fmt.Sprintf(`
		variable "name" {
			default = "%s"
		}
		resource "alibabacloudstack_cms_alarm_contact" "default" {
			alarm_contact_name = var.name
		    describe           = "For Test"
		    channels_mail      = "hello.uuuu@aaa.com"
			lifecycle {
				ignore_changes = [channels_mail]
  			}	
		}

`, name)
}
