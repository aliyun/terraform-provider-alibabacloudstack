---
subcategory: "GPDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_gpdb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-gpdb-dbinstance"
description: |- 
  编排图数据库实例
---

# alibabacloudstack_gpdb_dbinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_gpdb_instance`

使用Provider配置的凭证在指定的资源集下编排图数据库实例。

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
  engine                  = "gpdb"
  engine_version          = "4.3"
  instance_class          = "gpdb.group.segsdx2"
  instance_group_count    = "2"
  description            = "Terraform Test GPDB Instance"
  db_instance_description = "Terraform Test GPDB Instance Description"
  availability_zone      = data.alibabacloudstack_zones.default.zones[0].id
  vswitch_id             = alibabacloudstack_vswitch.default.id
  security_ip_list       = ["10.168.1.12", "100.69.7.112"]
}
```

## 参数说明

支持以下参数：

* `availability_zone` - (选填, 变更时重建) - DB实例将被创建的可用区。如果不指定，Terraform将使用`alibabacloudstack_zones`数据源中的第一个可用区。
* `instance_class` - (必填, 变更时重建) - GPDB实例的规格。例如，`gpdb.group.segsdx2`。更多详情请参阅 [实例规格](https://www.alibabacloud.com/help/doc-detail/86942.htm)。
* `db_instance_class` - (已废弃, 变更时重建) - 已废弃，请改用 `instance_class`。
* `instance_group_count` - (必填) - GPDB实例中的组数。有效值为 `[2, 4, 8, 16, 32]`。
* `instance_charge_type` - (选填, 变更时重建) - 实例的计费方式。有效值为 `PrePaid`(包年包月) 和 `PostPaid`(按量付费)。默认值为 `PostPaid`。
* `payment_type` - (选填, 变更时重建) - `instance_charge_type` 的别名。指定实例的支付类型。
* `description` - (选填) - GPDB实例的简要描述。最多可以是256个字符。
* `db_instance_description` - (选填) - `description` 的别名。提供GPDB实例的详细描述。
* `vswitch_id` - (选填, 变更时重建) - 将启动GPDB实例的交换机ID。如果要在VPC中启动实例，则此参数是必填的。
* `instance_inner_connection` - (选填, 变更时重建) - GPDB实例的内部连接端点。
* `instance_inner_port` - (选填, 变更时重建) - GPDB实例的内部端口。
* `port` - (选填, 变更时重建) - GPDB实例的公网端口号。默认值为 `5432`。
* `engine` - (必填, 变更时重建) - 数据库引擎。目前仅支持 `gpdb`。
* `engine_version` - (必填, 变更时重建) - GPDB引擎的版本。有效值包括 `4.3`、`6.0` 和 `7.0`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `availability_zone` - GPDB实例所在的可用区。
* `instance_class` - GPDB实例的规格。
* `db_instance_class` - `instance_class` 的别名。
* `instance_id` - GPDB实例的唯一标识符。
* `db_instance_id` - `instance_id` 的别名。
* `region_id` - GPDB实例所在的区域。
* `status` - GPDB实例的当前状态。
* `instance_network_type` - GPDB实例的网络类型。可能的值包括 `VPC`。
* `network_type` - `instance_network_type` 的别名。
* `instance_charge_type` - GPDB实例的计费方式。
* `payment_type` - `instance_charge_type` 的别名。
* `description` - GPDB实例的描述。
* `db_instance_description` - `description` 的别名。
* `vswitch_id` - 与GPDB实例关联的交换机ID。
* `instance_inner_connection` - GPDB实例的内部连接端点。
* `instance_inner_port` - GPDB实例的内部端口。
* `port` - GPDB实例的公网端口号。
* `instance_vpc_id` - 与GPDB实例关联的VPC ID。
* `vpc_id` - `instance_vpc_id` 的别名。
* `engine` - GPDB实例使用的数据库引擎。
* `engine_version` - GPDB实例使用的数据库引擎版本。