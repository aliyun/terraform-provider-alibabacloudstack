---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-instances"
description: |-
  查询阿里云MQTT实例列表
---

# alibabacloudstack_mqtt_instances

查询阿里云MQTT实例列表，可用于检索已创建的MQTT实例信息，包括实例配置、状态和访问端点等。

## 示例用法

```hcl

variable "name" {
  default = "12918"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max    = 500
  tps_send_max       = 500
  topic_capacity     = 50
  cluster            = "cluster1"
  independent_naming = "true"
  name               = "${var.name}MQ"
  remark             = "Ons_instance"
}

resource "alibabacloudstack_mqtt_instance" "default" {
  instance_name      = var.name
  remark             = "Mqtt"
  max_conn           = 1000
  max_sub            = 1000
  max_up_tps         = 1000
  max_down_tps       = 1000
  independent_naming = true
  store_instance_id  = alibabacloudstack_ons_instance.default.id
}

data "alibabacloudstack_mqtt_instances" "default" {
  name_regex = alibabacloudstack_mqtt_instance.default.instance_name
}


```

## 参数说明

以下参数支持过滤查询结果：

* `ids` (可选)：MQTT实例ID列表，用于精确匹配需要查询的实例。

* `name_regex` (可选)：实例名称正则表达式，用于模糊匹配实例名称。例如，设置为`"test.*"`将匹配所有以"test"开头的实例名称。

## 属性说明

以下属性被导出：

* `id` (字符串)：数据源ID，根据查询结果生成的唯一标识符。

* `ids` (列表)：查询到的MQTT实例ID列表。

* `instances` (列表)：查询到的MQTT实例详细信息列表，每个实例包含以下属性：
  * `create_time` (整数)：实例创建时间，Unix时间戳（毫秒）。
  * `endpoints` (映射)：实例的访问端点信息，包含内部和外部访问地址。
  * `independent_naming` (布尔)：是否使用独立命名空间。
  * `instance_id` (字符串)：MQTT实例的唯一标识符。
  * `instance_name` (字符串)：MQTT实例的名称。
  * `instance_status` (整数)：实例当前状态，5表示运行中。
  * `instance_type` (整数)：实例类型，0表示标准版。
  * `max_conn` (整数)：实例支持的最大连接数。
  * `max_down_tps` (整数)：实例支持的最大下行TPS（每秒事务数）。
  * `max_sub` (整数)：实例支持的最大订阅数。
  * `max_tps` (整数)：实例支持的最大总TPS（每秒事务数）。
  * `max_up_tps` (整数)：实例支持的最大上行TPS（每秒事务数）。
  * `namespace_rules_type` (布尔)：命名空间规则类型标识。
  * `order_id` (字符串)：关联的订单ID。
  * `sp_instance_type` (整数)：特殊实例类型标识。
  * `store_instance_id` (字符串)：关联的存储实例ID。
  * `store_type` (整数)：存储类型，1表示标准存储。

* `names` (列表)：查询到的MQTT实例名称列表。