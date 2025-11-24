---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topic"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-topic"
description: |-
  Provides a MQTT topic resource.
---

# alibabacloudstack_mqtt_topic

Manages MQTT topic resources in Alibaba Cloud MQTT service.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `store_instance_id` - (String, Required, ForceNew) The ID of the MQTT instance, in the format `MQ_INST_xxx`.
* `topic` - (String, Required, ForceNew) The topic name, used to identify the message publishing topic.
* `order_type` - (Integer, Optional, ForceNew) The order type. 1 indicates Professional Edition, 0 indicates Basic Edition.
* `remark` - (String, Optional, ForceNew) The remark for the topic, used to describe the topic's purpose.

> **Note**: This resource does not support modification operations. Any parameter change will cause the resource to be recreated.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - (String) The resource ID, in the format `{instance_id:topic}`.
* `channel_id` - (Integer) The channel ID, identifying the message transmission channel.
* `channel_name` - (String) The channel name, such as "ALIYUN".
* `create_time` - (Integer) The resource creation timestamp in milliseconds.
* `independent_naming` - (Boolean) Whether independent naming is enabled.
* `namespace_id` - (String) The namespace ID, which is the same as the instance ID.
* `relation` - (Integer) The resource relation identifier.
* `relation_name` - (String) The resource relation name, such as "Owner".
* `status` - (Integer) The resource status code. 0 indicates running normally.
* `status_name` - (String) The resource status name, such as "Running".
* `unit_flag` - (Boolean) Whether it is a unitized deployment identifier.
* `update_time` - (Integer) The last update timestamp of the resource in milliseconds.