---
subcategory: "表格存储 Tablestore"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ots-clusters"
description: |-
  获取阿里云表格存储（Tablestore）可用的集群列表。
---

# alibabacloudstack_ots_clusters

该数据源用于查询阿里云表格存储（Tablestore）服务中当前账号下可用的公共集群列表。您可以根据集群名称进行筛选，以获取特定的集群信息。

## 示例用法

### 获取所有公共集群
```hcl
data "alibabacloudstack_ots_clusters" "all" {}

output "cluster_names" {
  value = data.alibabacloudstack_ots_clusters.all.names
}
```

### 按名称正则表达式筛选集群
```hcl
data "alibabacloudstack_ots_clusters" "filtered" {
  name_regex = "^cn-"
}

output "matching_clusters" {
  value = data.alibabacloudstack_ots_clusters.filtered.clusters
}
```

### 指定具体集群名称列表
```hcl
data "alibabacloudstack_ots_clusters" "specific" {
  names = ["cn-hangzhou", "cn-shanghai"]
}

output "selected_clusters" {
  value = data.alibabacloudstack_ots_clusters.specific.clusters
}
```

## 参数说明

支持以下参数：

* `names` - (可选) 指定要查询的集群名称列表。仅返回名称在此列表中的集群。至少包含一个元素。
* `name_regex` - (可选) 使用正则表达式匹配集群名称。返回名称匹配该正则表达式的集群。必须是有效的正则表达式。

> **注意**：`names` 和 `name_regex` 可同时使用，此时将返回同时满足两个条件的集群（即交集）。

## 属性说明

除了上述参数外，该数据源还导出以下属性：

* `clusters` - 集群信息列表，每个元素包含以下字段：
  * `cluster_name` - 集群的唯一标识名称（如 `cn-hangzhou`）。
  * `cluster_type` - 集群类型（例如 `Public` 表示公共集群）。
  * `alias_name` - 集群的别名（通常为地域中文名称或描述性名称）。
  * `support_replica` - 是否支持副本功能（布尔值）。

* `names` - 所有匹配集群的名称列表（字符串列表），与 `clusters` 顺序一致。

## 注意事项

- 该数据源仅返回**公共集群**（Public Cluster），不包括私有集群。
- 查询结果依赖于当前账号在对应地域的权限和开通状态。
- `name_regex` 使用 Go 语言的正则表达式语法，需确保其有效性，否则 Terraform 会报错。