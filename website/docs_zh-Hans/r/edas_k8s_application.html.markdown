---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_k8s_application"
sidebar_current: "docs-alibabacloudstack-resource-edas-k8s-application"
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

  // set 'image_url' and 'repo_id' when package_type is 'image'
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
  logical_region_id     = cn-beijing
}
```

## 参数参考

以下参数受支持：

* `application_name` - (必填，变更时重建) 要创建的应用程序名称。必须以字母开头，支持数字、字母和连字符 (-)，最多支持36个字符。
* `cluster_id` - (必填，变更时重建) 要导入的阿里云容器服务 Kubernetes 集群 ID。您可以通过调用 ListCluster 操作查询。
* `package_type` - (必填，变更时重建) 应用包类型。可选参数值包括：FatJar、WAR 和 Image。
* `replicas` - (可选) 应用实例的数量。
* `image_url` - (可选) 镜像地址。当 `package_type` 设置为 'Image' 时，此参数项是必填的。
* `application_description` - (可选) 应用程序的描述。
* `package_url` - (可选) 部署包的 URL。通过 FatJar 或 WAR 包部署的应用需要配置它。
* `package_version` - (可选) 部署包的版本号。WAR 和 FatJar 类型需要此参数。请自定义其含义。
* `jdk` - (可选，变更时重建) 部署包依赖的 JDK 版本。可选参数值为 Open JDK 7 和 Open JDK 8。Image 不支持此参数。
* `web_container` - (可选，变更时重建) 部署包依赖的 Tomcat 版本。适用于通过 WAR 包部署的 Spring Cloud 和 Dubbo 应用。Image 不支持此参数。
* `edas_container_version` - (可选) 部署包依赖的 EDAS-Container 版本。Image 不支持此参数。

* `internet_target_port` - (可选，变更时重建) 公网 SLB 后端端口，也是应用的服务端口，范围为 1 到 65535。(已废弃， internet_service_port_infos 相关属性)
* `internet_slb_port` - (可选，变更时重建) 公网 SLB 前端端口，范围为 1~65535。(已废弃， internet_service_port_infos 相关属性)
* `internet_slb_protocol` - (可选，变更时重建) 公网 SLB 协议支持 TCP、HTTP 和 HTTPS 协议。(已废弃， internet_service_port_infos 相关属性)
* `internet_slb_id` - (可选，变更时重建) 公网 SLB ID。如果不配置，EDAS 将为用户自动购买一个新的 SLB。

* `intranet_target_port` - (可选，变更时重建) 内网 SLB 后端端口，也是应用的服务端口，范围为 1 到 65535。(已废弃， intranet_service_port_infos 相关属性)
* `intranet_slb_port` - (可选，变更时重建) 内网 SLB 前端端口，范围为 1~65535。(已废弃， intranet_service_port_infos 相关属性)
* `intranet_slb_protocol` - (可选，变更时重建) 内网 SLB 协议支持 TCP、HTTP 和 HTTPS 协议。(已废弃， intranet_service_port_infos 相关属性)
* `intranet_slb_id` - (可选，变更时重建) 内网 SLB ID。如果不配置，EDAS 将为用户自动购买一个新的 SLB。

* `limit_mem` - (可选) 应用运行期间实例的内存限制，单位：M。
* `requests_mem` - (可选) 创建应用时实例的内存限制，单位：M。设置为 0 表示不限制。
* `requests_m_cpu` - (可选) 创建应用时实例的 CPU 配额，单位：毫核，类似于 request_cpu。
* `limit_m_cpu` - (可选) 应用运行期间实例的 CPU 配额。单位：毫核，设置为 0 表示不限制，类似于 request_cpu。
* `command` - (可选) 设置的命令，如果设置，将在镜像启动时替换镜像中的启动命令。
* `command_args` - (可选) 与命令配合使用，命令的参数是一个 JsonArray 格式的字符串，格式为：`[{"argument":"-c"},{"argument":"test"}]`。其中，-c 和 test 是需要设置的两个参数。
* `envs` - (可选，变更时重建) 部署环境变量，格式必须符合 JSON 对象数组，例如：`{"name":"x","value":"y"},{"name":"x2","value":"y2"}`。如果要取消配置，需要设置一个空的 JSON 数组 "" 来表示无配置。
* `pre_stop` - (可选) 停止前执行脚本。
* `post_start` - (可选) 启动后执行脚本。
* `liveness` - (可选) 容器存活状态监控，格式如下：`{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1,"tcpSocket":{"host":"","port":8080} }`。
* `readiness` - (可选) 容器服务状态检查。如果检查失败，通过 K8s Service 的流量将不会转移到容器。格式如下：`{"failureThreshold": 3,"initialDelaySeconds": 5,"successThreshold": 1,"timeoutSeconds": 1, "httpGet": {"path": "/consumer","port": 8080,"scheme": "HTTP","httpHeaders": [{"name": "test","value": "testvalue"} ]}}`。
* `nas_id` - (可选) 挂载的 NAS 必须与集群在同一地域。必须有可用的挂载点创建配额，或者其挂载点必须在 VPC 的交换机上。如果不填写且 mountDescs 字段存在，默认会自动购买 NAS 并挂载到 VPC 的交换机上。
* `mount_descs` - (可选，变更时重建) 挂载配置描述，作为序列化的 JSON。例如：`[{"nasPath": "/k8s","mountPath": "/mnt"},{"nasPath": "/files","mountPath": "/app/files"}]`。其中，nasPath 表示文件存储路径；mountPath 表示容器内的挂载路径。
* `namespace` - (可选) K8s 集群的命名空间，它将决定您的应用程序部署在哪一个 K8s 命名空间中。默认为 'default'。
* `logical_region_id` - (可选) EDAS 命名空间对应的 ID，非默认命名空间必须填写。
* `config_mount_descs` - (可选) 配置 K8s ConfigMap 和 Secret 挂载，支持将 ConfigMaps 和 Secrets 挂载到指定的容器目录。ConfigMountDescs 的配置参数如下：
  * `name` - (必填)  ConfigMap 或 Secret 的名称。
  * `type` - (必填)  配置类型，支持 ConfigMap 和 Secret 类型。
  * `mount_path` - (必填)  挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
* `pvc_mount_descs` - (可选) 配置 K8s PVC (PersistentVolumeClaim) 挂载，支持将 K8s PVC 卷挂载到指定的容器目录。PvcMountDescs 的配置参数如下：
  * `pvc_name` - (必填)  PVC 卷的名称。PVC 卷必须已存在且处于 Bound 状态。
  * `mount_paths` - (必填)  挂载目录列表，支持配置多个挂载目录。每个挂载目录支持两个配置参数：
    * `mount_path` - (必填)  挂载路径，容器内的绝对路径，以斜杠 (/) 开头。
    * `read_only` - (必填)  挂载模式，true 表示只读，false 表示读写，默认为 false。
* `local_volume` - (可选) 主机文件挂载到容器目录的配置。
  * `node_path` - (必填)  主机上的路径。
  * `mount_path` - (必填)  容器内的路径。
  * `type` - (可选) 挂载类型。
* `update_type` - (可选) 部署类型。在使用批量部署或灰度部署时可以设置此参数。可选值：`BatchUpdate` 和 `GrayBatchUpdate`。
* `update_batch` - (可选) 部署批次数量。在使用批量部署时，需要设置部署的批次数量。
* `update_release_type` - (可选) 批量部署的发布类型。可选值：`auto` 和 `manual`。
* `update_batch_wait_time` - (可选) 批量部署的自动发布时间。当 `update_release_type` 设置为 `auto` 时，需要设置自动发布时间。
* `update_gray` - (可选) 灰度部署的批次数量。
* `host_aliases` - (可选) host映射配置
  * `ip` - (可选) Ip地址.
  * `hostnames` - (可选) hostname数组.
* `intranet_service_port_infos` - (可选) 内网SLB服务端口配置。
  * `port` - (可选) 内网SLB服务端口的端口号
  * `protocol` - (可选) 内网SLB服务端口的协议。
  * `target_port` - (可选) 内网SLB服务端口的目标端口号。
* `intranet_external_traffic_policy` - (可选) 内网SLB外部流量策略。
* `intranet_scheduler` - (可选) 内网SLB调度规则。
* `internet_service_port_infos` - (可选) 公网SLB服务端口配置。
  * `port` - (可选) 公网SLB服务端口。
  * `protocol` - (可选) 公网SLB服务端口的协议。
  * `target_port` - (可选) 公网SLB服务的目标端口号。
* `internet_external_traffic_policy` - (可选) 公网SLB外部流量策略。
* `internet_scheduler` - (可选) 服务的公网 SLB 调度规则。

## 属性参考

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