---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_backup_policy"
sidebar_current: "docs-alibabacloudstack-resource-kvstore-backup-policy"
description: |-
  编排Redis或Memcache实例的备份策略
---

# alibabacloudstack_kvstore_backup_policy

使用Provider配置的凭证在指定的资源集编排Redis或Memcache实例的备份策略。

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
  instance_id   = "${alibabacloudstack_kvstore_instance.default.id}"
  backup_period = ["Tuesday", "Wednesday"]
  backup_time   = "10:00Z-11:00Z"
}

```

## 参数参考

支持以下参数：

* `instance_id` - (必填，变更时重建) ApsaraDB for Redis 或 Memcache 实例的ID。
* `backup_time` - (可选) 备份时间，格式为HH:mmZ- HH:mm Z。
* `preferred_backup_time` - (可选) 首选备份时间，格式为HH:mmZ- HH:mm Z。
* `backup_period` - (可选) 备份周期。允许的值：Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday。
* `preferred_backup_period` - (可选) 首选备份周期。允许的值：Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday。

## 属性参考

导出以下属性：

* `id` - 备份策略的ID。
* `instance_id` - ApsaraDB for Redis 或 Memcache 实例的ID。
* `backup_time` - 备份时间，格式为HH:mmZ- HH:mm Z。
* `preferred_backup_time` - 首选备份时间，格式为HH:mmZ- HH:mm Z。
* `backup_period` - 备份周期。允许的值：Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday。
* `preferred_backup_period` - 首选备份周期。允许的值：Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday。

## 导入

KVStore备份策略可以使用ID导入，例如：

```bash
$ terraform import alibabacloudstack_kvstore_backup_policy.example r-abc12345678
```