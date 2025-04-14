---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_storageset"
sidebar_current: "docs-Alibabacloudstack-ecs-storageset"
description: |-
  编排云绑定服务器（Ecs）存储集资源
---

# alibabacloudstack_ecs_ebs_storage_set
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_storageset`

使用Provider配置的凭证在指定的资源集下编排云绑定服务器（Ecs）存储集资源。

## 示例用法
```
variable "name" {
	default = "tf-testAcc_storage_set4148"
}
data "alibabacloudstack_zones" "default" {}



resource "alibabacloudstack_ecs_ebs_storage_set" "default" {
  storage_set_name = "tf-testAcc_storage_set4148"
  maxpartition_number = "2"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
```

## 参数说明

支持以下参数：
  * `storage_set_name` - (必填，变更时重建) 存储集名称。
  * `maxpartition_number` - (可选，变更时重建) 存储集中分区的最大数量。此参数定义了存储集可以包含的分区上限。
  * `zone_id` - (可选，变更时重建) 可用区ID。指定存储集所在的可用区。
  * `storage_set_id` - (变更时重建) 存储集ID。用于唯一标识一个存储集资源。

## 属性说明

除了上述参数外，还导出了以下属性：
  * `storage_set_id` - 存储集ID。这是创建存储集后自动生成的唯一标识符，可用于后续操作中引用该存储集。