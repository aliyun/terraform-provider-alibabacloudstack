---
subcategory: "MongoDB" 
layout: "alibabacloudstack" 
page_title: "Alibabacloudstack: alibabacloudstack_mongodb_backup" 
sidebar_current: "docs-Alibabacloudstack-mongodb-backup" 
description: |- 
提供一个 MongoDB 备份资源。
---
# alibabacloudstack\_mongodb\_backup
提供一个 MongoDB 备份资源。

## 示例用法
```
hcl
variable "name" {
		default = "tfmongodb_backup24637"
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
  db_instance_id = "${alibabacloudstack_mongodb_instance.default.id}"
}
```
## 参数说明
支持以下参数：

* `backup_method` - (必填) - 备份方法。
* `db_instance_id` - (必填, 强制新建) - 数据库实例 ID。
## 导出属性
除了上述列出的参数外，还导出了以下属性：

* `backup_db_names` - 备份数据库名称。
* `backup_download_url` - 备份下载 URL。
* `backup_id` - 备份 ID。
* `backup_intranet_download_url` - 备份内网下载 URL。
* `backup_mode` - 备份模式。
* `backup_size` - 备份大小。
* `backup_type` - 备份类型。
* `end_time` - 备份结束时间。
* `start_time` - 备份开始时间。
* `status` - 资源状态。