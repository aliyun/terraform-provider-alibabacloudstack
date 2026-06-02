---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_instance_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-instance-types"
description: |-
  提供 PolarDB 集群实例类型列表。
---

# alibabacloudstack_polardb_cluster_instance_types

本数据源提供阿里云平台环境中的 PolarDB 集群实例类型列表。

## 示例用法

```hcl
data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type     = "MySQL"
  db_version  = "8.0"
  sorted_by   = "CPU"
  cpu_type    = "hygon"
  sub_category = "normal_exclusive"
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 实例类型 ID 列表。
* `db_version` - (可选) 数据库引擎版本。
* `cpu` - (可选) CPU 核心数量。
* `cpu_type` - (可选, 变更时强制重建) CPU 类型。有效值：`intel`、`arm64`、`hygon`。
* `memory` - (可选) 内存大小（GB）。
* `sorted_by` - (可选, 变更时强制重建) 排序方法。有效值：`CPU`、`Memory`。
* `sub_category` - (可选, 变更时强制重建) 子类别。有效值：`normal_general`、`normal_exclusive`。
* `db_type` - (可选, 变更时强制重建) 数据库引擎类型。有效值：`MySQL`、`PostgreSQL`、`Oracle`。

## 属性参考

导出以下属性：

* `ids` - 实例类型 ID 列表。
* `instance_types` - 实例类型列表。每个元素包含以下属性：
  * `id` - 实例类型的 ID。
  * `proxy_mem` - 代理内存。
  * `sub_category` - 子类别。
  * `cpu_type` - CPU 类型。
  * `memory` - 内存大小（GB）。
  * `gmt_modify` - 修改时间。
  * `spec_from` - 规格来源。
  * `proxy_cpu` - 代理 CPU。
  * `product_type` - 产品类型。
  * `db_node_class` - 数据库节点类别。
  * `product` - 产品。
  * `db_version` - 数据库版本。
  * `proxy_type` - 代理类型。
  * `db_node_num` - 数据库节点数量。
  * `db_type` - 数据库类型。
  * `cpu` - CPU 核心数量。
  * `uni_key` - 唯一密钥。
  * `proxy_class` - 代理类别。
  * `gmt_create` - 创建时间。
  * `region_id` - 区域 ID。
  * `status` - 状态。