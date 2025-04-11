---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_snapshots"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-snapshots"
description: |- 
  查询云服务器快照
---

# alibabacloudstack_ecs_snapshots
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_snapshots`

根据指定过滤条件列出当前凭证权限可以访问的云服务器快照列表。

## 示例用法

```hcl
variable "name" {
	default = "tf-testaccSnapshotDataSourceBasic73027"
}

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
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

resource "alibabacloudstack_snapshot" "default" {
  disk_id      = "${alibabacloudstack_ecs_instance.default.system_disk_id}"
  name         = "${var.name}"
  description  = "${var.name}"
}

data "alibabacloudstack_ecs_snapshots" "default" {
  ids            = ["${alibabacloudstack_snapshot.default.id}"]
  instance_id    = "${alibabacloudstack_ecs_instance.default.id}"
  disk_id       = "${alibabacloudstack_ecs_instance.default.system_disk_id}"
  name_regex    = "tf-testaccSnapshotDataSourceBasic.*"
  status        = "accomplished"
  type          = "user"
  source_disk_type = "system"
  usage         = "image"
}

output "snapshot_ids" {
  value = data.alibabacloudstack_ecs_snapshots.default.ids
}
```

## 参数说明

以下参数是支持的：

* `instance_id` - (选填, 变更时重建) - 与快照关联的实例ID。如果指定，数据源将返回与该实例相关的快照。
* `disk_id` - (选填, 变更时重建) - 与快照关联的磁盘ID。如果指定，数据源将返回与该磁盘相关的快照。
* `ids` - (选填, 变更时重建) - 快照ID列表。如果指定，数据源将返回与这些ID匹配的快照。
* `name_regex` - (选填, 变更时重建) - 用于按快照名称筛选结果的正则表达式字符串。
* `status` - (选填, 变更时重建) - 快照的状态。取值范围：
  * `progressing`: 正在创建的快照。
  * `accomplished`: 创建成功的快照。
  * `failed`: 创建失败的快照。
  * `all` (默认): 所有快照状态。
* `type` - (选填, 变更时重建) - 快照的类别。取值范围：
  * `auto`: 自动快照。
  * `user`: 手动快照。
  * `all` (默认): 包括自动和手动快照。
* `source_disk_type` - (选填, 变更时重建) - 快照来源磁盘的类型。取值范围：
  * `system`: 系统盘。
  * `data`: 数据盘。
  > 取值不区分大小写。
* `usage` - (选填, 变更时重建) - 指定快照是否已被用来创建自定义镜像或磁盘。取值范围：
  * `image`: 快照已被用来创建自定义镜像。
  * `disk`: 快照已被用来创建磁盘。
  * `image_disk`: 快照已被用来创建自定义镜像和数据盘。
  * `none`: 快照未被用来创建自定义镜像或磁盘。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 快照ID列表。
* `names` - 快照名称列表。
* `snapshots` - 快照列表。每个元素包含以下属性：
  * `id` - 快照的ID。
  * `name` - 快照的名称。
  * `description` - 快照的描述。长度为2~256个英文或中文字符，不能以`http://`或`https://`开头。默认值为空。
  * `progress` - 快照创建进度，单位为百分比。
  * `source_disk_id` - 来源磁盘的ID。
  * `source_disk_size` - 来源磁盘的大小，单位为GB。
  * `source_disk_type` - 来源磁盘的类型。取值范围：
    * `system`: 系统盘。
    * `data`: 数据盘。
  * `product_code` - 从镜像市场继承的产品代码。
  * `remain_time` - 快照创建任务的剩余时间，单位为秒。
  * `creation_time` - 快照的创建时间。遵循ISO8601标准并使用UTC时间。格式：`YYYY-MM-DDThh:mmZ`。
  * `status` - 快照的状态。取值范围：
    * `progressing`: 正在创建的快照。
    * `accomplished`: 创建成功的快照。
    * `failed`: 创建失败的快照。
    * `all`: 表示所有快照状态。
  * `usage` - 快照是否已被用来创建资源。取值范围：
    * `image`: 快照已被用来创建自定义镜像。
    * `disk`: 快照已被用来创建磁盘。
    * `image_disk`: 快照已被用来创建自定义镜像和数据盘。
    * `none`: 快照未被用来创建自定义镜像或磁盘。
* `tags` - 分配给快照的标签映射。