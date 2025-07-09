---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_invocation"
sidebar_current: "docs-Alibabacloudstack-ecs-invocation"
description: |-
  Provides a ecs Invocation resource.
---

# alibabacloudstack\_ecs\_invocation

Provides a ecs Invocation resource.

## Example Usage
```
variable "name" {
    default = "tf-testaccecsinvocation22834"
}

resource alibabacloudstack_ecs_command	"default" {
	description = "command description"
	command_content = "pwd"
	type = "RunShellScript"
	name = "tf-testaccecsinvocation22834"
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
  count = 8  # 遍历 1-8 核CPU配置

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






resource "alibabacloudstack_ecs_invocation" "default" {
  command_id = "${alibabacloudstack_ecs_command.default.id}"
  repeat_mode = "Once"
  username = "root"
  instance_ids = [
                   "${alibabacloudstack_ecs_instance.default.id}"
                 ]
}
```

## Argument Reference

The following arguments are supported:
  * `command_id` - (Required, ForceNew) - Command ID. You can query all available commandids by calling [DescribeCommands](~~ 64843 ~~).
  * `repeat_mode` - (Optional, ForceNew) - Sets how commands are executed. Value range:-Once: execute the command immediately.-Period: Execute commands regularly. When the value of this parameter is 'Period', you must specify both the 'Timed = true' parameter value and the 'requency' parameter.-Nextrabootonly: When the instance is started next time, the command is automatically executed.-EveryReboot: The command will be executed automatically every time the instance is started.Default value:-When the 'Timed = true' parameter value and the 'requency' parameter are not specified, the default value is' once '.-When the 'Timed = true' parameter value and the 'requency' parameter are specified, it will be processed as 'Period' regardless of whether the parameter value has been set.Precautions:-When the value of this parameter is 'Period', 'nextrabootonly', or 'EveryReboot', you can call [StopInvocation](~~ 64838 ~~) to stop the pending command or periodically executed command.-When the value of this parameter is 'Period' or 'EveryReboot', you can call [DescribeInvocationResults](~~ 64845 ~~) and specify IncludeHistory = true to view the history of the command cycle execution.
  * `username` - (Optional, ForceNew) - The name of the user who runs the command in the ECS instance.-The ECS instance of the Linux system. The command is run by the root user by default.-The ECS instance of the Windows System. By default, the System user executes commands.You can also specify other users that already exist in the instance to execute commands, making it safer for ordinary users to execute cloud assistant commands. For more information, see [Set up common users to execute cloud assistant commands](~~ 203771 ~~).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `create_time` - The creation time of the resource
  * `invocation_id` - Invocation ID
