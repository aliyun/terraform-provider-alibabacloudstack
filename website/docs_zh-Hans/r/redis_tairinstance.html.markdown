---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_redis_tairinstance"
sidebar_current: "docs-Alibabacloudstack-redis-tairinstance"
description: |-
  编排Redis实例
---

# alibabacloudstack_redis_tairinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_kvstore_instance`

使用Provider配置的凭证在指定的资源集编排Redis实例。

## 示例用法

以下是一个完整的示例，展示如何创建一个 Redis Tair 实例：

```hcl
variable "name" {
    default = "tf-testAccCheckAlibabacloudStackRKVInstances92773"
}

variable "kv_edition" {
    default = "community"
}

variable "kv_engine" {
    default = "Redis"
}

variable "password" {
    default = "1qaz@WSX"
}

data "alibabacloudstack_zones" "kv_zone" {
  available_resource_creation = "KVStore"
  enable_details = true
}

data "alibabacloudstack_kvstore_instance_classes" "default" {
  zone_id      = data.alibabacloudstack_zones.kv_zone.zones[0].id
  edition_type = var.kv_edition
  engine       = var.kv_engine
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

resource "alibabacloudstack_redis_tairinstance" "default" {
  tair_instance_name = var.name
  instance_class     = data.alibabacloudstack_kvstore_instance_classes.default.instance_classes.0.instance_class
  engine_version     = "5.0"
  zone_id           = data.alibabacloudstack_zones.kv_zone.zones.0.id
  instance_type     = "tair_rdb"
  vswitch_id        = alibabacloudstack_vswitch.default.id
  password          = var.password
  node_type         = "MASTER_SLAVE"
  architecture_type = "standard"
  maintain_start_time = "01:00Z"
  maintain_end_time   = "02:00Z"
  vpc_auth_mode      = "Open"
}
```

## 参数说明

支持以下参数：

* `tair_instance_name` - (可选) 资源的名称。它必须是2到256个字符长，并以字母或数字开头。它可以包含下划线、字母和数字。
* `password` - (可选) 用于连接实例的密码。密码必须为8到32个字符长，并至少包含以下三种字符类型：大写字母、小写字母、数字和特殊字符(! @ # $ % ^ & * ( ) _ + - =)。
* `kms_encrypted_password` - (可选) KMS加密的实例密码。如果指定了`password`，此字段将被忽略。
* `kms_encryption_context` - (可选) 用于解密`kms_encrypted_password`的加密上下文。仅当设置了`kms_encrypted_password`时有效。
* `instance_class` - (必填) 资源的实例类型。更多信息，请参见 [实例类型](https://www.alibabacloud.com/help/en/apsaradb-for-redis/latest/instance-types)。
* `engine_version` - (可选，变更时重建) 数据库版本。默认值为`5.0`。不同Tair产品类型的参数传递规则：
  - `tair_rdb`: 内存型兼容Redis 5.0和Redis 6.0协议，传入为`5.0`或`6.0`。
  - `tair_scm`: 持久内存兼容Redis 6.0协议，传入为`1.0`。
  - `tair_essd`: 磁盘(ESSD/SSD)兼容Redis 4.0和Redis 6.0协议，分别传入为`1.0`和`2.0`。
* `zone_id` - (可选，变更时重建) 实例所在的可用区ID。
* `availability_zone` - (可选，变更时重建) 实例的可用区。
* `instance_charge_type` - (可选) 实例的计费方式。有效值为`PrePaid`和`PostPaid`。默认为`PostPaid`。
* `instance_type` - (可选，变更时重建) 实例的存储介质。有效值：`tair_rdb`, `tair_scm`, `tair_essd`。
* `vswitch_id` - (可选，变更时重建) VSwitch的ID。
* `backup_id` - (可选) 如果实例是基于另一个实例的备份创建的，则为备份集的ID。
* `vpc_auth_mode` - (可选) VPC认证模式。有效值为`Open`(启用密码认证)和`Close`(禁用密码认证并启用免密访问)。
* `maintain_start_time` - (可选) 维护窗口的开始时间。格式为`HH:mmZ`(UTC时间)。
* `maintain_end_time` - (可选) 维护窗口的结束时间。格式为`HH:mmZ`(UTC时间)。
* `cpu_type` - (可选) 资源的CPU类型。有效值：`intel`。
* `node_type` - (可选) 节点类型。有效值：
  - `MASTER_SLAVE`: 高可用(双副本)
  - `STAND_ALONE`: 单副本
  - `double`: 双副本
  - `single`: 单副本
* `architecture_type` - (可选) 实例的架构类型。有效值：`cluster`, `standard`, `rwsplit`。
* `series` - (可选，变更时重建) 实例系列。
* `tde_status` - (可选) 透明数据加密(TDE)的状态。
* `encryption_key` - (可选) 用于加密实例数据的加密密钥。
* `role_arn` - (可选) RAM角色的ARN。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - Tair实例的ID。
* `connection_domain` - 实例的内部端点。
* `private_ip` - 实例的私有IP地址。
* `vpc_auth_mode` - VPC认证模式。有效值：`Open`(启用密码认证)，`Close`(禁用密码认证并启用免密访问)。
* `maintain_start_time` - 维护窗口的开始时间。格式为`HH:mmZ`(UTC时间)。
* `maintain_end_time` - 维护窗口的结束时间。格式为`HH:mmZ`(UTC时间)。
* `node_type` - 节点类型。有效值：
  - `MASTER_SLAVE`: 高可用(双副本)
  - `STAND_ALONE`: 单副本
  - `double`: 双副本
  - `single`: 单副本
* `architecture_type` - 实例的架构类型。有效值：`cluster`, `standard`, `rwsplit`。
* `series` - 实例系列。