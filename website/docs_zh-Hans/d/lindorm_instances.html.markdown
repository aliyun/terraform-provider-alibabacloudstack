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
  * `service_type` - `GetLindormInstanceList` 接口返回的原始服务类型值，例如 `"lindorm"` 或 `"lindorm_dual_zone"`。
  * `engine_type` - 归一化后的引擎类型，与 `alibabacloudstack_lindorm_instance` 的 `engine_type` 参数含义一致。多可用区实例（`lindorm_dual_zone`）统一报告为 `"lindorm"`。
  * `ali_uid` - 阿里云账户的UID。
  * `topology_type` - 由 `service_type` 派生的拓扑类型。有效值：`"1azone"`（单可用区）或 `"multizone"`（双可用区）。
  * `arch_version` - 架构版本。注意：`GetLindormInstanceList` 响应中并不包含此字段，因此该属性通常为空。
  * `standby_zone_id` - 多可用区实例的备可用区 ID（来源于接口返回的 `Idc2ZoneId`）。单可用区实例该字段为空。
  * `arbiter_zone_id` - 仲裁可用区 ID。列表接口当前不返回该字段，保留以备后续扩展。
  * `standby_vswitch_id` - 备可用区交换机 ID（来源于接口返回的 `Idc2VSwitchId`）。单可用区实例或接口未返回该字段时为空。
  * `arbiter_vswitch_id` - 仲裁可用区交换机 ID。列表接口当前不返回该字段，保留以备后续扩展。