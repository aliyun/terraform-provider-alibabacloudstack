---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_cdc_classes"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-cdc-classes"
description: |-
  查询PolarDBX实例可用的CDC（Change Data Capture）规格。
---

# alibabacloudstack_polardbx_cdc_classes

查询指定PolarDB-X实例可用的CDC（Change Data Capture，变更数据捕获）规格列表。CDC规格定义了可用于CDC节点的CPU和内存配置。

## Example Usage

```hcl
data "alibabacloudstack_polardbx_cdc_classes" "example" {
  instance_id = "pxc-example-instance"
  cpu         = 4
  memory      = 16
}

output "cdc_classes" {
  value = data.alibabacloudstack_polardbx_cdc_classes.example.cdc_classes
}
```

## 参数说明
以下参数是支持的：

* `instance_id` - （必选）PolarDB-X实例ID。
* `ids` - （可选）用于过滤CDC规格ID列表。
* `cpu` - （可选）用于过滤CDC规格的CPU核心数。
* `memory` - （可选）用于过滤CDC规格的内存大小（单位：GB）。
* `sorted_by` - （可选）返回的CDC规格的排序方式。取值：`CPU`、`Memory`。

## 属性说明
除了上述参数外，还导出以下属性：

* `ids` - CDC规格ID列表。
* `cdc_classes` - CDC规格配置列表，包含以下属性：
    * `id` - CDC规格唯一标识符（ClassCode）。
    * `cpu` - CPU核心数。
    * `memory` - 内存大小（单位：GB）。
