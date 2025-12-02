---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_services"
sidebar_current: "docs-Alibabacloudstack-datasource-api-gateway-v2-services"
description: |-
  Queries the list of services for Alibaba Cloud API Gateway v2.
---

# alibabacloudstack_api_gateway_v2_services

Queries the list of services for Alibaba Cloud API Gateway v2. This data source retrieves all service information under a specified gateway instance, supporting filtering by service ID or name regular expression.

## Example Usage

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

## Argument Reference

The following arguments support filtering the service list:

* `gw_instance_id` (Required): The ID of the gateway instance, used to specify the API Gateway instance to which the services belong.
* `ids` (Optional): A list of service IDs, used to precisely match the resource IDs of services to query (format: `{gwInstanceId:serviceId}`).
* `name_regex` (Optional): A regular expression for service names, used to fuzzy match service names.

## Attributes Reference

The following attributes are exported as data source results:

* `names`: A list of service names, consistent with the order of the returned service list.
* `services`: A list of services. Each service contains the following attributes:
  * `id`: The service resource ID, in the format `{gwInstanceId:serviceId}`.
  * `create_time`: The creation time of the service, in the format `YYYY-MM-DD HH:mm:ss`.
  * `description`: The description of the service.
  * `health_check_struct`: A list of health check configurations (typically contains a single item). Each configuration includes:
    * `health_interval`: The interval time for health checks (in seconds).
    * `health_path`: The path for health checks (e.g., `/check`).
    * `http_failures`: The threshold for the number of failed health checks.
    * `http_statuses`: The expected HTTP status codes (e.g., `200`).
    * `http_successes`: The threshold for the number of successful health checks.
    * `timeout`: The timeout for health checks (in milliseconds).
    * `type`: The type of health check (1 indicates HTTP).
    * `un_health_interval`: The interval time for unhealthy checks (in seconds).
  * `load_balance_type`: The type of load balancing (integer).
  * `name`: The name of the service.
  * `protocol`: The service protocol (e.g., `HTTP`).
  * `real_service_name`: The real service name (may be empty).
  * `service_group`: The group to which the service belongs (may be empty).
  * `service_id`: The unique identifier ID of the service.
  * `service_nodes`: A list of service nodes. Each node contains:
    * `enable`: Whether the node is enabled (boolean).
    * `ip`: The IP address of the node.
    * `port`: The port of the node.
    * `weight`: The weight of the node.
  * `service_version`: The version of the service (may be empty).
  * `source_id`: The source service ID (may be empty).
  * `upstream_type`: The type of upstream service (integer).