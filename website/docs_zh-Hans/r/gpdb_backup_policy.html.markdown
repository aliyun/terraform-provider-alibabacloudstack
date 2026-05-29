---
subcategory: "云原生数据仓库 AnalyticDB PostgreSQL版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_backup_policy"
sidebar_current: "docs-Alibabacloudstack-resource-gpdb-backup-policy"
description: |-
  配置AnalyticDB PostgreSQL版实例的备份策略
---

# alibabacloudstack_gpdb_backup_policy

配置AnalyticDB PostgreSQL版实例的备份策略。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tftest2982"
}

data "alibabacloudstack_zones" "gpdb" {
  available_resource_creation = "Gpdb"
}

data "alibabacloudstack_gpdb_instance_types" "default" {
  engine_version = "6.0"
}

data "alibabacloudstack_gpdb_instances" "default" {
  ids = [""]
}

resource "alibabacloudstack_gpdb_instance" "default" {
  count                    = length(data.alibabacloudstack_gpdb_instances.default.ids) > 0 ? 0 : 1
  engine                   = "gpdb"
  engine_version           = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.engine_version
  instance_class           = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.id
  db_instance_mode         = data.alibabacloudstack_gpdb_instance_types.default.instance_types.0.db_instance_mode
  db_instance_storage_type = "local_ssd"
  availability_zone        = data.alibabacloudstack_zones.gpdb.zones.0.id
  description              = var.name
  seg_node_num             = "2"
  cpu_type                 = "Intel"
}

locals {
  gpdb_instance_id = length(data.alibabacloudstack_gpdb_instances.default.ids) > 0 ? data.alibabacloudstack_gpdb_instances.default.ids.0 : alibabacloudstack_gpdb_instance.default.0.id
}




resource "alibabacloudstack_gpdb_backup_policy" "default" {
  preferred_backup_time   = "02:00Z-03:00Z"
  preferred_backup_period = "Monday,Wednesday,Friday"
  backup_retention_period = "7"
  enable_recovery_point   = "false"
  db_instance_id          = local.gpdb_instance_id
}
```

## 参数说明

支持以下参数：

* `db_instance_id` - (必填, 变更时重建) 实例ID。

* `backup_retention_period` - (可选) 数据备份保留天数。默认7天，最大值定为7天，取值范围1~7。
* `enable_recovery_point` - (可选) 是否开启自动恢复点。取值说明：
  * `true`：开启。
  * `false`：关闭。
* `preferred_backup_period` - (可选) 数据备份周期，多个取值用英文逗号（,）隔开。取值说明：
  * `Monday`：周一。
  * `Tuesday`：周二。
  * `Wednesday`：周三。
  * `Thursday`：周四。
  * `Friday`：周五。
  * `Saturday`：周六。
  * `Sunday`：周日。
* `preferred_backup_time` - (可选) 数据备份时间。格式：HH:mmZ-HH:mmZ（UTC时间）。
* `recovery_point_period` - (可选) 恢复点频次。取值说明：
  * `1`：每小时。
  * `2`：每两小时。
  * `4`：每四小时。
  * `8`：每八小时。

## 属性说明

导出以下属性：

* `id` - 资源ID，即实例ID。
* `backup_retention_period` - 数据备份保留天数。
* `enable_recovery_point` - 是否开启自动恢复点。
* `preferred_backup_period` - 数据备份周期。
* `preferred_backup_time` - 数据备份时间。
* `recovery_point_period` - 恢复点频次。

## Import

GPDB 备份策略可以使用 DBInstanceId 进行导入，例如：

```
$ terraform import alibabacloudstack_gpdb_backup_policy.example pgm-xxxxxxxxx
```