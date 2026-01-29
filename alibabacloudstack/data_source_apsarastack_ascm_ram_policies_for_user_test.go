package alibabacloudstack

import (
	"testing"
)

func TestAccAlibabacloudStackAscmRamPoliciesForUserDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ascm_ram_policies_for_user.default"

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, "", dataSourceAscmRamPoliciesForUserConfigDependence)

	loginNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"login_name": "${data.alibabacloudstack_account.current.login_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"login_name": "nonexistent_user",
		}),
	}

	namesConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"login_name": "${data.alibabacloudstack_account.current.login_name}",
			"names":      []string{"${data.alibabacloudstack_ascm_ram_policies_for_user.anyone.names.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"login_name": "${data.alibabacloudstack_account.current.login_name}",
			"names":      []string{"fake_policy_name"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"login_name": "${data.alibabacloudstack_account.current.login_name}",
			"names":      []string{"${data.alibabacloudstack_ascm_ram_policies_for_user.anyone.names.0}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"login_name": "nonexistent_user",
			"names":      []string{"fake_policy_name"},
		}),
	}

	var existAscmRamPoliciesForUserMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":                    CHECKSET,
			"policies.#":                 CHECKSET,
			"policies.0.policy_name":     CHECKSET,
			"policies.0.policy_type":     CHECKSET,
			"policies.0.attach_date":     CHECKSET,
			"policies.0.default_version": CHECKSET,
			"policies.0.policy_document": CHECKSET,
		}
	}

	var fakeAscmRamPoliciesForUserMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"names.#":    "0",
			"policies.#": "0",
		}
	}

	var ascmRamPoliciesForUserCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmRamPoliciesForUserMapFunc,
		fakeMapFunc:  fakeAscmRamPoliciesForUserMapFunc,
	}
	ascmRamPoliciesForUserCheckInfo.dataSourceTestCheck(t, rand, loginNameConf, namesConf, allConf)
}

func dataSourceAscmRamPoliciesForUserConfigDependence(name string) string {
	return `
	
data "alibabacloudstack_account" "current" {
}

data "alibabacloudstack_ascm_ram_policies_for_user" "anyone" {
  login_name = data.alibabacloudstack_account.current.login_name
}
`
}
