---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxies"
sidebar_current: "docs-alibabacloudstack-datasource-polardb-cluster-proxies"
description: |-
  Provides a list of PolarDB cluster proxies owned by an alibabacloudstack account.
---

# datasource: alibabacloudstack_polardb_cluster_proxies

This data source provides a list of PolarDB cluster proxies in an alibabacloudstack account according to the specified filters.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required) The ID of the PolarDB cluster.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the data source.
* `db_proxy_cluster_id` - The ID of the cluster proxy.
* `db_proxy_cluster_num` - The number of cluster proxy nodes.
* `proxy_instances` - The list of proxy instances.
  * `db_node_status` - The status of the proxy node.
  * `db_node_id` - The ID of the proxy node.
  * `db_node_class` - The specification of the proxy node.
