---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbconnection"
sidebar_current: "docs-Alibabacloudstack-resource-polardb-dbconnection"
description: |-
  Provides a PolarDB DB Connection resource.
---

# alibabacloudstack_polardb_dbconnection

Provides a PolarDB DB Connection resource.

## Example Usage

```hcl
resource "alibabacloudstack_polardb_dbconnection" "default" {
  instance_id       = "your_polardb_instance_id"
  connection_prefix = "your_connection_prefix"
  port              = "3306"
}
```

## Argument Reference
The following arguments are supported:

* `instance_id` - (Required, ForceNew) - The ID of the PolarDB instance.
* `connection_prefix` - (Optional, Computed) - The prefix of the connection string. It must be 1 to 31 characters in length and can contain numbers, letters, underscores (_), and hyphens (-). If not specified, the system will automatically generate a default value. Modifying this parameter will not recreate the resource; it can be updated in-place.
* `port` - (Optional) - The port number for the connection. Default is 3306. Valid values range from 1000 to 65534.


## Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `connection_string` - The connection string for the PolarDB instance.
* `ip_address` - The IP address for the connection.

## Import
PolarDB DB Connection can be imported using the id, e.g.

```sh
$ terraform import alibabacloudstack_polardb_dbconnection.example <instance_id>:<connection_prefix>
```

## Example:

```sh
$ terraform import alibabacloudstack_polardb_dbconnection.example polardb-instance-123456:my_connection_prefix
```
