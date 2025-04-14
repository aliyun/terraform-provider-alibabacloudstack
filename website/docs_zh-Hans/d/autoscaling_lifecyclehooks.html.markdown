---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_lifecyclehooks"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-lifecyclehooks"
description: |- 
  查询弹性伸缩生命周期事件钩子
---

# alibabacloudstack_autoscaling_lifecyclehooks
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_lifecycle_hooks`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩生命周期事件钩子列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccDataSourceLcHooks-843"
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

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size = 0
  max_size = 2
  default_cooldown = 20
  removal_policies = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ess_lifecycle_hook" "default" {
  scaling_group_id      = "${alibabacloudstack_ess_scaling_group.default.id}"
  name                  = "${var.name}-hook"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}

data "alibabacloudstack_ess_lifecycle_hooks" "default" {
  name_regex = "${alibabacloudstack_ess_lifecycle_hook.default.name}"
  scaling_group_id = "${alibabacloudstack_ess_lifecycle_hook.default.scaling_group_id}"
}

output "lifecycle_hooks_info" {
  value = data.alibabacloudstack_ess_lifecycle_hooks.default.hooks
}
```

## 参数说明

以下参数是支持的：

* `scaling_group_id` - (可选) 伸缩组ID。用于筛选属于特定伸缩组的生命周期钩子。
* `name_regex` - (可选) 正则表达式字符串，用于通过生命周期钩子名称筛选结果。
* `ids` - (可选) 生命周期钩子ID列表，用于进一步筛选结果。

## 属性说明

除了上述参数外，还导出以下属性：

* `hooks` - 生命周期钩子列表。每个元素包含以下属性：
  * `id` - 生命周期钩子的ID。
  * `name` - 生命周期钩子的名称。
  * `scaling_group_id` - 生命周期钩子所属的伸缩组的ID。
  * `default_result` - 定义当生命周期钩子超时后，伸缩组应采取的操作。它可以是 `CONTINUE` 或 `ABANDON`。
  * `heartbeat_timeout` - 定义在生命周期钩子超时之前可以经过的时间量（以秒为单位）。当生命周期钩子超时时，弹性伸缩执行 `default_result` 参数中定义的操作。
  * `lifecycle_transition` - 与生命周期钩子关联的伸缩活动类型。可能的值包括 `INSTANCE_LAUNCHING` 和 `INSTANCE_TERMINATING` 等。
  * `notification_arn` - 弹性伸缩在实例因生命周期钩子进入等待状态时将通知的ARN通知目标。
  * `notification_metadata` - 您希望弹性伸缩在向通知目标发送消息时包含的其他信息。

* `ids` - 生命周期钩子ID列表。
* `names` - 生命周期钩子名称列表。