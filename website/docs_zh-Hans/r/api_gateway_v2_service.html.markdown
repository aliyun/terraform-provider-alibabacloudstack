---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service"
sidebar_current: "docs-Alibabacloudstack-api_gateway-api_gateway_v2_service"
description: |-
  管理API网关V2版本的服务资源
---

# alibabacloudstack_api_gateway_v2_service

管理API网关V2版本的服务资源，用于创建、读取、更新和删除API网关服务。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "testtf-apigw-1656"
}
data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
}

resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
  cs_cluster_id    = local.k8s_cluster_id
  k8s_cluster_name = var.name
}


variable "existed_k8s_cluster_id" {
  default = ""
}


data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}


resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
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
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_zones.default.zones.0.id
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}


resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = alibabacloudstack_vpc_vpc.default.id
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_ecs_securitygroup.default.id
  cidr_ip           = "192.168.0.0/16"
}




resource "random_password" "password" {
  count            = 1
  length           = 12
  special          = true
  override_special = "!@#$^&*()_"
  min_lower        = 1
  min_upper        = 1
  min_numeric      = 1
}

data "alibabacloudstack_cs_kubernetes_clusters" "default" {
  ids = var.existed_k8s_cluster_id == "" ? [] : [var.existed_k8s_cluster_id]
}

locals {
  create_count = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? 0 : 1
}

resource "alibabacloudstack_cs_kubernetes" "default" {
  count                 = local.create_count
  name                  = var.name
  version               = "1.30.7-aliyun.1"
  os_type               = "linux"
  platform              = "AliyunLinux"
  num_of_nodes          = "3"
  master_count          = "3"
  master_vswitch_ids    = ["${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}", "${alibabacloudstack_vpc_vswitch.default.id}"]
  master_instance_types = ["ecs.n4v2.large", "ecs.n4v2.large", "ecs.n4v2.large"]
  master_disk_category  = "cloud_ssd"
  vpc_id                = alibabacloudstack_vpc_vpc.default.id
  worker_instance_types = ["ecs.n4v2.large"]
  worker_vswitch_ids    = ["${alibabacloudstack_vpc_vswitch.default.id}"]
  worker_disk_category  = "cloud_ssd"
  password              = random_password.password.0.result
  pod_cidr              = "172.20.0.0/16"
  service_cidr          = "172.21.0.0/20"
  worker_disk_size      = "40"
  master_disk_size      = "40"
  slb_internet_enabled  = "true"
  security_group_id     = alibabacloudstack_ecs_securitygroup.default.id
  runtime {
    name    = "containerd"
    version = "1.6.28"
  }
}

locals {
  k8s_cluster_id   = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.ids.0 : alibabacloudstack_cs_kubernetes.default.0.id
  k8s_cluster_name = length(data.alibabacloudstack_cs_kubernetes_clusters.default.ids) > 0 ? data.alibabacloudstack_cs_kubernetes_clusters.default.names.0 : alibabacloudstack_cs_kubernetes.default.0.name
}


resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name            = var.name
  node_number              = "1"
  instance_class           = "mini"
  broker_engine_type       = "SCG"
  deploy_mode              = "k8s"
  deploy_cluster_code      = alibabacloudstack_api_gateway_v2_k8s_cluster.default.id
  deploy_cluster_namespace = "${var.name}-namespace"
  ingress_class_name       = "${var.name}-class"
  sls_enabled              = true
  prometheus_enabled       = true
}




resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name              = var.name
  description       = var.name
  upstream_type     = "1"
  load_balance_type = "1"
  protocol          = "HTTP"
  service_nodes {
    port   = "80"
    weight = "100"
    enable = "true"
    ip     = "127.0.0.1"
  }

  health_check_struct {
    http_failures      = "0"
    type               = "1"
    health_path        = "/check"
    http_statuses      = "200"
    timeout            = "20000"
    health_interval    = "30"
    un_health_interval = "30"
    http_successes     = "1"
  }

  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
}
```
### Ai网关服务用法
```hcl
resource "alibabacloudstack_api_gateway_v2_service" "default" {
  name = "${var.name}"
  service_source_type = "ip"
  protocol = "HTTP"
  service_nodes {
	ip = "192.168.1.1"
	port = 80
  }
  gw_instance_id = "${alibabacloudstack_api_gateway_v2_instance.default.id}"
}
```


## 参数说明

支持以下参数：

* `gw_instance_id` - (必填, 变更时重建) API网关实例ID。
* `name` - (必填, 变更时重建) 服务名称。
* `description` - (可选) 服务描述信息。
* `upstream_type` - (可选, 可回读) 上游类型。
* `load_balance_type` - (可选, 可回读) 负载均衡类型。
* `protocol` - (可选, 可回读) 服务协议类型，如HTTP。
* `real_service_name` - (可选) 真实服务名称。
* `service_group` - (可选) 服务分组。
* `service_version` - (可选) 服务版本。
* `source_id` - (可选) 源ID。
* `source_group` - (可选, 可回读) 源分组。
* `service_source_type` - (可选) 服务源类型。可选值：`dns`、`ip`。创建AI网关服务时必填。
* `health_check_struct` - (可选, 最多1项) 健康检查配置结构。
  * `type` - (必填) 健康检查类型。
  * `health_path` - (可选) 健康检查路径。
  * `http_statuses` - (可选) 健康检查的HTTP状态码。
  * `timeout` - (可选) 健康检查超时时间（毫秒）。
  * `health_interval` - (可选) 健康检查间隔时间（秒）。
  * `un_health_interval` - (可选) 不健康检查间隔时间（秒）。
  * `http_successes` - (可选) 判定为健康的HTTP成功次数。
  * `http_failures` - (可选) 判定为不健康的HTTP失败次数。
* `service_nodes` - (可选) 服务节点列表。
  * `ip` - (必填) 节点IP地址。
  * `port` - (必填) 节点端口。
  * `weight` - (可选, 可回读) 节点权重。
  * `enable` - (可选, 可回读) 节点是否启用。
* `sql_input_parameters` - (可选) SQL输入参数配置。
  * `original_name` - (必填) 原始参数名称。
  * `target_name` - (必填) 目标参数名称。
  * `isoptional` - (可选) 是否可选参数。
  * `description` - (必填) 参数描述。
  * `sample` - (必填) 参数示例值。
* `sql_output_parameters` - (可选) SQL输出参数配置。
  * `original_name` - (必填) 原始参数名称。
  * `target_name` - (必填) 目标参数名称。
  * `param_type` - (可选) 参数类型。默认值：`java.lang.String`。
  * `isoptional` - (可选) 是否可选参数。
  * `description` - (必填) 参数描述。
  * `sample` - (必填) 参数示例值。

## 属性说明

以下属性会从API网关服务资源导出：

* `id` - 资源ID，格式为 `<gw_instance_id>^<service_id>`。
* `service_id` - 服务ID。

## 导入

API网关V2服务可以通过组合ID（gw_instance_id^service_id）导入，例如：

```
$ terraform import alibabacloudstack_api_gateway_v2_service.example gw-inst-123456^svc-789012
```