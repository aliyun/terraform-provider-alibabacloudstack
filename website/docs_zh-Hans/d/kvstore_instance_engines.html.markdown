---
subcategory: "Redis And Memcache (KVStore)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_kvstore_instance_engines"
sidebar_current: "docs-alibabacloudstack-datasource-kvstore-instance-engines"
description: |-
   查询KVStore实例引擎信息
---

# alibabacloudstack_kvstore_instance_engines

根据指定过滤条件列出当前凭证权限可以访问的KVStore实例引擎资源信息。

## 示例用法

```tf
data "alibabacloudstack_zones" "resources" {
  available_resource_creation = "KVStore"
}

data "alibabacloudstack_kvstore_instance_engines" "resources" {
  zone_id              = "${data.alibabacloudstack_zones.resources.zones.0.id}"
  instance_charge_type = "PrePaid"
  engine               = "Redis"
  engine_version       = "5.0"
  output_file          = "./engines.txt"
}

output "first_kvstore_instance_class" {
  value = "${data.alibabacloudstack_kvstore_instance_engines.resources.instance_engines.0.engine}"
}
```

## 参数说明

以下是支持的参数：

* `zone_id` - (必填) 启动KVStore实例的可用区。
* `instance_charge_type` - (可选) 通过付费类型过滤结果。有效值：`PrePaid`（预付费）和 `PostPaid`（按量付费）。默认为 `PrePaid`。
* `engine` - (可选) 数据库类型。选项为 `Redis` 和 `Memcache`。默认为 `Redis`。
* `engine_version` - (可选) 用户所需的数据库版本。Redis 的可选项可以参考最新文档 [详细信息](https://www.alibabacloud.com/help/doc-detail/60873.htm) 中的 `EngineVersion`。对于 `Memcache`，该值应为空。
* `output_file` - (可选) 将查询结果保存到本地文件的路径。

## 属性说明

除了上述参数外，还导出以下属性：

* `instance_engines` - KVStore 可用实例引擎的列表。每个元素包含以下属性：
    * `zone_id` - 启动 KVStore 实例的可用区 ID。
    * `engine` - 数据库类型，例如 `Redis` 或 `Memcache`。
    * `engine_version` - KVStore 实例的版本号，例如 Redis 的 `5.0` 或其他支持的版本。