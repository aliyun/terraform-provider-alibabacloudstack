package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPrometheusV2ContactsDataSource(t *testing.T) {
	randInt := getAccTestRandInt(10000, 99999)
	testAcc := dataSourceAttr{
		resourceId: "data.alibabacloudstack_prometheus_v2_contacts.default",
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

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_contact.default.username}"`,
		}),
		fakeConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"name_regex": `"^fake-name$"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"ids": `["${alibabacloudstack_prometheus_v2_contact.default.id}"]`,
		}),
		fakeConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_contact.default.username}"`,
			"ids":        `["${alibabacloudstack_prometheus_v2_contact.default.id}"]`,
		}),
		fakeConfig: resourcePrometheusV2ContactBasicDependenceNew(randInt, map[string]string{
			"name_regex": `"${alibabacloudstack_prometheus_v2_contact.default.username}"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, randInt, nameRegexConf, idsConf, allConf)
}

// resourcePrometheusV2ContactBasicDependenceNew creates a basic dependency template for prometheus v2 contact tests.
func resourcePrometheusV2ContactBasicDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tfacc_prometheus%d"
}
resource "alibabacloudstack_prometheus_v2_contact" "default" {
  username = "tfacc-${var.name}"
  mobile   = "13812345678"
  mail     = "test@example.com"
}

data "alibabacloudstack_prometheus_v2_contacts" "default" {
  %s
}
`, rand, strings.Join(pairs, "\n  "))
}
