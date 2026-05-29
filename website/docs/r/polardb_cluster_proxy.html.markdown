---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_proxy"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-cluster-proxy"
description: |-
  Provides a PolarDB cluster proxy resource.
---

# alibabacloudstack_polardb_cluster_proxy

Provides a PolarDB cluster proxy resource.

## Example Usage

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
  db_type    = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by  = "CPU"
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

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster to which the proxy belongs. Changing this parameter will force a recreation of the resource.
* `db_proxy_cluster_class` - (Required) The specification of the proxy cluster node.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the PolarDB cluster (same as `db_cluster_id`).
* `db_cluster_id` - The ID of the PolarDB cluster.
* `db_proxy_cluster_id` - The ID of the proxy cluster.
* `db_proxy_cluster_class` - The specification of the proxy cluster node.
* `proxy_instances` - A list of proxy instances. Each element contains:
  * `db_node_id` - The ID of the proxy node.
  * `db_node_status` - The status of the proxy node.
  * `db_node_class` - The specification of the proxy node.

## Import

PolarDB cluster proxy can be imported using the cluster ID, e.g.

```
$ terraform import alibabacloudstack_polardb_cluster_proxy.example pc-xxxxxxxxxxxx
```
