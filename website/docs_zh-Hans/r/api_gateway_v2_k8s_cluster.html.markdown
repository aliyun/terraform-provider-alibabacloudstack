---
subcategory: "API 网关（API Gateway）V2 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_api_gateway_v2_k8s_cluster"
sidebar_current: "docs-Alibabacloudstack-resource-api-gateway-v2-k8s-cluster"
description: |-
  导入K8s集群到API网关V2版本
---

# alibabacloudstack_api_gateway_v2_k8s_cluster

导入K8s集群到API网关V2版本，支持容器服务集群和自建集群。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tf_testAccApiGatewayv2cluster_2641257"
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


data "alibabacloudstack_cs_kubernetes_clusters_kubeconfig" "k8s_clusters_kubeconfig" {
  cluster_id = local.k8s_cluster_id
}



resource "alibabacloudstack_api_gateway_v2_k8s_cluster" "default" {
  cs_cluster_id    = local.k8s_cluster_id
  k8s_cluster_name = var.name
}
```

## 参数说明

支持以下参数：

* `k8s_cluster_name` - (必填, 变更时重建) K8s集群名称。名称长度为1-128个字符，可包含字母、数字、短划线(-)和下划线(_)。
* `cs_cluster_id` - (可选, 变更时重建) 容器服务集群ID。如果提供了此参数，将自动查询集群名称和配置内容。
* `vpc_id` - (可选, 变更时重建) VPC网络ID。如果提供了此参数，将使用内网SLB类型；否则使用公网SLB类型。
* `config_content` - (可选, 变更时重建) K8s集群的配置内容。当未提供`cs_cluster_id`时，此参数为必填，用于自建集群的导入。

## 属性说明

以下属性会从API网关服务端导出：

* `id` - K8s集群的ID，由API网关服务生成的唯一标识符。
* `cluster_type` - K8s集群类型，可能的值包括"container-service"(容器服务集群)或"self-built"(自建集群)。