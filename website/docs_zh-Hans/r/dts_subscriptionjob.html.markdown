---
subcategory: "Data Transmission Service (DTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dts_subscriptionjob"
sidebar_current: "docs-Alibabacloudstack-dts-subscriptionjob"
description: |- 
  编排数据传输服务（Dts）订阅任务
---

# alibabacloudstack_dts_subscription_job
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dts_subscriptionjob`

使用Provider配置的凭证在指定的资源集下编排数据传输服务（Dts）订阅任务。

## 示例用法

```hcl
variable "name" {
  default = "tf-testaccdtstf-testaccdtssubscriptionjob30638"
}

variable "creation" {
  default = "Rds"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = var.creation
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name       = var.name
  cidr_block     = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = alibabacloudstack_vpc.default.id
  cidr_block        = "172.16.0.0/24"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_db_instance" "instance" {
  engine           = "MySQL"
  engine_version   = "5.6"
  instance_type    = "rds.mysql.s2.large"
  instance_storage = "30"
  vswitch_id       = alibabacloudstack_vswitch.default.id
  instance_name    = var.name
  storage_type     = "local_ssd"
}

resource "alibabacloudstack_db_database" "db" {
  count        = 2
  instance_id  = alibabacloudstack_db_instance.instance.id
  name         = "tfaccountpri_${count.index}"
  description  = "from terraform"
  character_set = "UTF8"
}

resource "alibabacloudstack_db_account" "account" {
  instance_id  = alibabacloudstack_db_instance.instance.id
  name         = "tftestprivilege"
  password     = "inputYourCodeHere"
  description  = "from terraform"
}

resource "alibabacloudstack_db_account_privilege" "privilege" {
  instance_id  = alibabacloudstack_db_instance.instance.id
  account_name = alibabacloudstack_db_account.account.name
  privilege    = "ReadWrite"
  db_names     = alibabacloudstack_db_database.db.*.name
}

resource "alibabacloudstack_vpc" "default1" {
  vpc_name       = var.name
  cidr_block     = "10.0.0.0/8"
}

resource "alibabacloudstack_vswitch" "default1" {
  vpc_id            = alibabacloudstack_vpc.default1.id
  cidr_block        = "10.1.0.0/16"
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  name              = var.name
}

resource "alibabacloudstack_dts_subscription_job" "default" {
  source_endpoint_database_name = "tfaccountpri_0"
  source_endpoint_engine_name   = "MySQL"
  source_endpoint_instance_type = "RDS"
  source_endpoint_region        = "cn-hangzhou"
  source_endpoint_user_name     = "tftestprivilege"
  subscription_instance_vpc_id   = alibabacloudstack_vpc.default1.id
  subscription_instance_vswitch_id = alibabacloudstack_vswitch.default1.id
  subscription_instance_network_type = "vpc"
  dts_job_name                 = "tf-testAccCase"
  source_endpoint_instance_id  = alibabacloudstack_db_instance.instance.id
  source_endpoint_password     = "inputYourCodeHere"
  payment_type                = "PayAsYouGo"
  db_list                     = "{\"tfaccountpri_0\":{\"name\":\"tfaccountpri_0\",\"all\":true,\"state\":\"normal\"}}"
}
```

## 参数说明

支持以下参数：

* `checkpoint` - (选填) 订阅启动时间，格式为Unix时间戳。
* `compute_unit` - (选填) ETL规格。单位是计算单元 ComputeUnit (CU)，1CU=1vCPU+4 GB 内存。取值范围是大于等于 2 的整数。
* `database_count` - (选填) PolarDB-X 下用户自定义 RDS 实例个数，默认值为 1。仅当 `source_endpoint_engine_name` 等于 `drds` 时需要传递该参数。
* `db_list` - (选填) 订阅对象，JSON字符串格式。具体定义请参考迁移、同步或订阅对象 [文档](https://help.aliyun.com/document_detail/209545.html)。
* `delay_notice` - (选填) 是否监控延迟状态。有效值：`true`, `false`。
* `delay_phone` - (选填) 延迟告警联系人的手机号码。多个手机号码以英文逗号 `,` 分隔。此参数目前仅支持中国站，且仅支持大陆手机号码，最多可以传入 10 个手机号码。
* `delay_rule_time` - (选填) 当 `delay_notice` 设置为 `true` 时，必须传入此参数。触发延迟告警的阈值，单位为秒，需为整数。阈值可以根据业务需求设置，建议设置在 10 秒以上，避免因网络和数据库负载导致的延迟波动。
* `destination_endpoint_engine_name` - (选填) 目标端引擎名称。有效值：`ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`。
* `destination_region` - (选填) 目标区域。[支持的区域列表](https://help.aliyun.com/document_detail/141033.html)。
* `dts_instance_id` - (计算后, 变更时重建) 订阅实例ID。
* `dts_job_name` - (选填) 订阅任务名称。
* `error_notice` - (选填) 是否监控异常状态。有效值：`true`, `false`。
* `error_phone` - (选填) 异常告警联系人的手机号码。多个手机号码以英文逗号 `,` 分隔。此参数目前仅支持中国站，且仅支持大陆手机号码，最多可以传入 10 个手机号码。
* `instance_class` - (选填) 实例规格。有效值：`large`, `medium`, `micro`, `small`, `xlarge`, `xxlarge`。
* `payment_type` - (必填, 变更时重建) 资源的计费类型。有效值：`Subscription`, `PayAsYouGo`。
* `payment_duration_unit` - (选填) 计费时长单位。有效值：`Month`, `Year`。当 `payment_type` 为 `Subscription` 时，此参数有效且必须传入。
* `payment_duration` - (选填) 预付费实例购买时长。当 `payment_type` 为 `Subscription` 时，此参数有效且必须传入。
* `reserve` - (选填) DTS保留参数，格式为JSON字符串，可以传入该参数完成源库和目标库信息(如目标Kafka数据库的数据存储格式、云企业网CEN实例ID)。更多详情请参考 [Reserve参数描述](https://help.aliyun.com/document_detail/176470.html)。
* `source_endpoint_database_name` - (选填) 要订阅的数据库名称。
* `source_endpoint_engine_name` - (选填) 源数据库类型，值为 MySQL 或 Oracle。有效值：`MySQL`, `Oracle`。
* `source_endpoint_instance_id` - (选填) 源实例ID。仅当源数据库实例类型为 RDS MySQL、PolarDB-X 1.0、PolarDB MySQL 时，此参数可用且必须设置。
* `source_endpoint_instance_type` - (选填, 变更时重建) 源实例类型。有效值：`RDS`, `PolarDB`, `DRDS`, `LocalInstance`, `ECS`, `Express`, `CEN`, `dg`。
* `source_endpoint_ip` - (选填) 源端IP地址。
* `source_endpoint_oracle_sid` - (选填) Oracle数据库的SID。当源数据库为自建Oracle且Oracle数据库为非RAC实例时，此参数可用且必须传入。
* `source_endpoint_owner_id` - (选填) 源实例所属的阿里云账号ID。仅在跨阿里云账号配置数据订阅时该参数可用且必须传入。
* `source_endpoint_user_name` - (选填) 源数据库实例账号用户名。
* `source_endpoint_password` - (选填) 源数据库实例账号密码。
* `source_endpoint_port` - (选填) 源数据库端口。
* `source_endpoint_region` - (选填) 源数据库所在地域。
* `source_endpoint_role` - (选填) 授权角色。当源实例与配置订阅任务的阿里云账号不同时需要传入该参数，指定源的授权角色，允许配置订阅任务的阿里云账号访问源的源实例信息。
* `subscription_data_type_ddl` - (选填) 是否订阅DDL类型的数据。有效值：`true`, `false`。
* `subscription_data_type_dml` - (选填) 是否订阅DML类型的数据。有效值：`true`, `false`。
* `subscription_instance_network_type` - (选填) 订阅任务网络类型值：经典网络 classic。虚拟私有云 (vpc): vpc。有效值：`classic`, `vpc`。
* `subscription_instance_vpc_id` - (选填) 订阅VPC实例ID。当 `subscription_instance_network_type` 值为vpc时，此参数可用且必须传入。
* `subscription_instance_vswitch_id` - (选填) 订阅VSwitch实例ID。当 `subscription_instance_network_type` 值为vpc时，此参数可用且必须传入。
* `sync_architecture` - (选填) 同步架构。有效值：`bidirectional`, `oneway`。
* `synchronization_direction` - (选填) 同步方向。有效值：`Forward`, `Reverse`。当数据同步实例的拓扑类型为双向时，可以传入该参数反向启动反向同步链路。
* `status` - (选填) 任务状态。有效值：`Normal`, `Abnormal`。当任务创建时，任务处于 `NotStarted` 状态。可以将此状态指定为 `Normal` 来启动任务，并将状态指定为 `Abnormal` 来停止任务。**注意：我们将状态 `Starting` 视为 `Normal` 状态，并认为这两种状态在用户侧是一致的。**
* `tags` - (选填, Map) 订阅任务的标签。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform中订阅任务的资源ID。
* `checkpoint` - 订阅启动时间，格式为Unix时间戳。
* `dts_instance_id` - 订阅实例ID。
* `status` - 资源状态。
* `subscription_data_type_ddl` - 是否订阅DDL类型的数据。
* `subscription_data_type_dml` - 是否订阅DML类型的数据。