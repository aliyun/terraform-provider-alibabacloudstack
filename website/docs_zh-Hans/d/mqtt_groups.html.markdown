---
subcategory: "物联网平台(连接版)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-groups"
description: |-
  提供MQTT Group列表。
---

# alibabacloudstack_mqtt_groups

该数据源用于获取专有云中可用的MQTT Group列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_mqtt_groups" "example" {
  instance_id = "MQ_INST_xxx"
  name_regex  = "^GID_test.*"
}

output "mqtt_groups" {
  value = data.alibabacloudstack_mqtt_groups.example.groups
}
```

## 参数说明

以下参数支持配置：

* `instance_id` - (必选) MQTT实例的ID。
* `ids` - (可选) 用于过滤结果的Group ID列表。
* `name_regex` - (可选) 用于按Group名称过滤的正则表达式字符串。

## 属性参考

以下属性会被导出：

* `ids` - Group ID列表，格式为 `{instance_id}:{group_id}`。
* `names` - Group名称列表。
* `groups` - MQTT Group列表。每个元素包含以下属性：
  * `id` - 资源ID，格式为 `{instance_id}:{group_id}`。
  * `group_id` - Group ID。
  * `instance_id` - MQTT实例的ID。
  * `create_time` - Group创建时间戳（毫秒）。
  * `update_time` - Group最后更新时间戳（毫秒）。
  * `independent_naming` - 是否启用独立命名。
  * `channel_name` - 通道名称。
  * `region_name` - 区域名称。