---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpcclusters"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-hpcclusters"
description: |- 
  查询云服务器高性能计算集群
---

# alibabacloudstack_ecs_hpcclusters
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_hpc_clusters`

根据指定过滤条件列出当前凭证权限可以访问的云服务器高性能计算集群（HPC）列表。

## 示例用法

```terraform
resource "alibabacloudstack_ecs_hpc_cluster" "default" {
    name        = "tf-testAccEcsHpcClustersTest65057"
    description = "For Terraform Test"
}

data "alibabacloudstack_ecs_hpc_clusters" "default" {
    ids        = [alibabacloudstack_ecs_hpc_cluster.default.id]
    name_regex = "tf-testAcc.*"
    output_file = "hpc_clusters_output.txt"
}

output "first_hpc_cluster_id" {
    value = data.alibabacloudstack_ecs_hpc_clusters.default.clusters[0].id
}
```

## 参数说明

以下参数是支持的：

* `ids` - （可选，变更时重建）HPC 集群 ID 列表。用于通过指定的 HPC 集群 ID 过滤结果。
* `name_regex` - （可选，变更时重建）用于通过 HPC 集群名称过滤结果的正则表达式字符串。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - HPC 集群名称列表。
* `clusters` - ECS HPC 集群列表。每个元素包含以下属性：
    * `description` - ECS HPC 集群的描述信息。
    * `id` - HPC 集群的 ID，等同于 `hpc_cluster_id`。
    * `hpc_cluster_id` - HPC 集群的唯一标识符。
    * `name` - ECS HPC 集群的名称。