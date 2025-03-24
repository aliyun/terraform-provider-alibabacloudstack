---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_account_privilege"
sidebar_current: "docs-alibabacloudstack-resource-db-account-privilege"
description: |-
  Provides an RDS account privilege resource.
---

# alibabacloudstack_db_account_privilege

Provides an RDS account privilege resource and used to grant several database some access privilege. A database can be granted by multiple account.

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbaccountprivilegebasic"
}

variable "password" {
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s1.small"
  instance_storage = "10"
  vswitch_id       = "${alibabacloudstack_vswitch.default.id}"
  instance_name    = "${var.name}"
}

resource "alibabacloudstack_db_database" "db" {
  count       = 2
  instance_id = "${alibabacloudstack_db_instance.instance.id}"
  name        = "tfaccountpri_${count.index}"
  description = "from terraform"
}

resource "alibabacloudstack_db_account" "account" {
  instance_id = "${alibabacloudstack_db_instance.instance.id}"
  name        = "tftestprivilege"
  password    = var.password
  description = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege" {
  instance_id  = "${alibabacloudstack_db_instance.instance.id}"
  account_name = "${alibabacloudstack_db_account.account.name}"
  privilege    = "ReadOnly"
  db_names     = "${alibabacloudstack_db_database.db.*.name}"
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
* `data_base_instance_id` - (Optional, ForceNew) The ID of the database instance. 

## Attributes Reference

The following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID, account name and privilege with format `<instance_id>:<name>:<privilege>`.
* `instance_id` - The ID of the instance where the account belongs. 
* `data_base_instance_id` - The ID of the database instance. 