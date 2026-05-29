---
subcategory: "云原生数据库 PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_cluster_backup_policy"
sidebar_current: "docs-alibabacloudstack-resource-polardb-cluster-backup-policy"
description: |-
  提供 PolarDB 集群备份策略资源。
---

# alibabacloudstack_polardb_cluster_backup_policy

提供 PolarDB 集群备份策略资源。

## 示例用法

```hcl
variable "name" {
  default = "tfaccpolicy-datasource-test"
}

variable "db_type" {
  default = "MySQL"
}

variable "db_version" {
  default = "8.0"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
  lifecycle {
      ignore_changes = [
		secondary_cidr_blocks,
        tags
      ]
  }
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  vswitch_name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.1.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
  lifecycle {
      ignore_changes = [
        tags
      ]
  }
}

data "alibabacloudstack_polardb_cluster_instance_types" "default" {
  db_type = "${var.db_type}"
  db_version = "${var.db_version}"
  sorted_by = "CPU"
  cpu_type = "hygon"
  sub_category = "normal_exclusive"
}

resource "alibabacloudstack_polardb_cluster_instance" "default" {
	db_cluster_description 	=  "${var.name}"
	zone_id 				= "${data.alibabacloudstack_zones.default.zones.0.id}"
	db_type 				= "${var.db_type}"
	db_version 				= "${var.db_version}"
	storage_space 			= "20"
	vpc_id 					= "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id				= "${alibabacloudstack_vpc_vswitch.default.id}"
	db_node_class 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.id}"
	sub_category 			= "${data.alibabacloudstack_polardb_cluster_instance_types.default.instance_types.0.sub_category}"
	storage_type 			= "ESSDPL1"
}

resource "alibabacloudstack_polardb_cluster_backup_policy" "example" {
  db_cluster_id = "${alibabacloudstack_polardb_cluster_instance.default.id}"
  # 备份策略的可选参数
  data_level1_backup_period            = "Monday,Wednesday,Friday"
  data_level1_backup_time              = "02:00Z-03:00Z"
  data_level1_backup_retention_period  = 7
  log_backup_retention_period          = 7
}
```

## 参数参考

支持以下参数：

* `db_cluster_id` - (必需, ForceNew) PolarDB 集群的 ID。
* `data_level1_backup_period` - (可选) 一级数据备份周期（天）。
* `data_level1_backup_time` - (可选) 一级数据备份时间。
* `data_level1_backup_retention_period` - (可选) 一级数据备份保留周期（天）。
* `log_backup_retention_period` - (可选) 日志备份保留周期（天）。

## 属性参考

导出以下属性：

* `data_level1_backup_frequency` - 一级数据备份频率。
* `data_level1_backup_period` - 一级数据备份周期（天）。
* `data_level1_backup_time` - 一级数据备份时间。
* `data_level1_backup_retention_period` - 一级数据备份保留周期（天）。
* `data_level2_backup_retention_period` - 二级数据备份保留周期（天）。
* `backup_retention_policy_on_cluster_deletion` - 集群删除时的备份保留策略。
* `log_backup_retention_period` - 日志备份保留周期（天）。
* `preferred_backup_time` - 首选备份时间。
* `preferred_backup_period` - 首选备份周期（天）。
* `backup_retention_period` - 备份保留周期（天）。
* `backup_frequency` - 备份频率。
* `preferred_next_backup_time` - 下次首选备份时间。
* `data_level2_backup_period` - 二级数据备份周期（天）。
* `data_level2_backup_another_region_region` - 另一个区域的二级备份区域。
* `data_level2_backup_another_region_retention_period` - 另一个区域的二级备份保留周期（天）。
* `log_backup_another_region_region` - 另一个区域的日志备份区域。
* `log_backup_another_region_retention_period` - 另一个区域的日志备份保留周期（天）。
* `enable_backup_log` - 是否启用日志备份。
