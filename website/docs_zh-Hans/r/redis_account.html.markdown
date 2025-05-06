---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_account"
sidebar_current: "docs-Alibabacloudstack-redis-account"
description: |- 
  编排Redis账户
---

# alibabacloudstack_redis_account
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_kvstore_account`

使用Provider配置的凭证在指定的资源集编排Redis账户。

## 示例用法

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
	default_kv_instance_classes = "redis.master.small.default"
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
	password       = "1qaz@WSX"
	vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
}

resource "alibabacloudstack_redis_account" "default" {
  instance_id         = "${alibabacloudstack_kvstore_instance.default.id}"
  account_name       = "rdk_test_name_01"
  account_password   = "1qaz@WSX"
  description        = "This is a test Redis account."
  account_privilege  = "RoleReadWrite"
  account_type       = "Normal"
}
```

## 参数说明

支持以下参数：
* `instance_id` - (必填, 变更时重建) - 实例ID，该账户所属的实例。
* `account_name` - (必填, 变更时重建) - 账号名称。它必须以字母开头，并且可以包含小写字母、数字和下划线 (`_`)。最大长度为 16 个字符。
* `account_password` - (选填, 敏感信息) - 账户密码。它必须在 6 到 32 个字符之间，并且可以包括大写和小写字母、数字以及特殊字符如 `_`, `@`, 和 `!`。必须指定 `account_password` 或 `kms_encrypted_password`。
* `kms_encrypted_password` - (选填) - 使用 KMS 加密的账户密码。如果提供了 `account_password`，此字段将被忽略。
* `kms_encryption_context` - (选填) - 用于在创建或更新账户之前解密 `kms_encrypted_password` 的加密上下文。仅当设置了 `kms_encrypted_password` 时有效。
* `account_type` - (选填, 变更时重建) - 账号类型。有效值：
  * `Normal`: 普通权限。
  默认值为 `Normal`。
* `account_privilege` - (选填) - 数据库权限列表。有效值：
  * `RoleReadOnly`: 只读访问。
  * `RoleReadWrite`: 读写访问。
  * `RoleRepl`: 读、写和复制命令(`SYNC` / `PSYNC`)访问。仅适用于引擎版本为 4.0 或更高版本且架构类型为标准的 Redis 实例。
  默认值为 `RoleReadWrite`。
* `description` - (选填) - 账号备注信息。它必须以中文字符或英文字母开头，并且可以包括中文字符、英文字母、下划线 (`_`)、连字符 (`-`) 和数字。长度必须在 2 到 256 个字符之间。
* `account_description` - (选填) - 账号备注信息（与 `description` 相同）。

## 属性说明

除了上述所有参数外，还导出了以下属性：
* `id` - 账户的唯一标识符。它由实例 ID 和账户名称组成，格式为 `<instance_id>:<account_name>`。
* `description` - 账号备注信息。
* `account_description` - 账号备注信息（与 `description` 相同）。