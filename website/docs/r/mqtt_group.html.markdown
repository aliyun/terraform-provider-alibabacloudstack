---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_group"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-mqtt_group"
description: |-
  Manages Alibaba Cloud MQTT Group resources.
---

# alibabacloudstack_mqtt_group

Creates and manages Group resources in the specified MQTT instance using credentials configured by the Provider.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `group_id` - (Required, ForceNew) MQTT Group ID. Must start with `GID_` or `GID-`, used to uniquely identify the Group.
* `instance_id` - (Required, ForceNew) MQTT instance ID. Specifies the MQTT instance to which the Group belongs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - Resource ID in the format `{instance_id}:{group_id}`.
* `channel_name` - Channel name (e.g., `ALIYUN`).
* `create_time` - Creation time, Unix timestamp in milliseconds.
* `independent_naming` - Whether independent naming is enabled (boolean).
* `update_time` - Last update time, Unix timestamp in milliseconds.