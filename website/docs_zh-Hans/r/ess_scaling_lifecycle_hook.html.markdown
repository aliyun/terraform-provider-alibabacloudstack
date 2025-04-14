---
subcategory: "Auto Scaling(ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_lifecycle_hook"
sidebar_current: "docs-alibabacloudstack-resource-ess-lifecycle-hook"
description: |-
  编排弹性伸缩生命周期挂钩
---

# alibabacloudstack_ess_lifecycle_hook
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_autoscaling_lifecyclehook`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩生命周期挂钩(ESS lifecycle hook)资源。

## 示例用法

```
data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  name       = "testAccEssScalingGroup_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "foo" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vswitch" "bar" {
  vpc_id            = "${alibabacloudstack_vpc.foo.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ess_scaling_group" "foo" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "testAccEssScaling_group"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = ["${alibabacloudstack_vswitch.foo.id}", "${alibabacloudstack_vswitch.bar.id}"]
}

resource "alibabacloudstack_ess_lifecycle_hook" "foo" {
  scaling_group_id      = "${alibabacloudstack_ess_scaling_group.foo.id}"
  name                  = "testAccEssLifecycle_hook"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}
```

## 参数说明

支持以下参数：

* `scaling_group_id` - (必填，变更时重建) 要为其分配生命周期挂钩的弹性伸缩组的ID。
* `name` - (可选，变更时重建) 生命周期挂钩的名称，必须包含2-64个字符(英文或中文)，以数字、英文字母或中文字符开头，可以包含数字、下划线 `_`、连字符 `-` 和小数点 `.`。如果不指定此参数值，默认值为生命周期挂钩的ID。
* `lifecycle_transition` - (必填) 与生命周期挂钩关联的伸缩活动类型。支持的值：`SCALE_OUT`, `SCALE_IN`。
* `heartbeat_timeout` - (可选) 定义生命周期挂钩超时之前可以经过的时间量(以秒为单位)。当生命周期挂钩超时时，弹性伸缩将执行默认结果参数中定义的操作。默认值：600。
* `default_result` - (可选) 定义当生命周期挂钩超时后，弹性伸缩组应采取的操作。适用值：`CONTINUE`, `ABANDON`，默认值：`CONTINUE`。
* `notification_arn` - (可选) 通知目标的Arn。
* `notification_metadata` - (可选) 当弹性伸缩向通知目标发送消息时，您希望包含的其他信息。

## 属性说明

导出以下属性：

* `id` - 生命周期挂钩的ID。
* `scaling_group_id` - 生命周期挂钩所属的弹性伸缩组ID。
* `name` - 生命周期挂钩的名称。
* `default_result` - 当生命周期挂钩超时时，弹性伸缩组应采取的操作。
* `heartbeat_timeout` - 生命周期挂钩超时前可以经过的时间量(以秒为单位)。
* `lifecycle_transition` - 与生命周期挂钩关联的伸缩活动类型。
* `notification_metadata` - 将发送到通知目标的其他信息。
* `notification_arn` - 通知目标的Arn。