package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackCRNamespacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(1000000, 9999999)
	resourceId := "data.alibabacloudstack_cr_namespaces.default"
	name := fmt.Sprintf("testacc-cr-namespace%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceCRNamespacesConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_cr_namespace.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-nonexistent-namespace",
		}),
	}

	// Since the schema doesn't have an 'ids' filter parameter, we only test name_regex scenarios
	// But we can still create a basic config without filters
	basicConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{}),
		// fakeConfig not applicable since empty config returns all namespaces
	}

	var existCRNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "1",
			"names.#":      "1",
			"namespaces.#": "1",
			"namespaces.0.name":                name,
			"namespaces.0.auto_create":         "false",
			"namespaces.0.default_visibility":  "PUBLIC",
		}
	}

	var fakeCRNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}

	var crNamespacesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existCRNamespacesMapFunc,
		fakeMapFunc:  fakeCRNamespacesMapFunc,
	}
	crNamespacesCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, basicConf)
}

func dataSourceCRNamespacesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_cr_namespace" "default" {
  name               = var.name
  auto_create        = false
  default_visibility = "PUBLIC"
}

`, name)
}
