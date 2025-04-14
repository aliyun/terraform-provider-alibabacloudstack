---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_rds_dbinstance"
sidebar_current: "docs-Alibabacloudstack-rds-dbinstance"
description: |- 
  编排RDS实例
---

# alibabacloudstack_rds_dbinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_db_instance`

使用Provider配置的凭证在指定的资源集编排RDS实例。
## 示例用法

### 创建一个 RDS MySQL 实例

```hcl
variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default" {
  instance_name           = "${var.name}"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.s2.large"
  db_instance_storage     = 20
  maintain_time           = "03:00Z-04:00Z"
  security_ip_mode        = "normal"
  role_arn                = "acs:ram::123456789012:role/example-role"
}
```

### 创建一个具有特定参数的 RDS MySQL 实例

```hcl
resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "vpc-123456"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "vpc-123456"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default1" {
  instance_name           = "tf-testAccDBInstanceConfig1"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.t1.small"
  db_instance_storage     = 10
  encryption_key          = "f23ed1c9-b91f-..."
  tde_status              = false
  enable_ssl             = false
  zone_id_slave1         = "${data.alibabacloudstack_zones.default.zones.0.id}"
  zone_id                = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_rds_dbinstance" "default2" {
  instance_name           = "tf-testAccDBInstanceConfig2"
  vswitch_id              = "${alibabacloudstack_vpc_vswitch.default.id}"
  storage_type            = "local_ssd"
  engine                  = "MySQL"
  engine_version          = "5.6"
  db_instance_class       = "rds.mysql.t1.small"
  db_instance_storage     = 10
}
```

## 参数说明

支持以下参数：

* `engine` - (必填，变更时重建) 数据库类型。返回值：
  * **MySQL**
  * **PostgreSQL**
  * **SQLServer**
  * **MariaDB**

* `engine_version` - (必填，变更时重建) 数据库版本。例如，对于MySQL，有效版本包括 `5.6`, `5.7`, 和 `8.0`。

* `zone_id_slave1` - (可选) 第一个备用实例的可用区ID。

* `zone_id_slave2` - (可选) 第二个备用实例的可用区ID。

* `tde_status` - (可选) 启用透明数据加密(TDE)以用于RDS实例。

* `enable_ssl` - (可选，变更时重建) 指定是否为RDS实例启用SSL加密。

* `storage_type` - (可选，变更时重建) 实例使用的存储介质类型。有效值：
  * **local_ssd**: 本地SSD磁盘。
  * **ephemeral_ssd**: 临时SSD磁盘。
  * **cloud_ssd**: 云SSD磁盘。
  * **cloud_essd**: 云ESSD磁盘。

* `db_instance_storage_type` - (可选，变更时重建) 实例的存储类型。有效值：
  * **local_ssd**, **ephemeral_ssd**: 本地SSD磁盘。
  * **cloud_ssd**: 云SSD磁盘。
  * **cloud_essd**: 云ESSD磁盘。

* `encryption_key` - (可选) 同一区域中磁盘加密的密钥ID。此参数表示启用了云盘加密，并且在启用后无法关闭。您可以在密钥管理服务(KMS)控制台中查看密钥ID或创建新密钥。

* `encryption` - (可选，变更时重建) 指定是否启用加密。

* `db_instance_class` - (可选) 实例类型。更多信息，请参见[实例类型表](https://www.alibabacloud.com/help/doc-detail/26312.htm)。

* `db_instance_storage` - (可选) 用户定义的DB实例存储空间。取值范围取决于数据库类型和版本。按5 GB递增。

* `instance_charge_type` - (可选) 实例的计费方式。有效值：
  * **Prepaid**: 预付费计费。
  * **Postpaid**: 后付费计费。
  * **Serverless**: Serverless计费(仅支持MySQL实例)。

* `payment_type` - (可选) 实例的支付类型。有效值：
  * **PayAsYouGo**: 按量付费。
  * **Subscription**: 包年包月。
  * **Serverless**: Serverless付费类型(仅支持MySQL实例)。

* `period` - (可选) 实例的订阅时长。有效值：
  * Year: 年度订阅。
  * Month: 月度订阅。

* `monitoring_period` - (可选) 监控频率，单位为秒。有效值：`5`, `10`, `60`, `300`。默认值为`300`。

* `auto_renew` - (可选) 指定实例是否自动续费。仅在创建包年包月实例时传递。有效值：
  * **true**
  * **false**

* `auto_renew_period` - (可选) 实例的自动续费周期。有效值：`1`~`12`个月。

* `zone_id` - (可选，变更时重建) 启动DB实例的可用区。

* `vswitch_id` - (可选，变更时重建) 启动DB实例所在的虚拟交换机ID。

* `instance_name` - (可选) DB实例的名称。长度必须为2到256个字符。

* `db_instance_description` - (可选) DB实例的描述。

* `security_ip_mode` - (可选) 指定安全IP模式。有效值：
  * **normal**: 普通模式。
  * **safety**: 高安全模式。

* `maintain_time` - (可选) 实例的维护时间段，UTC时间。格式：`HH:MMZ-HH:MMZ`。

* `role_arn` - (可选) 授权RDS云服务账号访问KMS的全局资源标识符(ARN)。您可以使用[CheckCloudResourceAuthorized](~~ 446261 ~~)接口查看ARN信息。

* `force_restart` - (可选) 指定是否强制重启实例。

* `tags` - (可选) 要分配给资源的标签映射。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - RDS实例的ID。

* `port` - RDS实例的连接端口。

* `connection_string` - RDS实例的连接字符串。

* `storage_type` - 实例使用的存储介质类型。

* `db_instance_storage_type` - 实例的存储类型。有效值：
  * **local_ssd**, **ephemeral_ssd**: 本地SSD磁盘。
  * **cloud_ssd**: 云SSD磁盘。
  * **cloud_essd**: 云ESSD磁盘。

* `db_instance_class` - 实例类型。

* `db_instance_storage` - 实例的存储大小。

* `instance_charge_type` - 实例的计费方式。

* `payment_type` - 实例的支付类型。

* `monitoring_period` - 监控频率，单位为秒。

* `zone_id` - 实例的可用区ID。

* `instance_name` - DB实例的名称。

* `db_instance_description` - DB实例的描述。

* `maintain_time` - 实例的维护时间段。

* `role_arn` - 授权RDS云服务账号访问KMS的全局资源标识符(ARN)。