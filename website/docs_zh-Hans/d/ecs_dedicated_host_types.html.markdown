---
subcategory: "云服务器 ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_dedicated_host_types"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-dedicated-host-types"
description: |-
  查询专有宿主机规格类型
---

# alibabacloudstack_ecs_dedicated_host_types

-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_dedicatedhost_types`

根据指定过滤条件查询可用的专有宿主机规格类型列表。

## 示例用法

```hcl
# 查询所有可用的专有宿主机规格
data "alibabacloudstack_ecs_dedicated_host_types" "default" {
}

# 输出第一个专有宿主机规格ID
output "first_ddh_type_id" {
  value = "${data.alibabacloudstack_ecs_dedicated_host_types.default.ddh_types.0.id}"
}

# 根据可用区过滤专有宿主机规格
data "alibabacloudstack_ecs_dedicated_host_types" "filtered" {
  availability_zone = "cn-hangzhou-a"
}
```

## 参数说明

以下参数是支持的：

* `availability_zone` - (选填, 变更时重建) - 可用区ID。用于过滤指定可用区内可用的专有宿主机规格。
* `ids` - (选填, 变更时重建, 支持查询) - 专有宿主机规格ID列表，用于过滤结果。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 专有宿主机规格ID列表。
* `ddh_types` - 专有宿主机规格类型列表。每个元素包含以下属性：
  * `id` - 专有宿主机规格ID，例如 `ddh.g5`、`ddh.c5` 等。
  * `availability_zones` - 该规格可用的可用区ID列表。
