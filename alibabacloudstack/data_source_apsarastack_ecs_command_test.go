package alibabacloudstack

import (
	"fmt"
	"strings"
	"testing"

	
)

func TestAccAlibabacloudStackEcsCommandsDataSource(t *testing.T) {

	rand := getAccTestRandInt(10000, 99999)

	idsConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
		}),
	}

	command_idConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_commands.default.id}"]`,
			"command_id": `"${alibabacloudstack_ecs_commands.default.CommandId}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":        `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
			"command_id": `"${alibabacloudstack_ecs_commands.default.CommandId}_fake"`,
		}),
	}

	command_nameConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_commands.default.id}"]`,
			"command_name": `"${alibabacloudstack_ecs_commands.default.CommandName}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":          `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
			"command_name": `"${alibabacloudstack_ecs_commands.default.CommandName}_fake"`,
		}),
	}

	descriptionConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_commands.default.id}"]`,
			"description": `"${alibabacloudstack_ecs_commands.default.Description}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":         `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
			"description": `"${alibabacloudstack_ecs_commands.default.Description}_fake"`,
		}),
	}

	latestConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}"]`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
		}),
	}

	typeConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_ecs_commands.default.id}"]`,
			"type": `"${alibabacloudstack_ecs_commands.default.Type}"`,
		}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids":  `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,
			"type": `"${alibabacloudstack_ecs_commands.default.Type}_fake"`,
		}),
	}

	allConf := dataSourceTestAccConfig{
		existConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}"]`,

			"command_id":   `"${alibabacloudstack_ecs_commands.default.CommandId}"`,
			"command_name": `"${alibabacloudstack_ecs_commands.default.CommandName}"`,
			"description":  `"${alibabacloudstack_ecs_commands.default.Description}"`,
			"type":         `"${alibabacloudstack_ecs_commands.default.Type}"`}),
		fakeConfig: testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand, map[string]string{
			"ids": `["${alibabacloudstack_ecs_commands.default.id}_fake"]`,

			"command_id":   `"${alibabacloudstack_ecs_commands.default.CommandId}_fake"`,
			"command_name": `"${alibabacloudstack_ecs_commands.default.CommandName}_fake"`,
			"description":  `"${alibabacloudstack_ecs_commands.default.Description}_fake"`,
			"type":         `"${alibabacloudstack_ecs_commands.default.Type}_fake"`}),
	}

	AlibabacloudstackEcsCommandsCheckInfo.dataSourceTestCheck(t, rand, idsConf, command_idConf, command_nameConf, descriptionConf, latestConf, typeConf, allConf)
}

var existAlibabacloudstackEcsCommandsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"commands.#":    "1",
		"commands.0.id": CHECKSET,
	}
}

var fakeAlibabacloudstackEcsCommandsMapFunc = func(rand int) map[string]string {
	return map[string]string{
		"commands.#": "0",
	}
}

var AlibabacloudstackEcsCommandsCheckInfo = dataSourceAttr{
	resourceId:   "data.alibabacloudstack_ecs_commands.default",
	existMapFunc: existAlibabacloudstackEcsCommandsMapFunc,
	fakeMapFunc:  fakeAlibabacloudstackEcsCommandsMapFunc,
}

func testAccCheckAlibabacloudstackEcsCommandsSourceConfig(rand int, attrMap map[string]string) string {
	var pairs []string
	for k, v := range attrMap {
		pairs = append(pairs, k+" = "+v)
	}
	config := fmt.Sprintf(`
variable "name" {
	default = "tf-testAlibabacloudstackEcsCommands%d"
}






data "alibabacloudstack_ecs_commands" "default" {
%s
}
`, rand, strings.Join(pairs, "\n   "))
	return config
}
