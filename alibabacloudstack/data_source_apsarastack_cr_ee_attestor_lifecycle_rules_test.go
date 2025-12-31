package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackCreeAttestorLifecycleRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	idsConf := dataSourceTestAccConfig{
		existConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id": `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"ids":         `["${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}"]`,
		}),
		fakeConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id": `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"ids":         `["${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}_fake"]`,
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":     `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"namespace_regex": `"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}"`,
		}),
		fakeConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":     `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"namespace_regex": `"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}_fake"`,
		}),
	}

	enableDeleteTagConf := dataSourceTestAccConfig{
		existConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":       `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"enable_delete_tag": "true",
		}),
		fakeConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":       `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"enable_delete_tag": "false",
		}),
	}

	enableDeleteUntaggedManifestConf := dataSourceTestAccConfig{
		existConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":                     `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"enable_delete_untagged_manifest": "false",
		}),
		fakeConfig: resourceCrEEArtifactLifecycleRuleDependenceNew(rand, map[string]string{
			"instance_id":                     `"${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"`,
			"enable_delete_untagged_manifest": "true",
		}),
	}

	var existMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"rules.#":                   "1",
			"rules.0.namespace_name":    fmt.Sprintf("tf-testacc-cree-rule-%d", rand),
			"rules.0.enable_delete_tag": "true",
		}
	}

	var fakeMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
			"rules.#": "0",
		}
	}

	dataSourceAttr := dataSourceAttr{
		resourceId:   "data.alibabacloudstack_cr_ee_attestor_lifecycle_rules.default",
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	dataSourceAttr.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, enableDeleteTagConf, enableDeleteUntaggedManifestConf)
}

// Generate dependency resource template for CreeAttestorLifecycleRules
func resourceCrEEArtifactLifecycleRuleDependenceNew(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable "name" {
  default = "tf-testacc-cree-rule-%d"
}

data "alibabacloudstack_cr_ee_instances" "default" {
}

resource "alibabacloudstack_cr_ee_namespace" "default" {
	instance_id = "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"
	name = "${var.name}"
	auto_create	= false
	default_visibility = "PRIVATE"
}

resource "alibabacloudstack_cr_ee_attestor_lifecycle_rule" "default" {
	scope = "NAMESPACE"
	retention_tag_count = "30"
	tag_regexp = "release-v.*"
	enable_delete_tag = "true"
	namespace_name = "${alibabacloudstack_cr_ee_namespace.default.name}"
}

data "alibabacloudstack_cr_ee_attestor_lifecycle_rules" "default" {
    %s
}
`, rand, strings.Join(pairs, "\n   "))
}
