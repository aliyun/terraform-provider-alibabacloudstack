package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackPolardbxAccountsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	name := fmt.Sprintf("tf_acc_polardbx_account_%d", rand)
	instanceIdRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_polardbx_account.default.instance_id}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"drdsusrztw1cfake"`,
		}),
	}
	namesRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_polardbx_account.default.instance_id}"`,
			"names":       `["${var.name}",]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackPolardbxAccountsSourceConfig(name, map[string]string{
			"instance_id": `"${alibabacloudstack_polardbx_account.default.instance_id}"`,
			"names":       `["${var.name}_fake",]`,
		}),
	}

	var exisMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#":              CHECKSET,
			"accounts.0.account_name": CHECKSET,
			"accounts.0.instance_id":  CHECKSET,
		}
	}
	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"accounts.#": "0",
			"ids.#":      "0",
			"names.#":    "0",
		}
	}

	var CheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_polardbx_accounts.default",
		existMapFunc: exisMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}
	preCheck := func() {
	}
	CheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, instanceIdRegexConf, namesRegexConf,)
}

func testAccCheckAlibabacloudStackPolardbxAccountsSourceConfig(name string, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "%s"
}

%s

%s

%s


resource "alibabacloudstack_polardbx_account" "default" {
	instance_id  = local.polardbx_instance.id
	account_name = var.name
	password     = "${random_password.password.0.result}"
	description  = var.name
}
	
data "alibabacloudstack_polardbx_accounts" "default" {
  %s
}
`, name, RandomPasswordTestCase(12,1), VSwitchCommonTestCase, PolardbxReadOrCreateCommonTestCase(), strings.Join(pairs, "\n  "))
	return config
}
