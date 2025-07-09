---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_invocations"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-invocations"
description: |-
  提供当前 AlibabacloudStack 账户下所有 ECS 执行记录（Invocation）的列表。
---

# alibabacloudstack\_ecs\_invocations

该数据源用于根据指定过滤条件获取当前账户下所有 ECS 实例命令执行记录（Invocation）信息。

## 示例用法

```hcl
variable "name" {
		default = "tf-testAccEcsInvocationsTest78815"
	}

resource "alibabacloudstack_ecs_command" "default" {
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
	instance_ids = ["${alibabacloudstack_ecs_instance.default.id}"]
	username = "root"
	repeat_mode = "Once"
}

data "alibabacloudstack_ecs_invocations" "default" {
  invocation_id = "${alibabacloudstack_ecs_invocation.default.id}"
}
```
## 参数说明
以下参数支持配置：

* `ids` - (可选) 指定查询的调用 ID 列表。
* `command_id` - (可选) 命令 ID。你可以通过调用 [DescribeCommands](~~ 64843 ~~) 查询可用的命令 ID。
* `command_name` - (可选) 命令名称。
* `instance_id` - (可选) 指定要执行命令的实例 ID 列表，最多支持 50 个实例 ID。
* `timed` - (可选) 是否周期性地执行命令，默认值为 false。
* `invocation_id` - (可选) 执行记录 ID。
## 属性说明
除上述参数外，还将导出以下属性：

* `invocations` - 所有执行记录的列表
* `id` - 执行记录 ID
* `command_id` - 命令 ID，可通过 [DescribeCommands](~~ 64843 ~~) 查询。
* `create_time` - 资源创建时间
* `frequency` - 周期性任务的执行周期，请注意：
两次定时任务之间的间隔不能小于 10 秒；
定时任务的最小执行间隔不能小于任务超时时间；
任务超时时间是调用 CreateCommand 创建命令时设置的 Timeout 参数；
该参数遵循 Cron 表达式格式，详见 [设置定时执行命令](~~ 64769 ~~)。
* `invocation_id` - 执行记录的唯一标识 ID
* `repeat_mode` - 设置命令执行方式，取值范围：
* `Once`: 立即执行一次命令；
* `Period`: 周期性执行命令（需配合 Timed = true 和 Frequency 使用）；
* `NextRebootOnly`: 下次实例启动时执行一次；
* `EveryReboot`: 每次实例启动时都自动执行。 默认值：
若未指定 Timed = true 和 Frequency，默认为 Once；
若指定了这两个参数，则视为 Period。 注意事项：
当值为 Period, NextRebootOnly, 或 EveryReboot 时，可调用 [StopInvocation](~~ 64838 ~~) 停止待执行或周期性执行的任务；
当值为 Period 或 EveryReboot 时，可调用 [DescribeInvocationResults](~~ 64845 ~~) 并设置 IncludeHistory = true 查看历史执行记录。
parameters - 启用自定义参数功能后，在执行命令时传入的键值对。自定义参数限制如下：
键长度最大 64 字符，不能为空；
值可以为空字符串；
自定义参数和原始命令内容经 Base64 编码后的总长度不得超过 16KB；
自定义参数名称集合必须是创建命令时定义的参数集的子集；
对于未传入的参数，可以用空字符串代替。
* `resource_type` - 资源类型定义，取值为 invocation。
* `tags` - 标签列表。
* `timed` - 是否为定时任务，默认值为 false。
* `username` - 在 ECS 实例中运行命令的用户名：
Linux 实例默认以 root 用户运行；
Windows 实例默认以 System 用户运行；
可指定其他已存在的用户来执行命令，提升安全性，详见 [设置普通用户执行云助手命令](~~ 203771 ~~)。
windows_password_name - 在 Windows 实例中以非默认用户（System）执行命令时，需要同时传入 Username 和此参数。为降低密码泄露风险，密码需托管在 O&M 编排服务的密钥仓库中，此处仅传递密码名称。详见 [加密参数](~~ 186828 ) 和 [设置普通用户执行云助手命令]( 203771 ~~)。
如果使用 Linux 实例的 root 用户或 Windows 实例的 System 用户执行命令，则无需传入此参数。