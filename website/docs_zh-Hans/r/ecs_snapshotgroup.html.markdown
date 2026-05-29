---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshot_group"
sidebar_current: "docs-Alibabacloudstack-ecs-snapshotgroup"
description: |-
  提供 ECS 快照一致性组资源。
---

# alibabacloudstack\_ecs\_snapshot\_group

> **Note:** 该资源也可以使用以下别名引用：
> - `alibabacloudstack_ecs_snapshotgroup`

提供 ECS 快照一致性组资源。快照一致性组包含一个或多个云盘对应的快照，这些云盘可以属于同一 ECS 实例或同一可用区内的多个实例。

## 示例用法
```
variable "name" {
  default = "tf-testaccebs-snapshot-group41329"
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
  //name_regex  = "arm_centos_7_6_20G_20211110.raw"
  //name_regex  = "^arm_centos_7"
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
  zone_id    		   = data.alibabacloudstack_zones.default.zones.0.id
  is_outdated          = false
  data_disks {
      name                 = "disk1"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
      delete_with_instance = true
    }
  data_disks {
      name                 = "disk2"
      category             = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
      size                 = 20
	  delete_with_instance = true
    }
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

data "alibabacloudstack_ecs_disks" "disks" {
	instance_id = "${alibabacloudstack_ecs_instance.default.id}"
}

resource "alibabacloudstack_ecs_snapshot_group" "default" {
  instant_access_retention_days = "7"
  snapshot_group_name = "${var.name}"
  disk_ids = [
               "${data.alibabacloudstack_ecs_disks.disks.disks.0.id}",
               "${data.alibabacloudstack_ecs_disks.disks.disks.1.id}",
               "${data.alibabacloudstack_ecs_disks.disks.disks.2.id}"
             ]
  description = "${var.name}"
  instance_id = "${alibabacloudstack_ecs_instance.default.id}"
  instant_access = "true"
}
```

## 参数参考

支持以下参数：

  * `disk_ids` - (选填) 需要创建快照一致性组的云盘 ID。在同可用区内支持跨实例设置多个云盘 ID。N 的取值范围为 1~16，即一个快照一致性组内最多支持设置 16 块总大小不超过 32 TiB 的云盘。该参数不可与 `exclude_disk_ids` 同时设置。如果您设置了 `instance_id`，则该参数只能设置指定实例内已挂载的云盘。
  * `exclude_disk_ids` - (选填) 实例中不需要创建快照的云盘 ID。指定云盘 ID 后，创建的快照一致性组将不包含该云盘对应的快照。N 的取值范围为 1~16。默认值：空，表示为实例中的所有云盘创建快照。该参数不可与 `disk_ids` 同时设置。
  * `instance_id` - (选填) ECS 实例 ID。您可以设置实例 ID 为实例内的指定云盘创建快照一致性组。
  * `snapshot_group_name` - (选填) 快照一致性组名称。长度为 2~128 个英文或中文字符。必须以大小写字母或中文开头，不能以 `http://` 或 `https://` 开头，可以包含数字、半角句号（.）、下划线（_）、短划线（-）或者半角冒号（:）。
  * `description` - (选填) 快照一致性组描述。长度为 2～256 个字符，不能以 `http://` 或 `https://` 开头。
  * `instant_access` - (选填，3.18+ 版本支持) 是否开启快照极速可用。取值范围：`true`（开启）、`false`（关闭）。默认值为 `false`。
  * `instant_access_retention_days` - (选填，3.18+ 版本支持) 设置快照极速可用的使用时间。单位：天，取值范围：1~65535。仅当 `instant_access=true` 时，该参数生效。到期后自动关闭快照极速功能。默认值：空，表示和快照释放时间一致。

## 属性参考

除了上述参数外，还导出了以下属性：

  * `id` - 快照一致性组 ID。
  * `snapshot_group_id` - 快照一致性组 ID。
  * `status` - 快照一致性组的状态。取值范围：`progressing`（进行中）、`accomplished`（已完成）、`failed`（失败）。
  * `create_time` - 快照一致性组的创建时间。
  * `disk_ids` - 快照一致性组中包含的云盘 ID。

## Import

ECS 快照一致性组可以使用快照一致性组 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_ecs_snapshot_group.example ssg-j6ciyh3k52qp7ovm****
```
