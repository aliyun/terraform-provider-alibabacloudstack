---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxies"
sidebar_current: "docs-alibabacloudstack-datasource-polardb-cluster-proxies"
description: |-
  提供由 alibabacloudstack 账户拥有的 PolarDB 集群代理列表。
---

# datasource: alibabacloudstack_polardb_cluster_proxies

该数据源根据指定的筛选条件提供 alibabacloudstack 账户中的 PolarDB 集群代理列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testacc-example"
}

variable "db_type" {
  default = "PostgreSQL"
}

variable "db_version" {
  default = "14"
}

data "alibabacloudstack_polardb_cluster_proxy_types" "types" {
  db_type    = var.db_type
  db_version = var.db_version
}

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type      = var.db_type
  db_version   = var.db_version
  sorted_by    = "CPU"
  sub_category = "normal_exclusive"
}

resource "alibabacloudstack_polardb_cluster_instance" "instance" {
  db_cluster_description = var.name
  db_type                = var.db_type
  db_version             = var.db_version
  storage_type           = "ESSDPL1"
  storage_space          = 20
  db_node_class          = data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id
  zone_id                = data.alibabacloudstack_zones.default.zones.0.id
  vswitch_id             = alibabacloudstack_vpc_vswitch.default.id
  sub_category           = data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category
}

resource "alibabacloudstack_polardb_cluster_proxy" "default" {
  db_cluster_id      = alibabacloudstack_polardb_cluster_instance.instance.id
  db_proxy_cluster_class = data.alibabacloudstack_polardb_cluster_proxy_types.types.proxy_classes.0.id
}

data "alibabacloudstack_polardb_cluster_proxies" "default" {
  db_cluster_id = alibabacloudstack_polardb_cluster_instance.instance.id
}
```

## 参数说明

支持以下参数：

* `db_cluster_id` - （必填）PolarDB 集群的 ID。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - 数据源的 ID。
* `db_proxy_cluster_id` - 集群代理的 ID。
* `db_proxy_cluster_num` - 集群代理节点的数量。
* `proxy_instances` - 代理实例列表。
  * `db_node_status` - 代理节点的状态。
  * `db_node_id` - 代理节点的 ID。
  * `db_node_class` - 代理节点的规格。
