---
subcategory: "物联网平台(连接版)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topics"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-topics"
description: |-
  提供MQTT主题列表。
---

# alibabacloudstack_mqtt_topics

该数据源用于获取专有云中可用的MQTT主题列表。

-> **注意:** 适用于专有云环境。

## 示例

```hcl
data "alibabacloudstack_mqtt_topics" "example" {
  store_instance_id = "MQ_INST_xxx"
  name_regex        = "^my-topic.*"
}

output "mqtt_topics" {
  value = data.alibabacloudstack_mqtt_topics.example.topics
}
```

## 参数说明

以下参数支持配置：

* `store_instance_id` - (必选) MQTT实例的ID。
* `ids` - (可选) 用于过滤结果的主题ID列表。
* `name_regex` - (可选) 用于按主题名称过滤的正则表达式字符串。

## 属性参考

以下属性会被导出：

* `ids` - 主题ID列表。
* `topics` - MQTT主题列表。每个元素包含以下属性：
  * `store_instance_id` - MQTT实例的ID。
  * `topic` - 主题名称。
  * `remark` - 主题的备注信息。
  * `order_type` - 顺序类型。1表示专业版，0表示基础版。
  * `independent_naming` - 是否启用独立命名。
  * `update_time` - 最后更新时间戳（毫秒）。
  * `relation` - 资源关系标识符。
  * `relation_name` - 资源关系名称。
  * `create_time` - 创建时间戳（毫秒）。
  * `namespace_id` - 命名空间ID。
  * `unit_flag` - 是否为单元化部署标识。
  * `status_name` - 资源状态名称。
  * `channel_name` - 通道名称。
  * `channel_id` - 通道ID。
  * `status` - 资源状态码。0表示正常运行。