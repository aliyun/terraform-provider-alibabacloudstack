---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cluster_quotas"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-cluster-quotas"
description: |-
  查询Max Compute集群配额
---

# alibabacloudstack_maxcompute_cluster_quotas 数据源

根据指定过滤条件列出当前凭证权限可以访问的Max Compute集群配额信息


## 示例用法

```hcl
data "alibabacloudstack_maxcompute_cluster_quotas" "example" {
  cluster = "example-cluster"
}

output "cluster_quotas" {
  value = data.alibabacloudstack_maxcompute_cluster_quotas.example
}
```

## 参数说明
支持以下参数：

* `cluster` - (必填) 要获取配额的 MaxCompute 集群名称。此参数用于指定需要查询配额的具体 MaxCompute 集群。

## 属性说明
导出以下属性：

* `cluster` - MaxCompute 集群的名称，表示当前配额所属的集群。
* `cu_total` - 分配给集群的总计算单元 (CUs) 数量，表示该集群可以使用的最大计算资源。
* `disk_available` - 集群中的可用磁盘空间，表示当前未被占用的磁盘存储容量。
* `cu_available` - 集群中可用的计算单元 (CUs) 数量，表示当前未被占用的计算资源。
* `disk_total` - 分配给集群的总磁盘空间，表示该集群可以使用的最大磁盘存储容量。