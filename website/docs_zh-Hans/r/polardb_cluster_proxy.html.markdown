---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxy"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-proxy"
description: |-
  提供 PolarDB 集群代理资源。
---

# alibabacloudstack_polardb_cluster_proxy

提供 PolarDB 集群代理资源。

## 示例用法

```hcl
variable "name" {
  default = "tf-proxy-test"
}

variable "db_type" {
  default = "PostgreSQL"
}

variable "db_version" {
  default = "14"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
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
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.1.0/24"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
    ignore_changes = [
      tags
    ]
  }
}

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type      = "${var.db_type}"
  db_version   = "${var.db_version}"
  sorted_by    = "CPU"
  sub_category = "normal_exclusive"
}

data "alibabacloudstack_polardb_cluster_proxy_types" "types" {
  db_type    = "${var.db_type}"
  db_version = "${var.db_version}"
}

resource "alibabacloudstack_polardb_cluster_instance" "instance" {
  db_cluster_description = "${var.name}"
  db_type                = "${var.db_type}"
  db_version             = "${var.db_version}"
  storage_type           = "ESSDPL1"
  storage_space          = 20
  db_node_class          = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id             = "${alibabacloudstack_vpc_vswitch.default.id}"
  sub_category           = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
}

resource "alibabacloudstack_polardb_cluster_proxy" "default" {
  db_cluster_id          = "${alibabacloudstack_polardb_cluster_instance.instance.id}"
  db_proxy_cluster_class = "${data.alibabacloudstack_polardb_cluster_proxy_types.types.proxy_classes.0.id}"
}
```

## 参数说明

以下参数为支持的配置项：

* `db_cluster_id` - （必填，变更后重建）PolarDB 集群的 ID。更改此参数将强制重新创建资源。
* `db_proxy_cluster_class` - （必填）代理集群节点的规格。

## 属性说明

以下属性为导出项：

* `id` - PolarDB 集群的 ID（与 `db_cluster_id` 相同）。
* `db_cluster_id` - PolarDB 集群的 ID。
* `db_proxy_cluster_id` - 代理集群的 ID。
* `db_proxy_cluster_class` - 代理集群节点的规格。
* `proxy_instances` - 代理实例列表。每个元素包含：
  * `db_node_id` - 代理节点的 ID。
  * `db_node_status` - 代理节点的状态。
  * `db_node_class` - 代理节点的规格。

## 导入

PolarDB 集群代理可以通过集群 ID 导入，例如：

```
$ terraform import alibabacloudstack_polardb_cluster_proxy.example pc-xxxxxxxxxxxx
```
