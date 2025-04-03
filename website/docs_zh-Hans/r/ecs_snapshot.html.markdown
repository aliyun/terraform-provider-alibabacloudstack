---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshot"
sidebar_current: "docs-Alibabacloudstack-ecs-snapshot"
description: |- 
  编排云服务器（Ecs）快照
---

# alibabacloudstack_ecs_snapshot
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_snapshot`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）快照。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccecssnapshot10982"
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
  cpu_core_count   = 1
  memory_size      = 1
  instance_type_family = "ecs.n4"
  sorted_by        = "Memory"
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

resource "alibabacloudstack_ecs_snapshot" "default" {
  description  = "rdk_test_description"
  snapshot_name = "rdk_test_name"
  disk_id       = alibabacloudstack_ecs_instance.default.system_disk_id
}
```

## 参数参考

支持以下参数：
  * `disk_id` - (必填, 变更时重建) - 要为其创建快照的磁盘的 ID。
  * `snapshot_name` - (选填, 变更时重建) - 快照的显示名称。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(:)、下划线(_)或者短划线(-)。为防止和自动快照的名称冲突，不能以`auto`开头。
  * `description` - (选填, 变更时重建) - 快照的描述。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。默认值：空。
  * `tags` - (可选) - 要分配给资源的标签映射。

### 超时时间

* `create` - (默认为 2 分钟)用于创建快照(直到它达到初始 `SnapshotCreatingAccomplished` 状态)。
* `delete` - (默认为 2 分钟)用于终止快照。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - 快照的 ID。
  * `snapshot_name` - 快照的显示名称。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`或`https://`开头。可以包含数字、半角冒号(:)、下划线(_)或者短划线(-)。为防止和自动快照的名称冲突，不能以`auto`开头。
  * `description` - 快照的描述。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。