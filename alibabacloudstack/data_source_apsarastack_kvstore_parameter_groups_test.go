package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackKVStoreParameterGroupsDataSource(t *testing.T) {
	rand := getAccTestRandInt(10000, 99999)
	resourceId := "data.alibabacloudstack_kvstore_parameter_groups.default"
	name := fmt.Sprintf("tf-kvparamgroup%d", rand)

	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceKVStoreParameterGroupsConfigDependence)

	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_kvstore_parameter_group.default.parameter_group_name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "${alibabacloudstack_kvstore_parameter_group.default.parameter_group_name}_fake",
		}),
	}

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_kvstore_parameter_group.default.id}_fake"},
		}),
	}

	characterTypeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"character_type": "${alibabacloudstack_kvstore_parameter_group.default.character_type}",
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"character_type": "${alibabacloudstack_kvstore_parameter_group.default.character_type}_fake",
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
		}),
	}

	engineVersionConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"engine_version": "${alibabacloudstack_kvstore_parameter_group.default.engine_version}",
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"engine_version": "${alibabacloudstack_kvstore_parameter_group.default.engine_version}_fake",
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}"},
			"name_regex":     "${alibabacloudstack_kvstore_parameter_group.default.parameter_group_name}",
			"character_type": "${alibabacloudstack_kvstore_parameter_group.default.character_type}",
			"engine_version": "${alibabacloudstack_kvstore_parameter_group.default.engine_version}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids":            []string{"${alibabacloudstack_kvstore_parameter_group.default.id}_fake"},
			"name_regex":     "${alibabacloudstack_kvstore_parameter_group.default.parameter_group_name}_fake",
			"character_type": "${alibabacloudstack_kvstore_parameter_group.default.character_type}_fake",
			"engine_version": "${alibabacloudstack_kvstore_parameter_group.default.engine_version}_fake",
		}),
	}

	var existKVStoreParameterGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                         "1",
			"names.#":                       "1",
			"groups.#":                      "1",
			"groups.0.parameter_group_name": name,
			"groups.0.parameter_group_desc": CHECKSET,
			"groups.0.engine_version":       CHECKSET,
			"groups.0.character_type":       CHECKSET,
			"groups.0.parameters.#":         "3",
		}
	}

	var fakeKVStoreParameterGroupsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":    "0",
			"names.#":  "0",
			"groups.#": "0",
		}
	}

	var kvStoreParameterGroupsCheckInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existKVStoreParameterGroupsMapFunc,
		fakeMapFunc:  fakeKVStoreParameterGroupsMapFunc,
	}
	kvStoreParameterGroupsCheckInfo.dataSourceTestCheck(t, rand, nameRegexConf, idsConf, characterTypeConf, engineVersionConf, allConf)
}

func dataSourceKVStoreParameterGroupsConfigDependence(name string) string {
	return fmt.Sprintf(`
variable "name" {
  default = "%s"
}

resource "alibabacloudstack_kvstore_parameter_group" "default" {
	character_type = "logic"
	parameter_group_name = "${var.name}"
	engine_version = "7.0"
	parameter_group_desc = "${var.name}"
	parameters {
		param_name = "resp_version"
		value = "3"
	}
	parameters {
		param_name = "rt_threshold_ms"
		value = "400"
	}
	parameters {
		param_name = "#no_loose_check-whitelist-always"
		value = "yes"
	}

}

`, name)
}
