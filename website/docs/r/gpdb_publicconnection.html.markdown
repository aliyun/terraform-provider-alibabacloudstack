---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_publicconnection"
sidebar_current: "docs-Alibabacloudstack-gpdb-publicconnection"
description: |- 
  Provides a gpdb Publicconnection resource.
---

# alibabacloudstack_gpdb_publicconnection
-> **NOTE:** Alias name has: `alibabacloudstack_gpdb_connection`

Provides a gpdb Publicconnection resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "Gpdb"
}

variable "name" {
  default = "tf-testAccGpdbInstance"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "testing"
  cidr_block = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "10.1.0.0/16"
  name              = "apsara_vswitch"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

resource "alibabacloudstack_gpdb_instance" "default" {
  vswitch_id           = alibabacloudstack_vswitch.default.id
  engine               = "gpdb"
  engine_version       = "4.3"
  instance_class       = "gpdb.group.segsdx2"
  instance_group_count = "2"
  description          = var.name
}

resource "alibabacloudstack_gpdb_connection" "default" {
  instance_id       = alibabacloudstack_gpdb_instance.default.id
  connection_prefix = "tf-testacc10623"
  port              = 3306
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the GPDB instance for which the public connection will be created.
* `connection_prefix` - (Optional, ForceNew) The prefix of the public connection string. It must start with a letter and can only contain lowercase letters, numbers, and underscores (`_`). The length cannot exceed 30 characters. If not specified, it defaults to `<instance_id>-tf`.
* `port` - (Optional) The port number for the public connection. Valid values range from `3200` to `3999`. Default value is `3306`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Default `10 mins`) Used when creating the public connection.
* `update` - (Default `10 mins`) Used when updating the public connection.
* `delete` - (Default `10 mins`) Used when deleting the public connection.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the GPDB public connection resource. It is composed of the instance ID and the connection prefix in the format `<instance_id>:<connection_prefix>`.
* `connection_string` - The complete connection string for accessing the GPDB instance via the public network.
* `ip_address` - The public IP address associated with the public connection string.