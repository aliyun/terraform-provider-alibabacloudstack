---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_instance"
sidebar_current: "docs-alibabacloudstack-resource-db-instance"
description: |-
  Provides an RDS instance resource.
---

# alibabacloudstack\_db\_instance

Provides an RDS instance resource. A DB instance is an isolated database
environment in the cloud. A DB instance can contain multiple user-created
databases.

## Example Usage

### Create a RDS MySQL instance

```
variable "name" {
  default = "dbInstanceconfig"
}
variable "creation" {
  default = "Rds"
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
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  storage_type     = "local_ssd"
  instance_name        = "${var.name}"
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
  encryption_key="f23ed1c9-b91f-......"
  tde_status=false
  enable_ssl=false
  zone_id_slave1="${data.alibabacloudstack_zones.default.zones.0.id}"
  zone_id="${data.alibabacloudstack_zones.default.zones.0.id}"
}
```

### Create a RDS MySQL instance with specific parameters

```
resource "alibabacloudstack_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "vpc-123456"
}

resource "alibabacloudstack_db_instance" "default1" {
  engine              = "MySQL"
  engine_version      = "5.6"
  instance_type   = "rds.mysql.t1.small"
  instance_storage = "10"
  vswitch_id          = "${alibabacloudstack_vswitch.default.id}"
  storage_type     = "local_ssd"
  encryption_key="f23ed1c9-b91f-......"
  zone_id_slave1="${data.alibabacloudstack_zones.default.zones.0.id}"
  zone_id="${data.alibabacloudstack_zones.default.zones.0.id}"
  tde_status=false
  enable_ssl=false
}

resource "alibabacloudstack_db_instance" "default2" {
  engine              = "MySQL"
  engine_version      = "5.6"
  instance_type   = "rds.mysql.t1.small"
  storage_type     = "local_ssd"
  instance_storage = "10"
  
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required,ForceNew) Database type. Value options: MySQL, SQLServer, PostgreSQL, and PPAS.
* `engine_version` - (Required,ForceNew) Database version. Value options can refer to the latest docs [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) `EngineVersion`.
* `instance_type` - (Required) DB Instance type. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
* `instance_storage` - (Required) User-defined DB instance storage space. Value range:
    - [5, 2000] for MySQL/PostgreSQL/PPAS HA dual node edition;
    - [20,1000] for MySQL 5.7 basic single node edition;
    - [10, 2000] for SQL Server 2008R2;
    - [20,2000] for SQL Server 2012 basic single node edition
    Increase progressively at a rate of 5 GB. For details, see [Instance type table](https://www.alibabacloud.com/help/doc-detail/26312.htm).
    Note: There is extra 5 GB storage for SQL Server Instance and it is not in specified `instance_storage`.
* `storage_type` - (Required) The type of storage media that is used for the instance.
* `instance_name` - (Optional) The name of DB instance. It a string of 2 to 256 characters.
* `zone_id` - (ForceNew) The Zone to launch the DB instance.
* `encryption_key` - (Optional) Add encryptionkey to the DBInstance.
* `zone_id_slave1` - (Optional) The zone ID of the secondary instance.
* `zone_id_slave2` - (Optional) The zone ID of the secondary instance.
* `tde_status` - (Optional) Enables the Transparent Data Encryption (TDE) function for an ApsaraDB for RDS instance.
* `encryption` - (Optional) Enable encryption or not.
* `enable_ssl` - (Optional) To enable the SSL encryption of an ApsaraDB RDS instance.If it is a multi-zone and `vswitch_id` is specified, the vswitch must in the one of them.The multiple zone ID can be retrieved by setting `multi` to "true" in the data source `alibabacloudstack_zones`.
* `vswitch_id` - (ForceNew) The virtual switch ID to launch DB instances in one VPC.
* `security_ips` - (Optional) List of IP addresses allowed to access all databases of an instance. The list contains up to 1,000 IP addresses, separated by commas. Supported formats include 0.0.0.0/0, 10.23.12.24 (IP), and 10.23.12.24/24 (Classless Inter-Domain Routing (CIDR) mode. /24 represents the length of the prefix in an IP address. The range of the prefix length is [1,32]).
* `role_arn` - (ForceNew, Optional) The Alibaba Cloud Resource Name (ARN) of the RAM role.
* `period` - (Optional) The duration that you will buy DB instance (in month). It is valid when instance_charge_type is PrePaid. Valid values: [1~9], 12, 24, 36.
* `security_ip_mode` - (Optional) Valid values are normal, safety, Default to normal. support safety switch to high security access mode.
* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values are 5, 10, 60, 300. Defaults to 300.
* `maintain_time` - (Optional) Maintainable time period format of the instance: HH:MMZ-HH:MMZ (UTC time)
* `auto_renew` - (Optional) Whether to renewal a DB instance automatically or not. It is valid when instance_charge_type is PrePaid. Default to false.
* `auto_renew_period` - (Optional) Auto-renewal period of an instance, in the unit of the month. It is valid when instance_charge_type is PrePaid. Valid value:[1~12], Default to 1.
* `instance_charge_type` - (Optional) Valid values are Prepaid, Postpaid, Serverless, Default to Postpaid. Currently, the resource only supports PostPaid to PrePaid.
* `tags` - (Optional) A mapping of tags to assign to the resource.
    - Key: It can be up to 64 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It cannot be a null string.
    - Value: It can be up to 128 characters in length. It cannot begin with "aliyun", "acs:", "http://", or "https://". It can be a null string.
* `parameters` - (Optional) Set of parameters needs to be set after DB instance was launched. Available parameters can refer to the latest docs [View database parameter templates](https://www.alibabacloud.com/help/doc-detail/26284.htm) . See [`parameters`](#parameters) below.    
  * `name` - (Required) The parameter name.
  * `value` - (Required) The parameter value.
-> **NOTE:** Because of data backup and migration, change DB instance type and storage would cost 15~20 minutes. Please make full preparation before changing them.

## Attributes Reference

The following attributes are exported:

* `id` - The RDS instance ID.
* `port` - RDS database connection port.
* `connection_string` - RDS database connection string.