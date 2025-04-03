---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_dbinstance"
sidebar_current: "docs-Alibabacloudstack-rds-dbinstance"
description: |- 
  Provides a rds Dbinstance resource.
---

# alibabacloudstack_rds_dbinstance
-> **NOTE:** Alias name has: `alibabacloudstack_db_instance`

Provides a rds Dbinstance resource.

## Example Usage

### Create a RDS MySQL instance

```hcl
variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default" {
  instance_name           = "${var.name}"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.s2.large"
  db_instance_storage     = 20
  maintain_time           = "03:00Z-04:00Z"
  security_ip_mode        = "normal"
  role_arn                = "acs:ram::123456789012:role/example-role"
}
```

### Create a RDS MySQL instance with specific parameters

```hcl
resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "vpc-123456"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default1" {
  instance_name           = "tf-testAccDBInstanceConfig1"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.t1.small"
  db_instance_storage     = 10
  encryption_key          = "f23ed1c9-b91f-..."
  tde_status              = false
  enable_ssl             = false
  zone_id_slave1         = "${data.alibabacloudstack_zones.default.zones.0.id}"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default2" {
  instance_name           = "tf-testAccDBInstanceConfig2"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.t1.small"
  db_instance_storage     = 10
}
```

## Argument Reference

The following arguments are supported:

* `engine` - (Required, ForceNew) Database type. Return value:
  * **MySQL**
  * **PostgreSQL**
  * **SQLServer**
  * **MariaDB**

* `engine_version` - (Required, ForceNew) Database version. For example, for MySQL, valid versions include `5.6`, `5.7`, and `8.0`.

* `zone_id_slave1` - (Optional) The zone ID of the first standby instance.

* `zone_id_slave2` - (Optional) The zone ID of the second standby instance.

* `tde_status` - (Optional) Enables Transparent Data Encryption (TDE) for the RDS instance.

* `enable_ssl` - (Optional, ForceNew) Specifies whether to enable SSL encryption for the RDS instance.

* `storage_type` - (Optional, ForceNew) The type of storage media used for the instance. Valid values:
  * **local_ssd**: Local SSD disk.
  * **ephemeral_ssd**: Ephemeral SSD disk.
  * **cloud_ssd**: Cloud SSD disk.
  * **cloud_essd**: Cloud ESSD disk.

* `db_instance_storage_type` - (Optional, ForceNew) The storage type of the instance. Valid values:
  * **local_ssd**, **ephemeral_ssd**: Local SSD disk.
  * **cloud_ssd**: Cloud SSD disk.
  * **cloud_essd**: Cloud ESSD disk.

* `encryption_key` - (Optional) The key ID of the disk encryption in the same region. This parameter indicates that cloud disk encryption is enabled and cannot be turned off after it is enabled. You can view the key ID in the Key Management Service (KMS) console or create a new key.

* `encryption` - (Optional, ForceNew) Specifies whether to enable encryption.

* `db_instance_class` - (Optional) Instance type. For more information, see [Instance Type Table](https://www.alibabacloud.com/help/doc-detail/26312.htm).

* `db_instance_storage` - (Optional) User-defined DB instance storage space. Value range depends on the database type and edition. Increase progressively at a rate of 5 GB.

* `instance_charge_type` - (Optional) The billing method of the instance. Valid values:
  * **Prepaid**: Prepaid billing.
  * **Postpaid**: Postpaid billing.
  * **Serverless**: Serverless billing (only supports MySQL instances).

* `payment_type` - (Optional) The payment type of the instance. Valid values:
  * **PayAsYouGo**: Pay-as-you-go.
  * **Subscription**: Subscription.
  * **Serverless**: Serverless paid type (only supports MySQL instances).

* `period` - (Optional) The subscription duration of the instance. Valid values:
  * Year: Annual subscription.
  * Month: Monthly subscription.

* `monitoring_period` - (Optional) The monitoring frequency in seconds. Valid values: `5`, `10`, `60`, `300`. Defaults to `300`.

* `auto_renew` - (Optional) Specifies whether the instance is automatically renewed. It is only passed in when creating a subscription instance. Valid values:
  * **true**
  * **false**

* `auto_renew_period` - (Optional) The automatic renewal period of the instance. Valid values: `1`~`12` months.

* `zone_id` - (Optional, ForceNew) The Zone to launch the DB instance.

* `vswitch_id` - (Optional, ForceNew) The virtual switch ID to launch DB instances in one VPC.

* `instance_name` - (Optional) The name of the DB instance. It must be 2 to 256 characters in length.

* `db_instance_description` - (Optional) The description of the DB instance.

* `security_ip_mode` - (Optional) Specifies the security IP mode. Valid values:
  * **normal**: Normal mode.
  * **safety**: High-security mode.

* `maintain_time` - (Optional) The maintenance time period of the instance, which is the UTC time. Format: `HH:MMZ-HH:MMZ`.

* `role_arn` - (Optional) The Global Resource Descriptor (ARN) that authorizes the RDS cloud service account to access KMS. You can view the ARN information by using the [CheckCloudResourceAuthorized](~~ 446261 ~~) interface.

* `force_restart` - (Optional) Specifies whether to force restart the instance.

* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the RDS instance.

* `port` - The connection port of the RDS instance.

* `connection_string` - The connection string of the RDS instance.

* `storage_type` - The type of storage media used for the instance.

* `db_instance_storage_type` - The storage type of the instance. Valid values:
  * **local_ssd**, **ephemeral_ssd**: Local SSD disk.
  * **cloud_ssd**: Cloud SSD disk.
  * **cloud_essd**: Cloud ESSD disk.

* `db_instance_class` - The instance type.

* `db_instance_storage` - The storage size of the instance.

* `instance_charge_type` - The billing method of the instance.

* `payment_type` - The payment type of the instance.

* `monitoring_period` - The monitoring frequency in seconds.

* `zone_id` - The Zone ID of the instance.

* `instance_name` - The name of the DB instance.

* `db_instance_description` - The description of the DB instance.

* `maintain_time` - The maintenance time period of the instance.

* `role_arn` - The Global Resource Descriptor (ARN) that authorizes the RDS cloud service account to access KMS.