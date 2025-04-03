---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_commands"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-commands"
description: |- 
  查询云服务器命令
---

# alibabacloudstack_ecs_commands

根据指定过滤条件列出当前凭证权限可以访问的云服务器命令列表。

## 示例用法

### 基础用法：

```terraform
data "alibabacloudstack_ecs_commands" "example" {
  ids        = ["E2RY53-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_command_id" {
  value = data.alibabacloudstack_ecs_commands.example.commands.0.command_id
}
```

高级用法带过滤器：

```terraform
resource "alibabacloudstack_ecs_command" "default" {
    name              = "tf-testAccEcsCommandsTest51954"
    command_content   = "bHMK"
    description       = "For Terraform Test"
    type              = "RunShellScript"
    working_dir       = "/root"
}

data "alibabacloudstack_ecs_commands" "default" {
  ids               = [alibabacloudstack_ecs_command.default.id]
  content_encoding  = "Base64"
  description       = "Test command"
  name              = "test-command"
  type              = "RunShellScript"
}

output "command_ids" {
  value = data.alibabacloudstack_ecs_commands.default.ids
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - （可选，变更时重建）用于通过命令名称筛选结果的正则表达式字符串。
* `ids` - （可选，变更时重建）用于筛选结果的命令ID列表。
* `content_encoding` - （可选，变更时重建）命令内容（CommandContent）的编码方式。取值范围：
  * `PlainText`：不编码，采用明文传输。
  * `Base64`：Base64编码。
  默认值：`Base64`。错填该取值会当作`Base64`处理。
* `description` - （可选，变更时重建）命令的描述。
* `name` - （可选，变更时重建）命令的名称。
* `command_provider` - （可选，变更时重建）命令的提供者。
* `type` - （可选，变更时重建）命令的类型。取值范围：
  * `RunBatScript`
  * `RunPowerShellScript`
  * `RunShellScript`

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 命令名称列表。
* `commands` - ECS命令列表。每个元素包含以下属性：
  * `command_content` - 命令内容，以 Base64 编码后传输。
  * `command_id` - 命令 ID。
  * `description` - 命令的描述。
  * `enable_parameter` - 是否在命令中使用自定义参数。
  * `id` - 命令的 ID（与`command_id`相同）。
  * `name` - 命令的名称。
  * `parameter_names` - 在创建命令时从命令内容中解析出的自定义参数名称列表。
  * `timeout` - 命令在ECS实例上运行的超时时间（以秒为单位）。
  * `type` - 命令的类型。
  * `working_dir` - 命令在ECS实例中的执行路径。