package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsCommandsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ecs_commands.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsCommandsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsCommandsDependence)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_command.default.id}"},
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_command.default.id}-fake"},
		}),
	}
	nameRegexConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name_regex": "^${alibabacloudstack_ecs_command.default.name}$",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name_regex": "fake-name",
		}),
	}
	nameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"name": "${alibabacloudstack_ecs_command.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"name": "fake-name",
		}),
	}
	typeConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_command.default.id}"},
			"type"              : "RunShellScript",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"ids": []string{"${alibabacloudstack_ecs_command.default.id}"},
			"type":"RunBatScript",
		}),
	}
	var existEcsCommandsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"ids.0":                       CHECKSET,
			"names.#":                     "1",
			"names.0":                     name,
			"commands.#":                  "1",
			"commands.0.id":               CHECKSET,
			"commands.0.command_content":  "bHMK",
			"commands.0.command_id":       CHECKSET,
			"commands.0.description":      name,
			"commands.0.enable_parameter": "false",
			"commands.0.name":             name,
			"commands.0.type":             "RunShellScript",
			"commands.0.working_dir":      "/root",
		}
	}

	var fakeEcsCommandsMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":      "0",
			"names.#":    "0",
			"commands.#": "0",
		}
	}

	var EcsCommandsInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsCommandsMapFunc,
		fakeMapFunc:  fakeEcsCommandsMapFunc,
	}

	EcsCommandsInfo.dataSourceTestCheck(t, 0, idsConf, nameRegexConf, nameConf, typeConf)
}

func dataSourceEcsCommandsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}
	resource "alibabacloudstack_ecs_command" "default" {
		name              = var.name
		command_content   = "bHMK"
		description       = var.name
		type              = "RunShellScript"
		working_dir       = "/root"
	}`, name)
}
