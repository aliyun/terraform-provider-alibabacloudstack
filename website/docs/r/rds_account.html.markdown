---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_account"
sidebar_current: "docs-Alibabacloudstack-rds-account"
description: |- 
  Provides a rds Account resource.
---

# alibabacloudstack_rds_account
-> **NOTE:** Alias name has: `alibabacloudstack_db_account`

Provides a rds Account resource.

## Example Usage

```hcl
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

variable "creation" {
  default = "Rds"
}

variable "name" {
  default = "tf-testAccdbaccount-52957"
}

variable "password" {
}

resource "alibabacloudstack_db_instance" "instance" {
  engine               = "MySQL"
  engine_version       = "5.6"
  instance_type        = "rds.mysql.s2.large"
  instance_storage     = "30"
  vswitch_id           = "${alibabacloudstack_vswitch.default.id}"
  instance_name        = "${var.name}"
  storage_type         = "local_ssd"
}

resource "alibabacloudstack_rds_account" "default" {
  data_base_instance_id = "${alibabacloudstack_db_instance.instance.id}"
  account_name          = "tftestnormal"
  password              = var.password
  account_type          = "Normal"
  account_description   = "This is a test account."
}
```

## Argument Reference

The following arguments are supported:

* `data_base_instance_id` - (Required, ForceNew) The instance ID. You can call the `DescribeDBInstances` operation to query the instance ID.
* `account_name` - (Required, ForceNew) Database account name. Consists of lowercase letters, numbers, or underscores. For MySQL, uppercase letters are also supported. It must start with a letter and end with a letter or number.  
  - **Length**: 
    - MySQL 8.0 and 5.7: 2 to 32 characters.
    - MySQL 5.6: 2 to 16 characters.
    - SQL Server: 2 to 64 characters.
    - PostgreSQL cloud disk version: 2 to 63 characters.
    - PostgreSQL local disk version: 2 to 16 characters.
    - MariaDB: 2 to 16 characters.
  - **Note**: The common account name and the high-privilege account name cannot be similar to each other. For example, if the high-privilege account name is `Test1`, the common account name cannot be `test1`.
* `password` - (Required) Operation password. It may consist of letters, digits, or underlines, with a length of 6 to 32 characters.
* `kms_encrypted_password` - (Optional) An KMS encrypts password used to a db account. If the `password` is filled in, this field will be ignored.
* `kms_encryption_context` - (Optional) An KMS encryption context used to decrypt `kms_encrypted_password` before creating or updating a db account with `kms_encrypted_password`. See [Encryption Context](https://www.alibabacloud.com/help/doc-detail/42975.htm). It is valid when `kms_encrypted_password` is set.
* `account_type` - (Optional, ForceNew) Account type. Valid values:
  - `Normal`: Normal account (default).
  - `Super`: High-privilege account.
  - `Sysadmin`: A database account with SA permissions (only RDS SQL Server instances are supported).
  - **Note**: Before creating a database account with SA permissions, check whether the instance meets the prerequisites. For more information, see [create a database account with SA permissions](https://www.alibabacloud.com/help/doc-detail/122334.htm).
* `account_description` - (Optional) Account description. It can be 2 to 256 characters in length. It starts with a Chinese or English letter and can contain numbers, Chinese, English, underscores (`_`), and hyphens (`-`).  
  - **Note**: Cannot start with `http://` or `https://`.
* `instance_id` - (Optional, ForceNew) Deprecated field, use `data_base_instance_id` instead.
* `name` - (Optional, ForceNew) Deprecated field, use `account_name` instead.
* `type` - (Optional, ForceNew) Deprecated field, use `account_type` instead.
* `description` - (Optional) Deprecated field, use `account_description` instead.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The current account resource ID. Composed of instance ID and account name with format `<instance_id>:<account_name>`.
* `data_base_instance_id` - The instance ID. You can call the `DescribeDBInstances` operation to query the instance ID.
* `account_name` - Database account name. Consists of lowercase letters, numbers, or underscores. For MySQL, uppercase letters are also supported. It must start with a letter and end with a letter or number.
* `account_type` - Account type. Valid values: `Normal`, `Super`, or `Sysadmin`.
* `account_description` - Account description. It can be 2 to 256 characters in length. It starts with a Chinese or English letter and can contain numbers, Chinese, English, underscores (`_`), and hyphens (`-`).
* `instance_id` - Deprecated field, use `data_base_instance_id` instead.
* `name` - Deprecated field, use `account_name` instead.
* `type` - Deprecated field, use `account_type` instead.
* `description` - Deprecated field, use `account_description` instead.