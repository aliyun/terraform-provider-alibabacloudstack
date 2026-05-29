---
subcategory: "分布式关系型数据库"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_drds_rds_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-drds-rds-instances"
description: |-
  查询云原生分布式数据库（Drds）实例关联的私有定制 RDS 实例列表
---

# alibabacloudstack_drds_rds_instances

根据指定过滤条件，列出与 DRDS 实例关联的所有私有定制 RDS 实例。

> **注意:** 该资源也可以使用以下别名引用：
> - `apsarastack_drds_rds_instances`

## 示例用法

```hcl
data "alibabacloudstack_drds_rds_instances" "example" {
  drds_instance_id = "drds-xxxxxxxxxxxx"
  ids              = ["rm-xxxxxxxxxxxx"]
}

output "rds_instance_ids" {
  value = data.alibabacloudstack_drds_rds_instances.example.ids
}
```

## 参数说明

支持以下参数：

* `drds_instance_id` - (必选) DRDS 实例 ID。您可以调用 `DescribeDrdsInstances` API 查询 DRDS 实例 ID。
* `ids` - (可选) RDS 实例 ID 列表。通过该参数可以过滤查询结果，仅返回指定的 RDS 实例。

## 属性说明

除了上述参数外，还导出以下属性：

* `ids` - RDS 实例 ID 列表。
* `rds_instances` - DRDS RDS 实例列表。每个实例包含以下属性：
  * `rds_instance_id` - RDS 实例 ID。
  * `db_instance_storage` - RDS 实例的存储空间，单位：GB。
  * `create_time` - RDS 实例的过期时间。
