---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_backup"
sidebar_current: "docs-Alibabacloudstack-rds-backup"
description: |-
提供一个 rds 备份资源。
---

# alibabacloudstack\_rds\_backup
提供一个 rds 备份资源。

## 示例用法
```
variable "name" {
		default = "tfrds_backup77672"
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
  instance_id = "${alibabacloudstack_db_instance.default.id}"
}
```
## 参数说明
支持以下参数：

* `backup_download_url` - (可选) - 备份下载地址。
* `backup_method` - (必填) - 备份方式。
* `instance_id` - (必填) - 实例 ID。
## 属性输出
除了上述参数外，还导出以下属性：

* `backup_db_names` - 被备份的数据库名称。
* `backup_download_url` - 备份下载地址。
* `backup_id` - 备份 ID。
* `backup_intranet_download_url` - 内网环境下的备份下载地址。
* `backup_mode` - 备份模式。
* `backup_size` - 备份大小。
* `backup_type` - 备份类型。
* `end_time` - 查询结束时间，必须晚于开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `start_time` - 查询开始时间，格式为 <I> yyyy-MM-dd </I> T <I> HH:mm </I> Z (UTC 时间)。
* `status` - 备份状态。