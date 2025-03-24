---
subcategory: "PolarDB"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_polardb_dbinstances"
sidebar_current: "docs-Alibabacloudstack-datasource-polardb-dbinstances"
description: |- 
  查询polardb数据库实例
---

# alibabacloudstack_polardb_dbinstances
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_polardb_instances`

根据指定过滤条件列出当前凭证权限可以访问的polardb数据库实例列表。

## 示例用法
```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "${var.creation}"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

variable "name" {
  default = "tf-testAccDBInstanceConfig"
}

variable "creation" {
  default = "PolarDB"
}

resource "alibabacloudstack_polardb_instance" "default" {
  engine                 = "MySQL"
  engine_version         = "5.7"
  instance_name          = "${var.name}"
  db_instance_storage_type = "local_ssd"
  db_instance_storage    = 5
  db_instance_class      = "rds.mysql.t1.small"
  zone_id               = "${data.alibabacloudstack_zones.default.zones.0.id}"
  vswitch_id            = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_polardb_instances" "default" {
  ids                   = ["${alibabacloudstack_polardb_instance.default.id}"]
  network_type          = "VPC"
  db_instance_class     = "${alibabacloudstack_polardb_instance.default.db_instance_class}"
  status               = "Running"
  region_id            = "${data.alibabacloudstack_zones.default.zones.0.region_id}"
  payment_type         = "PayAsYouGo"
  engine               = "MySQL"
  output_file          = "output.txt"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 用于过滤的数据库实例 ID 列表。
  * `network_type` - (选填) - 实例的网络类型，取值：* **Classic**：经典网络。* **VPC**：专有网络。
  * `db_instance_type` - (选填) - 数据库实例的类型。
  * `vswitch_id` - (选填) - VSwitch 的 ID。
  * `vpc_id` - (选填) - 专有网络VPCID。
  * `db_instance_class` - (选填) - 实例规格，详情请参见[实例规格表](~~26312~~)。
  * `payment_type` - (选填) - 实例付费方式，取值：**PayAsYouGo**：按量付费。**Subscription**：包年包月。**Serverless**：Serverless付费类型，仅支持MySQL实例。
  * `status` - (选填) - 资源状态。
  * `db_instance_id` - (选填) - 实例ID。
  * `engine_version` - (选填) - 数据库版本。
  * `resource_group_id` - (选填) - 资源组ID。
  * `region_id` - (必填) - 代表region的资源属性字段。
  * `engine` - (选填) - 数据库类型。返回值：* MySQL* PostgreSQL* SQLServer* MariaDB。
  
## 属性参考
除了上述参数外，还导出以下属性：
  * `db_instances` - 数据库实例列表。列表中的每个元素是一个包含以下键的映射：
    * `id` - 数据库实例的 ID。
    * `auto_pay` - 是否自动支付。取值范围：- **true**：自动支付。您需要确保账户余额充足。- **false**：只生成订单不扣费。> 默认值为 true。如果您的支付方式余额不足，可以将参数 AutoPay 设置为 false，此时会生成未支付订单，您可以登录 POLARDB 管理控制台自行支付。
    * `auto_renew` - 实例是否自动续费，仅在创建包年包月实例时传入，取值：- **true**- **false**> * 按月购买，则自动续费周期为1个月。* 按年购买，则自动续费周期为1年。
    * `auto_upgrade_minor_version` - 实例升级小版本的方式，取值：* **Auto**：自动升级小版本。* **Manual**：不自动升级，仅在当前版本下线时才强制升级。
    * `business_info` - 业务扩展参数。
    * `category` - 实例系列，取值：* **Basic**：基础版。* **HighAvailability**：高可用版。* **AlwaysOn**：集群版。* **Finance**：三节点企业版。* **serverless_basic**：Serverless基础版。
    * `classic_expired_days` - 经典网络连接地址保留天数，取值：**1~120**，单位：天。
    * `commodity_code` - 商品编码。
    * `connection_mode` - 实例的访问模式，取值：* **Standard**：标准访问模式。* **Safe**：数据库代理模式。默认为 POLARDB 系统分配。> SQL Server 2012、2016、2017 只支持标准访问模式。
    * `connection_string_prefix` - 只读地址前缀名，不可重复，由小写字母和中划线组成，需以字母开头，长度不超过30个字符。> 默认以“实例名+rw”字符串组成前缀。
    * `connection_string_type` - 连接地址类型，取值：* **Normal**：普通连接* **ReadWriteSplitting**：读写分离连接，默认返回所有连接。
    * `current_connection_string` - 实例当前的某个连接地址，可以是内外网连接地址，或者混访模式下的经典网络连接地址。
    * `db_instance_description` - 数据库实例的描述。
    * `db_instance_storage` - 数据库实例的存储容量。
    * `db_instance_type` - 数据库实例的类型。
    * `db_instance_class` - 实例规格，详情请参见[实例规格表](~~26312~~)。
    * `db_instance_id` - 实例ID。
    * `db_instance_net_type` - 网络类型，取值：* **Intranet**：内网。
    * `db_instance_storage_type` - 实例储存类型，取值：* **local_ssd**、**ephemeral_ssd**：本地SSD盘。* **cloud_ssd**：SSD云盘。* **cloud_essd**：ESSD云盘。
    * `distribution_type` - 读权重分配模式，取值：- **Standard**：按规格权重自动分配- **Custom**：自定义分配权重。
    * `effective_time` - 生效时间，取值：* **Immediate**：立即生效。* **MaintainTime**：在可运维时间段内生效，请参见[ModifyDBInstanceMaintainTime](~~26249~~)。默认值：**Immediate**。
    * `encryption_key` - 同地域内的云盘加密的密钥ID。传入此参数表示开启云盘加密(开启后无法关闭)，并且需要传入**RoleARN**。您可以在密钥管理服务控制台查看密钥ID，也可以创建新的密钥。详情请参见[创建密钥](~~181610~~)。
    * `engine` - 数据库类型。返回值：* MySQL* PostgreSQL* SQLServer* MariaDB。
    * `engine_version` - 数据库版本。
    * `expire_time` - 到期时间。<i>yyyy-MM-dd</i>T<i>HH:mm:ss</i>Z(UTC时间)。> 按量付费实例无到期时间。
    * `lock_mode` - 实例锁定模式，取值：* **Unlock**：正常。* **ManualLock**：手动触发锁定。* **LockByExpiration**：实例过期自动锁定。* **LockByRestoration**：实例回滚前的自动锁定。* **LockByDiskQuota**：实例空间满自动锁定。* **LockReadInstanceByDiskQuota**：只读实例空间满自动锁定。
    * `lock_reason` - 锁定原因。
    * `maintain_time` - 实例可维护时间段，是UTC时间，+8小时才是控制台上显示的可维护时间段。
    * `master_instance_id` - 主实例的ID，如果没有返回此参数则表示该实例是主实例。
    * `max_delay_time` - 延迟阈值，范围是0~7200，单位：秒，默认为30。> 当只读实例延迟超过该阈值时，读取流量不发往该实例。
    * `network_type` - 实例的网络类型，取值：* **Classic**：经典网络。* **VPC**：专有网络。
    * `payment_type` - 实例付费方式，取值：**PayAsYouGo**：按量付费。**Subscription**：包年包月。**Serverless**：Serverless付费类型，仅支持MySQL实例。
    * `period` - 指定预付费实例为包年或者包月类型，取值：Year：包年。Month：包月。
    * `port` - 连接端口。
    * `private_ip_address` - 无需配置，表示目标实例的内网IP。系统默认通过VPCId和vSwitchId自动分配。
    * `read_write_splitting_classic_expired_days` - 读写分离的经典网络地址保留的天数，取值**1-120**，单位：天。默认值：**7**。> 当实例存在经典网络类型的读写分离地址，且**RetainClassic**=**True**，本参数有效。
    * `read_write_splitting_private_ip_address` - 设置实例的内网读写分离地址的IP，需要在指定交换机的IP地址范围内。系统默认通过**VPCId**和**VSwitchId**自动分配。> 当前实例存在经典网络类型的读写分离地址时，该值有效。
    * `record_total` - 总记录数。
    * `region_id` - 代表region的资源属性字段。
    * `resource_group_id` - 资源组ID。
    * `resource_type` - 资源类型定义。唯一取值：**INSTANCE**。
    * `retain_classic` - 是否保留经典网络地址，取值：* **True**：保留* **False**：不保留，默认值：**False**。
    * `role_arn` - 主账号授权POLARDB云服务账号访问KMS权限的全局资源描述符(ARN)。您可以通过[CheckCloudResourceAuthorized](~~446261~~)接口查看ARN信息。
    * `security_ip_list` - 安全IP地址列表。
    * `security_ip_mode` - 安全IP地址模式。
    * `sql_collector_status` - 开启或关闭SQL洞察(SQL审计)，取值：**Enable | Disabled**。
    * `status` - 资源状态。
    * `table_meta` - 指定恢复的库表。格式：```[{"type":"db","name":"<数据库1名称>","newname":"<新数据库1名称>","tables":[{"type":"table","name":"<数据库1内的表1名称>","newname":"<新的表1名称>"},{"type":"table","name":"<数据库1内的表2名称>","newname":"<新的表2名称>"}]},{"type":"db","name":"<数据库2名称>","newname":"<新数据库2名称>","tables":[{"type":"table","name":"<数据库2内的表3名称>","newname":"<新的表3名称>"},{"type":"table","name":"<数据库2内的表4名称>","newname":"<新的表4名称>"}]}]```
    * `tags` - 代表资源标签的资源属性字段。
    * `temp_db_instance_id` - 临时 DB 实例的 ID。
    * `used_time` - 指定购买时长，取值：当参数Period=Year时，UsedTime取值为1~5。当参数Period=Month时，UsedTime取值为1~11。
    * `vswitch_id` - VSwitch 的 ID。
    * `vpc_cloud_instance_id` - 专有网络实例ID。
    * `vpc_id` - 专有网络VPCID。
    * `weight` - 读权重分配，即传入主实例和只读实例的读请求比例。以100进行递增，最大值为10000。* POLARDB实例格式：`{"<只读实例ID>":<权重>,"master":<权重>,"slave":<权重>}`* MyBASE实例格式：`[{"instanceName":"<主实例ID>","weight":<权重>,"role":"master"},{"instanceName":"<主实例ID>","weight":<权重>,"role":"slave"},{"instanceName":"<只读实例ID>","weight":<权重>,"role":"master"}]`> - 当**DistributionType**为**Custom**时，必须传入该参数。> - 当**DisrtibutionType**为**Standard**时，传入该参数无效。
    * `zone_id_slave_one` - 仅当原实例为高可用版时，该参数可配置，表示目标实例备可用区ID。POLARDB PostgreSQL支持升级后将新的备实例配置到与原实例同一地域的其他可用区。可以通过[DescribeRegions](~~26243~~)接口查看可用区ID。
    * `zone_id_slave_two` - 备可用区2。