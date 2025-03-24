---
subcategory: "MongoDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_instance"
sidebar_current: "docs-Alibabacloudstack-mongodb-instance"
description: |- 
  集编排mongodb实例
---

# alibabacloudstack_mongodb_instance

使用Provider配置的凭证在指定的资源集编排mongodb实例。

## 示例用法

### 基础用法：

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "MongoDB"
}

variable "password" {}

resource "alibabacloudstack_mongodb_instance" "example" {
  engine_version      = "3.4"
  db_instance_class   = "dds.mongo.mid"
  db_instance_storage = 10
  zone_id             = data.alibabacloudstack_zones.default.zones[0].id
  backup_period       = ["Monday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"]
  preferred_backup_time = "20:00Z-21:00Z"
  name                = "testMongoDB"
  security_ip_list    = ["10.168.1.12", "100.69.7.112"] 
  ssl_action          = "Open"
  tde_status          = "Enabled"
  replication_factor  = 3
  storage_engine      = "WiredTiger"
  instance_charge_type= "PostPaid"
  vswitch_id          = "vsw-abc123"
  security_group_id   = "sg-abc123"
  account_password    = var.password
}
```

## 参数参考

支持以下参数：

* `engine_version` - (必填，变更时重建) 实例的数据库引擎版本。有效值包括：`3.4`, `4.0` 等。
* `db_instance_class` - (必填) 实例类型。例如：`dds.mongo.s.small`, `dds.mongo.mid`。
* `db_instance_storage` - (必填) 实例的存储容量。有效值：10 到 3000。值必须是 10 的倍数。单位：GB。
* `zone_id` - (选填，变更时重建) 实例所在的可用区 ID。如果不指定，系统将默认选择一个。
* `backup_period` - (选填)MongoDB 实例的备份周期。当设置了 `preferred_backup_time` 时为必填项。有效值：`[Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]`。默认值：`[Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday]`。
* `preferred_backup_time` - (选填)MongoDB 实例的备份时间窗口。格式为 `HH:mmZ-HH:mmZ`。时间设置间隔为一小时。如果未设置，默认返回类似 `23:00Z-24:00Z`。
* `name` - (选填)DB 实例的名称。长度为 2 到 256 个字符的字符串。
* `security_ip_list` - (选填)允许访问 MongoDB 实例的 IP 地址列表。每个 IP 地址最多可以有 256 个字符。默认为空列表。
* `ssl_action` - (选填)在 SSL 功能上执行的操作。有效值：`Open`：打开 SSL 加密；`Close`：关闭 SSL 加密；`Update`：更新 SSL 证书。
* `tde_status` - (选填，变更时重建) 透明数据加密(TDE)状态。有效值：`Enabled`, `Disabled`。
* `replication_factor` - (选填)副本集节点数量。有效值：`1`, `3`, `5`, `7`。默认值：`3`。
* `storage_engine` - (选填，变更时重建) 实例的存储引擎。有效值：`WiredTiger`, `RocksDB`。系统默认值：`WiredTiger`。
* `instance_charge_type` - (选填)实例的计费类型。有效值：`PrePaid`, `PostPaid`。默认值：`PostPaid`。
* `period` - (选填)购买 DB 实例的时长(以月为单位)。当 `instance_charge_type` 为 `PrePaid` 时有效。有效值：`[1~9], 12, 24, 36`。默认值：`1`。
* `vswitch_id` - (选填，变更时重建) 启动 DB 实例所在 VPC 的虚拟交换机 ID。
* `security_group_id` - (选填)ECS 的安全组 ID。一个实例最多可以绑定 10 个 ECS 安全组。
* `account_password` - (选填)root 账户的密码。它是一个由字母、数字和下划线组成的 6 到 32 个字符的字符串。
* `kms_encrypted_password` - (选填)用于创建或更新实例的 KMS 加密密码。如果提供了 `account_password`，则此字段将被忽略。
* `kms_encryption_context` - (选填)用于在创建或更新实例之前解密 `kms_encrypted_password` 的 KMS 加密上下文。
* `maintain_start_time` - (选填)维护窗口的开始时间。指定 UTC 时间格式为 `HH:mmZ`。
* `maintain_end_time` - (选填)维护窗口的结束时间。指定 UTC 时间格式为 `HH:mmZ`。
* `tags` - (选填，映射)要分配给资源的标签映射。

## 属性参考

除了上述参数外，还导出以下属性：

* `retention_period` - 实例日志备份保留天数。
* `replica_set_name` - mongo 副本集名称。
* `maintain_start_time` - 维护窗口的开始时间。
* `maintain_end_time` - 维护窗口的结束时间。
* `ssl_status` - SSL 功能的状态。`Open`：SSL 已打开；`Closed`：SSL 已关闭。
* `zone_id` - 实例所在的可用区 ID。
* `vswitch_id` - 启动 DB 实例所在 VPC 的虚拟交换机 ID。
* `security_ip_list` - 允许访问 MongoDB 实例的 IP 地址列表。
* `security_group_id` - ECS 的安全组 ID。
* `backup_period` - MongoDB 实例的备份周期。