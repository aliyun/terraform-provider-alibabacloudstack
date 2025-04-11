---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbinstance"
sidebar_current: "docs-Alibabacloudstack-polardb-dbinstance"
description: |- 
  编排polardb数据库实例
---

# alibabacloudstack_polardb_dbinstance

使用Provider配置的凭证在指定的资源集编排polardb数据库实例。

## 示例用法

```hcl
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
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

variable "name" {
  default = "tf-testaccdbinstanceconfig"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_polardb_dbinstance" "default" {
  instance_storage = "5"
  instance_name = "${var.name}"
  vswitch_id = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type = "local_ssd"
  engine = "MySQL"
  engine_version = "5.7"
  instance_type = "rds.mysql.t1.small"
  db_instance_class = "polardb.mysql.x4.large"
  payment_type = "PayAsYouGo"
  maintain_time = "02:00Z-03:00Z"
  role_arn = "acs:ram::123456789012:role/aliyunpolardbdefaultrole"
}
```

## 参数参考

支持以下参数：
  * `engine` - (必填, 变更时重建) - 数据库类型。返回值：* MySQL* PostgreSQL* SQLServer* MariaDB
  * `engine_version` - (必填, 变更时重建) - 数据库版本。
  * `zone_id_slave1` - (选填) - 第一从实例所在的可用区ID。
  * `zone_id_slave2` - (选填) - 第二从实例所在的可用区ID。
  * `tde_status` - (选填) - 实例的透明数据加密(TDE)状态。
  * `enable_ssl` - (选填, 变更时重建) - 是否为实例启用SSL。
  * `storage_type` - (选填, 变更时重建) - 字段`storage_type`已被废弃，将在未来的版本中移除，请改用新字段`db_instance_storage_type`。
  * `db_instance_storage_type` - (选填, 变更时重建) - 实例储存类型，取值：* **local_ssd**、**ephemeral_ssd**：本地SSD盘。* **cloud_ssd**：SSD云盘。* **cloud_essd**：ESSD云盘。
  * `encryption_key` - (选填) - 同地域内的云盘加密的密钥ID。传入此参数表示开启云盘加密(开启后无法关闭)，并且需要传入**RoleARN**。您可以在密钥管理服务控制台查看密钥ID，也可以创建新的密钥。详情请参见[创建密钥](~~181610~~)。
  * `encryption` - (选填, 变更时重建) - 是否为实例启用加密。
  * `instance_type` - (选填) - 字段`instance_type`已被废弃，将在未来的版本中移除，请改用新字段`db_instance_class`。
  * `db_instance_class` - (选填) - 实例规格，详情请参见[实例规格表](~~26312~~)。
  * `instance_storage` - (选填) - 字段`instance_storage`已被废弃，将在未来的版本中移除，请改用新字段`db_instance_storage`。
  * `db_instance_storage` - (选填) - 数据库实例的存储容量。
  * `instance_charge_type` - (选填) - 字段`instance_charge_type`已被废弃，将在未来的版本中移除，请改用新字段`payment_type`。
  * `payment_type` - (选填) - 实例付费方式，取值：**PayAsYouGo**：按量付费。**Subscription**：包年包月。**Serverless**：Serverless付费类型，仅支持MySQL实例。
  * `period` - (选填) - 指定预付费实例为包年或者包月类型，取值：Year：包年。Month：包月。
  * `monitoring_period` - (选填) - 实例的监控周期。
  * `auto_renew` - (选填) - 实例是否自动续费，仅在创建包年包月实例时传入，取值：- **true**- **false**> * 按月购买，则自动续费周期为1个月。* 按年购买，则自动续费周期为1年。
  * `auto_renew_period` - (选填) - 实例的自动续费周期。
  * `zone_id` - (选填, 变更时重建) - 实例所在的可用区ID。
  * `vswitch_id` - (选填, 变更时重建) - 实例的VSwitch ID。
  * `instance_name` - (选填) - 字段`instance_name`已被废弃，将在未来的版本中移除，请改用新字段`db_instance_description`。
  * `db_instance_description` - (选填) - 数据库实例的描述。
  * `security_ip_mode` - (选填) - 实例的安全IP模式。
  * `maintain_time` - (选填) - 实例可维护时间段，是UTC时间，+8小时才是控制台上显示的可维护时间段。
  * `role_arn` - (选填) - 主账号授权POLARDB云服务账号访问KMS权限的全局资源描述符(ARN)。您可以通过[CheckCloudResourceAuthorized](~~446261~~)接口查看ARN信息。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `storage_type` - 实例的存储类型。有效值：* **local_ssd**, **ephemeral_ssd**: 本地SSD磁盘.* **cloud_ssd**: SSD磁盘.* **cloud_essd**: ESSD云盘。
  * `db_instance_storage_type` - 实例储存类型，取值：* **local_ssd**、**ephemeral_ssd**：本地SSD盘。* **cloud_ssd**：SSD云盘。* **cloud_essd**：ESSD云盘。
  * `instance_type` - 实例类型。
  * `db_instance_class` - 实例规格，详情请参见[实例规格表](~~26312~~)。
  * `instance_storage` - 实例的存储容量。
  * `db_instance_storage` - 数据库实例的存储容量。
  * `instance_charge_type` - 实例的计费类型。
  * `payment_type` - 实例付费方式，取值：**PayAsYouGo**：按量付费。**Subscription**：包年包月。**Serverless**：Serverless付费类型，仅支持MySQL实例。
  * `monitoring_period` - 实例的监控周期。
  * `zone_id` - 实例所在的可用区ID。
  * `instance_name` - 实例名称。
  * `db_instance_description` - 数据库实例的描述。
  * `connection_string` - 实例的连接字符串。
  * `port` - 连接端口。
  * `maintain_time` - 实例可维护时间段，是UTC时间，+8小时才是控制台上显示的可维护时间段。
  * `role_arn` - 主账号授权POLARDB云服务账号访问KMS权限的全局资源描述符(ARN)。您可以通过[CheckCloudResourceAuthorized](~~446261~~)接口查看ARN信息。