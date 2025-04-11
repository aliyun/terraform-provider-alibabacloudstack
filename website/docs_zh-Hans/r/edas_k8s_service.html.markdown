---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_service"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-service"
description: |-
  编排绑定企业级分布式应用服务（Edas）k8s服务
---

# alibabacloudstack_edas_k8s_service
使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）k8s服务。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_edas_k8s_application" "default" {
  // package type 是 Image / FatJar / War
  package_type            = "Image"
  application_name        = "DemoApplication"
  application_description = "This is description of application"
  cluster_id              = var.cluster_id
  replicas                = 2

  // 当 package_type 是 'image' 时设置 'image_url' 和 'repo_id'
  image_url = "registry-vpc.cn-beijing.aliyuncs.com/edas-demo-image/consumer:1.0"

  // 当 package_type 不是 'image' 时设置 'package_url','package_version' 和 'jdk'
  package_url     = var.package_url
  package_version = var.package_version
  jdk             = var.jdk

  // 当 package_type 是 'war' 时设置 'web_container' 和 'edas_container'
  web_container          = var.web_container
  edas_container_version = var.edas_container_version

  internet_target_port  = var.internet_target_port
  internet_slb_port     = var.internet_slb_port
  internet_slb_protocol = var.internet_slb_protocol
  internet_slb_id       = var.internet_slb_id
  limit_mem             = 2048
  requests_mem          = 0
  requests_m_cpu        = 0
  limit_m_cpu           = 4000
  command               = var.command
  command_args          = var.command_args
  envs                  = var.envs
  pre_stop              = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  post_start            = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  liveness              = var.liveness
  readiness             = var.readiness
  nas_id                = var.nas_id
  mount_descs           = var.mount_descs
  local_volume          = var.local_volume
  namespace             = "default"
  logical_region_id     = "cn-beijing"
}

resource "alibabacloudstack_edas_k8s_service" "default" {
  app_id = alibabacloudstack_edas_k8s_application.default.id
  service_name = "tf-testAccEdasK8sService"
  type = "NodePort"
  external_traffic_policy = "Local"
  port_mappings {
    protocol = "TCP"
    service_port = 80
    target_port = 8080
  }
}
```

## 参数参考

支持以下参数：

* `app_id` - (必填，变更时重建) Edas k8s 应用程序的 id。
* `service_name` - (必填，变更时重建) Edas k8s 服务的名称。
* `type` - (必填，变更时重建) Edas K8s 服务类型，有效值为 `ClusterIP`, `NodePort`, `LoadBalancer`。
* `port_mappings` - (可选) K8s 服务端口映射表，需要符合 JsonArray 格式。支持的参数如下：
  * `protocol` - (必填) 服务协议，支持 TCP 和 UDP。
  * `service_port` - (必填) 前端服务端口，取值范围为 1~65535。
  * `target_port` - (必填) 后端容器端口，取值范围为 1~65535。
* `external_traffic_policy` - (可选) 当服务类型为 `NodePort` 或 `LoadBalancer` 时，设置外部流量管理策略。有效值为 `Local`, `Cluster`，默认为 `Local`。
* `annotations` - (可选) 服务的注解映射。
* `labels` - (可选) 服务的标签映射。
* `allow_edit` - (可选) 表示该服务是否允许编辑。

## 属性说明

导出以下属性：

* `cluster_ip` - Kubernetes 的集群 IP。
* `inner_endpointer` - 服务的内部终结点。
* `namespace` - K8s 集群的命名空间。
* `nodeip_list` - 服务的节点 IP 列表。
* `allow_edit` - 表示该服务是否允许编辑。

## 导入

EDAS k8s 应用程序可以按以下方式导入，例如：

```bash
$ terraform import alibabacloudstack_edas_k8s_service.new_k8s_service "app_id:name"
```