---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_invocations"
description: |-
  Provides a list of ecs invocations owned by an alibabacloudstack account.
---

# alibabacloudstack\_ecs\_invocations

This data source provides a list of ecs invocations in an alibabacloudstack account according to the specified filters.

## Example Usage
```
variable "name" {
		default = "tf-testAccEcsInvocationsTest78815"
	}

	resource alibabacloudstack_ecs_command	"default" {
		description = "command description"
		command_content = "pwd"
		type = "RunShellScript"
		name = "${var.name}"
	}

	
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}


resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  	type = "ingress"
  	ip_protocol = "tcp"
  	nic_type = "intranet"
  	policy = "accept"
  	port_range = "22/22"
  	priority = 1
  	security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  	cidr_ip = "192.168.0.0/16"
}


data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
  most_recent = true
  owners      = "system"
}


data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  sorted_by         = "CPU"
}

data "alibabacloudstack_instance_types" "default" {
  count = 8  

  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = count.index + 1  # 1-8
  sorted_by            = "Memory"
}

locals {
  filtered_default = [for d in data.alibabacloudstack_instance_types.default : d if length(d.ids) > 0]
  fallback_all     = length(data.alibabacloudstack_instance_types.all.ids) > 0 ? data.alibabacloudstack_instance_types.all.ids : []
  
  default_instance_type_id = coalesce(
    try(local.filtered_default[0].ids[0], null),
    try(local.fallback_all[0], null),
    "no-available-instance-type"
  )
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}



	resource alibabacloudstack_ecs_invocation "default" {
		command_id = "${alibabacloudstack_ecs_command.default.id}"
		instance_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
		username = "root"
		repeat_mode = "Once"
	}

	

data "alibabacloudstack_ecs_invocations" "default" {
  invocation_id = "${alibabacloudstack_ecs_invocation.default.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `ids` - (Optional) - the ids of the ecs invocations.
  * `command_id` - (Optional) - Command ID. You can query all available commandids by calling [DescribeCommands](~~ 64843 ~~).
  * `command_name` - (Optional) - - Invocation name
  * `instance_id` - (Optional) - The list of instances to execute the command. You can specify up to 50 instance IDs. The value range of N is 1 to 50.
  * `timed` - (Optional) - Whether the command is periodically executed.The default value is false.
  * `invocation_id` - (Optional) - Invocation ID

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `invocations` - the list of all invocations
    * `id` - the invocation ID
    * `command_id` - Command ID. You can query all available commandids by calling [DescribeCommands](~~ 64843 ~~).
    * `create_time` - The creation time of the resource
    * `frequency` - The execution cycle of a periodic execution command.You need to note:-The command interval between two periodic executions cannot be less than 10 seconds, and the minimum execution interval of periodic commands cannot be less than the task execution timeout. The task execution timeout is The 'Timeout' parameter that you set when you call the CreateCommand to create the command. The value of the 'Timeout' parameter can be queried by [DescribeCommands](~~ 64843 ~~).-When the value of the parameter' time' is' true', 'requency' is a required parameter.-The value of this parameter follows the Cron expression. For more information, see [Set Timed Execution Command](~~ 64769 ~~).
    * `invocation_id` - Invocation ID
    * `repeat_mode` - Sets how commands are executed. Value range:-Once: execute the command immediately.-Period: Execute commands regularly. When the value of this parameter is 'Period', you must specify both the 'Timed = true' parameter value and the 'requency' parameter.-Nextrabootonly: When the instance is started next time, the command is automatically executed.-EveryReboot: The command will be executed automatically every time the instance is started.Default value:-When the 'Timed = true' parameter value and the 'requency' parameter are not specified, the default value is' once '.-When the 'Timed = true' parameter value and the 'requency' parameter are specified, it will be processed as 'Period' regardless of whether the parameter value has been set.Precautions:-When the value of this parameter is 'Period', 'nextrabootonly', or 'EveryReboot', you can call [StopInvocation](~~ 64838 ~~) to stop the pending command or periodically executed command.-When the value of this parameter is 'Period' or 'EveryReboot', you can call [DescribeInvocationResults](~~ 64845 ~~) and specify IncludeHistory = true to view the history of the command cycle execution.
    * `parameters` - When the custom parameter function is enabled, the key-value pair of the custom parameter passed in when the command is executed. The number of custom parameters ranges from 0 to 10.The key of-Map is not allowed to be an empty string and supports up to 64 characters.The value of-Map is allowed to be an empty string.-The combined length of the custom parameter and the original command content cannot exceed 16KB after Base64 encoding.-The set of custom parameter names must be a subset of the parameter set defined when the command was created. For parameters that are not passed in, you can use an empty string instead.You can disable custom parameters by canceling this parameter.
    * `resource_type` - Resource type definition. Valued invocation.
    * `tags` - List of labels.
    * `timed` - Whether the command is periodically executed.The default value is false.
    * `username` - The name of the user who runs the command in the ECS instance.-The ECS instance of the Linux system. The command is run by the root user by default.-The ECS instance of the Windows System. By default, the System user executes commands.You can also specify other users that already exist in the instance to execute commands, making it safer for ordinary users to execute cloud assistant commands. For more information, see [Set up common users to execute cloud assistant commands](~~ 203771 ~~).
    * `windows_password_name` - The password name of the user who executed the command in the Windows instance.When you want to execute a command in a Windows instance as a non-default user (System), you need to pass in both 'Username' and this parameter. To reduce the risk of password leakage, you must host the password in plain text in the parameter warehouse of the O & M orchestration service. Only the name of the password is passed in here. For more information, see [Encryption Parameters](~~ 186828 ~~) and [Set Common User to Execute Cloud Assistant Command](~~ 203771 ~~).> When you use the root user of the Linux instance or the System user of the Windows instance to execute the command, you do not need to pass this parameter.
