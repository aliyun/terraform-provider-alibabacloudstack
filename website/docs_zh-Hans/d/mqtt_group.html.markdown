---
subcategory: "物联网平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_group"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-group"
description: |-
  用于创建和管理MQTT Group
---

# alibabacloudstack_mqtt_group

用于创建和管理MQTT Group资源。

## 示例用法

```hcl

variable "name" {
  default = "testacc-16879"
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

resource "alibabacloudstack_mqtt_group" "default" {
  group_id    = "GID_${var.name}"
  instance_id = alibabacloudstack_mqtt_instance.default.id
}

data "alibabacloudstack_mqtt_groups" "default" {
  instance_id = alibabacloudstack_mqtt_group.default.instance_id
  name_regex  = "GID_testacc.*$"
}

```

## 参数说明
以下参数支持配置：

* `group_id` (必填, 变更时重建)：MQTT Group ID，格式为"GID_xxx"。
* `instance_id` (必填, 变更时重建)：MQTT实例ID。

## 属性说明
以下属性被导出：

* `id` (字符串)：资源ID，格式为"{instance_id}:{group_id}"。
* `channel_name` (字符串)：渠道名称。
* `create_time` (整数)：Group创建时间，Unix时间戳格式。
* `independent_naming` (布尔)：是否启用独立命名空间。
* `region_name` (字符串)：区域名称。
* `update_time` (整数)：Group最后更新时间，Unix时间戳格式。