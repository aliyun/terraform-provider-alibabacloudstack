---
subcategory: "Enterprise Distributed Application Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_cluster"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-cluster"
description: |-
  编排绑定企业级分布式应用服务（Edas）k8s集群
---

# alibabacloudstack_edas_k8s_cluster

使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）k8s集群。
有关 EDAS K8s 集群的更多信息以及如何使用它，请参阅[什么是 EDAS K8s 集群](https://www.alibabacloud.com/help/en/doc-detail/85108.htm)。



## 示例用法

### 基础用法

```
resource "alibabacloudstack_edas_k8s_cluster" "default" {
  cs_cluster_id = "xxxx-xxx-xxx"
}
```

## 参数说明

支持以下参数：

* `cs_cluster_id` - (必填，变更时重建) 要导入的容器服务 Kubernetes 集群的 ID。您可以调用 [GetK8sCluster](https://www.alibabacloud.com/help/en/doc-detail/85108.htm) 接口查询集群 ID。
* `namespace_id` - (可选，变更时重建) 您要导入的命名空间的 ID。您可以调用 [ListUserDefineRegion](https://www.alibabacloud.com/help/en/doc-detail/149377.htm) 接口查询命名空间 ID。

## 属性说明

导出以下属性：

* `id` - EDAS K8s 集群的 ID。
* `cluster_name` - 集群名称。
* `cluster_type` - 集群类型。有效值：`5`：容器服务 K8s 集群或 Serverless K8s 集群。
* `network_mode` - 集群的网络类型。有效值：`1`：经典网络。`2`：VPC。
* `vpc_id` - 集群的虚拟私有云（VPC）ID。
* `cluster_import_status` - 集群的导入状态。有效值：
    * `1`：成功。
    * `2`：失败。
    * `3`：正在导入。
    * `4`：已删除。
* `cs_cluster_id` - 要导入的容器服务 Kubernetes 集群的 ID。
* `namespace_id` - 要导入的命名空间的 ID。

## 导入

EDAS K8s 集群可以使用集群 ID（导入后返回的 EDAS 内部集群 ID）导入，例如：

```bash
$ terraform import alibabacloudstack_edas_k8s_cluster.example 81453e4b-4df0-4592-xxxx-b835a2eexxxx
```