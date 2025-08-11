---
subcategory: "PolarDBX"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardbx_backup_policys"
sidebar_current: "docs-Alibabacloudstack-datasource-polardbx-backup-policys"
description: |-
  提供阿里云栈账户拥有的PolarDBX备份策略列表。
---

# alibabacloudstack\_polardbx\_backup\_policys

此数据源根据指定的过滤条件提供阿里云栈账户中的PolarDBX备份策略列表。

## 使用示例

```hcl
variable "name" {
  default = "tf-testAccPolardbxInstancesDataSource-3625795"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details              = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = alibabacloudstack_vpc_vpc.default.id
  cidr_block = "172.16.1.0/24"
  zone_id    = data.alibabacloudstack_zones.default.zones.0.id
}

resource "alibabacloudstack_polardbx_instance" "default" {
  description      = "testtf1111"
  series           = "enterprise"
  topology_type    = "1azone"
  zone_id          = data.alibabacloudstack_zones.default.zones.0.id
  engine_version   = "5.7"
  storage          = "50"
  network_type     = "vpc"
  vpc_id           = alibabacloudstack_vpc_vpc.default.id
  vswitch_id       = alibabacloudstack_vpc_vswitch.default.id
  cn_node_class    = "polarx.x4.medium.2e"
  cn_node_count    = "2"
  dn_node_class    = "mysql.n4.medium.25"
  dn_node_count    = "2"
}

data "alibabacloudstack_polardbx_backup_policys" "default" {
  db_instance_id = alibabacloudstack_polardbx_instance.default.id
}
```

## 参数参考
支持以下参数：

* `db_instance_id` - (必填, ForceNew) PolarDBX实例的ID。

## 属性参考
除了上述参数外，还导出以下属性：

* `backup_period` - 备份周期。您可以选择一周中的多天。有效值：Monday, Tuesday, Wednesday, Thursday, Friday, Saturday, Sunday。多个值用逗号分隔。
* `backup_set_retention` - 备份集的保留周期。有效值：7到730天。
* `backup_plan_begin` - 备份开始时间。时间格式为HH:mmZ，例如03:00Z。
* `remove_log_retention` - 事务日志的保留周期。有效值：7到730天。
* `cold_data_backup_interval` - 冷数据备份间隔。
* `local_log_retention_number` - 本地保留的事务日志数量。有效值：6到100。
* `cold_data_backup_retention` - 冷数据备份保留周期。
* `force_clean_on_high_space_usage` - 指定当存储空间使用率高时是否强制删除备份。有效值：0（否），1（是）。
* `backup_way` - 备份方式。
* `local_log_retention` - 本地保留的事务日志保留周期。单位：天。
* `backup_type` - 备份类型。
* `log_local_retention_space` - 本地保留的事务日志占用的最大空间。单位：%。有效值：0到100。