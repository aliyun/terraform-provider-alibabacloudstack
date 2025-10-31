---
subcategory: "API GatewayV2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: apsarastack_api_gateway_v2_instance"
sidebar_current: "docs-alibabacloudstack-resource-api-gateway-v2-instance"
description: |-
    Provides a Alibabacloudstack Api Gateway V2 Instance Resource.
---
提供一个API网关资源。

## 示例用法

### 基本用法

```hcl
variable "name" {
  default = "%s"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
	sorted_by = "CPU"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_cs_kubernetes" "default" {
	name						= var.name
	version						= "1.30.7-aliyun.1"
	os_type						= "linux"
	platform					= "AliyunLinux"
	num_of_nodes				= "3"
	master_count				= "3"
	master_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
	master_instance_types		= ["ecs.n4v2.large","ecs.n4v2.large","ecs.n4v2.large"]
	master_disk_category		= "cloud_ssd"
	vpc_id						= "${alibabacloudstack_vpc_vpc.default.id}"
	worker_instance_types		= ["ecs.n4v2.large"]
	worker_vswitch_ids			= ["${alibabacloudstack_vpc_vswitch.default.id}"]
	worker_disk_category		= "cloud_ssd"
	password					= random_password.password.0.result
	pod_cidr					= "172.20.0.0/16"
	service_cidr				= "172.21.0.0/20"
	worker_disk_size			= "40"
	master_disk_size			= "40"
	slb_internet_enabled		= "true"
	security_group_id			= alibabacloudstack_ecs_securitygroup.default.id
	runtime	 {
		name	= "containerd"
		version	= "1.6.28"
	}
}

resource "alibabacloudstack_edas_k8s_cluster" "default" {
	cs_cluster_id	= alibabacloudstack_cs_kubernetes.default.id
}

variable "region" {
  default = "cn-hangzhou"
}

resource "alibabacloudstack_edas_namespace" "default" {
  debug_enable         = false
  description          = var.name
  namespace_logical_id = "${var.region}:${var.name}"
  namespace_name       = var.name
}

resource "apsarastack_api_gateway_v2_instance" "example" {
  instance_name = "example-instance"
  deploy_mode   = "edas"
  broker_engine_type = "SCG"
  instance_class = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
  node_number    = 1
    
  edas_app_infos {
    edas_namespace = "${alibabacloudstack_edas_namespace.default.id}"
    k8s_cluster_id = "${alibabacloudstack_edas_k8s_cluster.default.id}"
  }
}
```

### 自定义部署模式

```hcl
resource "apsarastack_api_gateway_v2_instance" "edas_example" {
  instance_name      = "example-instance"
  deploy_mode        = "custom"
  broker_engine_type = "SCG"
  instance_class     = "${data.alibabacloudstack_api_gateway_v2_instance_types.default.instance_types[0].id}"
  node_number        = 1
}
```

### AI网关模式

```hcl
resource "apsarastack_api_gateway_v2_instance" "custom_example" {
  instance_name         = "custom-instance"
  deploy_mode           = "custom"
  broker_engine_type    = "HIGRESS"
  deploy_cluster_code   = "cluster-code"
  deploy_cluster_namespace = "testnamespace"
  ingress_class_name    = "testingress"
  prometheus_enabled    = true
  sls_enabled           = true
}
```

## 参数参考

支持以下参数：

* `instance_name` - (必需) API 网关实例的名称。
* `deploy_mode` - (可选，强制新建) 实例的部署模式。有效值：`k8s`、`edas`、`custom`。
* `broker_engine_type` - (可选，强制新建) broker 引擎类型。有效值：`HIGRESS`(AI网关)、`SCG`(API网关)。
* `instance_class` - (可选，强制新建) 实例类别/规格。
* `node_number` - (可选，强制新建) 实例的节点数量。
* `edas_namespace_id` - (可选，强制新建) EDAS 命名空间 ID。
* `deploy_cluster_code` - (可选) 部署集群代码。
* `deploy_cluster_namespace` - (可选，强制新建) 部署集群命名空间。
* `ingress_class_name` - (可选，强制新建) ingress 类名。
* `prometheus_enabled` - (可选，强制新建) 是否启用 Prometheus 监控。
* `sls_enabled` - (可选，强制新建) 是否启用 SLS 日志记录。
* `edas_app_infos` - (可选) EDAS 应用信息块集合。每个块包含：
  * `edas_namespace` - (可选) EDAS 命名空间。
  * `k8s_cluster_id` - (可选) Kubernetes 集群 ID。
  * `k8s_namespace` - (可选) Kubernetes 命名空间。默认为 `default`。

## 属性参考

导出以下属性：

* `id` - 实例的 ID。
* `broker_engine_version` - broker 引擎版本。
* `access_mode` - 实例的访问模式。
* `tid` - 租户 ID。
* `create_time` - 实例的创建时间。
* `k8s_cluster_id` - Kubernetes 集群的 ID。
* `edas_app_id` - EDAS 应用程序 ID。
* `status` - 实例的状态。
* `broker_latest_engine_version` - 最新的 broker 引擎版本。
* `shared_instance` - 实例是否为共享实例。
* `custom_deploy_config` - 自定义部署配置（仅适用于自定义部署模式）。

## 导入

可以使用实例 ID 导入 API 网关 V2 实例：

```shell
terraform import apsarastack_api_gateway_v2_instance.example <instance_id>
```