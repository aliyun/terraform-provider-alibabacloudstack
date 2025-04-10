---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_ebs_storage_sets"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-ebs-storage-sets"
description: |- 
  查询云服务器EBS存储集
---

# alibabacloudstack_ecs_ebs_storage_sets

根据指定过滤条件列出当前凭证权限可以访问的云服务器EBS存储集的列表。

## 示例用法

```hcl
data "alibabacloudstack_ecs_ebs_storage_sets" "example" {
  storage_set_name = "example-storage-set"
  zone_id          = "cn-hangzhou-e"
}

output "storages" {
  value = data.alibabacloudstack_ecs_ebs_storage_sets.example.storages
}
```

## 参数说明

以下参数被支持：

* `storage_set_name` - (可选) 用于过滤结果的存储集名称。通过此参数可以筛选出特定名称的存储集。
* `maxpartition_number` - (可选) 存储集的最大分区数量。通过此参数可以限制返回的存储集的分区数量范围。
* `zone_id` - (可选) 存储集所在的可用区ID。通过此参数可以筛选出特定可用区内的存储集。
* `storage_set_id` - (可选) 用于过滤结果的存储集ID。通过此参数可以直接定位到特定的存储集。

## 属性说明

以下属性被导出：

* `ids` - 存储集的ID列表。每个ID唯一标识一个存储集。
* `names` - 存储集的名称列表。每个名称对应一个存储集。
* `storages` - 存储集列表。每个元素包含以下属性：
    * `storage_set_id` - 存储集的唯一标识符，用于在系统中唯一识别该存储集。
    * `storage_set_name` - 存储集的名称，用于描述该存储集。
    * `storage_set_partition_number` - 存储集的分区数量，表示该存储集中的分区数目。