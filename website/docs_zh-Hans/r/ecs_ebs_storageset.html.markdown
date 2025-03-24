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

## 参数参考

支持以下参数：
  * `storage_set_name` - (必填，变更时重建) 存储集名称。
  * `maxpartition_number` - (可选，变更时重建) 存储集中分区的最大数量。
  * `zone_id` - (可选，变更时重建) 可用区ID。
  * `storage_set_id` - (变更时重建) 存储集ID。

## 属性参考

除了上述参数外，还导出了以下属性：
  * `storage_set_id` - 存储集ID