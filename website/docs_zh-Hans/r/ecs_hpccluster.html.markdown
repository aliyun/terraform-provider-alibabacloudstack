---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpc_cluster"
sidebar_current: "docs-Alibabacloudstack-ecs-hpccluster"
description: |- 
  编排云服务器（ECS）高性能计算集群（HPC）
---

# alibabacloudstack_ecs_hpc_cluster
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_hpccluster`

使用Provider配置的凭证在指定的资源集下编排云服务器（ECS）高性能计算集群（HPC）。

有关 ECS 高性能计算集群的更多信息以及如何使用它，请参阅 [什么是高性能计算集群](https://www.alibabacloud.com/help/en/doc-detail/109138.htm)。

## 示例用法

### 基础用法

```terraform
variable "name" {
    default = "tf-testaccecshpc_cluster26697"
}

resource "alibabacloudstack_ecs_hpc_cluster" "example" {
  name        = var.name
  description = "For Terraform Test"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) 高性能计算集群的名称。长度为2~128个英文或中文字符。必须以大小写字母或中文开头，不能以`http://`和`https://`开头。可以包含数字、英文句号（.）、下划线（_）或者短划线（-）。此名称在同一地域内必须唯一。
* `description` - (可选) 高性能计算集群的描述信息。长度为2~256个英文或中文字符，不能以`http://`和`https://`开头。可以包含大写/小写字母、数字、句点(.)、冒号(:)、下划线(_)、连字符(-)和at符号(@)。默认值：空。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - 高性能计算集群的ID。
* `hpc_cluster_id` - (输出) HPC集群ID。
* `name` - (输出) HPC集群名称。
* `description` - (输出) HPC集群的描述信息。

## Import

ECS HPC Cluster 可以通过 HpcClusterId 导入，例如：

```
$ terraform import alibabacloudstack_ecs_hpc_cluster.example hpc-bp1a5zr3u7nq9cx****
```