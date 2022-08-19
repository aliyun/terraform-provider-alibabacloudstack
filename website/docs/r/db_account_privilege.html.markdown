---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_account_privilege"
sidebar_current: "docs-apsarastack-resource-db-account-privilege"
description: |-
  Provides an RDS account privilege resource.
---

# apsarastack\_db\_account\_privilege

Provides an RDS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountprivilegebasic"
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
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${apsarastack_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "apsarastack_db_database" "db" {
  count       = 2
  instance_id = "${apsarastack_db_instance.instance.id}"
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "apsarastack_db_account" "account" {
  instance_id = "${apsarastack_db_instance.instance.id}"
  name        = "tftestprivilege"
  password    = "Test12345"
  description = "from terraform"
}

resource "apsarastack_db_account_privilege" "privilege" {
  instance_id  = "${apsarastack_db_instance.instance.id}"
  account_name = "${apsarastack_db_account.account.name}"
  privilege    = "ReadOnly"
  db_names     = "${apsarastack_db_database.db.*.name}"
}
```

## Argument Reference
 
The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance in which account belongs.
* `account_name` - (Required, ForceNew) A specified account name.
* `privilege` - The privilege of one account access database. Valid values: 
    - ReadOnly: This value is only for MySQL, MariaDB and SQL Server
    - ReadWrite: This value is only for MySQL, MariaDB and SQL Server
     
   Default to "ReadOnly". 
* `db_names` - (Required) List of specified database name.

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<instance_id>:<name>:<privilege>`.
