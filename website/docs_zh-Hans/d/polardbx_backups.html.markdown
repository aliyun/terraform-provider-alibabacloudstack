---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "阿里云栈：alibabacloudstack_polardbx_backups"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-backups"
description: |-
  提供阿里云栈账户拥有的polardbx备份列表。
---

# alibabacloudstack\_polardbx\_backups

此数据源根据指定的过滤条件提供阿里云栈账户中的polardbx备份列表。

## 使用示例
```
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
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

resource "alibabacloudstack_polardbx_instance" "default" {
  description = "testtf1111"
	series = "enterprise"
	topology_type = "1azone"
	zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
	engine_version = "5.7"
	storage = "50"
	network_type = "vpc"
	vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
	vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
	cn_node_class = "polarx.x4.medium.2e"
	cn_node_count = "2"
	dn_node_class = "mysql.n4.medium.25"
	dn_node_count = "2"
}

data "alibabacloudstack_polardbx_backups" "default" {
  db_instance_id = "${alibabacloudstack_polardbx_instance.default.id}"
}
```

## 参数参考

支持以下参数：
  * `ids` - (可选) - 备份ID列表。
  * `db_instance_id` - (必填, ForceNew) - 要查询备份的PolarDBX集群的ID。
  * `end_time` - (可选) - 此备份的结束时间（UTC时间）。
  * `start_time` - (可选) - 备份开始时间（UTC时间）。

## 属性参考

除了上述参数外，还导出以下属性：
  * `backups` - PolarDBX备份列表。
    * `id` - 备份的ID。
    * `backup_method` - 数据备份方法。仅支持快照备份。值固定为**Snapshot**。
    * `backup_model` - 备份模式
    * `backup_set_size` - 备份大小。
    * `backup_type` - 备份类型。
    * `backup_set_id` - 备份集ID。
    * `status` - 备份的状态。
    * `end_time` - 此备份的结束时间（UTC时间）。
    * `begin_time` - 备份开始时间（UTC时间）。