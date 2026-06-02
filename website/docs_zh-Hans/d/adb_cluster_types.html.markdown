---
subcategory: "云原生数据库 MySQL版 (ADB)"
layout: "alibabacloudstack"
page_title: "ApsaraStack: alibabacloudstack_adb_cluster_types"
sidebar_current: "docs-apsarastack-data-source-adb-cluster-types"
description: |-
  提供ADB集群类型列表。
---

# alibabacloudstack_adb_cluster_types

该数据源用于获取专有云中可用的ADB（AnalyticDB）集群类型。


## 示例

```hcl
data "alibabacloudstack_adb_cluster_types" "example" {
  cpu_type = "Intel"
  cpu      = 16
  memory   = 64
}

output "adb_cluster_types" {
  value = data.alibabacloudstack_adb_cluster_types.example.instance_types
}
```

## 参数说明

以下参数支持配置：

* `cpu_type` - (可选) ADB集群的CPU类型。
* `cpu` - (可选) CPU核心数。
* `memory` - (可选) 内存大小（单位：GB）。
* `sorted_by` - (可选, ForceNew) 排序字段。有效值：`CPU`、`Memory`。
* `status` - (可选) 集群类型的状态。
* `cluster_type` - (可选) 集群类型。
* `ids` - (可选, ForceNew) 集群类型ID列表。

## 属性参考

以下属性会被导出：

* `ids` - 集群类型ID列表。
* `instance_types` - ADB集群类型列表。每个元素包含以下属性：
  * `id` - 集群类型的ID。
  * `cpu_type` - CPU类型。
  * `cpu` - CPU核心数。
  * `memory` - 内存大小（单位：GB）。
  * `storage_type` - 存储类型。
  * `storage_min` - 最小存储空间（单位：GB）。
  * `storage_max` - 最大存储空间（单位：GB）。
  * `mode` - 集群模式。
  * `node_min` - 最小节点数。
  * `node_max` - 最大节点数。
  * `cluster_type` - 集群类型。
  * `status` - 集群类型的状态。
  * `series` - 系列。
  * `cluster_category` - 集群类别。
