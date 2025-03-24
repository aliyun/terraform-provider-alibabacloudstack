---
subcategory: "Data Transmission Service (DTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dts_synchronization_job"
sidebar_current: "docs-alibabacloudstack-resource-dts-synchronization-job"
description: |-
  编排数据传输服务（Dts）数据同步任务
---

# alibabacloudstack_dts_synchronization_job

使用Provider配置的凭证在指定的资源集下编排数据传输服务（Dts）数据同步任务。

有关 DTS 数据同步作业及其使用方法的更多信息，请参阅 [什么是数据同步作业](https://www.alibabacloud.com/product/data-transmission-service)。


## 示例用法

### 基础用法

```terraform
variable "password" {
}

resource "alibabacloudstack_dts_synchronization_instance" "default" {
  payment_type                     = "PayAsYouGo"
  source_endpoint_engine_name      = "PolarDB"
  source_endpoint_region           = "cn-hangzhou"
  destination_endpoint_engine_name = "ADB30"
  destination_endpoint_region      = "cn-hangzhou"
  instance_class                   = "small"
  sync_architecture                = "oneway"
}

resource "alibabacloudstack_dts_synchronization_job" "default" {
  dts_instance_id                    = alibabacloudstack_dts_synchronization_instance.default.id
  dts_job_name                       = "tf-testAccCase1"
  source_endpoint_instance_type      = "PolarDB"
  source_endpoint_instance_id        = "pc-xxxxxxxx"
  source_endpoint_engine_name        = "PolarDB"
  source_endpoint_region             = "cn-hangzhou"
  source_endpoint_database_name      = "tf-testacc"
  source_endpoint_user_name          = "root"
  source_endpoint_password           = var.password
  destination_endpoint_instance_type = "ads"
  destination_endpoint_instance_id   = "am-xxxxxxxx"
  destination_endpoint_engine_name   = "ADB30"
  destination_endpoint_region        = "cn-hangzhou"
  destination_endpoint_database_name = "tf-testacc"
  destination_endpoint_user_name     = "root"
  destination_endpoint_password      = var.password
  db_list                            = "{\"tf-testacc\":{\"name\":\"tf-test\",\"all\":true,\"state\":\"normal\"}}"
  structure_initialization           = "true"
  data_initialization                = "true"
  data_synchronization               = "true"
  status                             = "Synchronizing"
}
```

## 参数参考

以下参数受支持：

* `dts_instance_id` - (必填，变更时重建) 同步实例 ID。`alibabacloudstack_dts_synchronization_instance` 的 ID。
* `synchronization_direction` - (可选，变更时重建) 同步方向。有效值：`Forward`，`Reverse`。仅当 `alibabacloudstack_dts_synchronization_instance` 的属性 `sync_architecture` 为 `bidirectional` 时，此参数应传递，否则不应指定此参数。
* `dts_job_name` - (可选) 同步作业名称。
* `dts_job_id` - (可选，变更时重建) 同步实例的作业 ID。
* `instance_class` - (可选) 实例类。有效值：`large`，`medium`，`micro`，`small`，`xlarge`，`xxlarge`。您只能升级配置，不能降级配置。如果您要降级实例，请[提交工单](https://selfservice.console.aliyun.com/ticket/category/dts/today)。
* `checkpoint` - (可选，计算，变更时重建) 以 Unix 时间戳格式表示的开始时间。
* `data_initialization` - (必填，变更时重建) 是否执行 DTS 支持的架构迁移、全量数据迁移或全量数据初始化。值包括：
* `data_synchronization` - (必填，变更时重建) 是否对迁移类型或同步执行增量数据迁移。值包括：
* `structure_initialization` - (必填，变更时重建) 是否对数据库表结构进行迁移或初始化。值包括：
* `db_list` - (必填，变更时重建) 迁移对象，JSON 字符串格式。有关详细定义说明，请参阅 [迁移、同步或订阅对象的描述](https://help.aliyun.com/document_detail/209545.html)。
* `reserve` - (可选，变更时重建) DTS 预留参数，格式为 JSON 字符串，您可以在此参数中完成源和目标数据库信息(例如目标 Kafka 数据库的数据存储格式、云企业网 CEN 实例 ID)。有关更多说明，请参阅参数 [Reserve 参数的描述](https://help.aliyun.com/document_detail/273111.html)。
* `source_endpoint_instance_type` - (必填，变更时重建) 源实例类型。有效值：`CEN`，`DG`，`DISTRIBUTED_DMSLOGICDB`，`ECS`，`EXPRESS`，`MONGODB`，`OTHER`，`PolarDB`，`POLARDBX20`，`RDS`。
* `source_endpoint_engine_name` - (必填，变更时重建) 源数据库类型。有效值：`AS400`，`DB2`，`DMSPOLARDB`，`HBASE`，`MONGODB`，`MSSQL`，`MySQL`，`ORACLE`，`PolarDB`，`POLARDBX20`，`POLARDB_O`，`POSTGRESQL`，`TERADATA`。
* `source_endpoint_instance_id` - (可选，变更时重建) 源实例 ID。
* `source_endpoint_region` - (可选，变更时重建) 源实例区域。
* `source_endpoint_ip` - (可选，变更时重建) 源端点 IP。
* `source_endpoint_port` - (可选，变更时重建) 源端点端口。
* `source_endpoint_oracle_sid` - (可选，变更时重建) Oracle 数据库的 SID。
* `source_endpoint_database_name` - (可选，变更时重建) 迁移的数据库名称。
* `source_endpoint_user_name` - (可选，变更时重建) 数据库账户的用户名。
* `source_endpoint_password` - (可选) 数据库账户的密码。
* `source_endpoint_owner_id` - (可选，变更时重建) 源实例所属的阿里云账户 ID。
* `source_endpoint_role` - (可选，变更时重建) 源实例所属云账户配置的角色名称。
* `destination_endpoint_instance_type` - (必填，变更时重建) 目标实例类型。有效值：`ads`，`CEN`，`DATAHUB`，`DG`，`ECS`，`EXPRESS`，`GREENPLUM`，`MONGODB`，`OTHER`，`PolarDB`，`POLARDBX20`，`RDS`。
* `destination_endpoint_engine_name` - (必填，变更时重建) 目标数据库类型。有效值：`ADB20`，`ADB30`，`AS400`，`DATAHUB`，`DB2`，`GREENPLUM`，`KAFKA`，`MONGODB`，`MSSQL`，`MySQL`，`ORACLE`，`PolarDB`，`POLARDBX20`，`POLARDB_O`，`PostgreSQL`。
* `destination_endpoint_instance_id` - (可选，变更时重建) 目标实例 ID。
* `destination_endpoint_region` - (可选，变更时重建) 目标实例区域。
* `destination_endpoint_ip` - (可选，变更时重建) 目标端点 IP。
* `destination_endpoint_port` - (可选，变更时重建) 目标端点端口。
* `destination_endpoint_database_name` - (可选，变更时重建) 迁移的数据库名称。
* `destination_endpoint_user_name` - (可选，变更时重建) 数据库账户的用户名。
* `destination_endpoint_password` - (可选，变更时重建) 数据库账户的密码。
* `destination_endpoint_oracle_sid` - (可选，变更时重建) Oracle 数据库的 SID。
* `delay_notice` - (可选，变更时重建) 延迟通知。有效值：`true`，`false`。
* `delay_phone` - (可选，变更时重建) 延迟电话。延迟报警联系人的手机号码。多个手机号码用英文逗号 `,` 分隔。该参数目前仅支持中国站点，并且仅支持大陆手机号码，最多可以传入 10 个手机号码。
* `delay_rule_time` - (可选，变更时重建) 延迟规则时间。当 `delay_notice` 设置为 `true` 时，必须传入此参数。触发延迟报警的阈值。单位为秒，需要为整数。阈值可以根据业务需求设置。建议设置在 10 秒以上，以避免因网络和数据库负载引起的延迟波动。
* `error_notice` - (可选，变更时重建) 错误通知。有效值：`true`，`false`。
* `error_phone` - (可选，变更时重建) 错误电话。错误报警联系人的手机号码。多个手机号码用英文逗号 `,` 分隔。该参数目前仅支持中国站点，并且仅支持大陆手机号码，最多可以传入 10 个手机号码。
* `status` - (可选) 资源状态。有效值：`Synchronizing`，`Suspending`。通过指定 `Suspending` 可以停止任务，通过指定 `Synchronizing` 可以启动任务。
* `source_endpoint_password` - (可选，变更时重建) 源数据库账户的密码。
* `destination_endpoint_password` - (可选，变更时重建) 目标数据库账户的密码。

-> **注意:** 从 `NotStarted` 状态到 `Synchronizing` 状态，资源会经历 `Prechecking` 和 `Initializing` 阶段。由于 `Initializing` 阶段耗时较长，一旦资源进入 `Prechecking` 状态，就可以认为任务可以正常执行。因此，我们将 `Initializing` 状态视为等同于 `Synchronizing` 状态。

-> **注意:** 如果您想通过属性 `instance_class` 升级同步作业规格，还必须修改其实例的属性 `instance_class`，以保持它们的一致性。

## 注意事项

1. 年费月付类型的作业暂停后，到期时间无法更改；
2. 按量付费类型作业暂停后，您的作业配置费用仍会被收取；
3. 如果任务暂停超过 6 小时，任务将无法成功启动。
4. 暂停任务只会停止写入目标库，但仍然会继续获取源的增量日志，以便取消暂停后任务能够快速恢复。因此，在此期间，源库的一些资源，如带宽资源，将继续被占用。
5. 在任务暂停期间，费用将继续产生。如果需要停止计费，请释放实例。
6. 当 DTS 实例暂停超过 7 天，实例无法恢复，状态将从暂停变为失败。

## 属性参考

导出以下属性：

* `id` - DTS 同步作业的 Terraform 资源 ID。
* `dts_job_name` - 同步作业名称。
* `checkpoint` - 以 Unix 时间戳格式表示的开始时间。
* `instance_class` - 实例类。
* `status` - 资源状态。有效值：`Synchronizing`，`Suspending`。

## 导入

可以通过 id 导入 DTS 同步作业，例如：

```bash
$ terraform import alibabacloudstack_dts_synchronization_job.example <id>
```