package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackNasLifecyclePolicies_DataSource(t *testing.T) {
	rand := getAccTestRandInt(100000, 999999)
	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_file_system.default.id}"`,
		}),
		fakeConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"file_system_id": `"eeeeeeeeee"`,
		}),
	}
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_lifecycle_policy.default.id}"]`,
		}),
		fakeConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_nas_lifecycle_policy.default.id}_fake"]`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"file_system_id": `"${alibabacloudstack_nas_file_system.default.id}"`,
			"ids":            `["${alibabacloudstack_nas_lifecycle_policy.default.id}"]`,
		}),
		fakeConfig: testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand, map[string]string{
			"file_system_id": `"eeeeeeeeee"`,
			"ids":            `["${alibabacloudstack_nas_lifecycle_policy.default.id}_fake"]`,
		}),
	}

	LifecyclePoliciesCheckInfo.dataSourceTestCheck(t, rand, descriptionConf, idsConf, allConf)
}

func testAccAlibabacloudStackNasLifecyclePoliciesDataSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testacc-nas-liecycle-datasource%d"
}

%s

resource "alibabacloudstack_oss_bucket" "default" {
  bucket = "${var.name}"
  acl    = "public-read"
}

resource "alibabacloudstack_nas_lifecycle_policy" "default" {
	lifecycle_policy_name = "${var.name}"
	file_system_id        = "${alibabacloudstack_nas_file_system.default.id}"
	path                  = "/"
	recursive             = "false"
	lifecycle_rule_name   = "DEFAULT_ATIME_14"
	oss_bucket            = "${alibabacloudstack_oss_bucket.default.id}"

}
data "alibabacloudstack_nas_lifecycle_policies" "default" {
	depends_on = [
		alibabacloudstack_nas_lifecycle_policy.default
	]
	%s
}`, rand, NasCommonTestCase, strings.Join(pairs, "\n  "))
	return config
}

var existLifecyclePoliciesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lifecycle_policies.#":                     "1",
		"lifecycle_policies.0.id":                  CHECKSET,
		"lifecycle_policies.0.path":                "/",
		"lifecycle_policies.0.recursive":           "false",
		"lifecycle_policies.0.lifecycle_rule_name": "DEFAULT_ATIME_14",
		"ids.#": "1",
		"ids.0": CHECKSET,
	}
}

var fakeLifecyclePoliciesMapCheck = func(rand int) map[string]string {
	return map[string]string{
		"lifecycle_policies.#": "0",
	}
}

var LifecyclePoliciesCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_nas_lifecycle_policies.default",
	existMapFunc: existLifecyclePoliciesMapCheck,
	fakeMapFunc:  fakeLifecyclePoliciesMapCheck,
}
