---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_mcpserver"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-mcpserver"
description: |-
  管理API网关V2版本的MCP服务器
---

# alibabacloudstack_api_gateway_v2_mcpserver

管理API网关V2版本的MCP服务器，支持OPEN_API、DATABASE和DIRECT_ROUTE三种类型的服务配置。

## 示例用法

### 创建OPEN_API类型的MCP服务器

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

### 创建DATABASE类型的MCP服务器

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

### 创建DIRECT_ROUTE类型的MCP服务器

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

## 参数说明

支持以下参数：

* `name` - (必填) MCP服务器的名称。
* `gw_instance_id` - (必填) API网关实例ID。
* `type` - (必填) MCP服务器类型，可选值：`OPEN_API`、`DATABASE`、`DIRECT_ROUTE`。
* `service` - (必填) 服务名称，格式为Kubernetes服务地址。
* `description` - (可选) MCP服务器的描述信息。
* `domains` - (可选) 绑定的域名列表。
* `consumer_auth` - (可选) 是否启用消费者认证。
* `direct_route_path` - (可选) 直连路由路径，当`type`为`DIRECT_ROUTE`时必填。
* `direct_route_type` - (可选) 直连路由类型，当`type`为`DIRECT_ROUTE`时必填，可选值：`streamable`、`sse`。
* `db_host` - (可选) 数据库主机地址，当`type`为`DATABASE`时必填。
* `db_port` - (可选) 数据库端口，当`type`为`DATABASE`时必填。
* `db_username` - (可选) 数据库用户名，当`type`为`DATABASE`时必填。
* `db_password` - (可选) 数据库密码，当`type`为`DATABASE`时必填。
* `db_name` - (可选) 数据库名称，当`type`为`DATABASE`时必填。
* `other_params` - (可选) 数据库其他参数，当`type`为`DATABASE`时必填。
* `db_type` - (可选) 数据库类型，当`type`为`DATABASE`时必填，可选值：`MYSQL`、`POSTGRESQL`、`CLICKHOUSE`。

## 属性说明

以下属性会从API中导出：

* `id` - MCP服务器的资源ID，格式为`{gwInstanceId}:{name}`。
* `raw_configurations` - MCP服务器的原始配置信息。
* `services` - 服务列表，包含以下属性：
  * `name` - 服务名称。
  * `port` - 服务端口。
  * `version` - 服务版本。
  * `weight` - 服务权重。