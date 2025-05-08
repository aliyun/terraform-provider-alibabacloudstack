---
subcategory: "AnalyticDB for MySQL (ADB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_connection"
sidebar_current: "docs-Alibabacloudstack-adb-connection"
description: |- 
  Provides a adb Connection resource.
---

# alibabacloudstack_adb_connection

Provides a adb Connection resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testaccadbconnection96904"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alibabacloudstack_adb_db_cluster" "cluster" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode                = "reserver"
  vswitch_id          = alibabacloudstack_vswitch.default.id
  description         = var.name
  cluster_type        = "analyticdb"
  cpu_type            = "intel"
}

resource "alibabacloudstack_adb_connection" "default" {
  db_cluster_id         = alibabacloudstack_adb_db_cluster.cluster.id
  connection_prefix     = "testabc"
}
```

## Argument Reference

The following arguments are supported:

* `db_cluster_id` - (Required, ForceNew) The ID of the ADB cluster that can run the database. This field is immutable and cannot be changed after creation.
* `connection_prefix` - (Optional, ForceNew) Prefix of the cluster's public endpoint. The prefix must be 6 to 30 characters in length and can contain lowercase letters, digits, and hyphens (-). It must start with a letter and end with a digit or letter. If not specified, it defaults to `<db_cluster_id> + tf`.
* `port` - (Optional, Computed) The port number used for connecting to the ADB cluster.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier for the ADB connection resource. It is composed of the cluster ID and connection string in the format `<db_cluster_id>:<connection_prefix>`.
* `connection_prefix` - The prefix used for the connection string.
* `port` - The port number used for connecting to the ADB cluster.
* `connection_string` - The full connection string for accessing the ADB cluster.
* `ip_address` - The IP address associated with the connection string. 

## Import

ADB connection can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_adb_connection.example am-12345678
```