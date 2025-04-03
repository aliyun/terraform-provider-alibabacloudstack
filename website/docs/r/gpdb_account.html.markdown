---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_account"
sidebar_current: "docs-Alibabacloudstack-gpdb-account"
description: |- 
  Provides a gpdb Account resource.
---

# alibabacloudstack_gpdb_account

Provides a gpdb Account resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tftest1124"
}
variable "password" {
}
data "alibabacloudstack_gpdb_zones" "default" {}
data "alibabacloudstack_zones" "default" {}
data "alibabacloudstack_vpcs" "default" {
  name_regex = "default-NODELETING"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_gpdb_instance" "default" {
  availability_zone      = data.alibabacloudstack_zones.default.zones.0.id
  engine                = "gpdb"
  engine_version        = "4.3"
  instance_class        = "gpdb.group.segsdx2"
  instance_group_count  = 2
  description          = "tf-testAccGpdbInstance_new"
  vswitch_id           = alibabacloudstack_vswitch.default.id
}

resource "alibabacloudstack_gpdb_account" "default" {
  account_name         = "tftest1124"
  account_password     = var.password
  account_description  = "tftest1124"
  db_instance_id       = alibabacloudstack_gpdb_instance.default.id
}
```

## Argument Reference

The following arguments are supported:

* `account_description` - (Optional, ForceNew) The description of the account.  
  * Starts with a letter.
  * Does not start with `http://` or `https://`.
  * Contains letters, underscores (_), hyphens (-), or digits.
  * Be 2 to 256 characters in length.

* `account_name` - (Required, ForceNew) The name of the account. The account name must be unique and meet the following requirements:
  * Starts with a letter.
  * Contains only lowercase letters, digits, or underscores (_).
  * Be up to 16 characters in length.
  * Contains no reserved keywords.

* `account_password` - (Required) The password of the account. The password must be 8 to 32 characters in length and contain at least three of the following character types: uppercase letters, lowercase letters, digits, and special characters. Special characters include `! @ # $ % ^ & * ( ) _ + - =`.

* `db_instance_id` - (Required, ForceNew) The ID of the instance.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID of Account. The value formats as `<db_instance_id>:<account_name>`.
* `status` - The status of the account. Valid values: `Active`, `Creating`, and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 5 mins) Used when create the Account.

## Import

GPDB Account can be imported using the id, e.g.

```bash
$ terraform import alicloud_gpdb_account.example <db_instance_id>:<account_name>
```