---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_services"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-services"
description: |-
  查询阿里云API网关v2版本的服务列表
---

# alibabacloudstack_api_gateway_v2_services

查询阿里云API网关v2版本的服务列表。该数据源用于检索指定网关实例下的所有服务信息，支持通过服务ID或名称正则表达式进行过滤。

## 示例用法

```hcl



variable "name" {
  default = "testtf-apigw-15604"
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
  protocol          = "HTTP"
  upstream_type     = "1"
  load_balance_type = "1"
  gw_instance_id    = alibabacloudstack_api_gateway_v2_instance.default.id
  service_nodes {
    ip     = "127.0.0.1"
    port   = "80"
    weight = "100"
    enable = "true"
  }
  health_check_struct {
    type               = "1"
    health_path        = "/check"
    http_statuses      = "200"
    timeout            = "20000"
    health_interval    = "30"
    un_health_interval = "30"
    http_successes     = "1"
    http_failures      = "0"
  }
}


data "alibabacloudstack_api_gateway_v2_services" "default" {
  gw_instance_id = alibabacloudstack_api_gateway_v2_service.default.gw_instance_id
  name_regex     = "testtf-apigw-*"
}

```

## 参数说明

以下参数支持过滤服务列表：

* `gw_instance_id` (必填)：网关实例ID，用于指定要查询的服务所属的API网关实例。
* `ids` (可选)：服务ID列表，用于精确匹配需要查询的服务资源ID（格式为`{gwInstanceId:serviceId}`）。
* `name_regex` (可选)：服务名称的正则表达式，用于模糊匹配服务名称。

## 属性说明

以下属性导出为数据源结果：

* `names`：服务名称列表，与返回的服务列表顺序一致。
* `services`：服务列表。每个服务包含以下属性：
  * `id`：服务资源ID，格式为`{gwInstanceId:serviceId}`。
  * `create_time`：服务创建时间，格式为`YYYY-MM-DD HH:mm:ss`。
  * `description`：服务描述信息。
  * `health_check_struct`：健康检查配置列表（通常包含单个配置项）。每个配置包含：
    * `health_interval`：健康检查间隔时间（秒）。
    * `health_path`：健康检查路径（如`/check`）。
    * `http_failures`：健康检查失败次数阈值。
    * `http_statuses`：期望的HTTP状态码（如`200`）。
    * `http_successes`：健康检查成功次数阈值。
    * `timeout`：健康检查超时时间（毫秒）。
    * `type`：健康检查类型（1表示HTTP）。
    * `un_health_interval`：不健康检查间隔时间（秒）。
  * `load_balance_type`：负载均衡类型（整数）。
  * `name`：服务名称。
  * `protocol`：服务协议（如`HTTP`）。
  * `real_service_name`：真实服务名称（可能为空）。
  * `service_group`：服务所属组（可能为空）。
  * `service_id`：服务唯一标识ID。
  * `service_nodes`：服务节点列表。每个节点包含：
    * `enable`：节点是否启用（布尔值）。
    * `ip`：节点IP地址。
    * `port`：节点端口。
    * `weight`：节点权重。
  * `service_version`：服务版本（可能为空）。
  * `source_id`：源服务ID（可能为空）。
  * `upstream_type`：上游服务类型（整数）。