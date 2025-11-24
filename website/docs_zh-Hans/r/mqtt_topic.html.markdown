---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topic"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-topic"
description: |-
  管理阿里云MQTT主题
---

# alibabacloudstack_mqtt_topic

管理阿里云MQTT主题资源。该资源用于在指定的MQTT实例中创建和管理主题。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "tftestacc612"
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
  order_type        = "1"
  remark            = "test"
  store_instance_id = alibabacloudstack_mqtt_instance.default.store_instance_id
}
```

## 参数说明

支持以下参数：

* `store_instance_id` - (必填, 变更时重建) MQTT的存储实例ID。格式为`MQ_INST_xxxxxxx_xxxxxx`。
* `topic` - (必填, 变更时重建) MQTT主题名称。主题名称长度限制为3-64个字符，只能包含字母、数字、短横线(-)、下划线(_)和斜杠(/)。
* `remark` - (必填, 变更时重建) 主题备注信息。备注信息长度限制为1-128个字符。
* `order_type` - (必填, 变更时重建) 主题类型。有效值：
  * `0`：普通主题
  * `1`：有序主题

## 属性说明

以下属性会从API响应中导出：

* `id` - 资源ID，格式为`{InstanceId:topic}`。
* `channel_id` - 通道ID。
* `channel_name` - 通道名称。
* `create_time` - 主题创建时间（时间戳，单位毫秒）。
* `independent_naming` - 是否启用独立命名空间标识。
* `namespace_id` - 命名空间ID，与store_instance_id相同。
* `relation` - 关系标识。
* `relation_name` - 关系名称。
* `status` - 主题状态标识。
* `status_name` - 主题状态名称。
* `unit_flag` - 单元标志。
* `update_time` - 主题最后更新时间（时间戳，单位毫秒）。