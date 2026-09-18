package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackDnsRecursorAclsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_dns_recursor_acls.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tfacc%d", rand),
		dataSourceDnsRecursorAclsConfigDependence)

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
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_recursor_acl.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^fake-name.*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_dns_recursor_acl.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"fake-id"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_dns_recursor_acl.default.name}",
			"ids":        []string{"${alibabacloudstack_dns_recursor_acl.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^fake-name.*",
			"ids":        []string{"fake-id"},
		}),
	}

	testAcc.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, allConf)
}

func dataSourceDnsRecursorAclsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
    default = "%s"
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
`, name)
}
