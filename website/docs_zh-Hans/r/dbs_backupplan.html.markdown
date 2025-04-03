---
subcategory: "DBS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dbs_backupplan"
sidebar_current: "docs-Alibabacloudstack-dbs-backupplan"
description: |- 
  编排Dbs备份计划
---

# alibabacloudstack_dbs_backupplan
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_dbs_backup_plan`

使用Provider配置的凭证在指定的资源集下编排Dbs备份计划。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccdbsbackupplan74295"
}

resource "alibabacloudstack_dbs_backup_plan" "default" {
  backup_method     = "logical"
  database_type     = "MySQL"
  instance_class    = "large"
  backup_plan_name  = var.name
  database_region   = "cn-hangzhou"
  storage_region    = "cn-hangzhou"
  instance_type     = "RDS"
  from_app          = "OpenAPI"
}
```

## 参数参考

支持以下参数：

* `backup_plan_id` - (变更时重建) - 备份计划的唯一标识符。此 ID 在创建时自动生成，之后无法修改。
* `backup_method` - (必填) - 要使用的备份方法。有效值包括：
  * `logical`: 逻辑备份。
  * `physical`: 物理备份。
* `database_type` - (必填) - 要备份的数据库类型。有效值包括：
  * `MySQL`
  * `MSSQL`
  * `Oracle`
  * `MongoDB`
  * `Redis`
* `instance_class` - (必填) - 备份实例的类别。有效值包括：
  * `small`: 小型实例。
  * `large`: 大型实例。
* `backup_plan_name` - (选填) - 备份计划的名称。如果不指定，Terraform 将自动生成一个名称。
* `database_region` - (选填) - 源数据库所在的区域。例如，`cn-hangzhou`。
* `storage_region` - (选填) - 备份数据将存储的区域。它可以与 `database_region` 相同或不同。
* `instance_type` - (选填) - 数据库实例的类型。有效值包括：
  * `RDS`: 关系型数据库服务。
  * `PolarDB`: 极速数据库。
  * `DDS`: 文档数据库服务。
  * `Kvstore`: 键值存储。
  * `Other`: 其他类型。
* `from_app` - (选填) - 表示请求的来源。默认值为 `OpenAPI`。通常无需手动设置。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `backup_plan_id` - 备份计划的唯一标识符。
* `backup_plan_name` - 备份计划的名称。
