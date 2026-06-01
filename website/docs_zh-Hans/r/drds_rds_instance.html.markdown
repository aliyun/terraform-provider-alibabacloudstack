---
subcategory: "Distributed Relational Database Service(DRDS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_rds_instance"
sidebar_current: "docs-Alibabacloudstack-resource-drds-rds-instance"
description: |-
  提供 DRDS RDS 实例资源。
---

# alibabacloudstack_drds_rds_instance

提供 DRDS RDS 实例资源。

DRDS RDS 实例是与 DRDS 实例关联的私有 RDS 实例，用于存储 DRDS 数据库的数据。

## 示例

```hcl
variable "name" {
  default = "tf-testacc-drds-rds"
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

resource "alibabacloudstack_drds_instance" "default" {
  description          = "${var.name}"
  zone_id              = "${alibabacloudstack_vpc_vswitch.default.availability_zone}"
  instance_series      = "drds.sn2.4c16g"
  instance_charge_type = "PostPaid"
  vswitch_id           = "${alibabacloudstack_vpc_vswitch.default.id}"
  specification        = "drds.sn2.4c16g.8C32G"
}

resource "alibabacloudstack_drds_rds_instance" "default" {
  storage_type        = "local_ssd"
  category            = "HighAvailability"
  db_instance_class   = "rds.mysql.s1.small"
  drds_instance_id    = "${alibabacloudstack_drds_instance.default.id}"
  zone_id             = "${data.alibabacloudstack_zones.default.zones.0.id}"
  db_instance_storage = "20"
}
```

## 参数说明

支持以下参数：

* `storage_type` - （必填，ForceNew）RDS 实例的存储类型。有效值：`local_ssd`、`cloud_ssd`、`cloud_essd`。修改此参数会强制重新创建资源。
* `category` - （必填，ForceNew）RDS 实例的系列。有效值：`HighAvailability`（高可用版）、`Finance`（三节点企业版）。修改此参数会强制重新创建资源。
* `drds_instance_id` - （必填，ForceNew）RDS 实例所属的 DRDS 实例 ID。修改此参数会强制重新创建资源。
* `zone_id` - （必填，ForceNew）RDS 实例所在的可用区 ID。修改此参数会强制重新创建资源。
* `db_instance_class` - （必填）RDS 实例的规格（实例等级）。示例值：`rds.mysql.s1.small`。
* `db_instance_storage` - （必填）RDS 实例的存储容量，单位：GB。
* `force_remove` - （可选）销毁时是否强制删除 RDS 实例。有效值：`true`、`false`。默认值为 `false`。

### 超时

`timeouts` 块允许您为某些操作指定[超时](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `create` - （默认 20 分钟）创建 DRDS RDS 实例时使用（直到达到运行状态）。
* `update` - （默认 20 分钟）更新 DRDS RDS 实例规格时使用。
* `delete` - （默认 20 分钟）删除 DRDS RDS 实例时使用。

## 属性说明

除上述参数外，还导出以下属性：

* `id` - 资源的 ID。格式为 `<drds_instance_id>:<rds_instance_id>`。
* `rds_instance_id` - RDS 实例的 ID。
* `create_time` - RDS 实例的创建时间。
* `status` - RDS 实例的状态。有效值：`Creating`（创建中）、`Running`（运行中）、`Changing Specifications`（变更规格中）、`Deleting`（删除中）、`Stopping`（停止中）、`Stopped`（已停止）。

## Import

DRDS RDS 实例可以使用组合 ID 进行导入，格式为 `<drds_instance_id>:<rds_instance_id>`，例如：

```bash
$ terraform import alibabacloudstack_drds_rds_instance.example drds-abc123456:rm-xyz789012
```
