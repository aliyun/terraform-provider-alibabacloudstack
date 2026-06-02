---
subcategory: "云原生分布式数据库PolarDB-X 2.0"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup"
sidebar_current: "docs-Alibabacloudstack-resource-polardbx-backup"
description: |-
  提供一个 PolarDB-X 备份资源。
---

# alibabacloudstack\_polardbx\_backup

提供一个 PolarDB-X 备份资源。该资源允许您创建和管理 PolarDB-X 实例的备份。

-> **注意：** 该资源不支持更新操作。修改任何参数都会强制创建新资源。

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

  * `instance_id` - (必填，变更后重建) PolarDB-X 实例的 ID。修改此参数会强制重新创建资源。
  * `backup_type` - (可选，变更后重建) 备份类型。有效值：`0`（物理备份）。默认值为 `0`。修改此参数会强制重新创建资源。

## 属性参考

除了上述参数外，还导出以下属性：

  * `id` - 备份的 ID。格式为 `<instance_id>:<backup_set_id>`。
  * `backup_model` - 备份模式。有效值：`0`（物理备份），`1`（逻辑备份）。
  * `backup_set_size` - 备份集的大小，单位为字节。
  * `backup_set_id` - 备份集 ID。
  * `status` - 备份的状态。有效值：`0`（创建中），`1`（成功），`2`（失败）。
  * `end_time` - 备份结束时间，UTC 格式。
  * `begin_time` - 备份开始时间，UTC 格式。

## Import

PolarDB-X 备份可以使用实例 ID 和备份集 ID（用冒号分隔）导入，例如：

```
$ terraform import alibabacloudstack_polardbx_backup.example pc-xxxxxxxxxxxxx:1234567890
```