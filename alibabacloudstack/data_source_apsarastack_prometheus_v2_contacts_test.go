package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2ContactsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_prometheus_v2_contacts.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc_prometheus%d", rand),
		dataSourcePrometheusV2ContactsConfigDependence)

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"contacts.#":          "1",
				"contacts.0.username": CHECKSET,
				"contacts.0.mobile":   CHECKSET,
				"contacts.0.mail":     CHECKSET,
				"contacts.0.groups.#": "0",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"contacts.#": "0",
			}
		},
	}

	testAcc.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_contact.default.username}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "^fake-name$",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_prometheus_v2_contact.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"fake-id"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_contact.default.username}",
				"ids":        []string{"${alibabacloudstack_prometheus_v2_contact.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_prometheus_v2_contact.default.username}",
				"ids":        []string{"fake-id"},
			}),
		},
	)
}

func dataSourcePrometheusV2ContactsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_prometheus_v2_contact" "default" {
	username = "tfacc-${var.name}"
	mobile   = "13812345678"
	mail     = "test@example.com"
}
`, name)
}
