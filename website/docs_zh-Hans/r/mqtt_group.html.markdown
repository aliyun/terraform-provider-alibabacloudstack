---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_group"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-mqtt_group"
description: |-
  管理阿里云MQTT Group资源
---

# alibabacloudstack_mqtt_group

使用Provider配置的凭证在指定的MQTT实例中创建和管理Group资源。

## 示例用法

### 基础用法

```hcl

variable "name" {
  default = "GID_tftestacc745"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max    = 500
  tps_send_max       = 500
  topic_capacity     = 50
  independent_naming = true
  cluster            = "cluster1"
  name               = "${var.name}ONS"
  remark             = "Ons_instance"
}

resource "alibabacloudstack_mqtt_instance" "default" {
  instance_name     = var.name
  remark            = "MqttInstance"
  max_conn          = 1000
  max_sub           = 1000
  max_up_tps        = 1000
  max_down_tps      = 1000
  store_instance_id = alibabacloudstack_ons_instance.default.id
}


resource "alibabacloudstack_mqtt_group" "default" {
  group_id    = var.name
  instance_id = alibabacloudstack_mqtt_instance.default.id
}
```

## 参数说明

支持以下参数：

* `group_id` - (必填, 变更时重建) MQTT Group ID。必须以`GID_`或`GID-`开头，用于唯一标识Group。
* `instance_id` - (必填, 变更时重建) MQTT实例ID。指定Group所属的MQTT实例。

## 属性说明

以下属性导出为资源属性：

* `id` - 资源ID，格式为`{instance_id}:{group_id}`。
* `channel_name` - 通道名称（例如：`ALIYUN`）。
* `create_time` - 创建时间，Unix时间戳（毫秒）。
* `independent_naming` - 是否启用独立命名空间（布尔值）。
* `update_time` - 最后更新时间，Unix时间戳（毫秒）。