package alibabacloudstack

import (
	"fmt"
	"testing"
)

func TestAccAlibabacloudStackEcsInvocationsDataSource(t *testing.T) {
	resourceId := "data.alibabacloudstack_ecs_invocations.default"
	rand := getAccTestRandInt(10000, 99999)
	name := fmt.Sprintf("tf-testAccEcsInvocationsTest%d", rand)
	testAccConfig := dataSourceTestAccConfigFunc(resourceId, name, dataSourceEcsInvocationsDependence)

	InvocationIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"invocation_id": "${alibabacloudstack_ecs_invocation.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"invocation_id": "${alibabacloudstack_ecs_invocation.default.id}_fake",
		}),
	}
	commandNameConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"command_name": "${alibabacloudstack_ecs_command.default.name}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"command_name": "${alibabacloudstack_ecs_command.default.name}_fake",
		}),
	}
	commandIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"command_id": "${alibabacloudstack_ecs_command.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"command_id": "${alibabacloudstack_ecs_command.default.id}_fake",
		}),
	}

	// 平台可能会自行启动一些Agent安装的命令，因此需要组合查找
	instanceIdConf := dataSourceTestAccConfig{
		existConfig: testAccConfig(map[string]interface{}{
			"command_id": "${alibabacloudstack_ecs_command.default.id}",
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}",
		}),
		fakeConfig: testAccConfig(map[string]interface{}{
			"command_id": "${alibabacloudstack_ecs_command.default.id}",
			"instance_id": "${alibabacloudstack_ecs_instance.default.id}_fake",
		}),
	}
	var existEcsInvocationMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":                       "1",
			"invocations.#":               "1",
			"invocations.0.command_id":    CHECKSET,
			"invocations.0.invocation_id": CHECKSET,
			"invocations.0.repeat_mode":   CHECKSET,
			"invocations.0.timed":         CHECKSET,
			"invocations.0.username":      CHECKSET,
		}
	}

	var fakeEcsInvocationMapFunc = func(rand int) map[string]string {
		return map[string]string{
			"ids.#":         "0",
			"names.#":       "0",
			"invocations.#": "0",
		}
	}

	var EcsInvocationInfo = dataSourceAttr{
		resourceId:   resourceId,
		existMapFunc: existEcsInvocationMapFunc,
		fakeMapFunc:  fakeEcsInvocationMapFunc,
	}

	EcsInvocationInfo.dataSourceTestCheck(t, 0, InvocationIdConf, commandNameConf, commandIdConf, instanceIdConf)
}

func dataSourceEcsInvocationsDependence(name string) string {
	return fmt.Sprintf(`
	variable "name" {
		default = "%s"
	}

	resource alibabacloudstack_ecs_command	"default" {
		description = "command description"
		command_content = "pwd"
		type = "RunShellScript"
		name = "${var.name}"
	}

	%s

	resource alibabacloudstack_ecs_invocation "default" {
		command_id = "${alibabacloudstack_ecs_command.default.id}"
		instance_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
		username = "root"
		repeat_mode = "Once"
	}

	`, name, ECSInstanceCommonTestCase)
}
