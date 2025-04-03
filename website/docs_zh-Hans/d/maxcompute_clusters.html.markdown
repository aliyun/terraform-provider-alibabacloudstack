---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_clusters"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-clusters"
description: |-
  查询Max Compute集群
---

# alibabacloudstack_maxcompute_clusters

根据指定过滤条件列出当前凭证权限可以访问的Max Compute集群列表。[什么是Cluster](https://www.alibabacloud.com/help/en/maxcompute)


## 示例用法

```hcl
data "alibabacloudstack_maxcompute_clusters" "example" {
  name_regex = "example-cluster"
}

output "clusters" {
  value = data.alibabacloudstack_maxcompute_clusters.example.clusters
}
```

## 参数参考
支持以下参数：

* `ids` - (可选) 用于过滤结果的集群ID列表。
* `name_regex` - (可选) 用于按名称过滤集群的正则表达式模式。

## 属性参考
导出以下属性：

* `ids` - 集群的ID列表。
* `clusters` - 集群列表。每个元素包含以下属性：
    * `cluster` - 集群名称。
    * `core_arch` - 集群的核心架构。
    * `project` - 与集群关联的项目。
    * `region` - 集群所在的区域。