---
subcategory: "RDS"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_db_readonly_instance"
sidebar_current: "docs-apsarastack-resource-db-readonly-instance"
description: |-
  Provides an RDS readonly instance resource.
---

# apsarastack\_db\_readonly\_instance

Provides an RDS readonly instance resource. 

## Example Usage

```
variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "dbInstancevpc"
}

data "apsarastack_zones" "default" {
  available_resource_creation = var.creation
}

resource "apsarastack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "apsarastack_vswitch" "default" {
  vpc_id            = apsarastack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.apsarastack_zones.default.zones[0].id
  name              = var.name
}

resource "apsarastack_db_instance" "default" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.t1.small"
  instance_storage     = "20"
  instance_name        = var.name
  vswitch_id           = apsarastack_vswitch.default.id
}

resource "apsarastack_db_connection" "connection" {
  instance_id       = apsarastack_db_instance.default.id
  connection_prefix = var.connection_prefix
}

resource "apsarastack_db_readonly_instance" "default" {
  master_db_instance_id = apsarastack_db_instance.default.id
  zone_id               = apsarastack_db_instance.default.zone_id
  engine_version        = apsarastack_db_instance.default.engine_version
  instance_type         = apsarastack_db_instance.default.instance_type
  instance_storage      = "30"
  instance_name         = "${var.name}ro"
  vswitch_id            = apsarastack_vswitch.default.id
  db_instance_storage_type= "local_ssd"
  depends_on = [apsarastack_db_connection.connection]
}
```

## Argument Reference

The following arguments are supported:

* `engine_version` - (Required, ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `master_db_instance_id` - (Required) ID of the master instance.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range: [5, 2000] for MySQL/SQL Server HA dual node edition. Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm).
* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.
* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `db_instance_storage_type` - (Required) The storage type of the instance. Valid values:
    local_ssd: specifies to use local SSDs. This value is recommended.
    cloud_ssd: specifies to use standard SSDs.
    cloud_essd: specifies to use enhanced SSDs (ESSDs).
    cloud_essd2: specifies to use enhanced SSDs (ESSDs).
    cloud_essd3: specifies to use enhanced SSDs (ESSDs).
    
-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `engine` - Database type.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.

### Timeouts


The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 60 mins) Used when creating the db instance (until it reaches the initial `Running` status). 
* `update` - (Defaults to 30 mins) Used when updating the db instance (until it reaches the initial `Running` status). 
* `delete` - (Defaults to 20 mins) Used when terminating the db instance. 

## Import

RDS readonly instance can be imported using the id, e.g.

```
$ terraform import apsarastack_db_readonly_instance.example rm-abc12345678
```
