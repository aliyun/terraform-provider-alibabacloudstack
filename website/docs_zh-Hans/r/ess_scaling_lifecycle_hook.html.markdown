---
subcategory: "弹性伸缩 ESS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_lifecycle_hook"
sidebar_current: "docs-alibabacloudstack-resource-ess-lifecycle-hook"
description: |-
  编排弹性伸缩生命周期挂钩
---

# alibabacloudstack_ess_lifecycle_hook

-> **NOTE：** 该资源在专有云上不受支持，将在 3.21.0 版本中移除。

-> **NOTE：** 该资源等效别名有: `alibabacloudstack_autoscaling_lifecyclehook`

使用 Provider 配置的凭证在指定的资源集下编排弹性伸缩生命周期挂钩（ESS Lifecycle Hook）资源。

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
  lifecycle_hook_name   = "testAccEssLifecycle_hook"
  lifecycle_transition  = "SCALE_OUT"
  heartbeat_timeout     = 400
  notification_metadata = "helloworld"
}
```

## 参数说明

支持以下参数：

* `scaling_group_id` - （必填，ForceNew）要为其分配生命周期挂钩的弹性伸缩组的 ID。

* `lifecycle_hook_name` - （可选，ForceNew）生命周期挂钩的名称。名称长度为 2~128 个字符，必须以字母开头，不能以 http:// 或 https:// 开头，可以包含字母、数字、下划线（_）、连字符（-）和半角句号（.）。该参数与 `name` 参数互斥。

* `lifecycle_transition` - （必填）生命周期挂钩的类型。取值：`SCALE_IN`（弹性收缩）、`SCALE_OUT`（弹性扩张）。

* `heartbeat_timeout` - （可选）生命周期挂钩超时前的最大等待时间。取值范围：30~21600，单位：秒。默认值：600。

* `default_result` - （可选）生命周期挂钩超时后弹性伸缩采取的操作。取值：`CONTINUE`（继续执行）、`ABANDON`（放弃执行）。默认值：`CONTINUE`。

* `notification_arn` - （可选）生命周期挂钩被触发时发送通知的事件总线 ARN。

* `notification_metadata` - （可选）弹性伸缩向通知目标发送消息时包含的附加信息。

* `name` - （可选，ForceNew，已弃用）生命周期挂钩的名称。该字段已弃用，请使用 `lifecycle_hook_name` 替代。如果不指定此参数值，默认值为生命周期挂钩的 ID。该参数与 `lifecycle_hook_name` 参数互斥。

## 属性说明

导出以下属性：

* `id` - 生命周期挂钩的 ID。

* `scaling_group_id` - 生命周期挂钩所属的弹性伸缩组 ID。

* `lifecycle_hook_name` - 生命周期挂钩的名称。

* `name` - （已弃用）生命周期挂钩的名称。

* `lifecycle_transition` - 生命周期挂钩的类型。

* `heartbeat_timeout` - 生命周期挂钩超时前的最大等待时间。

* `default_result` - 生命周期挂钩超时后弹性伸缩采取的操作。

* `notification_arn` - 发送通知的事件总线 ARN。

* `notification_metadata` - 生命周期挂钩被触发时发送的附加信息。

## Import

弹性伸缩生命周期挂钩可以使用生命周期挂钩 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_ess_lifecycle_hook.example lch-bp1234567890abcdef
```
