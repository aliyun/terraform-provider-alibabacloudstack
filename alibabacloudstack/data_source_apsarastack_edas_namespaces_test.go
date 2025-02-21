package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

)

func TestAccAlibabacloudStackEdasNamespacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(100, 999)
	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_edas_namespace.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"ids": `["${alibabacloudstack_edas_namespace.default.id}_fake"]`,
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_edas_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_edas_namespace.default.namespace_name}_fake"`,
		}),
	}
	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_edas_namespace.default.id}"]`,
			"name_regex": `"${alibabacloudstack_edas_namespace.default.namespace_name}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand, map[string]string{
			"ids":        `["${alibabacloudstack_edas_namespace.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_edas_namespace.default.namespace_name}_fake"`,
		}),
	}
	var existAlibabacloudStackEdasNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                             "1",
			"names.#":                           "1",
			"namespaces.#":                      "1",
// 			"namespaces.0.debug_enable":         "false",
			"namespaces.0.description":          fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.namespace_logical_id": fmt.Sprintf("%s:tftest%d", defaultRegionToTest, rand),
			"namespaces.0.namespace_name":       fmt.Sprintf("tf-testAccNamespace-%d", rand),
			"namespaces.0.user_id":              CHECKSET,
			"namespaces.0.belong_region":        CHECKSET,
		}
	}
	var fakeAlibabacloudStackEdasNamespacesDataSourceNameMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":   "0",
			"names.#": "0",
		}
	}
	var alibabacloudstackEdasNamespacesCheckInfo = dataSourceAttr{
		resourceId:   "data.alibabacloudstack_edas_namespaces.default",
		existMapFunc: existAlibabacloudStackEdasNamespacesDataSourceNameMapFunc,
		fakeMapFunc:  fakeAlibabacloudStackEdasNamespacesDataSourceNameMapFunc,
	}

	preCheck := func() {
		testAccPreCheck(t)
	}
	alibabacloudstackEdasNamespacesCheckInfo.dataSourceTestCheckWithPreCheck(t, rand, preCheck, idsConf, nameRegexConf, allConf)
}
func testAccCheckAlibabacloudStackEdasNamespacesDataSourceName(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}

	config := fmt.Sprintf(`

variable "name" {	
	default = "tf-testAccNamespace-%d"
}

variable "logical_id" {
  default = "%s:tftest%d"
}

resource "alibabacloudstack_edas_namespace" "default" {
	//debug_enable = false
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

data "alibabacloudstack_edas_namespaces" "default" {	
	%s
}
`, rand, defaultRegionToTest, rand, strings.Join(pairs, " \n "))
	return config
}
