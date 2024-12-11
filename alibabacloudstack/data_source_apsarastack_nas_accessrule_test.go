package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackAlibabacloudstackNasAccessRulesDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_rules.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_rules.default.id}_fake"]`,
		}),
	}

	access_group_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_nas_access_rules.default.id}"]`,
			"access_group_name": `"${alibabacloudstack_nas_access_rules.default.AccessGroupName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":               `["${alibabacloudstack_nas_access_rules.default.id}_fake"]`,
			"access_group_name": `"${alibabacloudstack_nas_access_rules.default.AccessGroupName}_fake"`,
		}),
	}

	access_rule_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_nas_access_rules.default.id}"]`,
			"access_rule_id": `"${alibabacloudstack_nas_access_rules.default.AccessRuleId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":            `["${alibabacloudstack_nas_access_rules.default.id}_fake"]`,
			"access_rule_id": `"${alibabacloudstack_nas_access_rules.default.AccessRuleId}_fake"`,
		}),
	}

	file_system_typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_nas_access_rules.default.id}"]`,
			"file_system_type": `"${alibabacloudstack_nas_access_rules.default.FileSystemType}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids":              `["${alibabacloudstack_nas_access_rules.default.id}_fake"]`,
			"file_system_type": `"${alibabacloudstack_nas_access_rules.default.FileSystemType}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_rules.default.id}"]`,

			"access_group_name": `"${alibabacloudstack_nas_access_rules.default.AccessGroupName}"`,
			"access_rule_id":    `"${alibabacloudstack_nas_access_rules.default.AccessRuleId}"`,
			"file_system_type":  `"${alibabacloudstack_nas_access_rules.default.FileSystemType}"`}),
		fakeConfig: testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_access_rules.default.id}_fake"]`,

			"access_group_name": `"${alibabacloudstack_nas_access_rules.default.AccessGroupName}_fake"`,
			"access_rule_id":    `"${alibabacloudstack_nas_access_rules.default.AccessRuleId}_fake"`,
			"file_system_type":  `"${alibabacloudstack_nas_access_rules.default.FileSystemType}_fake"`}),
	}

	AlibabacloudstackNasAccessRulesCheckInfo.dataSourceTestCheck(t, rand, idsConf, access_group_nameConf, access_rule_idConf, file_system_typeConf, allConf)
}

var existAlibabacloudstackNasAccessRulesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":    "1",
		"rules.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackNasAccessRulesMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
	}
}

var AlibabacloudstackNasAccessRulesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_access_rules.default",
	existMapFunc: existAlibabacloudstackNasAccessRulesMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackNasAccessRulesMapFunc,
}

func testAccCheckAlibabacloudstackNasAccessRulesSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackNasAccessRules%d"
}






data "alibabacloudstack_nas_access_rules" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
