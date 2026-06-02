---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_invocation"
sidebar_current: "docs-Alibabacloudstack-resource-ecs-invocation"
description: |-
  提供一个 ECS 命令执行记录（Invocation）资源。
---

# alibabacloudstack\_ecs\_invocation

提供一个 ECS 命令执行记录（Invocation）资源。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccecsinvocation22834"
}

resource "alibabacloudstack_ecs_command" "default" {
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
  zone_id              = data.alibabacloudstack_zones.default.zones.0.id
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
## 参数说明
以下参数支持配置：

* `command_id` - (必需, 变更时强制重建) 命令 ID。你可以通过调用 [DescribeCommands](~~ 64843 ~~) 查询所有可用的命令 ID。

* `repeat_mode` - (可选, 变更时强制重建) 设置命令的执行方式，取值范围：

* `Once`: 立即执行一次命令。
* `Period`: 定期执行命令。当该参数值为 Period 时，必须同时指定 Timed = true 和 frequency 参数。
NextRebootOnly: 下次实例启动时自动执行一次。
EveryReboot: 每次实例启动时都自动执行。
默认值：

如果未指定 Timed = true 和 frequency，默认为 Once。
如果指定了这两个参数，则无论其值是否设置，均视为 Period。
注意事项：

当此参数值为 Period、NextRebootOnly 或 EveryReboot 时，可以调用 [StopInvocation](~~ 64838 ~~) 来停止待执行或周期性执行的任务。
当此参数值为 Period 或 EveryReboot 时，可以调用 [DescribeInvocationResults](~~ 64845 ~~) 并设置 IncludeHistory = true 来查看历史执行记录。
username - (可选, 变更时强制重建) 在 ECS 实例中运行命令的用户名：

Linux 实例默认由 root 用户执行。
Windows 实例默认由 System 用户执行。
你也可以指定实例中已存在的其他用户来执行命令，以提升普通用户执行云助手命令的安全性。详见 [设置普通用户执行云助手命令](~~ 203771 ~~)。
## 属性说明
除上述参数外，还将导出以下属性：

* `create_time` - 资源创建时间。
* `invocation_id` - 执行记录的唯一标识 ID。