---
subcategory: "消息队列 RocketMQ 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_clusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ons-clusters"
description: |-
    查询消息队列集群列表
---

# alibabacloudstack_ons_clusters

根据指定过滤条件列出当前可用区支持的消息队列集群列表。

## 示例用法

```hcl
data "alibabacloudstack_ons_clusters" "example" {
  ids = ["cluster1"]
}

output "first_cluster_id" {
  value = data.alibabacloudstack_ons_clusters.example.clusters[0].id
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 用于过滤结果的集群ID列表。
* `name_regex` - (可选) 用于通过集群名称过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - 集群ID列表。
* `clusters` - 集群列表。每个元素包含以下属性：
  * `id` - 集群ID。
  * `name` - 集群名称。
  * `cpu_brand` - 集群节点的CPU品牌。
  * `cpu_arch` - 集群节点的CPU架构。
