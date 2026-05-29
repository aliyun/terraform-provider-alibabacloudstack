---
subcategory: "云数据库 RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-rds-backups"
description: |-
提供一个由 AlibabacloudStack 账户拥有的 rds 备份列表。
---

# alibabacloudstack\_rds\_backups
该数据源根据指定的过滤条件，提供一个 Alibabacrdstack 账户下的 rds 备份列表。

## 示例用法
```
variable "name" {
	default = "tf-testAlibabacloudstackrdsBackups21486"
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


resource "alibabacloudstack_db_instance" "default" {
	engine               = "MySQL"
	engine_version       = "5.6"
	instance_type        = "rds.mysql.s2.large"
	instance_storage     = "20"
	instance_name        = "${var.name}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	storage_type         = "local_ssd"
  }

resource "alibabacloudstack_rds_backup" "default" {
	backup_method = "Physical"
	instance_id = alibabacloudstack_db_instance.default.id
}

data "alibabacloudstack_rds_backups" "default" {
   backup_ids = ["${alibabacloudstack_rds_backup.default.id}"]
   instance_id = "${alibabacloudstack_rds_backup.default.instance_id}"
   start_time = "${alibabacloudstack_rds_backup.default.start_time}"
   end_time = "${alibabacloudstack_rds_backup.default.end_time}"
}
```
## 参数说明
支持以下参数：

* `backup_ids` - (可选) - 备份的 ID 列表。
* `start_time` - (可选) - 查询开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `backup_id` - (可选) - 备份 ID。
* `end_time` - (可选) - 查询结束时间，必须晚于开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `instance_id` - (必填) - 实例 ID。
属性输出
除了上述参数外，还导出以下属性：

* `backups` - 备份列表
* `id` - 备份 ID
* `backup_db_names` - 被备份的数据库名称
* `backup_download_url` - 备份文件的下载地址
* `backup_id` - 备份 ID
* `backup_intranet_download_url` - 内网环境下的备份文件下载地址
* `backup_method` - 备份方式
* `backup_mode` - 备份模式
* `backup_size` - 备份文件大小
* `backup_type` - 备份类型
* `instance_id` - 实例 ID
* `end_time` - 查询结束时间，必须晚于开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `start_time` - 查询开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `status` - 备份状态