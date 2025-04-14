---
subcategory: "RDS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_db_readonly_instance"
sidebar_current: "docs-alibabacloudstack-resource-db-readonly-instance"
description: |-
  编排RDS只读实例
---

# alibabacloudstack_db_readonly_instance
使用Provider配置的凭证在指定的资源集下编排RDS只读实例。

## 参数说明

以下参数被支持：

* `engine_version` - (必填，变更时重建) 数据库版本。可选值请参考最新文档 [CreateDBInstance](https://www.alibabacloud.com/help/doc-detail/26228.htm) 中的 `EngineVersion`。
* `master_db_instance_id` - (必填) 主实例的ID。
* `instance_type` - (必填) 数据库实例类型。详情请参阅 [实例类型表](https://www.alibabacloud.com/help/doc-detail/26312.htm)。
* `instance_storage` - (必填) 用户定义的数据库实例存储空间。取值范围：对于MySQL/SQL Server HA双节点版为[5, 2000]。以5GB为单位递增。详情请参阅 [实例类型表](https://www.alibabacloud.com/help/doc-detail/26312.htm)。
* `instance_name` - (可选) 数据库实例名称。长度为2到256个字符的字符串。
* `parameters` - (可选) 在数据库实例启动后需要设置的参数集合。可用参数请参考最新文档 [查看数据库参数模板](https://www.alibabacloud.com/help/doc-detail/26284.htm)。
  * `name` - (必填) 参数名称。
  * `value` - (必填) 参数值。
* `zone_id` - (可选，变更时重建) 启动数据库实例所在的可用区。
* `vswitch_id` - (可选，变更时重建) 用于在一个VPC中启动数据库实例的虚拟交换机ID。
* `tags` - (可选) 分配给资源的标签映射。
    - Key：最多可以是64个字符长度。不能以“aliyun”、“acs:”、“http://”或“https://”开头。不能是空字符串。
    - Value：最多可以是128个字符长度。不能以“aliyun”、“acs:”、“http://”或“https://”开头。可以是空字符串。
* `db_instance_storage_type` - (必填) 实例的存储类型。有效值：
    - local_ssd：表示使用本地SSD。推荐使用此值。
    - cloud_ssd：表示使用标准SSD。
    - cloud_essd：表示使用增强型SSD(ESSD)。
    - cloud_essd2：表示使用增强型SSD(ESSD)。
    - cloud_essd3：表示使用增强型SSD(ESSD)。
* `db_instance_class` - (可选) 数据库实例类。
* `db_instance_storage` - (可选) 数据库实例存储。
* `master_instance_id` - (可选，变更时重建) 主实例的ID。
* `db_instance_description` - (可选) 数据库实例的描述。

-> **注意：** 由于数据备份和迁移，更改数据库实例类型和存储需要花费15~20分钟。请在更改前做好充分准备。

## 属性说明

以下属性将被导出：

* `id` - RDS实例ID。
* `engine` - 数据库类型。
* `port` - RDS数据库连接端口。
* `connection_string` - RDS数据库连接字符串。
* `db_instance_description` - 数据库实例的描述。
* `engine` - 数据库引擎类型。