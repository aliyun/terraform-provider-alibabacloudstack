---
subcategory: "云数据库 MongoDB 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-mongodb-backups"
description: |-
  提供由 Alibabacloudstack 账户拥有的 MongoDB 备份列表。
---

# alibabacloudstack\_mongodb\_backups

该数据源根据指定的过滤条件，提供一个 Alibabacloudstack 账户中的 MongoDB 备份列表。

## 示例用法

```hcl
variable "name" {
	default = "tf-testAlibabacloudstackMongodbBackups39113"
}

data "alibabacloudstack_zones" default {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_mongodb_instance" "default" {
	vswitch_id          = alibabacloudstack_vpc_vswitch.default.id
	engine_version      = "3.0"
	db_instance_class   = "dds.mongo.mid"
	db_instance_storage = "10"
	name                = "${var.name}"
	storage_engine      = "WiredTiger"
	instance_charge_type = "PostPaid"
	replication_factor = "3"
}

resource "alibabacloudstack_mongodb_backup" "default" {
	backup_method = "Physical"
	db_instance_id = alibabacloudstack_mongodb_instance.default.id
}

data "alibabacloudstack_mongodb_backups" "default" {
	ids = ["${alibabacloudstack_mongodb_backup.default.id}"]
	db_instance_id = "${alibabacloudstack_mongodb_backup.default.db_instance_id}"
	start_time = "${alibabacloudstack_mongodb_backup.default.start_time}"
	end_time = "${alibabacloudstack_mongodb_backup.default.end_time}"
}
```

## 参数说明

支持以下参数：

* `db_instance_id` - (必填) MongoDB 实例 ID。
* `start_time` - (必填) 备份开始时间，格式为 `YYYY-MM-DDTHH:MMZ`。
* `end_time` - (必填) 备份结束时间，格式为 `YYYY-MM-DDTHH:MMZ`。
* `backup_id` - (可选) 备份 ID。
* `ids` - (可选) 备份 ID 列表。

## 导出属性

除了上述列出的参数外，还导出了以下属性：

* `ids` - 备份 ID 列表。
* `backups` - 备份列表，包含以下字段：
  * `id` - 备份 ID，格式为 `<backup_id>&<db_instance_id>&<start_time>`。
  * `backup_id` - 备份 ID。
  * `backup_db_names` - 备份数据库名称。
  * `backup_download_url` - 备份下载 URL。
  * `backup_intranet_download_url` - 备份内网下载 URL。
  * `backup_method` - 备份方法。
  * `backup_mode` - 备份模式。
  * `backup_size` - 备份大小。
  * `backup_type` - 备份类型。
  * `backup_status` - 备份状态。
  * `db_instance_id` - MongoDB 实例 ID。
  * `start_time` - 备份开始时间。
  * `end_time` - 备份结束时间。