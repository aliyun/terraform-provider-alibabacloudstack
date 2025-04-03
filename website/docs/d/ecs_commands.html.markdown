---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_commands"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-commands"
description: |- 
  Provides a list of ecs commands owned by an AlibabacloudStack account.
---

# alibabacloudstack_ecs_commands

This data source provides a list of ECS Commands in an AlibabacloudStack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_ecs_commands" "example" {
  ids        = ["E2RY53-xxxx"]
  name_regex = "tf-testAcc"
}

output "first_ecs_command_id" {
  value = data.alibabacloudstack_ecs_commands.example.commands.0.command_id
}
```

Advanced Usage with Filters

```terraform
data "alibabacloudstack_ecs_commands" "example" {
  content_encoding = "Base64"
  description      = "Test command"
  name             = "test-command"
  type             = "RunShellScript"
}

output "command_ids" {
  value = data.alibabacloudstack_ecs_commands.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by Command name.
* `ids` - (Optional, ForceNew) A list of Command IDs used to filter results.
* `content_encoding` - (Optional, ForceNew) The encoding method of the command content. Valid values:
  * `PlainText`: No encoding, using clear text transmission.
  * `Base64`: Base64 encoding.
  Default value: `Base64`. If an invalid value is provided, it will be treated as `Base64`.
* `description` - (Optional, ForceNew) The description of the command.
* `name` - (Optional, ForceNew) The name of the command.
* `command_provider` - (Optional, ForceNew) The provider of the command.
* `type` - (Optional, ForceNew) The type of the command. Valid values:
  * `RunBatScript`
  * `RunPowerShellScript`
  * `RunShellScript`

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Command names.
* `commands` - A list of ECS Commands. Each element contains the following attributes:
  * `command_content` - The Base64-encoded content of the command.
  * `command_id` - The ID of the command.
  * `description` - The description of the command.
  * `enable_parameter` - Specifies whether to use custom parameters in the command.
  * `id` - The ID of the command (same as `command_id`).
  * `name` - The name of the command.
  * `parameter_names` - A list of custom parameter names parsed from the command content when the command was created.
  * `timeout` - The timeout period (in seconds) for the command to run on ECS instances.
  * `type` - The type of the command.
  * `working_dir` - The execution path of the command in the ECS instance.