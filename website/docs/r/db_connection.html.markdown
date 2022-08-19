---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_connection"
sidebar_current: "docs-apsarastack-resource-db-connection"
description: |-
  Provides an RDS instance connection resource.
---

# apsarastack\_db\_connection

Provides an RDS connection resource to allocate an Internet connection string for RDS instance.

-> **NOTE:** Each RDS instance will allocate a intranet connnection string automatically and its prifix is RDS instance ID.
 To avoid unnecessary conflict, please specified a internet connection prefix before applying the resource.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbconnectionbasic"
}

data "apsarastack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "apsarastack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = "${apsarastack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.apsarastack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "apsarastack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.t1.small"
  instance_storage = "10"
  vswitch_id       = "${apsarastack_vswitch.default.id}"
  instance_name    = "${var.name}"
}
 
resource "apsarastack_db_connection" "foo" {
  instance_id       = "${apsarastack_db_instance.instance.id}"
  connection_prefix = "testabc"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `connection_prefix` - (ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'tf'.
* `port` - (Optional) Internet connection port. Valid value: [3001-3999]. Default to 3306.

## Attributes Reference

The following attributes are exported:

* `id` - The current instance connection resource ID. Composed of instance ID and connection string with format `<instance_id>:<connection_prefix>`.
* `connection_prefix` - Prefix of a connection string.
* `port` - Connection instance port.
* `connection_string` - Connection instance string.
* `ip_address` - The ip address of connection string.
