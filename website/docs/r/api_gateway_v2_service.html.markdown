---
subcategory: "API Gateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_service"
description: |-
  Manages API Gateway V2 service resources
---

# alibabacloudstack_api_gateway_v2_service

Manages API Gateway V2 service resources for creating, reading, updating, and deleting API Gateway services.

## Example Usage

### Basic Usage

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
### AiGw service  Usage
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

## Argument Reference

The following arguments are supported:

* `gw_instance_id` - (Required, ForceNew) The ID of the API Gateway instance.
* `name` - (Required, ForceNew) The name of the service.
* `description` - (Optional) The description of the service.
* `upstream_type` - (Optional, Computed) The upstream type.
* `load_balance_type` - (Optional, Computed) The load balancing type.
* `protocol` - (Optional, Computed) The service protocol type, such as HTTP.
* `real_service_name` - (Optional) The real service name.
* `service_group` - (Optional) The service group.
* `service_version` - (Optional) The service version.
* `source_id` - (Optional) The source ID.
* `source_group` - (Optional, Computed) The source group.
* `service_source_type` - (Optional) The service source type. Valid values: `dns`, `ip`. Required when creating AI Gateway service.
* `health_check_struct` - (Optional, Max Items: 1) The health check configuration structure.
  * `type` - (Required) The health check type.
  * `health_path` - (Optional) The health check path.
  * `http_statuses` - (Optional) The HTTP status codes for health check.
  * `timeout` - (Optional) The health check timeout in milliseconds.
  * `health_interval` - (Optional) The health check interval in seconds.
  * `un_health_interval` - (Optional) The unhealthy check interval in seconds.
  * `http_successes` - (Optional) The number of HTTP successes to determine healthy status.
  * `http_failures` - (Optional) The number of HTTP failures to determine unhealthy status.
* `service_nodes` - (Optional) The list of service nodes.
  * `ip` - (Required) The IP address of the node.
  * `port` - (Required) The port of the node.
  * `weight` - (Optional, Computed) The weight of the node.
  * `enable` - (Optional, Computed) Whether the node is enabled.
* `sql_input_parameters` - (Optional) The SQL input parameter configuration.
  * `original_name` - (Required) The original parameter name.
  * `target_name` - (Required) The target parameter name.
  * `isoptional` - (Optional) Whether the parameter is optional.
  * `description` - (Required) The parameter description.
  * `sample` - (Required) The sample value of the parameter.
* `sql_output_parameters` - (Optional) The SQL output parameter configuration.
  * `original_name` - (Required) The original parameter name.
  * `target_name` - (Required) The target parameter name.
  * `param_type` - (Optional) The parameter type. Default: `java.lang.String`.
  * `isoptional` - (Optional) Whether the parameter is optional.
  * `description` - (Required) The parameter description.
  * `sample` - (Required) The sample value of the parameter.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format of `<gw_instance_id>^<service_id>`.
* `service_id` - The ID of the service.

## Import

API Gateway V2 Service can be imported using the composite ID (gw_instance_id^service_id), e.g.

```
$ terraform import alibabacloudstack_api_gateway_v2_service.example gw-inst-123456^svc-789012
```