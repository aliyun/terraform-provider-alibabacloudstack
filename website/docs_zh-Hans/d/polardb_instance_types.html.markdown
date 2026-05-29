---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-instance-types"
description: |-
  查询polardb数据库实例规格。
---

# alibabacloudstack_polardb_instance_types

根据指定过滤条件列出当前凭证权限可以访问的polardb数据库实例规格列表。

## Example Usage

```hcl
data "alibabacloudstack_polardb_instance_types" "example" {
  engine        = "MySQL"
  engine_version = "8.0"
  cpu           = 4
  memory        = 8
}

output "instance_types" {
  value = data.alibabacloudstack_polardb_instance_types.example.instance_types
}
```
## 参数说明
以下参数是支持的：

  * `ids` - (Optional) 用于通过ID列表过滤实例类型
  * `engine` - (Optional) 指定数据库引擎类型（支持PolarDB_PG、MySQL、PolarDB_PPAS）
  * `engine_version` - (Optional) 数据库引擎版本，当指定engine时必填
  * `cpu` - (Optional) CPU核心数
  * `cpu_type` - (Optional) CPU架构（支持intel、arm64）
  * `memory` - (Optional) 内存大小（GB）
  * `sorted_by` - (Optional) 排序方式（支持按CPU或Memory排序）
  * `series` - (Optional) 实例系列（支持dual_ha、read_only）

## 属性说明
除了上述参数外，还导出以下属性：

  * `instance_types` - 实例类型详细信息列表，包含以下属性
    * `id` - 实例类型唯一标识符
    * `cpu` - CPU核心数
    * `memory` - 内存大小（GB）
    * `engine` - 数据库引擎
    * `engine_version` - 数据库引擎版本
    * `cpu_type` - CPU架构
    * `series` - 实例系列
    * `connections` - 最大并发连接数
    * `storage_type` - 存储类型
    * `storage_min` - 最小存储容量
    * `storage_max` - 最大存储容量