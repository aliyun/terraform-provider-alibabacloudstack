---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_read_write_splitting_connection"
sidebar_current: "docs-alibabacloudstack-resource-db-read-write-splitting-connection"
description: |-
  Provides an RDS instance read write splitting connection resource.
---

# alibabacloudstack_db_read_write_splitting_connection

Provides an RDS read write splitting connection resource to allocate an Intranet connection string for RDS instance.

## Example Usage

```
variable "creation" {
  default = "RDS"
}

variable "name" {
  default = "dbInstancevpc"
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

resource "alibabacloudstack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.t1.small"
  instance_storage     = "20"
  instance_charge_type = "Postpaid"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
  security_ips         = ["10.168.1.12", "100.69.7.112"]
}

resource "alibabacloudstack_db_readonly_instance" "default" {
  master_db_instance_id = "${alibabacloudstack_db_instance.default.id}"
  zone_id               = "${alibabacloudstack_db_instance.default.zone_id}"
  engine_version        = "${alibabacloudstack_db_instance.default.engine_version}"
  instance_type         = "${alibabacloudstack_db_instance.default.instance_type}"
  instance_storage      = "30"
  instance_name         = "${var.name}ro"
  vswitch_id            = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_db_read_write_splitting_connection" "default" {
  instance_id       = "${alibabacloudstack_db_instance.default.id}"
  connection_prefix = "t-con-123"
  distribution_type = "Standard"
  depends_on = ["alibabacloudstack_db_readonly_instance.default"]
}
```

-> **NOTE:** Resource `alibabacloudstack_db_read_write_splitting_connection` should be created after `alibabacloudstack_db_readonly_instance`, so the `depends_on` statement is necessary.

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The Id of instance that can run database.
* `distribution_type` - (Required) Read weight distribution mode. Values are as follows: `Standard` indicates automatic weight distribution based on types, `Custom` indicates custom weight distribution. 
* `connection_prefix` - (Optional, ForceNew) Prefix of an Internet connection string. It must be checked for uniqueness. It may consist of lowercase letters, numbers, and underlines, and must start with a letter and have no more than 30 characters. Default to <instance_id> + 'rw'.
* `port` - (Optional) Intranet connection port. Valid value: [3001-3999]. Default to 3306.
* `max_delay_time` - (Optional) Delay threshold, in seconds. The value range is 0 to 7200. Default to 30. Read requests are not routed to the read-only instances with a delay greater than the threshold.  
* `weight` - (Optional) Read weight distribution. Read weights increase at a step of 100 up to 10,000. Enter weights in the following format: {"Instanceid":"Weight","Instanceid":"Weight"}. This parameter must be set when distribution_type is set to Custom. 
* `connection_string` - (Optional)  Connection instance string.

## Attributes Reference

The following attributes are exported:

* `id` - The Id of DB instance.
* `connection_string` - Connection instance string.
* `port` -  Intranet connection port. Valid value: [3001-3999]. Default to 3306.