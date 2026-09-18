package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackFlinkNamespacesDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_flink_namespaces.default"
	testAccConfig := dataSourceTestAccConfigFunc(resourceId,
		fmt.Sprintf("tf-testacc-flink-ns-data-%d", rand),
		dataSourceFlinkNamespacesConfigDependenceNew)

	existMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"namespaces.#":      "1",
			"namespaces.0.name": CHECKSET,
		}
	}
	fakeMapFunc := func(rand int) map[string]string {
		return map[string]string{
			"namespaces.#": "0",
		}
	}
	checkInfo := dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existMapFunc,
		fakeMapFunc:  fakeMapFunc,
	}

	checkInfo.dataSourceTestCheck(t, rand,
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_flink_namespace.default.id}"},
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids": []string{"${alibabacloudstack_flink_namespace.default.id}_fake"},
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_flink_namespace.default.name}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"name_regex": "${alibabacloudstack_flink_namespace.default.name}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"owner_uid": "${alibabacloudstack_flink_namespace.default.owner_uid}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"owner_uid": "${alibabacloudstack_flink_namespace.default.owner_uid}_fake",
			}),
		},
		dataSourceTestAccConfig{
			existConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_flink_namespace.default.id}"},
				"name_regex": "${alibabacloudstack_flink_namespace.default.name}",
				"owner_uid":  "${alibabacloudstack_flink_namespace.default.owner_uid}",
			}),
			fakeConfig: testAccConfig(map[string]interface{}{
				"ids":        []string{"${alibabacloudstack_flink_namespace.default.id}_fake"},
				"name_regex": "${alibabacloudstack_flink_namespace.default.name}_fake",
				"owner_uid":  "${alibabacloudstack_flink_namespace.default.owner_uid}_fake",
			}),
		},
	)
}

func dataSourceFlinkNamespacesConfigDependenceNew(name string) string {
	return fmt.Sprintf(`
variable "name" {
	default = "%s"
}

resource "alibabacloudstack_ascm_user" "default" {
	display_name       = var.name
	mobile_nation_code = "86"
	login_name         = var.name
	login_policy_id    = "1"
	role_ids           = ["8", "9"]
	cellphone_number   = "13612345678"
	email              = "${var.name}@gmail.com"
}

resource "alibabacloudstack_flink_namespace" "default" {
	owner_uid = "${alibabacloudstack_ascm_user.default.user_uid}"
	name      = "${var.name}"
	cu        = "1"
	cpu_type  = "Intel"
}
`, name)
}
