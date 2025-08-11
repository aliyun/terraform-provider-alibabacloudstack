---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup"
sidebar_current: "docs-Alibabacloudstack-polardbx-backup"
description: |-
  提供一个polardbx备份资源。
---

# alibabacloudstack\_polardbx\_backup

提供一个polardbx备份资源。

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

resource "alibabacloudstack_polardbx_backup" "default" {
  instance_id = "${alibabacloudstack_polardbx_instance.default.id}"
}
```

## 参数参考

支持以下参数：
  * `instance_id` - (必填) - PolarDBX实例的ID。
  * `backup_type` - (可选) - 备份类型。目前仅支持"0"。

## 属性参考

除了上述参数外，还导出以下属性：
  * `backup_mode` - 备份模式。目前仅支持"0"。
  * `backup_set_size` - 备份的大小。
  * `backup_set_id` - 备份集ID。
  * `status` - 备份的状态。
  * `end_time` - 此备份的结束时间（UTC时间）。
  * `begin_time` - 备份开始时间（UTC时间）。