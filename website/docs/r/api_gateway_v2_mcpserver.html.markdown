---
subcategory: "API Gateway V2"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_mcpserver"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-mcpserver"
description: |-
  Manage MCP servers for API Gateway V2
---

# alibabacloudstack_api_gateway_v2_mcpserver

Manage MCP servers for API Gateway V2, supporting three types of service configurations: OPEN_API, DATABASE, and DIRECT_ROUTE.

## Example Usage

### Create an OPEN_API type MCP server

```hcl

variable "name" {
  default = "testtf"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
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


resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
  cs_cluster_id    = local.k8s_cluster_id
  k8s_cluster_name = var.name
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name            = "${var.name}-apigw"
  node_number              = 1
  instance_class           = "mini"
  broker_engine_type       = "HIGRESS"
  deploy_mode              = "k8s"
  deploy_cluster_code      = alibabacloudstack_api_gateway_v2_k8s_cluster.default.id
  deploy_cluster_namespace = "${var.name}-namespace"
  ingress_class_name       = "${var.name}-class"
  sls_enabled              = "true"
  prometheus_enabled       = "true"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain      = "${var.name}.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
}


resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  consumer_auth  = "true"
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  name           = var.name
  description    = "testddd"
  type           = "OPEN_API"
  service        = "kubernetes.default.svc.cluster.local"
  domains = [
    "${alibabacloudstack_api_gateway_v2_domain.default.id}"
  ]
}
```

### Create a DATABASE type MCP server

```hcl

variable "name" {
  default = "testtf"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
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


resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
  cs_cluster_id    = local.k8s_cluster_id
  k8s_cluster_name = var.name
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name            = "${var.name}-apigw"
  node_number              = 1
  instance_class           = "mini"
  broker_engine_type       = "HIGRESS"
  deploy_mode              = "k8s"
  deploy_cluster_code      = alibabacloudstack_api_gateway_v2_k8s_cluster.default.id
  deploy_cluster_namespace = "${var.name}-namespace"
  ingress_class_name       = "${var.name}-class"
  sls_enabled              = "true"
  prometheus_enabled       = "true"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain      = "${var.name}.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
}


resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  consumer_auth  = "true"
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  name           = var.name
  description    = "testddd"
  type           = "OPEN_API"
  service        = "kubernetes.default.svc.cluster.local"
  domains = [
    "${alibabacloudstack_api_gateway_v2_domain.default.id}"
  ]
}
```

### Create a DIRECT_ROUTE type MCP server

```hcl

variable "name" {
  default = "testtf"
}

data "alibabacloudstack_api_gateway_v2_instance_types" "default" {
  sorted_by = "CPU"
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


resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
  cs_cluster_id    = local.k8s_cluster_id
  k8s_cluster_name = var.name
}

resource "alibabacloudstack_api_gateway_v2_instance" "default" {
  instance_name            = "${var.name}-apigw"
  node_number              = 1
  instance_class           = "mini"
  broker_engine_type       = "HIGRESS"
  deploy_mode              = "k8s"
  deploy_cluster_code      = alibabacloudstack_api_gateway_v2_k8s_cluster.default.id
  deploy_cluster_namespace = "${var.name}-namespace"
  ingress_class_name       = "${var.name}-class"
  sls_enabled              = "true"
  prometheus_enabled       = "true"
}

resource "alibabacloudstack_api_gateway_v2_domain" "default" {
  domain      = "${var.name}.com"
  instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  protocol    = "HTTP"
}


resource "alibabacloudstack_api_gateway_v2_mcpserver" "default" {
  consumer_auth  = "true"
  gw_instance_id = alibabacloudstack_api_gateway_v2_instance.default.id
  name           = var.name
  description    = "testddd"
  type           = "OPEN_API"
  service        = "kubernetes.default.svc.cluster.local"
  domains = [
    "${alibabacloudstack_api_gateway_v2_domain.default.id}"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the MCP server.
* `gw_instance_id` - (Required) The ID of the API Gateway instance.
* `type` - (Required) The type of the MCP server. Valid values: `OPEN_API`, `DATABASE`, and `DIRECT_ROUTE`.
* `service` - (Required) The service name, formatted as a Kubernetes service address.
* `description` - (Optional) The description of the MCP server.
* `domains` - (Optional) A list of bound domain names.
* `consumer_auth` - (Optional) Whether to enable consumer authentication.
* `direct_route_path` - (Optional) The direct route path, required when `type` is `DIRECT_ROUTE`.
* `direct_route_type` - (Optional) The direct route type, required when `type` is `DIRECT_ROUTE`. Valid values: `streamable` and `sse`.
* `db_host` - (Optional) The host address of the database, required when `type` is `DATABASE`.
* `db_port` - (Optional) The port of the database, required when `type` is `DATABASE`.
* `db_username` - (Optional) The username of the database, required when `type` is `DATABASE`.
* `db_password` - (Optional) The password of the database, required when `type` is `DATABASE`.
* `db_name` - (Optional) The name of the database, required when `type` is `DATABASE`.
* `other_params` - (Optional) Other parameters of the database, required when `type` is `DATABASE`.
* `db_type` - (Optional) The type of the database, required when `type` is `DATABASE`. Valid values: `MYSQL`, `POSTGRESQL`, and `CLICKHOUSE`.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of the MCP server, formatted as `{gwInstanceId}:{name}`.
* `raw_configurations` - The raw configuration information of the MCP server.
* `services` - A list of services, containing the following attributes:
  * `name` - The name of the service.
  * `port` - The port of the service.
  * `version` - The version of the service.
  * `weight` - The weight of the service.