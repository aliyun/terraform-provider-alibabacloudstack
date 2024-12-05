package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccAlibabacloudStackAlibabacloudstackArmsAlertContactsDataSource(t *testing.T) {

	rand := acctest.RandIntRange(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackArmsAlertContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_arms_alert_contacts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackArmsAlertContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_arms_alert_contacts.default.id}_fake"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackArmsAlertContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_arms_alert_contacts.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackArmsAlertContactsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_arms_alert_contacts.default.id}_fake"]`,
		}),
	}

	AlibabacloudstackArmsAlertContactsCheckInfo.dataSourceTestCheck(t, rand, idsConf, allConf)
}

var existAlibabacloudstackArmsAlertContactsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#":    "1",
		"contacts.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackArmsAlertContactsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"contacts.#": "0",
	}
}

var AlibabacloudstackArmsAlertContactsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_arms_alert_contacts.default",
	existMapFunc: existAlibabacloudstackArmsAlertContactsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackArmsAlertContactsMapFunc,
}

func testAccCheckAlibabacloudstackArmsAlertContactsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackArmsAlertContacts%d"
}






data "alibabacloudstack_arms_alert_contacts" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
