---
subcategory: "云原生多模数据库 Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instances"
description: |-
  为用户提供Lindorm实例列表
---

# alibabacloudstack\_lindorm\_instances

该数据源为用户提供可用的Lindorm实例列表。

## 示例用法

### 基本用法

```terraform
data "alibabacloudstack_lindorm_instances" "example" {
  ids = ["i-12345678"]
}

output "first_instance_id" {
  value = data.alibabacloudstack_lindorm_instances.example.instances.0.id
}
```

### 通过描述正则表达式筛选

```terraform
data "alibabacloudstack_lindorm_instances" "example" {
  description_regex = "my-instance"
}

output "instance_ids" {
  value = data.alibabacloudstack_lindorm_instances.example.ids
}
```

## 参数参考

以下参数被支持：

* `ids` - (可选) 实例ID列表。
* `name_regex` - (可选, 已废弃) 用于通过实例描述筛选结果的正则表达式。字段'name_regex'已废弃，并将在未来版本中移除，请使用新字段'description_regex'替代。
* `description_regex` - (可选) 用于通过实例描述筛选结果的正则表达式。

## 属性参考

以下属性会被导出：

* `ids` - 实例ID列表。
* `instances` - Lindorm实例列表。每个元素包含以下属性：
  * `id` - 实例ID。
  * `instance_id` - 实例ID。
  * `cpu_brand` - 实例使用的CPU品牌。
  * `instance_storage` - 实例的存储容量。
  * `zone_id` - 部署实例的可用区ID。
  * `create_time` - 实例创建时间。
  * `ascm_create_user` - 创建实例的用户。
  * `instance_alias` - 实例别名。
  * `network_type` - 实例的网络类型。
  * `service_type` - 实例的服务类型。
  * `engine_type` - 实例的引擎类型。
  * `ali_uid` - 阿里云账户的UID。