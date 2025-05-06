---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_account"
sidebar_current: "docs-Alibabacloudstack-redis-account"
description: |- 
  Provides a redis Account resource.
---

# alibabacloudstack_redis_account
-> **NOTE:** Alias name has: `alibabacloudstack_kvstore_account`

Provides a redis Account resource.

## Example Usage

```hcl
variable "name" {
    default = "tf-testacc-redisaccount73052"
}

variable "kv_edition" {
    default = "enterprise"
}

variable "kv_engine" {
    default = "Redis"
}


data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details = true
}
 
data alibabacloudstack_kvstore_instance_classes "default" {
  zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
  edition_type = "${var.kv_edition}"
  engine = "${var.kv_engine}"
}

locals {
	# data alibabacloudstack_kvstore_instance_classes接口错误
	#default_kv_instance_classes = length(data.alibabacloudstack_kvstore_instance_classes.default.instance_classes) > 0 ? data.alibabacloudstack_kvstore_instance_classes.default.instance_classes[0] : "redis.master.small.default"
	default_kv_instance_classes = "redis.master.small.default"
}

variable "password" {
}

data "alibabacloudstack_zones" "default" {
	available_resource_creation = "VSwitch"
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

resource "alibabacloudstack_kvstore_instance" "default" {
	zone_id = data.alibabacloudstack_zones.kv_zone.zones[0].id
	instance_name  = var.name
	instance_type  = var.kv_engine
	instance_class = local.default_kv_instance_classes
	engine_version = "4.0"
	node_type = "double"
	architecture_type = "standard"
	password       = var.password
	vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_redis_account" "default" {
  instance_id         = "${alibabacloudstack_kvstore_instance.default.id}"
  account_name       = "rdk_test_name_01"
  account_password   = var.password
  description        = "This is a test Redis account."
  account_privilege  = "RoleReadWrite"
  account_type       = "Normal"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required, ForceNew) The ID of the database instance to which the account belongs.
* `account_name` - (Required, ForceNew) The name of the account. It must start with a letter and can consist of lowercase letters, numbers, and underscores (`_`). The maximum length is 16 characters.
* `account_password` - (Optional, Sensitive) The password for the account. It must be between 6 and 32 characters long and can include uppercase and lowercase letters, numbers, and special characters such as `_`, `@`, and `!`. You must specify either `account_password` or `kms_encrypted_password`.
* `kms_encrypted_password` - (Optional) An encrypted password for the account using KMS. If `account_password` is provided, this field will be ignored.
* `kms_encryption_context` - (Optional) The encryption context used to decrypt the `kms_encrypted_password` before creating or updating the account. This is valid only when `kms_encrypted_password` is set.
* `account_type` - (Optional, ForceNew) The type of the account. Valid values:
  * `Normal`: Common privilege.
  Default value is `Normal`.
* `account_privilege` - (Optional) The privilege level of the account. Valid values:
  * `RoleReadOnly`: Read-only access.
  * `RoleReadWrite`: Read and write access.
  * `RoleRepl`: Read, write, and replication commands (`SYNC` / `PSYNC`) access. This is only applicable to Redis instances with an engine version of 4.0 or higher and a standard architecture type.
  Default value is `RoleReadWrite`.
* `description` - (Optional) A description of the account. It must start with a Chinese character or an English letter and can include Chinese characters, English letters, underscores (`_`), hyphens (`-`), and numbers. The length must be between 2 and 256 characters.
* `account_description` - (Optional) The description of the account (same as `description`). 

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the account. It is composed of the instance ID and the account name in the format `<instance_id>:<account_name>`.
* `description` - The description of the account.
* `account_description` - The description of the account (same as `description`). 