---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scheduledtasks"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scheduledtasks"
description: |- 
  查询弹性伸缩定时任务
---

# alibabacloudstack_autoscaling_scheduledtasks
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_scheduled_tasks`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩定时任务列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccDataSourceScheduledtas"
}

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

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size = 0
  max_size = 2
  default_cooldown = 20
  removal_policies = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ess_scheduled_task" "default" {
  scheduled_action    = "${alibabacloudstack_ess_scaling_rule.default.ari}"
  launch_time         = "2025-03-08T17:50Z"
  scheduled_task_name = "${var.name}"
}

data "alibabacloudstack_ess_scheduled_tasks" "default" {
  ids = ["${alibabacloudstack_ess_scheduled_task.default.id}"]
  name_regex = "tf-testAccDataSourceScheduledtas.*"
  output_file = "scheduled_tasks_output.txt"
}

output "first_scheduled_task_id" {
  value = data.alibabacloudstack_ess_scheduled_tasks.default.tasks.0.id
}
```

## 参数说明

以下参数是支持的：

* `scheduled_task_id` - (可选) 定时任务的ID。
* `scheduled_action` - (可选) 定时任务触发时需要执行的操作。
* `name_regex` - (可选) 用于通过名称过滤定时任务的正则表达式字符串。
* `ids` - (可选) 用于过滤结果的定时任务ID列表。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 定时任务ID列表。
* `names` - 定时任务名称列表。
* `tasks` - 定时任务列表。每个元素包含以下属性：
  * `id` - 定时任务的ID。
  * `name` - 定时任务的名称。
  * `scheduled_action` - 定时任务触发时需要执行的操作。
  * `description` - 定时任务的描述信息。
  * `launch_expiration_time` - 定时任务触发操作失败后，在此时间内重试。单位为秒，取值范围：0~21600。
  * `launch_time` - 定时任务触发的时间点。
  * `min_value` - 当定时任务的缩放方法是指定伸缩组中的实例数时，伸缩组中的最小实例数。
  * `max_value` - 当定时任务的缩放方法是指定伸缩组中的实例数时，伸缩组中的最大实例数。
  * `recurrence_type` - 指定定时任务的重复类型。
  * `recurrence_value` - 指定定时任务的重复频率。
  * `recurrence_end_time` - 指定定时任务不再重复的结束时间。
  * `task_enabled` - 是否启动定时任务。默认值为 `true`。
  * `description` - (计算属性) 定时任务的描述信息。
