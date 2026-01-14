package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCRReposDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_cr_repos.default"
	name := fmt.Sprintf("acc-cr-repos-%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCRReposConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cr_repo.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	namespaceConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace": "${alibabacloudstack_cr_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace": "fake-namespace",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"namespace":  "${alibabacloudstack_cr_namespace.default.name}",
			"name_regex": "${alibabacloudstack_cr_repo.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"namespace":  "fake-namespace",
			"name_regex": "fake_*",
		}),
	}

	var existCRReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                     "1",
			"names.#":                   "1",
			"repos.#":                   "1",
			"repos.0.namespace":         name,
			"repos.0.name":              name,
			"repos.0.summary":           "OLD",
			"repos.0.repo_type":         "PUBLIC",
		}
	}

	var fakeCRReposMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":     "0",
			"names.#":   "0",
			"repos.#":   "0",
		}
	}

	var crReposCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCRReposMapFunc,
		fakeMapFunc:  fakeCRReposMapFunc,
	}
	crReposCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, namespaceConf, allConf)
}

func dataSourceCRReposConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_cr_namespace" "default" {
	name               = var.name
	auto_create        = false
	default_visibility = "PUBLIC"
}

resource "alibabacloudstack_cr_repo" "default" {
	namespace = alibabacloudstack_cr_namespace.default.name
	name      = var.name
	summary   = "OLD"
	repo_type = "PUBLIC"
	detail    = "OLD"
}

`, name)
}
