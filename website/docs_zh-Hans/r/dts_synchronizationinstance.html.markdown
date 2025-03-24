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

## 参数参考

支持以下参数：

* `compute_unit` - (选填) - ETL的规格。单位为计算单元ComputeUnit(CU)，1CU=1vCPU + 4GB内存。取值范围为大于等于2的整数。传入该参数，开通[ETL功能](https://help.aliyun.com/document_detail/212324.html)，进行数据清洗和转换。
* `database_count` - (选填) - PolarDB-X下的私有定制RDS实例的数量，默认值为**1**。仅当`source_endpoint_engine_name`为**drds**时需要传入该参数。
* `quantity` - (选填) - 购买实例数量。当前单次调用最多支持购买1个。
* `sync_architecture` - (选填, 变更时重建) - 同步拓扑，取值：
  - **oneway**：单向同步，为默认值。
  - **bidirectional**：双向同步。
* `destination_endpoint_engine_name` - (必填, 变更时重建) - 目标数据库引擎类型。有效值包括但不限于：
  - **MySQL**：MySQL数据库(包括RDS MySQL和自建MySQL)。
  - **PolarDB**：PolarDB MySQL。
  - **polardb_o**：PolarDB-O。
  - **polardb_pg**：PolarDB PostgreSQL。
  - **Redis**：Redis数据库(包括云数据库Redis和自建Redis)。
  - **DRDS**：云原生分布式数据库PolarDB-X 1.0和2.0。
  - **PostgreSQL**：自建PostgreSQL。
  - **odps**：MaxCompute。
  - **oracle**：自建Oracle。
  - **mongodb**：MongoDB数据库(包括云数据库MongoDB和自建MongoDB)。
  - **tidb**：TiDB数据库。
  - **ADS**：云原生数仓 AnalyticDB MySQL 2.0。
  - **ADB30**：云原生数仓 AnalyticDB MySQL 3.0。
  - **Greenplum**：云原生数仓 AnalyticDB PostgreSQL。
  - **MSSQL**：SQL Server数据库(包括RDS SQL Server和自建SQL Server)。
  - **kafka**：Kafka数据库(包括消息队列Kafka版和自建Kafka)。
  - **DataHub**：阿里云流式数据服务DataHub。
  - **clickhouse**：云数据库 ClickHouse。
  - **DB2**：自建DB2 LUW。
  - **as400**：AS/400。
  - **Tablestore**：表格存储Tablestore。
  > 默认取值为**MySQL**。关于支持的源库和目标库对应情况，请参见支持的[数据库、同步初始化类型和同步拓扑](https://help.aliyun.com/document_detail/130744.html), [支持的数据库和迁移类型](https://help.aliyun.com/document_detail/26618.html)。
* `destination_endpoint_region` - (必填, 变更时重建) - 目标实例所在区域。
* `source_endpoint_engine_name` - (必填, 变更时重建) - 源实例数据库引擎类型。有效值与`destination_endpoint_engine_name`相同。
* `source_endpoint_region` - (必填, 变更时重建) - 源实例所在区域。
* `instance_class` - (选填) - 迁移或同步实例的规格。迁移实例支持的规格：**xxlarge**、**xlarge**、**large**、**medium**、**small**。同步实例支持的规格：**large**、**medium**、**small**、**micro**。不同规格对应的性能说明，请参见[数据迁移链路规格说明](https://help.aliyun.com/document_detail/26606.html)和[数据同步链路规格说明](https://help.aliyun.com/document_detail/26605.html)。
* `payment_type` - (必填, 变更时重建) - 代表付费类型的资源属性字段。有效值：`Subscription`, `PayAsYouGo`。
* `payment_duration_unit` - (选填) - 付款时长单位。有效值：`Month`, `Year`。当`payment_type`为`Subscription`时，此参数有效且必须传递。
* `payment_duration` - (选填) - 预付费实例购买的时长。当`payment_type`为`Subscription`时，此参数有效且必须传递。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 同步实例的资源ID。
* `dts_job_id` - 同步实例的任务ID。
* `status` - 同步实例的状态。
* `instance_class` - 迁移或同步实例的规格。迁移实例支持的规格：**xxlarge**, **xlarge**, **large**, **medium**, **small**。同步实例支持的规格：**large**, **medium**, **small**, **micro**。关于不同规格的性能描述，请参见[数据迁移链路规格说明](https://help.aliyun.com/document_detail/26606.html)和[数据同步链路规格说明](https://help.aliyun.com/document_detail/26605.html)。
