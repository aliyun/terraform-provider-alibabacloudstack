package apsarastack

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
)

func TestAccApsaraStackNasAccessRuleDataSource(t *testing.T) {
	rand := acctest.RandIntRange(100000, 999999)
	ipConf := dataSourceTestAccConfig{
		existConfig: providerCommon + testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${apsarastack_nas_access_rule.default.source_cidr_ip}"`,
		}),
		fakeConfig: testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"source_cidr_ip":    `"${apsarastack_nas_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	RWAccessConf := dataSourceTestAccConfig{
		existConfig: providerCommon + testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${apsarastack_nas_access_rule.default.rw_access_type}"`,
		}),
		fakeConfig: testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"rw_access":         `"${apsarastack_nas_access_rule.default.rw_access_type}_fake"`,
		}),
	}
	UserAccessConf := dataSourceTestAccConfig{
		existConfig: providerCommon + testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${apsarastack_nas_access_rule.default.user_access_type}"`,
		}),
		fakeConfig: testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${apsarastack_nas_access_rule.default.user_access_type}_fake"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: providerCommon + testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"ids":               `["${apsarastack_nas_access_rule.default.access_rule_id}"]`,
		}),
		fakeConfig: testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"ids":               `["${apsarastack_nas_access_rule.default.access_rule_id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: providerCommon + testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${apsarastack_nas_access_rule.default.user_access_type}"`,
			"rw_access":         `"${apsarastack_nas_access_rule.default.rw_access_type}"`,
			"ids":               `["${apsarastack_nas_access_rule.default.access_rule_id}"]`,
			"source_cidr_ip":    `"${apsarastack_nas_access_rule.default.source_cidr_ip}"`,
		}),
		fakeConfig: testAccCheckApsaraStackAccessRuleDataSourceConfig(rand, map[string]string{
			"access_group_name": `"${apsarastack_nas_access_group.default.access_group_name}"`,
			"user_access":       `"${apsarastack_nas_access_rule.default.user_access_type}_fake"`,
			"rw_access":         `"${apsarastack_nas_access_rule.default.rw_access_type}_fake"`,
			"ids":               `["${apsarastack_nas_access_rule.default.access_rule_id}"]`,
			"source_cidr_ip":    `"${apsarastack_nas_access_rule.default.source_cidr_ip}_fake"`,
		}),
	}
	accessRuleCheckInfo.dataSourceTestCheck(t, rand, ipConf, RWAccessConf, UserAccessConf, idsConf, allConf)
}

func testAccCheckApsaraStackAccessRuleDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
        	default = "tf-testAccAccessGroupsdatasource-%d"
}
resource "apsarastack_nas_access_group" "default" {
        	access_group_name = "${var.name}"
	        access_group_type = "Vpc"
	        description = "tf-testAccAccessGroupsdatasource"
}
resource "apsarastack_nas_access_rule" "default" {
        	access_group_name = "${apsarastack_nas_access_group.default.access_group_name}"
	        source_cidr_ip = "168.1.1.0/16"
        	rw_access_type = "RDWR"
	        user_access_type = "no_squash"
	        priority = 2
}
data "apsarastack_nas_access_rules" "default" {
		%s
}`, rand, strings.Join(pairs, "\n  "))
	return config
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
		"ids.0":                  "1",
	}
}

var fakeAccessRuleMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"rules.#": "0",
		"ids.#":   "0",
	}
}

var accessRuleCheckInfo = dataSourceAttr{
	resourceId:   "data.apsarastack_nas_access_rules.default",
	existMapFunc: existAccessRuleMapCheck,
	fakeMapFunc:  fakeAccessRuleMapCheck,
}
