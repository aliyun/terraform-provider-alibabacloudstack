---
subcategory: "文件存储 NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_namespace_group"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-namespace-group"
description: |-
  查询阿里云NAS跨域挂载编排信息
---

# alibabacloudstack_nas_namespace_group

> NAS 跨域挂载编排

## 示例用法

```hcl

variable "name" {
  default = "tf-testnasng4165"
}

data "alibabacloudstack_nas_zones" "default" {
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  tags = {
    common_test = "terraform"
    filter      = var.name
  }
  enable_ipv6 = true
  lifecycle {
    ignore_changes = [
      secondary_cidr_blocks,
      tags
    ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  enable_ipv6  = true
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}



resource "alibabacloudstack_nas_namespace" "default" {
  zone_id       = data.alibabacloudstack_nas_zones.default.zones.0.zone_id
  cluster_id    = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.cluster_id
  description   = var.name
  storage_type  = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.storage_type
  protocol_type = data.alibabacloudstack_nas_zones.default.zones.0.clusters.0.instance_types.0.protocol_type
  encrypt_type  = "0"
}

resource "alibabacloudstack_nas_accessgroup" "default" {
  access_group_name = var.name
  access_group_type = "Vpc"
}

resource "alibabacloudstack_nas_namespace_mount_target" "default" {
  network_type      = "Vpc"
  access_group_name = alibabacloudstack_nas_accessgroup.default.access_group_name
  nas_namespace_id  = alibabacloudstack_nas_namespace.default.id
  vswitch_id        = alibabacloudstack_vpc_vswitch.default.id
}

resource "alibabacloudstack_nas_namespace_group" "default" {
  network_type        = "Vpc"
  nas_namespace_id    = alibabacloudstack_nas_namespace.default.id
  mount_target_domain = alibabacloudstack_nas_namespace_mount_target.default.mount_target_domain
  mapped_path         = var.name
}



data "alibabacloudstack_nas_namespace_groups" "default" {
  ids = [
    "${alibabacloudstack_nas_namespace_group.default.id}"
  ]
}
```

## 参数说明
以下参数可用于过滤查询结果：

* `mount_target_domain` (字符串, 可选)：用于按挂载目标域名过滤结果，用于挂载NAS文件系统的域名。

* `name_regex` (字符串, 可选)：用于按挂载点上的路径过滤结果的正则表达式。

* `nas_namespace_id` (字符串, 可选)：用于按NAS命名空间ID过滤结果。

* `network_type` (字符串, 可选)：用于按网络类型过滤结果，可以是Vpc（专有网络）或Classic（经典网络）。

## 属性说明
以下属性被导出：

* `id` (字符串)：数据源的唯一标识符，基于查询到的命名空间ID列表生成。

* `create_time` (字符串)：命名空间组的创建时间，遵循ISO 8601标准格式。

* `ids` (列表)：NAS命名空间组ID列表。

* `mapped_path` (字符串)：映射路径，表示NAS命名空间在挂载点上的路径。

* `member_id` (字符串)：成员ID，表示命名空间组中的成员标识。

* `mount_target_domain` (字符串)：挂载目标域名，用于挂载NAS文件系统，已去除末尾的点号。

* `nas_namespace_id` (字符串)：NAS命名空间ID。

* `names` (列表)：NAS命名空间组名称列表（使用NasNamespaceId作为名称）。

* `network_type` (字符串)：网络类型，可以是Vpc（专有网络）或Classic（经典网络）。

* `status` (字符串)：命名空间组的状态，可能的值包括：Enabled（启用）。

* `groups` (列表)：NAS命名空间组列表。每个元素包含命名空间组的属性。

`groups`列表中的每个对象包含以下属性：

* `create_time` (字符串)：命名空间组的创建时间。

* `mapped_path` (字符串)：映射路径。

* `member_id` (字符串)：成员ID。

* `mount_target_domain` (字符串)：挂载目标域名。

* `nas_namespace_id` (字符串)：NAS命名空间ID。

* `network_type` (字符串)：网络类型。

* `status` (字符串)：命名空间组的状态。