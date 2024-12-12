package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlicloudQuickBIUsersDataSource(t *testing.T) {
	//t.Skip()
	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_quick_bi_user.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_quick_bi_user.default.id}_fakeid"]`,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"keyword": `"${alibabacloudstack_quick_bi_user.default.nick_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"keyword": `"${alibabacloudstack_quick_bi_user.default.nick_name}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids":     `["${alibabacloudstack_quick_bi_user.default.id}"]`,
			"keyword": `"${alibabacloudstack_quick_bi_user.default.nick_name}"`,
		}),
		fakeConfig: testAccCheckAlicloudQuickBIUserDataSourceName(rand, map[string]string{
			"ids":     `["${alibabacloudstack_quick_bi_user.default.id}_fake"]`,
			"keyword": `"${alibabacloudstack_quick_bi_user.default.nick_name}"`,
		}),
	}

	var existDataAlicloudQuickBIUsersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"users.#":                 "1",
			"users.0.nick_name":       fmt.Sprintf("tf-testAccQuickBIUser%d", rand),
			"users.0.admin_user":      "false",
			"users.0.auth_admin_user": "false",
			"users.0.user_type":       "Developer",
		}
	}
	var fakeDataAlicloudQuickBIUsersSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}
	var alibabacloudstackQuickBIUserCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_quick_bi_users.default",
		existMapFunc: existDataAlicloudQuickBIUsersSourceNameMapFunc,
		fakeMapFunc:  fakeDataAlicloudQuickBIUsersSourceNameMapFunc,
	}
	alibabacloudstackQuickBIUserCheckInfo.dataSourceTestCheck(t, rand, idsConf, keywordConf, allConf)
}
func testAccCheckAlicloudQuickBIUserDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccQuickBIUser%d"
}

resource "alibabacloudstack_quick_bi_user" "default" {
  nick_name       = var.name
  account_name    = var.name
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}

data "alibabacloudstack_quick_bi_users" "default" {	
	%s
}
`, rand, strings.Join(pairs, " \n "))
	return config
}
