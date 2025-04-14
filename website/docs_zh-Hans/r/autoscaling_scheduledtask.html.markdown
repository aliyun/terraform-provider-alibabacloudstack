---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduledtask"
sidebar_current: "docs-Alibabacloudstack-autoscaling-scheduledtask"
description: |- 
  编排弹性伸缩定时任务
---

# alibabacloudstack_autoscaling_scheduledtask
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_scheduled_task`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩定时任务。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
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
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
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
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccEssScheduleConfig-227425"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  default_cooldown   = 20
  removal_policies   = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids        = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  image_id         = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type    = "ecs.e4.small"
  security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
  force_delete     = true
  active          = true
  enable          = true
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
}

resource "alibabacloudstack_ess_scaling_rule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type  = "TotalCapacity"
  adjustment_value = "1"
  cooldown        = 0
}

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action     = "${alibabacloudstack_ess_scaling_rule.default.ari}"
  launch_time          = "2025-03-08T17:50Z"
  scheduled_task_name  = "${var.name}"
  description         = "This is an example of a scheduled task."
  launch_expiration_time = 3600
  recurrence_type      = "Daily"
  recurrence_value     = "1"
  recurrence_end_time  = "2025-04-08T17:50Z"
  task_enabled        = true
}
```

## 参数说明

支持以下参数：
  * `scheduled_action` - (必填) - 定时任务触发时需要执行的操作。必须是伸缩规则的唯一标识符。
  * `launch_time` - (必填) - 定时任务触发的时间点。按照ISO 8601标准格式指定时间，格式为`YYYY-MM-DDThh:mm:ssZ`。时间必须为UTC。您不能输入创建计划任务日期之后90天以上的某个时间点。如果设置了`recurrence_type`参数，则任务会根据`launch_time`反复执行；否则，任务只会在`launch_time`指定的日期和时间执行一次。
  * `scheduled_task_name` - (选填) - 定时任务的名称，长度为2-40个字符(英文或中文)。
  * `description` - (选填) - 定时任务的描述信息，长度为2-200个字符(英文或中文)。
  * `launch_expiration_time` - (选填) - 定时任务触发操作失败后，在此时间内重试。单位为秒，取值范围：0~21600，默认值：600。
  * `recurrence_type` - (选填) - 重复执行定时任务的类型。有效值：
    * `Daily`: 定时任务每隔指定天数执行一次。
    * `Weekly`: 定时任务在每周的指定天执行。
    * `Monthly`: 定时任务在每月的指定天执行。
    * `Cron`: 定时任务基于指定的cron表达式执行。
  * `recurrence_value` - (选填) - 重复执行定时任务的数值。有效值取决于`recurrence_type`：
    * `Daily`: 可以输入一个值。有效值：1到31。
    * `Weekly`: 可以输入多个值并用逗号(,)分隔。例如，值0到6分别对应星期日到星期六。
    * `Monthly`: 可以输入两个值，格式为A-B。有效值A和B：1到31。B的值必须大于或等于A的值。
    * `Cron`: 可以输入一个cron表达式，该表达式以UTC编写，并由五个字段组成：分钟、小时、日期、月份和星期几。表达式可以包含通配符，包括逗号(,)、问号(?)、连字符(-)、星号(*)、井号(#)、正斜杠(/)以及字母L和W。
  * `recurrence_end_time` - (选填) - 重复执行定时任务的结束时间。按照ISO 8601标准格式指定时间，格式为`YYYY-MM-DDThh:mm:ssZ`。时间必须为UTC。您不能输入创建计划任务日期之后365天以上的某个时间点。
  * `task_enabled` - (选填) - 是否启动定时任务。默认值：`true`。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - 定时任务的ID。
  * `description` - 定时任务的描述信息。
  * `recurrence_type` - 重复执行定时任务的类型。
  * `recurrence_value` - 重复执行定时任务的数值。
  * `recurrence_end_time` - 重复执行定时任务的结束时间。
  * `task_enabled` - 表明定时任务是否已启用。