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

## 参数说明
支持以下参数：

* `ids` - (可选) 用于过滤结果的集群ID列表。此参数允许用户通过指定一个或多个集群ID来筛选查询结果。
* `name_regex` - (可选) 用于按名称过滤集群的正则表达式模式。此参数允许用户通过匹配集群名称的正则表达式来筛选查询结果。

## 属性说明
导出以下属性：

* `ids` - 集群的ID列表。此属性包含所有符合过滤条件的集群ID。
* `clusters` - 集群列表。每个元素包含以下属性：
    * `cluster` - 集群名称。表示当前集群的唯一标识名称。
    * `core_arch` - 集群的核心架构。描述了集群所使用的基础硬件架构或计算平台。
    * `project` - 与集群关联的项目。表示该集群所属的MaxCompute项目。
    * `region` - 集群所在的区域。表示该集群部署的具体阿里云区域。