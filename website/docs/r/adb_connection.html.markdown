---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_adb_connection"
sidebar_current: "docs-apsarastack-resource-adb-connection"
description: |-
  Provides an ADB cluster connection resource.
---

# apsarastack\_adb\_connection

Provides an ADB connection resource to allocate an Internet connection string for ADB cluster.

-> **NOTE:** Each ADB instance will allocate a intranet connnection string automatically and its prifix is ADB instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```
variable "creation" {
  default = "ADB"
}

variable "name" {
  default = "adbaccountmysql"
}

data "apsarastack_zones" "default" {
  available_resource_creation = var.creation
}

resource "apsarastack_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = apsarastack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.apsarastack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "apsarastack_adb_db_cluster" "cluster" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode				  = "reserver"
  vswitch_id          = apsarastack_vswitch.default.id
  description         = var.name
  cluster_type      = "analyticdb"
  cpu_type          = "intel"
}

resource "apsarastack_adb_connection" "connection" {
  db_cluster_id     = apsarastack_adb_db_cluster.cluster.id
  connection_prefix = "testabc"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The Id of cluster that can run database.
* `connection_prefix` - (ForceNew) Prefix of the cluster public endpoint. The prefix must be 6 to 30 characters in length, and can contain lowercase letters, digits, and hyphens (-), must start with a letter and end with a digit or letter. Default to `<db_cluster_id> + tf`.

## Attributes Reference

The following attributes are exported:

* `id` - The current cluster connection resource ID. Composed of cluster ID and connection string with format `<db_cluster_id>:<connection_prefix>`.
* `connection_prefix` - Prefix of a connection string.
* `port` - Connection cluster port.
* `connection_string` - Connection cluster string.
* `ip_address` - The ip address of connection string.

## Import

ADB connection can be imported using the id, e.g.

```
$ terraform import apsarastack_adb_connection.example am-12345678
```
