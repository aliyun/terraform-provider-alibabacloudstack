package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackDnsRecursorAclsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_recursor_acls.default"

	testAcc := dataSourceAttr{
		resourceId: resourceId,
		existMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":                      "1",
				"recursor_acls.#":            "1",
				"recursor_acls.0.name":       fmt.Sprintf("tfacc%d", rand),
				"recursor_acls.0.policy":     "ALLOW",
				"recursor_acls.0.remark":     fmt.Sprintf("tfacc%d", rand),
				"recursor_acls.0.line_ids.#": "1",
			}
		},
		fakeMapFunc: func(rand int) map[string]string {
			return map[string]string{
				"ids.#":           "0",
				"recursor_acls.#": "0",
			}
		},
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_dns_recursor_acl.default.name}"`,
		}),
		fakeConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name.*"`,
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"ids": `["${alibabacloudstack_dns_recursor_acl.default.id}"]`,
		}),
		fakeConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"ids": `["fake-id"]`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_dns_recursor_acl.default.name}"`,
			"ids":        `["${alibabacloudstack_dns_recursor_acl.default.id}"]`,
		}),
		fakeConfig: AlibabacloudTestAccDnsRecursorAcldependenceNew(rand, map[string]string{
			"name_regex": `"^fake-name.*"`,
			"ids":        `["fake-id"]`,
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func AlibabacloudTestAccDnsRecursorAcldependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
    default = "tfacc%d"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_recursor_acl" "default" {
    name     = "${var.name}"
    remark   = "${var.name}"
    policy   = "ALLOW"
    line_ids = ["${alibabacloudstack_dns_line.ipv4.id}"]
}

data "alibabacloudstack_dns_recursor_acls" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n    "))
}
