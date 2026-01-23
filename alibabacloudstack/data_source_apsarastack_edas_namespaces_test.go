package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEDASNamespacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	resourceId := "data.alibabacloudstack_edas_namespaces.default"
	name := fmt.Sprintf("tfnsdata%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEdasNamespacesConfigDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_namespace.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_edas_namespace.default.id}_fake"},
		}),
	}

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_edas_namespace.default.namespace_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake_*",
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_namespace.default.id}"},
			"name_regex": "${alibabacloudstack_edas_namespace.default.namespace_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":        []string{"${alibabacloudstack_edas_namespace.default.id}_fake"},
			"name_regex": "fake_*",
		}),
	}

	var existEdasNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"namespaces.#":                  "1",
			"namespaces.0.id":               CHECKSET,
			"namespaces.0.namespace_id":     CHECKSET,
			"namespaces.0.description":      name,
			"namespaces.0.namespace_logical_id": CHECKSET,
			"namespaces.0.namespace_name":   name,
			"namespaces.0.user_id":          CHECKSET,
			"namespaces.0.belong_region":    CHECKSET,
		}
	}

	var fakeEdasNamespacesMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":        "0",
			"names.#":      "0",
			"namespaces.#": "0",
		}
	}

	var edasNamespacesCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEdasNamespacesMapFunc,
		fakeMapFunc:  fakeEdasNamespacesMapFunc,
	}
	edasNamespacesCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameRegexConf, allConf)
}

func dataSourceEdasNamespacesConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

%s

data "alibabacloudstack_account" "current" {}

locals {
  logical_id = "${data.alibabacloudstack_account.current.region}:${var.name}"
}

resource "alibabacloudstack_edas_namespace" "default" {
  description            = var.name
  namespace_logical_id   = substr(local.logical_id, 0, min(length(local.logical_id), 32))
  namespace_name         = var.name
}

`, name, DataZoneCommonTestCase)
}
