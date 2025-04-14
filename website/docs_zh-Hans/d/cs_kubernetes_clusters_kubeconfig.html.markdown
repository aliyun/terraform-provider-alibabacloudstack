---
subcategory: "Container Service (CS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cs_kubernetes_clusters_kubeconfig"
sidebar_current: "docs-alibabacloudstack-datasource-cs-kubernetes-clusters-kubeconfig"
description: |-
  查询K8s集群的配置信息
---

# alibabacloudstack_cs_kubernetes_clusters_kubeconfig

根据指定过滤条件列出当前凭证权限可以访问的容器服务中K8s集群的配置信息。


## 示例用法

```
# 声明数据源
data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
  cluster_id = "xx"
}
output "kubeconfig" {
  value = data.alibabacloudstack_cs_kubernetes_clusters_kubeconfig.k8s_clusters_kubeconfig.kubeconfig
}
```

## 参数说明

支持以下参数：

* `cluster_id` - (可选) 集群 ID，用于指定要查询的 Kubernetes 集群。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `kubeconfig` - Kubernetes 集群的 kubeconfig 内容。该内容包括以下详细信息：
  
  ### 集群(Clusters)
  - **name**: 集群名称，在上下文中引用。
  - **cluster**:
    - **server**: Kubernetes API 服务器的 URL。
    - **certificate-authority**: 用于验证 API 服务器身份的 CA 证书路径。
    - **insecure-skip-tls-verify**: 如果设置为 true，则跳过 TLS 验证（不推荐）。

  ### 用户(Users)
  - **name**: 用户名，在上下文中引用。
  - **user**:
    - **client-certificate**: 用于向 API 服务器进行身份验证的客户端证书路径。
    - **client-key**: 与客户端证书对应的私钥路径。
    - **token**: 如果使用基于令牌的身份验证，这是提供令牌字符串的位置。
    - **username/password**: 对于某些简单的身份验证机制，可能需要用户名和密码。

  ### 上下文(Contexts)
  - **name**: 上下文名称，可以任意命名。
  - **context**:
    - **cluster**: 关联的集群名称。
    - **user**: 用于访问集群的用户。
    - **namespace**: 默认命名空间。 