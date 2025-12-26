package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"
)

func TestAccAlibabacloudStackFlinkNamespacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"ids": `["${alibabacloudstack_flink_namespace.default.id}"]`,
		}),
		fakeConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"ids": `["${alibabacloudstack_flink_namespace.default.id}_fake"]`,
		}),
	}

	nameregex_idConf := dataSourceTestAccConfig{
		existConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_flink_namespace.default.name}"`,
		}),
		fakeConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"name_regex": `"${alibabacloudstack_flink_namespace.default.name}_fake"`,
		}),
	}

	oweruid_idConf := dataSourceTestAccConfig{
		existConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"owner_uid": `"${alibabacloudstack_flink_namespace.default.owner_uid}"`,
		}),
		fakeConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"owner_uid": `"${alibabacloudstack_flink_namespace.default.owner_uid}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"ids":        `["${alibabacloudstack_flink_namespace.default.id}"]`,
			"name_regex": `"${alibabacloudstack_flink_namespace.default.name}"`,
			"owner_uid":  `"${alibabacloudstack_flink_namespace.default.owner_uid}"`,
		}),
		fakeConfig: dataSourceFlinkNamespacesConfigDependence(rand, map[string]string{
			"ids":        `["${alibabacloudstack_flink_namespace.default.id}_fake"]`,
			"name_regex": `"${alibabacloudstack_flink_namespace.default.name}_fake"`,
			"owner_uid":  `"${alibabacloudstack_flink_namespace.default.owner_uid}_fake"`,
		}),
	}

	AlibabacloudstacFlinkNamesapcesDataCheckInfo.dataSourceTestCheck(t, rand, idsConf, nameregex_idConf, oweruid_idConf, allConf)
}

var existAlibabacloudstackFlinkNamesapcesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"namespaces.#":      "1",
		"namespaces.0.name": CHECKSET,
	}
}

var fakeAlibabacloudstackFlinkNamesapcesDataMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"namespaces.#": "0",
	}
}

var AlibabacloudstacFlinkNamesapcesDataCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_flink_namespaces.default",
	existMapFunc: existAlibabacloudstackFlinkNamesapcesDataMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackFlinkNamesapcesDataMapFunc,
}

func dataSourceFlinkNamespacesConfigDependence(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	return fmt.Sprintf(`
variable name{
 default = "tf-testacc-flink-ns-data-%d"
}

resource "alibabacloudstack_ascm_user" "default" {
  display_name = var.name
  mobile_nation_code = "86"
  login_name = var.name
  login_policy_id = "1"
  role_ids = [
               "8",
               "9"
             ]
  cellphone_number = "13612345678"
  email = "${var.name}@gmail.com"
}

resource "alibabacloudstack_flink_namespace" "default" {
  owner_uid = "${alibabacloudstack_ascm_user.default.user_uid}"
  name = "${var.name}"
  cu = "1"
  cpu_type = "Intel"
}

data "alibabacloudstack_flink_namespaces" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
}
