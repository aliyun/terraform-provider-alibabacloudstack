---
subcategory: "分布式关系型数据库"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_instance_types"
sidebar_current: "docs-alibabacloudstack-rds-instance-types"
description: |-
  查询RDS实例规格
---

# alibabacloudstack_rds_instance_types

根据指定过滤条件列出当前凭证权限可以访问的RDS实例规格列表。

## 示例用法

```
data "alibabacloudstack_rds_instance_types" "default" {
  engine               = "MySQL"
  engine_version       = "5.7"
  sorted_by            = "CPU"
  series               = "dual_ha"
}
```

## 参数说明

支持以下参数：

* `ids` - (可选，变更时重建) 指定实例规格ID范围。如果未指定，则返回所有可用区的实例规格。
* `engine` - (可选，变更时重建)筛选特定例引擎的结果。 有效值： `PostgreSQL`, `MySQL`，`POLARDB`。
* `engine_version` - (可选，变更时重建) 筛选特定例引擎版本的结果。
* `series` - (可选) 筛选特定例规格系列的结果。
* `cpu` - (可选) 筛选特定数量CPU核心的结果。
* `cpu_type` - (可选) 筛选特定数量CPU类型的结果。
* `memory` - (可选) 筛选特定内存大小(GB)的结果。
* `sorted_by` - (可选，强制更新) 排序模式，有效值：`CPU`, `Memory`。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 包含所有匹配条件的实例规格ID的列表。
* `names` - 包含所有匹配条件的实例规格名称的列表。
* `instance_types` - 实例规格的详细信息列表。每个元素包含以下属性：
  * `id` - 实例规格的唯一标识符。
  * `name` - 实例规格的名称。
  * `cpu` - CPU核心数量。
  * `memory` - 内存大小，以GB为单位。
  * `engine` - 引擎类型。
  * `engine_version` - 引擎版本。
  * `cpu_type` - CPU类型。
  * `series` - 隶属的实例规格系列的ID。
  * `connections` - 最大链接数。
  * `storage_type` - 存储类型。
  * `storage_min` - 最小存储规格。
  * `storage_max` - 最大存储规格。
