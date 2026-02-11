package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlicloudQuickBIUsersDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_quick_bi_users.default"
	name := fmt.Sprintf("tf-testAccQuickBIUser%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceQuickBIUsersConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_quick_bi_user.default.id}"},
			"enable_details": true,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_quick_bi_user.default.id}_fakeid"},
			"enable_details": true,
		}),
	}

	keywordConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alibabacloudstack_quick_bi_user.default.nick_name}",
			"enable_details": true,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"keyword":        "${alibabacloudstack_quick_bi_user.default.nick_name}_fake",
			"enable_details": true,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_quick_bi_user.default.id}"},
			"keyword":        "${alibabacloudstack_quick_bi_user.default.nick_name}",
			"enable_details": true,
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_quick_bi_user.default.id}_fake"},
			"keyword":        "${alibabacloudstack_quick_bi_user.default.nick_name}",
			"enable_details": true,
		}),
	}

	var existQuickBIUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                   "1",
			"users.#":                 "1",
			"users.0.id":              CHECKSET,
			"users.0.nick_name":       fmt.Sprintf("tf-testAccQuickBIUser%d", rand),
			"users.0.admin_user":      "false",
			"users.0.auth_admin_user": "false",
			"users.0.user_type":       "Developer",
		}
	}

	var fakeQuickBIUsersMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"users.#": "0",
		}
	}

	var quickBIUsersCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existQuickBIUsersMapFunc,
		fakeMapFunc:  fakeQuickBIUsersMapFunc,
	}
	quickBIUsersCheckInfo.dataSourceTestCheck(t, rand, idsConf, keywordConf, allConf)
}

func dataSourceQuickBIUsersConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_quick_bi_user" "default" {
  nick_name       = var.name
  account_name    = var.name
  admin_user      = "false"
  auth_admin_user = "false"
  user_type       = "Developer"
}
`, name)
}
