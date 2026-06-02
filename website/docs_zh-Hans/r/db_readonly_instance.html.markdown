---
subcategory: "云数据库 RDS 版"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_readonly_instance"
sidebar_current: "docs-Alibabacloudstack-resource-db-readonly-instance"
description: |-
  编排RDS只读实例
---

# alibabacloudstack_db_readonly_instance
使用Provider配置的凭证在指定的资源集下编排RDS只读实例。

## 参数说明

以下参数被支持：

* `engine_version` - (必填，变更时强制重建) 数据库版本。可选值请参考最新文档 [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) 中的 `EngineVersion`。
* `db_instance_storage_type` - (必填，变更时强制重建) 实例的存储类型。有效值：`local_ssd`、`cloud_ssd`、`cloud_essd`、`cloud_essd2`、`cloud_essd3`、`cloud_pperf`、`cloud_sperf`。
* `master_db_instance_id` - (可选，变更时强制重建，已废弃) 主实例的ID。请使用 `master_instance_id` 替代。
* `master_instance_id` - (可选，变更时强制重建) 主实例的ID。
* `instance_type` - (可选，已废弃) 数据库实例类型。详情请参阅 [实例类型表](https://www.alibabacloud.com/help/doc-detail/26312.htm)。请使用 `db_instance_class` 替代。`instance_type` 和 `db_instance_class` 至少需要指定一个。
* `db_instance_class` - (可选) 数据库实例规格。`instance_type` 和 `db_instance_class` 至少需要指定一个。
* `instance_storage` - (可选，已废弃) 用户定义的数据库实例存储空间。详情请参阅 [实例类型表](https://www.alibabacloud.com/help/doc-detail/26312.htm)。请使用 `db_instance_storage` 替代。`instance_storage` 和 `db_instance_storage` 至少需要指定一个。
* `db_instance_storage` - (可选) 数据库实例存储空间。`instance_storage` 和 `db_instance_storage` 至少需要指定一个。
* `instance_name` - (可选，已废弃) 数据库实例名称。长度为2到256个字符。请使用 `db_instance_description` 替代。
* `db_instance_description` - (可选) 数据库实例的描述。长度为2到256个字符。
* `zone_id` - (可选，变更时强制重建) 启动数据库实例所在的可用区。
* `vswitch_id` - (可选，变更时强制重建) 用于在一个VPC中启动数据库实例的虚拟交换机ID。
* `parameters` - (可选) 在数据库实例启动后需要设置的参数集合。可用参数请参考最新文档 [查看数据库参数模板](https://www.alibabacloud.com/help/doc-detail/26284.htm)。
    * `name` - (必填) 参数名称。
    * `value` - (必填) 参数值。
* `tags` - (可选) 分配给资源的标签映射。
    - Key：最多可以是64个字符长度。不能以"aliyun"、"acs:"、"http://"或"https://"开头。不能是空字符串。
    - Value：最多可以是128个字符长度。不能以"aliyun"、"acs:"、"http://"或"https://"开头。可以是空字符串。
* `force_restart` - (可选) 参数变更时是否强制重启实例。默认为 `false`。

-> **注意：** 由于数据备份和迁移，更改数据库实例类型和存储需要花费15~20分钟。请在更改前做好充分准备。

## 属性说明

以下属性将被导出：

* `id` - RDS实例ID。
* `engine` - 数据库类型。
* `engine_version` - 数据库引擎版本。
* `port` - RDS数据库连接端口。
* `connection_string` - RDS数据库连接字符串。
* `zone_id` - 实例所在的可用区。
* `vswitch_id` - 实例所在的虚拟交换机ID。
* `master_instance_id` - 主实例ID。
* `db_instance_class` - 实例规格。
* `db_instance_storage` - 实例存储空间。
* `db_instance_storage_type` - 实例存储类型。
* `db_instance_description` - 数据库实例的描述。