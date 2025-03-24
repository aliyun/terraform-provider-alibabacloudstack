---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_cluster"
sidebar_current: "docs-Alibabacloudstack-edas-cluster"
description: |- 
  编排企业级分布式应用服务（Edas）集群
---

# alibabacloudstack_edas_cluster

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas）集群。

## 示例用法

以下是一个完整的 HCL 配置示例，演示如何使用 `alibabacloudstack_edas_cluster` 资源及所需参数：

```hcl
resource "alibabacloudstack_edas_cluster" "default" {
  cluster_name      = var.cluster_name
  cluster_type      = var.cluster_type
  network_mode      = var.network_mode
  logical_region_id = var.logical_region_id
  vpc_id            = var.vpc_id
}
```

## 参数参考

支持以下参数：

* `cluster_name` - (必填，变更时重建) 集群名称。它必须在阿里巴巴云账户内唯一，并且创建后无法修改。
* `cluster_type` - (必填，变更时重建) 集群类型。有效值：
  * `0`: 普通 Docker 集群。
  * `1`: Swarm 集群。
  * `2`: ECS 集群。
  * `3`: EDAS 自建 Kubernetes 集群。
  * `4`: Pandora 自动注册应用集群类型。
  * `5`: 容器服务 Kubernetes 集群。
* `network_mode` - (必填，变更时重建) 网络类型。有效值：
  * `1`: 经典网络。
  * `2`: VPC。
* `logical_region_id` - (可选，变更时重建) 集群所在的逻辑区域 ID。您可以调用 `ListUserDefineRegion` 操作查询逻辑区域 ID。
* `vpc_id` - (可选，变更时重建) VPC 网络 ID。如果 `network_mode` 设置为 `2`(VPC)，此参数是必填的。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - EDAS 集群的唯一标识符。它被制定为 `<cluster_id>`。