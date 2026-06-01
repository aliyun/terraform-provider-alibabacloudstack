---
subcategory: "ApsaraDB for PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_readwrite_splitting_connection"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-readwrite-splitting-connection"
description: |-
  Provides a PolarDB Read/Write Splitting Connection resource.
---

# alibabacloudstack\_polardb\_readwrite\_splitting\_connection

Provides a PolarDB Read/Write Splitting Connection resource that allows you to configure read/write splitting for a PolarDB cluster.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf-testAccPolarDB"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  db_type              = "MySQL"
  db_version           = "8.0"
  pay_type             = "Postpaid"
  db_node_class        = "polar.mysql.x4.medium"
  zone_id              = data.alibabacloudstack_zones.default.zones.0.id
  vswitch_id           = alibabacloudstack_vswitch.default.id
  security_ip_list     = ["10.168.1.12", "100.69.7.112"]
  instance_network_type = "VPC"
}

resource "alibabacloudstack_polardb_readonly_instance" "default" {
  master_db_instance_id = alibabacloudstack_polardb_dbinstance.default.id
  zone_id               = alibabacloudstack_polardb_dbinstance.default.zone_id
  engine_version        = alibabacloudstack_polardb_dbinstance.default.engine_version
  instance_type         = alibabacloudstack_polardb_dbinstance.default.instance_type
  instance_storage      = alibabacloudstack_polardb_dbinstance.default.instance_storage
  instance_name         = var.name
  db_instance_storage_type = alibabacloudstack_polardb_dbinstance.default.storage_type
}

resource "alibabacloudstack_polardb_proxy" "default" {
  db_instance_id        = alibabacloudstack_polardb_dbinstance.default.id
  db_proxy_instance_num = 1
}

resource "alibabacloudstack_polardb_readwrite_splitting_connection" "default" {
  instance_id       = alibabacloudstack_polardb_dbinstance.default.id
  connection_id     = alibabacloudstack_polardb_proxy.default.db_proxy_endpoint_name
  distribution_type = "Standard"
  max_delay_time    = 30
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the PolarDB cluster.
* `connection_id` - (Required, ForceNew) The ID of the proxy endpoint.
* `distribution_type` - (Required) The distribution type of read requests. Valid values: `Standard`, `Custom`.
  - `Standard`: System automatically distributes read requests to read-only nodes based on their weights.
  - `Custom`: You can customize the weight of each read-only node.
* `weight` - (Optional, Computed) The weight of each read-only node. This parameter is required when `distribution_type` is set to `Custom`. The value is a map where the key is the read-only node ID and the value is the weight (0-100).
* `max_delay_time` - (Optional) The maximum delay time threshold for read-only nodes. Unit: seconds. Default value: 30. If the delay of a read-only node exceeds this threshold, read requests will not be distributed to that node.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource, which is the same as `instance_id`.
* `connection_string` - The connection string of the read/write splitting endpoint.
* `port` - The port number of the read/write splitting endpoint.

## Import

PolarDB Read/Write Splitting Connection can be imported using the instance ID, e.g.

```
$ terraform import alibabacloudstack_polardb_readwrite_splitting_connection.example pc-xxxxxxxxxxxxx
```
