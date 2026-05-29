---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-instances"
description: |-
    提供一个数据源来查询 API 网关 V2 实例。
---

# alibabacloudstack_api_gateway_v2_instances

API 网关 V2 实例数据源提供与请求参数匹配的 API 网关 V2 实例列表。

## 示例用法

### 基本用法

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "example" {
  instance_id = "i-example123"
}

output "instance_ids" {
  value = data.alibabacloudstack_api_gateway_v2_instances.example.ids
}
```

### 按部署模式过滤

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "edas_instances" {
  deploy_mode = "edas"
}

output "edas_instance_details" {
  value = data.alibabacloudstack_api_gateway_v2_instances.edas_instances.instances
}
```

### 按 Broker 引擎类型过滤

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "scg_instances" {
  broker_engine_type = "SCG"
}
```

### 使用正则表达式过滤

```hcl
data "alibabacloudstack_api_gateway_v2_instances" "name_filtered" {
  description_regex = "^test-.*"
}
```

## 参数参考

支持以下参数：

* `ids` - (可选) 用于过滤结果的实例 ID 列表。
* `instance_id` - (可选) 要检索的特定实例的 ID。
* `broker_engine_type` - (可选) 实例的 broker 引擎类型。有效值：`HIGRESS`、`SCG`。
* `deploy_mode` - (可选) 实例的部署模式。有效值：`k8s`、`edas`、`custom`。
* `name_regex` - (可选，已弃用) 用于按名称过滤实例的正则表达式。此字段已弃用，将在未来版本中移除。请改用 `description_regex`。
* `description_regex` - (可选) 用于按描述（名称）过滤实例的正则表达式。

-> **注意：** `name_regex` 和 `description_regex` 是互斥的。只能指定其中一个。

## 属性参考

导出以下属性：

* `ids` - 实例 ID 列表。
* `instances` - 实例列表。每个元素包含以下属性：
  * `id` - 实例的 ID。
  * `instance_id` - 实例的 ID。
  * `instance_name` - 实例的名称。
  * `broker_engine_type` - 实例的 broker 引擎类型。
  * `deploy_mode` - 实例的部署模式。
  * `instance_class` - 实例的类别/规格。
  * `status` - 实例的状态。
  * `k8s_cluster_id` - 实例部署的 Kubernetes 集群的 ID。
  * `k8s_namespace` - 实例部署的 Kubernetes 命名空间。
  * `edas_app_id` - 与实例关联的 EDAS 应用程序 ID。
  * `edas_namespace_id` - 与实例关联的 EDAS 命名空间 ID。
  * `shared_instance` - 实例是否为共享实例。
  * `node_number` - 实例中的节点数量。
  * `edas_app_infos` - EDAS 应用程序信息对象集合，每个包含：
    * `edas_namespace` - EDAS 命名空间 ID。
    * `app_id` - 应用程序 ID。
    * `k8s_cluster_id` - Kubernetes 集群 ID。
    * `k8s_namespace` - Kubernetes 命名空间。
  * `custom_deploy_config` - 自定义部署配置映射（仅当 `deploy_mode` 为 `custom` 时可用）。