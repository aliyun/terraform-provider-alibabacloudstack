---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxy_types"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-cluster-proxy-types"
description: |-
  提供 PolarDB 集群代理类型列表。
---

# alibabacloudstack_polardb_cluster_proxy_types

该数据源提供专有云环境中 PolarDB 集群代理类型列表。


## 示例

```hcl
data "alibabacloudstack_polardb_cluster_proxy_types" "default" {
  db_type    = "MySQL"
  db_version = "5.7"
}
```

## 参数说明

以下参数适用于该数据源：

* `db_type` -（可选，ForceNew）数据库引擎类型。有效值：`MySQL`、`PostgreSQL`、`Oracle`。
* `db_version` -（可选）数据库引擎版本。
* `ids` -（可选，ForceNew）用于过滤结果的代理类型 ID 列表。
* `core_count` -（可选）CPU 核数。

## 属性说明

以下属性由该数据源导出：

* `ids` - 代理类型 ID 列表。
* `proxy_classes` - 代理类型列表。每个元素包含以下属性：
  * `id` - 代理类型 ID。
  * `core_count` - CPU 核数。

## Import

数据源不支持 Import 操作。
