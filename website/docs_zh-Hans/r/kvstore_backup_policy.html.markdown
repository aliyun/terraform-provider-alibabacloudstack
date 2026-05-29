---
subcategory: "云数据库 Redis 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_backup_policy"
sidebar_current: "docs-Alibabacloudstack-resource-kvstore-backup-policy"
description: |-
  编排Redis或Memcache实例的备份策略
---

# alibabacloudstack_kvstore_backup_policy

使用Provider配置的凭证在指定的资源集编排Redis或Memcache实例的备份策略。

> **注意：** 当前资源也可通过以下别名引用：
> - `apsarastack_kvstore_backup_policy`

## 示例用法

### 基础用法

```
variable "creation" {
  default = "KVStore"
}
variable "multi_az" {
  default = "false"
}
variable "name" {
  default = "kvstorebackuppolicyvpc"
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
resource "alibabacloudstack_kvstore_instance" "default" {
  instance_class = "Memcache"
  instance_name  = "${var.name}"
  vswitch_id     = "${alibabacloudstack_vswitch.default.id}"
  private_ip     = "172.16.0.10"
  security_ips   = ["10.0.0.1"]
  instance_type  = "memcache.master.small.default"
  
}
resource "alibabacloudstack_kvstore_backup_policy" "default" {
  instance_id             = "${alibabacloudstack_kvstore_instance.default.id}"
  preferred_backup_period = ["Tuesday", "Wednesday"]
  preferred_backup_time   = "10:00Z-11:00Z"
}

```

## 参数说明

支持以下参数：

* `instance_id` - (必填，ForceNew) ApsaraDB for Redis 或 Memcache 实例的ID。修改此参数会强制重新创建资源。
* `preferred_backup_time` - (可选，Computed) 首选备份时间，格式为HH:mmZ-HH:mmZ。例如：`02:00Z-03:00Z` 表示每天凌晨2点到3点之间进行备份。此属性由 API 返回，无法手动设置。
* `preferred_backup_period` - (可选，Computed) 首选备份周期。允许的值为：`Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`。例如：`["Monday", "Wednesday"]` 表示每周一和周三进行备份。此属性由 API 返回，无法手动设置。
* `backup_time` - (可选，已弃用，Computed) **此参数已弃用**，将在未来版本中移除。请使用 `preferred_backup_time` 替代。与 `preferred_backup_time` 参数互斥。
* `backup_period` - (可选，已弃用，Computed) **此参数已弃用**，将在未来版本中移除。请使用 `preferred_backup_period` 替代。与 `preferred_backup_period` 参数互斥。

## 属性说明

导出以下属性：

* `id` - 备份策略的唯一标识符，与实例ID相同。
* `instance_id` - ApsaraDB for Redis 或 Memcache 实例的ID。
* `preferred_backup_time` - 当前配置的首选备份时间，格式为HH:mmZ-HH:mmZ。
* `preferred_backup_period` - 当前配置的首选备份周期。允许的值为：`Monday`, `Tuesday`, `Wednesday`, `Thursday`, `Friday`, `Saturday`, `Sunday`。
* `backup_time` - 当前配置的备份时间（已弃用参数）。
* `backup_period` - 当前配置的备份周期（已弃用参数）。

## 导入

KVStore备份策略可以使用ID导入，例如：

```bash
$ terraform import alibabacloudstack_kvstore_backup_policy.example r-abc12345678
```