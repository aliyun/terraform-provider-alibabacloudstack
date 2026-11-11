---
subcategory: "云数据库 Redis 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_instance"
sidebar_current: "docs-Alibabacloudstack-resource-kvstore-instance"
description: |-
  编排KVStore实例
---

# alibabacloudstack_kvstore_instance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_redis_tairinstance`

使用Provider配置的凭证在指定的资源集编排KVStore（ApsaraDB for Redis / Memcache）实例。

## 示例用法

基础用法

```hcl
variable "name" {
  default = "kvstoreinstanceexample"
}

data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details              = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  zone_id        = data.alibabacloudstack_zones.kv_zone.zones[0].id
  engine         = "Redis"
  engine_version = "5.0"
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.kv_zone.zones.0.id
  name              = var.name
}

resource "alibabacloudstack_kvstore_instance" "default" {
  tair_instance_name  = var.name
  instance_class      = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version      = "5.0"
  zone_id             = data.alibabacloudstack_zones.kv_zone.zones.0.id
  instance_type       = "Redis"
  vswitch_id          = alibabacloudstack_vswitch.default.id
  password            = var.password
  node_type           = "MASTER_SLAVE"
  maintain_start_time = "01:00Z"
  maintain_end_time   = "02:00Z"
  vpc_auth_mode       = "Open"
}
```

多可用区用法

```hcl
variable "name" {
  default = "kvstoreinstancemultizoneexample"
}

data "alibabacloudstack_kvstore_zones" "default" {
  multi = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  engine         = "Redis"
  engine_version = "5.0"
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id       = alibabacloudstack_vpc_vpc.default.id
  cidr_block   = "172.16.1.0/24"
  zone_id      = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[0]
}

resource "alibabacloudstack_kvstore_instance" "default" {
  tair_instance_name = var.name
  instance_class     = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version     = "5.0"
  instance_type      = "Redis"
  vswitch_id         = alibabacloudstack_vpc_vswitch.default.id
  zone_id            = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[0]
  secondary_zone_id  = data.alibabacloudstack_kvstore_zones.default.zones.0.multi_zone_ids[1]
  node_type          = "double"
  password           = var.password
}
```

## 参数说明

支持以下参数：

* `instance_class` - (必填) KVStore实例的实例规格。
* `tair_instance_name` - (可选) 实例的名称。长度必须为2到128个字符。与`instance_name`冲突。
* `instance_name` - (可选，已废弃) 'instance_name'字段已被废弃，并将在未来的版本中移除。请改用新字段`tair_instance_name`。与`tair_instance_name`冲突。
* `password` - (可选，敏感字段) 用于连接实例的密码。密码必须为8到32个字符长，并至少包含以下三种字符类型：大写字母、小写字母、数字和特殊字符(! @ # $ % ^ & * ( ) _ + - =)。如果未指定，则只有在设置密码后才能访问实例。
* `kms_encrypted_password` - (可选) KMS加密的实例密码。如果指定了`password`，此字段将被忽略。
* `kms_encryption_context` - (可选) 用于解密`kms_encrypted_password`的加密上下文。仅当设置了`kms_encrypted_password`时有效。
* `engine_version` - (可选，变更时重建) 数据库引擎版本。默认值为`5.0`。
* `zone_id` - (可选，变更时重建) 实例所在的主可用区ID。与`availability_zone`冲突。
* `secondary_zone_id` - (可选，变更时重建) 备可用区ID。如果指定该参数，将创建多可用区（双可用区）实例；否则创建单可用区实例。可通过数据源`alibabacloudstack_kvstore_zones`的`multi_zone_ids`属性获取。
* `availability_zone` - (可选，变更时重建，已废弃) 'availability_zone'字段已被废弃，并将在未来的版本中移除。请改用新字段`zone_id`。与`zone_id`冲突。
* `instance_charge_type` - (可选，已废弃) 'instance_charge_type'字段已被废弃，将在未来的版本中移除，当前已无任何作用。有效值：`PrePaid`、`PostPaid`。
* `instance_type` - (可选，变更时重建) 实例的引擎类型。有效值：`Redis`、`Memcache`。默认值为`Redis`。
* `vswitch_id` - (可选，变更时重建) VSwitch的ID。如果指定该参数，实例将部署在VPC网络中。
* `backup_id` - (可选) 如果实例是基于另一个实例的备份创建的，则为备份集的ID。
* `security_ips` - (可选) 允许访问该实例的IP地址集合。
* `vpc_auth_mode` - (可选) VPC认证模式。有效值：`Open`(启用密码认证)、`Close`(禁用密码认证并启用免密访问)。仅当网络类型为VPC时生效。
* `parameters` - (可选) 用于配置实例的参数集合。每个元素支持以下字段：
  * `name` - (必填) 参数的名称。
  * `value` - (必填) 参数的值。
* `maintain_start_time` - (可选) 维护窗口的开始时间。格式为`HH:mmZ`(UTC时间)。
* `maintain_end_time` - (可选) 维护窗口的结束时间。格式为`HH:mmZ`(UTC时间)。
* `cpu_type` - (可选) 实例的CPU类型，例如`intel`。
* `node_type` - (可选) 实例的节点类型。有效值：`MASTER_SLAVE`、`STAND_ALONE`、`double`、`readone`。
* `architecture_type` - (可选，已废弃) 'architecture_type'字段不能正常工作，请通过'instance_class'设置。有效值：`cluster`、`rwsplit`、`standard`。
* `series` - (可选，已废弃) 'series'字段不能正常工作，请通过'instance_class'设置。有效值：`community`、`enterprise`。
* `tde_status` - (可选) 是否开启透明数据加密(TDE)。有效值：`true`、`false`。TDE开启后无法关闭。
* `enable_ssl` - (可选) 是否为实例开启SSL。有效值：`true`、`false`。
* `encryption_key` - (可选，已废弃) 用于TDE的加密密钥。TDE不支持同时使用`encryption_key`和`role_arn`。与`role_arn`冲突。
* `role_arn` - (可选，已废弃) 用于TDE的RAM角色ARN。TDE不支持同时使用`encryption_key`和`role_arn`。与`encryption_key`冲突。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `id` - KVStore实例的ID。
* `connection_domain` - 实例的内部端点。
* `private_ip` - 实例的私有IP地址。

## 超时配置

* `create` - （默认45分钟）创建资源时使用。
* `update` - （默认30分钟）更新资源时使用。
* `delete` - （默认45分钟）删除资源时使用。

## 存量资源纳管

alibabacloudstack_kvstore_instance 可通过ID导入，例如：

```shell
$ terraform import alibabacloudstack_kvstore_instance.example <id>
```
