package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackAscmLogonPoliciesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 20000)
	resourceId := "data.alibabacloudstack_ascm_logon_policies.default"
	name := fmt.Sprintf("tf-testacc-ascmlogonpolicy%v", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceAscmLogonPoliciesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_ascm_logon_policy.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"description": "${alibabacloudstack_ascm_logon_policy.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"description": "fake_*",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ascm_logon_policy.default.policy_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"-1"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"${alibabacloudstack_ascm_logon_policy.default.policy_id}"},
			"name_regex":  "${alibabacloudstack_ascm_logon_policy.default.name}",
			"description": "${alibabacloudstack_ascm_logon_policy.default.description}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":         []string{"-1"},
			"name_regex":  "${alibabacloudstack_ascm_logon_policy.default.name}_fake",
			"description": "fake_*",
		}),
	}

	var existAscmLogonPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                      "1",
			"policies.#":                 "1",
			"policies.0.id":              CHECKSET,
			"policies.0.name":            name,
			"policies.0.description":     name,
			"policies.0.rule":            "ALLOW",
			"policies.0.login_policy_id": CHECKSET,
		}
	}

	var fakeAscmLogonPoliciesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"policies.#": "0",
		}
	}

	var ascmLogonPoliciesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existAscmLogonPoliciesMapFunc,
		fakeMapFunc:  fakeAscmLogonPoliciesMapFunc,
	}
	ascmLogonPoliciesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, descriptionConf, allConf)
}

func dataSourceAscmLogonPoliciesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_ascm_logon_policy" "default" {
  name        = var.name
  description = var.name
  rule        = "ALLOW"
}
`, name)
}
