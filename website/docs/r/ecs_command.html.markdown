---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_command"
sidebar_current: "docs-Alibabacloudstack-ecs-command"
description: |- 
  Provides a ecs Command resource.
---

# alibabacloudstack_ecs_command

Provides a ECS Command resource.

For information about ECS Command and how to use it, see [What is Command](https://www.alibabacloud.com/help/en/doc-detail/64844.htm).

## Example Usage

Basic Usage

```terraform
variable "name" {
    default = "tf-testaccecscommand49325"
}

resource "alibabacloudstack_ecs_command" "default" {
  command_name      = var.name
  command_content   = "bHMK" # Base64-encoded content: "ls\n"
  type              = "RunShellScript"
  description       = "Command for listing files in the root directory"
  enable_parameter  = false
  timeout           = 120
  working_dir       = "/root"
}
```

## Argument Reference

The following arguments are supported:

* `command_name` - (Required, ForceNew) The name of the command. It must be unique and can contain up to 128 characters.
* `command_content` - (Required, ForceNew) The Base64-encoded content of the command. This is the actual script or command that will be executed on the ECS instances.
* `type` - (Required, ForceNew) The type of the command. Valid values:
  * `RunShellScript`: For Linux systems, executes shell scripts.
  * `RunBatScript`: For Windows systems, executes batch scripts.
  * `RunPowerShellScript`: For Windows systems, executes PowerShell scripts.
* `description` - (Optional, ForceNew) A brief description of the command. It helps identify the purpose of the command.
* `enable_parameter` - (Optional, ForceNew) Specifies whether the command supports custom parameters. Default value is `false`. If set to `true`, you can pass parameters when invoking the command.
* `timeout` - (Optional, ForceNew) The timeout period for the command execution on ECS instances. Unit: seconds. Default value is `60`.
* `working_dir` - (Optional, ForceNew) The working directory where the command will be executed on the ECS instance. Default value is `/root` for Linux systems and `C:\Windows\system32` for Windows systems.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the ECS Command resource.
* `command_id` - The ID of the created ECS Command.
