---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_lifecyclehook"
sidebar_current: "docs-Alibabacloudstack-autoscaling-lifecyclehook"
description: |- 
  编排弹性伸缩的生命周期钩子
---

# alibabacloudstack_autoscaling_lifecyclehook
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_lifecycle_hook`

使用Provider配置的凭证在指定的资源集下编排弹性伸缩的生命周期钩子。

## 示例用法

```hcl
data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type                = "ingress"
  ip_protocol         = "tcp"
  nic_type           = "intranet"
  policy             = "accept"
  port_range         = "22/22"
  priority           = 1
  security_group_id  = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip            = "172.16.0.0/24"
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
  availability_zone    = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
}

locals {
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
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
  zone_id             = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated         = false

  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccEssLifecycleHook-81869"
}

resource "alibabacloudstack_vswitch" "default2" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name             = "${var.name}"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size          = 1
  max_size          = 1
  scaling_group_name = "${var.name}"
  removal_policies  = ["OldestInstance", "NewestInstance"]
  vswitch_ids       = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
}

resource "alibabacloudstack_ess_lifecycle_hook" "default" {
  scaling_group_id     = "${alibabacloudstack_ess_scaling_group.default.id}"
  name                 = "${var.name}"
  lifecycle_transition = "SCALE_OUT"
  heartbeat_timeout    = 600
  default_result       = "CONTINUE"
  notification_metadata = "helloworld"
}
```

## 参数参考

支持以下参数：
  * `scaling_group_id` - (必填, 变更时重建) - 伸缩组ID。
  * `name` - (选填, 变更时重建) - 生命周期挂钩的名称。如果不指定，则默认生成一个唯一名称。
  * `lifecycle_hook_name` - (选填, 变更时重建) - 生命周期挂钩名称。与`name`功能类似，优先使用`lifecycle_hook_name`。
  * `lifecycle_transition` - (必填) - 生命周期挂钩对应伸缩活动类型。取值范围：`SCALE_OUT`(扩容)或`SCALE_IN`(缩容)。
  * `heartbeat_timeout` - (选填) - 生命周期挂钩为伸缩组活动设置的等待时间(单位：秒)，等待状态超时后会执行下一步动作。默认值为300秒。
  * `default_result` - (选填) - 当伸缩组发生弹性收缩活动(SCALE_IN)并触发多个生命周期挂钩时，DefaultResult为`ABANDON`的生命周期挂钩触发的等待状态结束时，会提前将其它对应的等待状态提前结束。其他情况下，下一步动作均以最后一个结束等待状态的下一步动作为准。取值范围：`CONTINUE`(继续)或`ABANDON`(放弃)。
  * `notification_arn` - (选填) - 生命周期挂钩通知对象标识符。用于指定接收通知的对象。
  * `notification_metadata` - (选填) - 伸缩活动的等待状态的固定字符串信息。最大长度为4KB。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `name` - 生命周期挂钩的名称。
  * `lifecycle_hook_name` - 生命周期挂钩名称。
  * `notification_arn` - 生命周期挂钩通知对象标识符。
  * `notification_metadata` - 伸缩活动的等待状态的固定字符串信息。
```