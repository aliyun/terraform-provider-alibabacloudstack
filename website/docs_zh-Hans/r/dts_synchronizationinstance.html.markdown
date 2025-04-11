---
subcategory: "Data Transmission Service (DTS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dts_synchronization_instance"
sidebar_current: "docs-alibabacloudstack-resource-dts-synchronization-instance"
description: |-
  编排数据传输服务（Dts）数据同步实例
---

# alibabacloudstack_dts_synchronizationinstance
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dts_synchronization_instance`

使用Provider配置的凭证在指定的资源集下编排数据传输服务（Dts）数据同步实例。
更多信息，请参见[创建同步实例](https://help.aliyun.com/document_detail/209156.html)。

## 示例用法

以下是一个完整的示例，展示如何创建一个DTS同步实例：

```terraform
variable "name" {
    default = "tf-testaccdtssynchronization_instance61662"
}

resource "alibabacloudstack_dts_synchronizationinstance" "default" {
  source_endpoint_engine_name      = "MySQL"
  destination_endpoint_engine_name  = "MySQL"
  destination_endpoint_region       = "cn-hangzhou"
  source_endpoint_region           = "cn-hangzhou"
  payment_type                     = "PayAsYouGo"
  instance_class                   = "small"
  sync_architecture                = "oneway"
  compute_unit                    = 4
  database_count                  = 1
  quantity                        = 1
}
```

## 参数说明

支持以下参数：

* `payment_type` - (必填, 变更时重建) - 代表付费类型的资源属性字段。有效值：`Subscription`, `PayAsYouGo`。
* `source_endpoint_engine_name` - (必填, 变更时重建) - 源实例数据库引擎类型。有效值包括但不限于：`ADS`, `DB2`, `DRDS`, `DataHub`, `Greenplum`, `MSSQL`, `MySQL`, `PolarDB`, `PostgreSQL`, `Redis`, `Tablestore`, `as400`, `clickhouse`, `kafka`, `mongodb`, `odps`, `oracle`, `polardb_o`, `polardb_pg`, `tidb`。关于支持的源库和目标库对应情况，请参见支持的[数据库、同步初始化类型和同步拓扑](https://help.aliyun.com/document_detail/130744.html), [支持的数据库和迁移类型](https://help.aliyun.com/document_detail/26618.html)。
* `source_endpoint_region` - (必填, 变更时重建) - 源实例所在区域。
* `destination_endpoint_engine_name` - (必填, 变更时重建) - 目标数据库引擎类型。有效值与`source_endpoint_engine_name`相同。
* `destination_endpoint_region` - (必填, 变更时重建) - 目标实例所在区域。
* `instance_class` - (选填) - 迁移或同步实例的规格。有效值：`large`, `medium`, `small`, `micro`, `xlarge`, `xxlarge`。仅支持升级配置，不支持降级配置。如果需要降级实例，请[提交工单](https://selfservice.console.aliyun.com/ticket/category/dts/today)。
* `sync_architecture` - (选填, 变更时重建) - 同步拓扑。有效值：`oneway`, `bidirectional`。
* `compute_unit` - (选填) - ETL的规格。单位为计算单元ComputeUnit(CU)，1CU=1vCPU + 4GB内存。取值范围为大于等于2的整数。传入该参数并启用[ETL功能](https://help.aliyun.com/document_detail/212324.html)，进行数据清洗和转换。
* `database_count` - (选填) - PolarDB-X下的私有定制RDS实例的数量，默认值为**1**。仅当`source_endpoint_engine_name`为**drds**时需要传入该参数。
* `quantity` - (选填) - 购买实例数量。当前单次调用最多支持购买1个。
* `payment_duration_unit` - (选填) - 付款时长单位。有效值：`Month`, `Year`。当`payment_type`为`Subscription`时，此参数有效且必须传递。
* `payment_duration` - (必填，当`payment_type`为`Subscription`时) - 预付费实例购买的时长。当`payment_type`为`Subscription`时，此参数有效且必须传递。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 同步实例的资源ID。
* `dts_job_id` - 同步实例的任务ID。
* `status` - 同步实例的状态。
* `instance_class` - 迁移或同步实例的规格。迁移实例支持的规格：`xxlarge`, `xlarge`, `large`, `medium`, `small`。同步实例支持的规格：`large`, `medium`, `small`, `micro`。关于不同规格的性能描述，请参见[数据迁移链路规格说明](https://help.aliyun.com/document_detail/26606.html)和[数据同步链路规格说明](https://help.aliyun.com/document_detail/26605.html)。