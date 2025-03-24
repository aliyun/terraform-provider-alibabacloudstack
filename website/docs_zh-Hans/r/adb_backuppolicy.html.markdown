---
subcategory: "ADB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_adb_backuppolicy"
sidebar_current: "docs-Alibabacloudstack-adb-backuppolicy"
description: |- 
  编排adb备份规则
---

# alibabacloudstack_adb_backup_policy
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_adb_backuppolicy`

使用Provider配置的凭证在指定的资源集下编排adb备份规则。

## 示例用法

```hcl
variable "name" {
  default = "adbClusterconfig"
}

variable "creation" {
  default = "ADB"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = var.creation
}

resource "alibabacloudstack_vpc" "default" {
  name       = var.name
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  zone_id           = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_name      = var.name
}

resource "alibabacloudstack_adb_db_cluster" "default" {
  db_cluster_version  = "3.0"
  db_cluster_category = "Basic"
  db_node_class       = "C8"
  db_node_count       = 2
  db_node_storage     = 200
  mode                = "reserver"
  pay_type            = "PostPaid"
  description         = var.name
  vswitch_id          = alibabacloudstack_vswitch.default.id
  cluster_type        = "analyticdb"
  cpu_type            = "intel"
}

resource "alibabacloudstack_adb_backup_policy" "policy" {
  db_cluster_id           = alibabacloudstack_adb_db_cluster.default.id
  preferred_backup_period = ["Tuesday", "Thursday", "Saturday"]
  preferred_backup_time   = "10:00Z-11:00Z"
}
```

### 从配置中移除 `alibabacloudstack_adb_backup_policy`

`alibabacloudstack_adb_backup_policy` 资源允许您管理 ADB 集群的备份策略，但 Terraform 无法销毁它。从配置中移除此资源将从状态文件和管理中移除它，但不会销毁集群策略。您可以继续通过 ADB 控制台管理集群。

## 参数参考

支持以下参数：

* `db_cluster_id` - (必填，变更时重建) 需要配置备份策略的 ADB 集群的 ID。
* `preferred_backup_period` - (必填) 应进行 ADB 集群备份的天数。有效值包括：`Monday`、`Tuesday`、`Wednesday`、`Thursday`、`Friday`、`Saturday`、`Sunday`。
* `preferred_backup_time` - (必填) 应进行 ADB 集群备份的时间窗口，格式为 `HH:mmZ-HH:mmZ`。开始和结束时间之间的间隔为一小时。请注意，指定的时间是以 UTC 为准。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `id` - 当前备份策略资源 ID。它与 `db_cluster_id` 相同。
* `backup_retention_period` - 数据备份文件保留的天数。此值固定为 7 天，无法修改。
