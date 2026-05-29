---
subcategory: "PolarDB-X"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-instance-types"
description: |-
  查询PolarDBX数据库实例规格。
---

# alibabacloudstack_polardbx_instance_types

根据指定过滤条件列出当前凭证权限可以访问的PolarDB-X数据库实例规格列表。

## Example Usage

```hcl
data "alibabacloudstack_polardbx_instance_types" "example" {
  engine_version = "8.0"
  cpu            = 4
  memory         = 16
}

output "instance_types" {
  value = data.alibabacloudstack_polardbx_instance_types.example.instance_types
}
```


## 参数说明
以下参数是支持的：

* `ids` - (Optional) 用于通过ID列表过滤实例类型
* `spec_series` - (Optional) 规格系列, 支持SHARE (共享型)、SINGLE (独享型)
* `engine_version` - (Optional) 数据库引擎版本（支持5.7、8.0）
* `cpu` - (Optional) CPU核心数
* `cpu_type` - (Optional) CPU架构
* `spec_type` - (Optional) 规格类型, 支持DN (数据节点)、CN (计算节点)
* `memory` - (Optional) 内存大小（GB）
* `sorted_by` - (Optional) 排序方式（支持按CPU或Memory排序）
* `series` - (Optional) 实例系列，支持enterprise (企业版)、standard (标准版)

## 属性说明
除了上述参数外，还导出以下属性：

* `instance_types` - 实例类型详细信息列表，包含以下属性
    * `id` - 实例类型唯一标识符
    * `cpu` - CPU核心数
    * `memory` - 内存大小（GB）