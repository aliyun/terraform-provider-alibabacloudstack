---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application"
sidebar_current: "docs-Alibabacloudstack-resource-edas-k8s-application"
description: |-
  编排绑定企业级分布式应用服务（Edas）k8s应用程序
---

# alibabacloudstack_edas_k8s_application

使用Provider配置的凭证在指定的资源集下编排绑定企业级分布式应用服务（Edas）k8s应用程序。
有关 EDAS K8s 应用程序的详细信息和如何使用它，请参阅 [什么是 EDAS K8s 应用程序](https://www.alibabacloud.com/help/doc-detail/85029.htm)。


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_edas_k8s_application" "default" {
  // package type is Image / FatJar / War
  package_type            = "Image"
  application_name        = "DemoApplication"
  application_description = "This is description of application"
  cluster_id              = var.cluster_id
  replicas                = 2

  // set 'image_url' and 'cr_ee_repo_id' when package_type is 'image'
  image_url = "registry-vpc.cn-beijing.aliyuncs.com/edas-demo-image/consumer:1.0"

  // set 'package_url','package_version' and 'jdk' when package_type is not 'image'
  package_url     = var.package_url
  package_version = var.package_version
  jdk             = var.jdk

  // set 'web_container' and 'edas_container' when package_type is 'war'
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
```

### 亲和性用法
```terraform
resource "alibabacloudstack_edas_k8s_application" "default" {
  package_type            = "FatJar"
  application_name        = "terraform-test-fatjar"
  application_description = "This is description of description"
  cluster_id              = "xxxxxxxxxxxxxxxxxxx"
  replicas                = 2

  package_url     = "xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx"
  package_version = "2025-07-09 10:00:18"
  jdk             = "Open JDK 8"

  command      = "/bin/sh"
  command_args = ["-c", "sleep 1001"]
  pre_stop     = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  post_start   = "{\"exec\":{\"command\":[\"ls\",\"/\"]}}"
  namespace    = "default"

  // 容忍节点上的 disk-pressure 污点。
  custom_tolerations {
    key      = "node.kubernetes.io/disk-pressure"
    operator = "Exists"
    effect   = "NoSchedule"
  }

  // 必须满足的节点亲和性：排除控制面节点，调度到 worker 节点。
  custom_node_affinity_require {
    match_expressions {
      key      = "node-role.kubernetes.io/control-plane"
      operator = "DoesNotExist"
    }
  }

  // 尽量满足的节点亲和性：优先选择 Linux 节点。
  custom_node_affinity_preferred {
    weight = 100
    match_expressions {
      key      = "kubernetes.io/os"
      values   = ["linux"]
      operator = "In"
    }
  }

  // 尽量满足的 Pod 亲和性：优先调度到已运行 EDAS 应用 Pod 的节点。
  custom_pod_affinity_preferred {
    weight        = 1
    k8s_namespace = ["default"]
    topology_key  = "kubernetes.io/hostname"
    match_expressions {
      key      = "edas.component"
      values   = ["app"]
      operator = "In"
    }
    match_expressions {
      key      = "edas.controlplane"
      values   = ["edas-oam"]
      operator = "In"
    }
  }

  // 尽量满足的 Pod 反亲和性：将 Pod 分散调度，避免与同一应用的其他实例调度到同一节点。
  custom_pod_ant_affinity_preferred {
    weight        = 1
    k8s_namespace = ["default"]
    topology_key  = "kubernetes.io/hostname"
    match_expressions {
      key      = "edas.component"
      values   = ["app"]
      operator = "In"
    }
    match_expressions {
      key      = "edas.oam.acname"
      values   = ["${var.name}"]
      operator = "NotIn"
    }
  }
}
```

## 参数说明

以下参数受支持：

* `application_name` - (必填) 要创建的应用程序名称。必须以字母开头，支持数字、字母和连字符 (-)，最多支持 36 个字符。
* `cluster_id` - (必填，变更时重建) 要导入的阿里云容器服务 Kubernetes 集群 ID。您可以通过调用 ListCluster 操作查询。
* `package_type` - (可选，变更时重建) 应用包类型。可选值：`FatJar`、`War` 和 `Image`。默认值：`Image`。
* `replicas` - (可选) 应用实例的数量。默认值：1。
* `image_url` - (可选) 镜像地址。当 `package_type` 设置为 `Image` 时，此参数为必填。此属性同时为可回读属性（Computed），将由 API 返回。
* `application_description` - (可选) 应用程序的描述。
* `application_descriotion` - (可选，已废弃) `application_description` 的废弃拼写错误。请改用 `application_description`。
* `package_url` - (可选) 部署包的 URL。通过 FatJar 或 WAR 包部署的应用需要配置它。此属性同时为可回读属性。
* `package_version` - (可选) 部署包的版本号。WAR 和 FatJar 类型需要此参数。请自定义其含义。此属性同时为可回读属性。
* `jdk` - (可选) 部署包依赖的 JDK 版本。可选值为 `Open JDK 7` 和 `Open JDK 8`。Image 类型不支持此参数。
* `web_container` - (可选) 部署包依赖的 Tomcat 版本。适用于通过 WAR 包部署的 Spring Cloud 和 Dubbo 应用。Image 类型不支持此参数。
* `edas_container_version` - (可选) 部署包依赖的 EDAS-Container 版本。Image 类型不支持此参数。
* `cr_ee_repo_id` - (可选) 企业版容器镜像仓库的 Repository ID。
* `cr_instance_id` - (可选) 企业版容器镜像仓库实例的 ID。使用企业版容器镜像仓库时必填。

* `internet_target_port` - (可选，已废弃) 公网 SLB 后端端口，也是应用的服务端口，范围为 1 到 65535。已废弃，请使用 `internet_service_port_infos` 相关属性。
* `internet_slb_port` - (可选，已废弃) 公网 SLB 前端端口，范围为 1~65535。已废弃，请使用 `internet_service_port_infos` 相关属性。
* `internet_slb_protocol` - (可选，已废弃) 公网 SLB 协议，支持 TCP、HTTP 和 HTTPS 协议。已废弃，请使用 `internet_service_port_infos` 相关属性。
* `internet_slb_id` - (可选) 公网 SLB ID。如果不配置，EDAS 将为用户自动购买一个新的 SLB。此属性同时为可回读属性。
* `internet_external_traffic_policy` - (可选) 公网 SLB 外部流量策略。可选值：`Local`、`Cluster`。默认值：`Local`。
* `internet_scheduler` - (可选) 公网 SLB 调度算法。可选值：`rr`、`wrr`。默认值：`rr`。
* `internet_service_port_infos` - (可选) 公网 SLB 服务端口配置。此属性同时为可回读属性。与 `internet_target_port`、`internet_slb_port`、`internet_slb_protocol` 互斥。
  * `port` - (必填) 公网 SLB 前端端口。
  * `protocol` - (必填) 公网 SLB 协议。可选值：`TCP`、`HTTP`、`HTTPS`。
  * `target_port` - (必填) 公网 SLB 后端（目标）端口。

* `intranet_target_port` - (可选，已废弃) 内网 SLB 后端端口，范围为 1 到 65535。已废弃，请使用 `intranet_service_port_infos` 相关属性。
* `intranet_slb_port` - (可选，已废弃) 内网 SLB 前端端口，范围为 1~65535。已废弃，请使用 `intranet_service_port_infos` 相关属性。
* `intranet_slb_protocol` - (可选，已废弃) 内网 SLB 协议，支持 TCP、HTTP 和 HTTPS 协议。已废弃，请使用 `intranet_service_port_infos` 相关属性。
* `intranet_slb_id` - (可选) 内网 SLB ID。如果不配置，EDAS 将为用户自动购买一个新的 SLB。此属性同时为可回读属性。
* `intranet_external_traffic_policy` - (可选) 内网 SLB 外部流量策略。可选值：`Local`、`Cluster`。默认值：`Local`。
* `intranet_scheduler` - (可选) 内网 SLB 调度算法。可选值：`rr`、`wrr`。默认值：`rr`。
* `intranet_service_port_infos` - (可选) 内网 SLB 服务端口配置。此属性同时为可回读属性。与 `intranet_target_port`、`intranet_slb_port`、`intranet_slb_protocol` 互斥。
  * `port` - (必填) 内网 SLB 前端端口。
  * `protocol` - (必填) 内网 SLB 协议。可选值：`TCP`、`HTTP`、`HTTPS`。
  * `target_port` - (必填) 内网 SLB 后端（目标）端口。

* `limit_mem` - (可选) 应用运行期间实例的内存限制，单位：M。此属性同时为可回读属性。
* `requests_mem` - (可选) 创建应用时实例的内存限制，单位：M。设置为 0 表示不限制。此属性同时为可回读属性。
* `requests_m_cpu` - (可选) 创建应用时实例的 CPU 配额，单位：毫核，类似于 request_cpu。此属性同时为可回读属性。
* `limit_m_cpu` - (可选) 应用运行期间实例的 CPU 配额。单位：毫核，设置为 0 表示不限制，类似于 request_cpu。此属性同时为可回读属性。
* `command` - (可选) 设置的命令，如果设置，将在镜像启动时替换镜像中的启动命令。
* `command_args` - (可选) 与命令配合使用，命令的参数为字符串列表，例如：`["-c", "sleep 1001"]`。
* `envs` - (可选) 部署环境变量，为键值对 Map，例如：`{x = "y", x2 = "y2"}`。如需取消配置，请设置为空 Map。
* `pre_stop` - (可选) 停止前执行脚本。此属性同时为可回读属性。
* `post_start` - (可选) 启动后执行脚本。
* `liveness` - (可选) 容器存活状态监控（liveness probe），格式如下：`{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1,"tcpSocket":{"host":"","port":8080} }`。
* `readiness` - (可选) 容器服务状态检查（readiness probe）。如果检查失败，通过 K8s Service 的流量将不会转移到容器。格式如下：`{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1, "httpGet": {"path": "/consumer","port": 8080,"scheme": "HTTP","httpHeaders": [{"name": "test","value": "testvalue"} ]}}`。
* `nas_id` - (可选) 挂载的 NAS 必须与集群在同一地域。必须有可用的挂载点创建配额，或者其挂载点必须在 VPC 的交换机上。如果不填写且 mountDescs 字段存在，默认会自动购买 NAS 并挂载到 VPC 的交换机上。
* `mount_descs` - (可选) 挂载配置描述，作为序列化的 JSON。例如：`[{"nasPath": "/k8s","mountPath": "/mnt"},{"nasPath": "/files","mountPath": "/app/files"}]`。其中，nasPath 表示文件存储路径；mountPath 表示容器内的挂载路径。此属性同时为可回读属性。
* `namespace` - (可选) K8s 集群的命名空间，它将决定您的应用程序部署在哪一个 K8s 命名空间中。默认为 'default'。此属性同时为可回读属性。
* `logical_region_id` - (可选) EDAS 命名空间对应的 ID，非默认命名空间必须填写。
* `config_mount_descs` - (可选) 配置 K8s ConfigMap 和 Secret 挂载，支持将 ConfigMaps 和 Secrets 挂载到指定的容器目录。ConfigMountDescs 的配置参数如下：
  * `name` - (必填) ConfigMap 或 Secret 的名称。
  * `type` - (必填) 配置类型。可选值：`ConfigMap`、`Secret`。
  * `mount_path` - (必填) 挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
* `pvc_mount_descs` - (可选) 配置 K8s PVC (PersistentVolumeClaim) 挂载，支持将 K8s PVC 卷挂载到指定的容器目录。PvcMountDescs 的配置参数如下：
  * `pvc_name` - (必填) PVC 卷的名称。PVC 卷必须已存在且处于 Bound 状态。
  * `mount_paths` - (必填) 挂载目录列表，支持配置多个挂载目录。每个挂载目录支持两个配置参数：
    * `mount_path` - (必填) 挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
    * `read_only` - (可选) 挂载模式，true 表示只读，false 表示读写，默认为 false。
* `local_volume` - (可选) 主机文件挂载到容器目录的配置。
  * `node_path` - (必填) 主机上的路径。
  * `mount_path` - (必填) 容器内的路径。
  * `type` - (必填) 挂载类型。
* `update_type` - (可选) 部署类型。在使用批量部署或灰度部署时可以设置此参数。可选值：`BatchUpdate` 和 `GrayBatchUpdate`。
* `update_batch` - (可选) 部署批次数量。在使用批量部署时，需要设置部署的批次数量。
* `update_release_type` - (可选) 批量部署的发布类型。可选值：`auto` 和 `manual`。
* `update_batch_wait_time` - (可选) 批量部署的自动发布时间。当 `update_release_type` 设置为 `auto` 时，需要设置自动发布时间。
* `update_gray` - (可选) 灰度部署的批次数量。
* `host_aliases` - (可选) HostAliases 配置。此属性同时为可回读属性。
  * `ip` - (可选) IP 地址。
  * `hostnames` - (可选) hostname 列表。
* `custom_tolerations` - (可选) 污点容忍。
  * `key` - (必填) 节点污点的键。
  * `operator` - (必填) 操作符。可选值：`Equal`、`Exists`。
  * `value` - (可选) 节点污点的值。当 `operator` 为 `Equal` 时必填。此属性同时为可回读属性。
  * `effect` - (必填) 效果。可选值：`NoSchedule`、`NoExecute`、`PreferNoSchedule`。
  * `toleration_seconds` - (可选) 容忍时间（秒）。此属性同时为可回读属性。
* `custom_node_affinity_require` - (可选) 必须满足的节点亲和性（硬性）。
  * `match_expressions` - (必填) 节点亲和性规则列表。
    * `key` - (必填) 节点标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`、`Gt`、`Lt`。
    * `values` - (可选) 节点标签的值。此属性同时为可回读属性。
* `custom_node_affinity_preferred` - (可选) 尽量满足的节点亲和性（软性）。
  * `weight` - (可选) 权重，取值范围：1~100。默认值：1。
  * `match_expressions` - (必填) 节点亲和性规则列表。
    * `key` - (必填) 节点标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`、`Gt`、`Lt`。
    * `values` - (可选) 节点标签的值。此属性同时为可回读属性。
* `custom_pod_affinity_require` - (可选) 必须满足的 Pod 亲和性（硬性）。
  * `k8s_namespace` - (可选) K8s 集群的命名空间列表。
  * `topology_key` - (必填) 拓扑域。
  * `match_expressions` - (可选) Pod 亲和性规则列表。
    * `key` - (必填) Pod 标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`。
    * `values` - (可选) Pod 标签的值。此属性同时为可回读属性。
* `custom_pod_affinity_preferred` - (可选) 尽量满足的 Pod 亲和性（软性）。
  * `weight` - (可选) 权重，取值范围：1~100。默认值：1。
  * `k8s_namespace` - (可选) K8s 集群的命名空间列表。
  * `topology_key` - (必填) 拓扑域。
  * `match_expressions` - (可选) Pod 亲和性规则列表。
    * `key` - (必填) Pod 标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`。
    * `values` - (可选) Pod 标签的值。此属性同时为可回读属性。
* `custom_pod_ant_affinity_require` - (可选) 必须满足的 Pod 反亲和性（硬性）。
  * `k8s_namespace` - (可选) K8s 集群的命名空间列表。
  * `topology_key` - (必填) 拓扑域。
  * `match_expressions` - (可选) Pod 反亲和性规则列表。
    * `key` - (必填) Pod 标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`。
    * `values` - (可选) Pod 标签的值。此属性同时为可回读属性。
* `custom_pod_ant_affinity_preferred` - (可选) 尽量满足的 Pod 反亲和性（软性）。
  * `weight` - (可选) 权重，取值范围：1~100。默认值：1。
  * `k8s_namespace` - (可选) K8s 集群的命名空间列表。
  * `topology_key` - (必填) 拓扑域。
  * `match_expressions` - (可选) Pod 反亲和性规则列表。
    * `key` - (必填) Pod 标签的键。
    * `operator` - (必填) 操作符。可选值：`In`、`NotIn`、`Exists`、`DoesNotExist`。
    * `values` - (可选) Pod 标签的值。此属性同时为可回读属性。
## 属性说明

以下属性被导出：

* `application_name` - 要创建的应用程序名称。必须以字母开头，支持数字、字母和连字符 (-)，最多支持36个字符。
* `cluster_id` - 要导入的阿里云容器服务 Kubernetes 集群 ID。您可以通过调用 ListCluster 操作查询。
* `replicas` - 应用实例的数量。
* `package_type` - 应用包类型。可选参数值包括：FatJar、WAR 和 Image。
* `image_url` - 镜像地址。当 `package_type` 设置为 'Image' 时，此参数项可用。
* `update_type` - (可选) 部署类型。在使用批量部署或灰度部署时可以设置此参数。可选值：`BatchUpdate` 和 `GrayBatchUpdate`。
* `update_batch` - (可选) 部署批次数量。在使用批量部署时，需要设置部署的批次数量。
* `update_release_type` - (可选) 批量部署的发布类型。可选值：`auto` 和 `manual`。
* `update_batch_wait_time` - (可选) 批量部署的自动发布时间。当 `update_release_type` 设置为 `auto` 时，需要设置自动发布时间。
* `update_gray` - (可选) 灰度部署的批次数量。
* `config_mount_descs` - 配置 K8s ConfigMap 和 Secret 挂载，支持将 ConfigMaps 和 Secrets 挂载到指定的容器目录。ConfigMountDescs 的配置参数如下：
  * `name` - ConfigMap 或 Secret 的名称。
  * `type` - 配置类型，支持 ConfigMap 和 Secret 类型。
  * `mount_path` - 挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
* `pvc_mount_descs` - 配置 K8s PVC (PersistentVolumeClaim) 挂载，支持将 K8s PVC 卷挂载到指定的容器目录。PvcMountDescs 的配置参数如下：
  * `pvc_name` - PVC 卷的名称。PVC 卷必须已存在且处于 Bound 状态。
  * `mount_paths` - 挂载目录列表，支持配置多个挂载目录。每个挂载目录支持两个配置参数：
    * `mount_path` - 挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
    * `read_only` - 挂载模式，true 表示只读，false 表示读写，默认为 false。
* `local_volume` - 配置主机文件挂载到容器目录。
  * `node_path` - 主机上的路径。
  * `mount_path` - 容器内的路径。
  * `type` - 挂载类型。
* `package_version` - 部署包的版本号。

## 导入

EDAS k8s 应用程序可以通过以下方式导入，例如：

```bash
$ terraform import alibabacloudstack_edas_k8s_application.new_k8s_application application_id
```