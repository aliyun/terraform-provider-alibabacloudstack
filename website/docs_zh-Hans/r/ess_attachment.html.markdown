---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_attachment"
sidebar_current: "docs-alibabacloudstack-resource-ess-attachment"
description: |- 
  编排绑定多台ECS实例附加到指定的伸缩组
---

# alibabacloudstack_ess_attachment

使用Provider配置的凭证在指定的资源集下编排绑定多台ECS实例附加到指定的伸缩组。

-> **NOTE:** 只有当伸缩组处于活动状态且没有正在进行的伸缩活动时，才能附加或移除ECS实例。

-> **NOTE:** 在一个伸缩组中，有两种类型的ECS实例："AutoCreated" 和 "Attached"。它们的总数不能超过伸缩组的 "MaxSize"。

## 示例用法


```
variable "name" {
  default = "essattachmentconfig"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_security_group.default.id
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 0
  max_size           = 2
  scaling_group_name = var.name
  removal_policies   = ["OldestInstance", "NewestInstance"]
  vswitch_ids        = [alibabacloudstack_vswitch.default.id]
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id  = alibabacloudstack_ess_scaling_group.default.id
  image_id          = data.alibabacloudstack_images.default.images[0].id
  instance_type     = data.alibabacloudstack_instance_types.default.instance_types[0].id
  security_group_id = alibabacloudstack_security_group.default.id
  force_delete      = true
  active            = true
  enable            = true
}

resource "alibabacloudstack_instance" "default" {
  image_id                   = data.alibabacloudstack_images.default.images[0].id
  instance_type              = data.alibabacloudstack_instance_types.default.instance_types[0].id
  count                      = 2
  security_groups            = [alibabacloudstack_security_group.default.id]
  internet_charge_type       = "PayByTraffic"
  internet_max_bandwidth_out = "10"
  instance_charge_type       = "PostPaid"
  system_disk_category       = "cloud_efficiency"
  vswitch_id                 = alibabacloudstack_vswitch.default.id
  instance_name              = var.name
}

resource "alibabacloudstack_ess_attachment" "default" {
  scaling_group_id = alibabacloudstack_ess_scaling_group.default.id
  instance_ids     = [alibabacloudstack_instance.default[0].id, alibabacloudstack_instance.default[1].id]
  force            = true
}
```

## 参数说明

支持以下参数：

* `scaling_group_id` - (必填) 伸缩组ID。
* `instance_ids` - (必填) 要附加到伸缩组的ECS实例ID列表。最多可以输入20个ID。
* `force` - (可选) 是否强制删除“AutoCreated”ECS实例以释放伸缩组容量“MaxSize”，以便附加ECS实例。默认为`false`。

-> **NOTE:** “AutoCreated”ECS实例在从伸缩组中移除后将被删除，但“Attached”不会被删除。

-> **NOTE:** 附加ECS实例的限制：

- 要附加的ECS实例和伸缩组必须具有相同的区域和网络类型(`Classic` 或 `VPC`)。
- 要附加的ECS实例和具有活动伸缩配置的实例必须具有相同的实例类型。
- 要附加的ECS实例必须处于运行状态。
- 要附加的ECS实例尚未附加到其他伸缩组。
- 要附加的ECS实例支持包年包月和按量付费两种计费方式。

## 属性说明

导出以下属性：

* `id` - (必填，ForceNew) ESS附件资源ID。
* `instance_ids` - (必填) 已附加到伸缩组的“Attached”ECS实例的ID列表。
* `force` - 是否强制删除“AutoCreated”ECS实例以释放伸缩组容量“MaxSize”。

## 导入

ESS attachment 可以通过id或伸缩组id导入，例如：

```bash
$ terraform import alibabacloudstack_ess_attachment.example asg-abc123456
```