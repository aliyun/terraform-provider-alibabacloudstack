---
subcategory: "云原生数据仓库 AnalyticDB PostgreSQL版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-resource-gpdb-dbinstance"
description: |- 
  编排云原生数据仓库 AnalyticDB PostgreSQL版实例
---

# alibabacloudstack_gpdb_dbinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_gpdb_instance`

使用Provider配置的凭证在指定的资源集下编排云原生数据仓库 AnalyticDB PostgreSQL版实例。

## 示例用法

### 创建带有VPC和安全IP列表的GPDB实例

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "Gpdb"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  vpc_id           = alibabacloudstack_vpc.default.id
  cidr_block       = "172.16.0.0/24"
  name             = "vswitch-123456"
}

resource "alibabacloudstack_gpdb_dbinstance" "example" {
  db_instance_mode          = "Classic"
  db_instance_class         = "gpdb.group.segsdx2"
  instance_group_count      = "2"
  db_instance_description   = "Terraform Test GPDB Instance"
  availability_zone         = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_id                = alibabacloudstack_vswitch.default.id
  security_ip_list          = ["10.168.1.12", "100.69.7.112"]
  engine                    = "gpdb"
  engine_version            = "4.3"
}
```

## 参数说明

支持以下参数：

* `db_instance_mode` - (必填, 变更时重建) GPDB实例的部署模式。有效值为 `Classic` 和 `StorageReserver`。当设置为 `StorageReserver` 时，`db_instance_storage_type` 和 `seg_node_num` 变为必填参数。
* `db_instance_class` - (选填, 变更时重建) GPDB实例的规格。例如 `gpdb.group.segsdx2`。`db_instance_class` 或已废弃的 `instance_class` 必须指定其一。
* `instance_class` - (选填, 变更时重建, 已废弃) 已废弃，请改用 `db_instance_class`。
* `instance_group_count` - (选填) GPDB实例的节点组数量。有效值为 `2`、`4`、`8`、`16`、`32`。
* `instance_charge_type` - (选填, 变更时重建, 已废弃) 已废弃，请改用 `payment_type`。实例的计费方式。有效值：`PostPaid`（按量付费）。
* `payment_type` - (选填, 变更时重建) `instance_charge_type` 的别名。实例的计费方式。有效值：`PostPaid`（按量付费）。
* `db_instance_description` - (选填) GPDB实例的描述。长度为2~256个字符。
* `description` - (选填, 已废弃) 已废弃，请改用 `db_instance_description`。
* `vswitch_id` - (选填, 变更时重建) 交换机ID。如果指定了该参数，实例将创建在该交换机所在的VPC中。
* `instance_inner_connection` - (选填, 变更时重建) 实例的内网连接前缀。
* `engine` - (选填, 变更时重建) 数据库引擎。有效值：`gpdb`。
* `engine_version` - (选填, 变更时重建) 引擎版本。有效值包括 `4.3`、`6.0`、`7.0`。
* `availability_zone` - (选填, 变更时重建) 实例的可用区ID。如果不指定，系统将自动分配。
* `cpu_type` - (选填, 变更时重建) 实例的CPU架构类型。
* `security_ip_list` - (选填) 允许访问实例的IP地址列表。
* `db_instance_storage_type` - (选填, 变更时重建) 存储类型。当 `db_instance_mode` 为 `StorageReserver` 时为必填。
* `seg_node_num` - (选填) 计算节点数量。当 `db_instance_mode` 为 `StorageReserver` 时为必填。
* `instance_pay_type` - (选填) 实例的支付类型。
* `network_type` - (选填, 已废弃) 实例的网络类型。该字段将在 3.21.0 版本中废弃。它将根据是否配置了 `vswitch_id` 自动确定。
* `tags` - (选填) 要绑定到资源的标签映射。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源ID，格式为 `instance_id`。
* `db_instance_id` - GPDB实例的唯一标识符。
* `instance_id` - (已废弃) `db_instance_id` 的别名。
* `region_id` - 实例所在的区域ID。
* `status` - 实例的状态。
* `instance_network_type` - 实例的网络类型。
* `network_type` - (已废弃) `instance_network_type` 的别名。
* `instance_charge_type` - (已废弃) 计费方式。
* `payment_type` - 实例的计费方式。
* `description` - (已废弃) 实例描述。
* `db_instance_description` - 实例描述。
* `vswitch_id` - 交换机ID。
* `instance_inner_connection` - 内网连接端点。
* `instance_inner_port` - (已废弃) 内网端口。请改用 `port`。
* `port` - 实例的连接端口。
* `vpc_id` - 实例的VPC ID。
* `instance_vpc_id` - (已废弃) `vpc_id` 的别名。
* `engine` - 数据库引擎。
* `engine_version` - 引擎版本。
* `instance_class` - 实例规格。
* `db_instance_class` - 实例规格。
* `availability_zone` - 可用区。
* `instance_group_count` - 节点组数量。
* `port` - 连接端口。

## Import

GPDB实例可以使用实例ID导入，例如：

```
$ terraform import alibabacloudstack_gpdb_dbinstance.example gp-gs5xxxxxxxxxx
```