---
subcategory: "物联网平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topic"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-topic"
description: |-
  Provides a MQTT topic resource.
---

# alibabacloudstack_mqtt_topic

管理阿里云MQTT服务中的主题资源。

## 示例用法

```hcl

variable "name" {
  default = "tf_testAcck9386"
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

resource "alibabacloudstack_mqtt_topic" "default" {
  topic             = var.name
  order_type        = 1
  remark            = "test"
  store_instance_id = alibabacloudstack_mqtt_instance.default.store_instance_id
}


data "alibabacloudstack_mqtt_topics" "default" {
  name_regex        = alibabacloudstack_mqtt_topic.default.topic
  store_instance_id = alibabacloudstack_mqtt_topic.default.store_instance_id
}
```

## 参数说明

以下参数支持入参配置：

* `store_instance_id` (字符串, 必填, 变更时重建) - MQTT的存储实例ID。

* `topic` (字符串, 必填, 变更时重建) - 主题名称，用于标识消息发布的主题。

* `order_type` (整数, 可选, 变更时重建) - 订单类型，1表示专业版，0表示基础版。

* `remark` (字符串, 可选, 变更时重建) - 主题备注信息，用于描述主题用途。

> `注意`：该资源不支持修改操作，任何参数变更都会导致资源重建。

## 属性说明

以下属性导出为只读属性：

* `id` (字符串) - 资源ID，格式为`{instance_id:topic}`。

* `channel_id` (整数) - 通道ID，标识消息传输通道。

* `channel_name` (字符串) - 通道名称，如"ALIYUN"。

* `create_time` (整数) - 资源创建时间戳（毫秒）。

* `independent_naming` (布尔) - 是否启用独立命名空间。

* `namespace_id` (字符串) - 命名空间ID，与实例ID相同。

* `relation` (整数) - 资源关系标识。

* `relation_name` (字符串) - 资源关系名称，如"Owner"。

* `status` (整数) - 资源状态码，0表示正常运行。

* `status_name` (字符串) - 资源状态名称，如"Running"。

* `unit_flag` (布尔) - 是否为单元化部署标识。

* `update_time` (整数) - 资源最后更新时间戳（毫秒）。