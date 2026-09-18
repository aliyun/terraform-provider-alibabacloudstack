package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCrEEArtifactLifecycleRulesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000, 9999)
	resourceId := "data.alibabacloudstack_cr_ee_attestor_lifecycle_rules.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-cree-rule-%d", rand),
		dataSourceCrEEAttestorLifecycleRulesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"ids":         []string{"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id": "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"ids":         []string{"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"namespace_regex": "${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"namespace_regex": "${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}_fake",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"ids":             []string{"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}"},
			"namespace_regex": "${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"instance_id":     "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}",
			"ids":             []string{"${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.id}_fake"},
			"namespace_regex": "${alibabacloudstack_cr_ee_attestor_lifecycle_rule.default.namespace_name}_fake",
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
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	dataSourceAttr.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceCrEEAttestorLifecycleRulesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
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
    instance_id = "${data.alibabacloudstack_cr_ee_instances.default.instances.0.id}"
	scope = "NAMESPACE"
	retention_tag_count = "30"
	tag_regexp = "release-v.*"
	enable_delete_tag = "true"
	namespace_name = "${alibabacloudstack_cr_ee_namespace.default.name}"
	recent_pull_keep = 30
	recent_push_keep = 20
}
`, name)
}
