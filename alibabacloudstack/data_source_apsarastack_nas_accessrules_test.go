package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackNasAccessRuleDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := accessRuleCheckInfo.resourceId
	name := fmt.Sprintf("tf-testnaacls%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, testAccCheckAlibabacloudStackAccessRuleDataSourceConfig)
	ipConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"source_cidr_ip":    "${alibabacloudstack_nas_access_rule.default.source_cidr_ip}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"source_cidr_ip":    "${alibabacloudstack_nas_access_rule.default.source_cidr_ip}_fake",
		}),
	}
	RWAccessConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"rw_access":         "${alibabacloudstack_nas_access_rule.default.rw_access_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"rw_access":         "${alibabacloudstack_nas_access_rule.default.rw_access_type}_fake",
		}),
	}
	UserAccessConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"user_access":       "${alibabacloudstack_nas_access_rule.default.user_access_type}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"user_access":       "${alibabacloudstack_nas_access_rule.default.user_access_type}_fake",
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"ids":               []string{"${alibabacloudstack_nas_access_rule.default.access_rule_id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"ids":               []string{"${alibabacloudstack_nas_access_rule.default.access_rule_id}_fake"},
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"user_access":       "${alibabacloudstack_nas_access_rule.default.user_access_type}",
			"rw_access":         "${alibabacloudstack_nas_access_rule.default.rw_access_type}",
			"ids":               []string{"${alibabacloudstack_nas_access_rule.default.access_rule_id}"},
			"source_cidr_ip":    "${alibabacloudstack_nas_access_rule.default.source_cidr_ip}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"access_group_name": "${alibabacloudstack_nas_access_group.default.access_group_name}",
			"user_access":       "${alibabacloudstack_nas_access_rule.default.user_access_type}_fake",
			"rw_access":         "${alibabacloudstack_nas_access_rule.default.rw_access_type}_fake",
			"ids":               []string{"${alibabacloudstack_nas_access_rule.default.access_rule_id}"},
			"source_cidr_ip":    "${alibabacloudstack_nas_access_rule.default.source_cidr_ip}_fake",
		}),
	}
	accessRuleCheckInfo.dataSourceTestCheck(t, rand, ipConf, RWAccessConf, UserAccessConf, idsConf, allConf)
}

func testAccCheckAlibabacloudStackAccessRuleDataSourceConfig(name string) string {
	return  fmt.Sprintf(`
variable "name" {
        	default = "%s"
}
resource "alibabacloudstack_nas_access_group" "default" {
        	access_group_name = "${var.name}"
	        access_group_type = "Vpc"
	        description = "tf-testAccAccessGroupsdatasource"
}
resource "alibabacloudstack_nas_access_rule" "default" {
        	access_group_name = "${alibabacloudstack_nas_access_group.default.access_group_name}"
	        source_cidr_ip = "168.1.1.0/16"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
}`,  name, )
}

var existAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#":                "1",
		"rules.0.source_cidr_ip": "168.1.1.0/16",
		"rules.0.priority":       "2",
		"rules.0.access_rule_id": CHECKSET,
		"rules.0.user_access":    "no_squash",
		"rules.0.rw_access":      "RDWR",
		"ids.#":                  "1",
	}
}

var fakeAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
		"ids.#":   "0",
	}
}

var accessRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_access_rules.default",
	existMapFunc: existAccessRuleMapCheck,
	fakeMapFunc:  fakeAccessRuleMapCheck,
}
