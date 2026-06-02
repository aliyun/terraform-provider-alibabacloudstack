---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_database"
description: |-
  Provides a PolarDB cluster database resource.
---

# alibabacloudstack_polardb_cluster_database

Provides a PolarDB cluster database resource. This resource allows you to manage databases in PolarDB clusters.

## Example Usage

```hcl
variable "name" {
  default = "tf-polar-test"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name       = "${var.name}_vpc"
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vpc_id       = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block   = "172.16.0.0/24"
  zone_id      = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_name = "${var.name}_vsw"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
  db_cluster_description = "${var.name}"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_type                = "${var.db_type}"
  db_version             = "${var.db_version}"
  storage_space          = "20"
  vpc_id                 = "${alibabacloudstack_vpc_vpc.default.id}"
  vswitch_id             = "${alibabacloudstack_vpc_vswitch.default.id}"
  db_node_class          = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
  sub_category           = "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
  storage_type           = "ESSDPL1"
}

resource "alibabacloudstack_polardb_cluster_database" "default" {
  db_cluster_id      = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  db_name            = "${var.name}"
  character_set_name = "utf8"
  db_description     = "Test database for PolarDB cluster"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the PolarDB cluster.
* `db_name` - (Required, ForceNew) The name of the database. 
* `character_set_name` - (Required, ForceNew) The character set of the database.
* `db_description` - (Optional) The description of the database.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the database, formatted as `{DBClusterId}:{DBName}`.
* `engine` - The database engine.
* `db_status` - The status of the database.

## Import

PolarDB cluster database can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_polardb_cluster_database.example pc-12345678:test_database
```