---
subcategory: "资源访问管理"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ram_role_attachment"
sidebar_current: "docs-Alibabacloudstack-ecs-ramroleattachment"
description: |-
  编排绑定云服务器（ECS）实例和RAM角色
---

# alibabacloudstack_ram_role_attachment
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_ramroleattachment`

使用 Provider 配置的凭证在指定的资源集下编排绑定云服务器（ECS）实例和 RAM 角色。

## 示例用法

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
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
  default_instance_type_id = try(
    element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0),
    sort(data.alibabacloudstack_instance_types.all.ids)[0]
  )
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

variable "name" {
  default = "Test_ram_role_attachment"
}

resource "alibabacloudstack_vpc" "default" {
  name        = var.name
  cidr_block  = "192.168.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "192.168.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_security_group" "default" {
  name   = var.name
  vpc_id = alibabacloudstack_vpc.default.id
}

resource "alibabacloudstack_instance" "default" {
  image_id           = data.alibabacloudstack_images.default.images.0.id
  instance_type      = local.default_instance_type_id
  instance_name      = var.name
  security_groups    = [alibabacloudstack_security_group.default.id]
  availability_zone  = data.alibabacloudstack_zones.default.zones[0].id
  system_disk_category = "cloud_pperf"
  system_disk_size  = 100
  vswitch_id        = alibabacloudstack_vswitch.default.id
}

data "alibabacloudstack_ascm_ram_service_roles" "role" {
  product = "ecs"
}

resource "alibabacloudstack_ram_role_attachment" "default" {
  role_name    = data.alibabacloudstack_ascm_ram_service_roles.role.roles.0.name
  instance_ids = [alibabacloudstack_instance.default.id]
}
```

## 参数说明

支持以下参数：

* `role_name` - （必填，ForceNew）要附加的 RAM 角色名称。该名称必须在 1 到 64 个字符之间，并且只能包含字母数字字符或连字符（`-`、`_`）。它不能以连字符开头。
* `instance_ids` - （必填，ForceNew）要附加 RAM 角色的 ECS 实例 ID 列表。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 资源的 ID。格式为 `<role_name>:<instance_ids>`，其中 `instance_ids` 是用双引号括起来的逗号分隔的 ECS 实例 ID 列表。

## Import

RAM 角色附加资源可以使用 role_name 和 instance_ids（用冒号连接）进行导入，例如：

```
$ terraform import alibabacloudstack_ram_role_attachment.example my-ram-role:i-uf6257yd1w3k1example
```
