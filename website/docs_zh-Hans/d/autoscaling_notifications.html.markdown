---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_notifications"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-notifications"
description: |- 
  查询弹性伸缩消息通知
---

# alibabacloudstack_autoscaling_notifications
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ess_notifications`

根据指定过滤条件列出当前凭证权限可以访问的弹性伸缩消息通知列表。

## 示例用法

```hcl
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
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
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
  availability_zone     = data.alibabacloudstack_zones.default.zones[0].id
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
  is_outdated         = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccDataSourceEssNs-17"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = ["${alibabacloudstack_vswitch.default.id}"]
}

resource "alibabacloudstack_ess_notification" "default" {
  scaling_group_id   = "${alibabacloudstack_ess_scaling_group.default.id}"
  notification_types = ["AUTOSCALING:SCALE_OUT_SUCCESS"]
  notification_arn   = "acs:ess"
}

data "alibabacloudstack_ess_notifications" "default" {
  scaling_group_id = "${alibabacloudstack_ess_notification.default.scaling_group_id}"
  ids              = ["notification-id-1", "notification-id-2"]
  output_file      = "notifications_output.txt"
}

output "first_notification_arn" {
  value = data.alibabacloudstack_ess_notifications.default.notifications[0].notification_arn
}
```

## 参数参考

以下参数是支持的：

* `scaling_group_id` - (必填) 伸缩组的ID。此参数用于指定需要查询的通知所属的伸缩组。
* `ids` - (选填)通知ID列表。如果指定，数据源将仅返回与这些ID匹配的通知。

## 属性参考

除了上述参数外，还导出以下属性：

* `ids` - 通知ID列表。此属性包含所有符合条件的通知的唯一标识符。
* `notifications` - 自动伸缩通知列表。每个元素包含以下属性：
  * `id` - 通知的唯一标识符。此属性用于唯一标识每个通知。
  * `scaling_group_id` - 与通知关联的伸缩组ID。此属性指定了通知所属的伸缩组。
  * `notification_arn` - 通知对象的阿里云资源名称(ARN)。此属性提供了通知对象的完整标识符。
  * `notification_types` - 一类或多类自动伸缩事件及资源变化通知。此属性列出了触发该通知的事件类型，例如生命周期事件、伸缩活动等。
