---
subcategory: "EDAS"
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

* `cs_cluster_id` - (必填，变更时重建) 要导入的阿里云容器服务 Kubernetes 集群的 ID。
* `namespace_id` - (可选，变更时重建) 您要导入的命名空间的 ID。您可以调用 [ListUserDefineRegion](https://www.alibabacloud.com/help/en/doc-detail/149377.htm?spm=a2c63.p38356.879954.34.331054faK2yNvC#doc-api-Edas-ListUserDefineRegion) 操作查询命名空间 ID。
* `vpc_id` - (可选，变更时重建) 集群所属的虚拟私有云（VPC）的 ID。

## 属性说明

导出以下属性：

* `cluster_name` - 要创建的集群名称。
* `cluster_type` - 要创建的集群类型。有效值仅：5: K8s 集群。
* `network_mode` - 要创建的集群的网络类型。有效值：1: 经典网络。2: VPC。
* `region_id` - 区域 ID。
* `vpc_id` - 集群的虚拟私有云(VPC)ID。
* `cluster_import_status` - 集群的导入状态：
    * `1`: 成功。
    * `2`: 失败。
    * `3`: 正在导入。
    * `4`: 已删除。
* `cs_cluster_id` - 要导入的阿里云容器服务 Kubernetes 集群的 ID。
* `namespace_id` - 要导入的命名空间的 ID。

## 导入

EDAS 集群可以使用 id 导入，例如：

```bash
$ terraform import alibabacloudstack_edas_k8s_cluster.cluster cluster_id
```