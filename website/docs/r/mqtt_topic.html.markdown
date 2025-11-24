---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topic"
sidebar_current: "docs-Alibabacloudstack-resource-mqtt-topic"
description: |-
  Manage Alibaba Cloud MQTT topics
---

# alibabacloudstack_mqtt_topic

Manages Alibaba Cloud MQTT topic resources. This resource is used to create and manage topics in a specified MQTT instance.

## Example Usage

### Basic Usage

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

## Argument Reference

The following arguments are supported:

* `store_instance_id` - (Required, ForceNew) The MQTT store instance ID.
* `topic` - (Required, ForceNew) The MQTT topic name. The topic name must be 3-64 characters in length and can contain letters, digits, hyphens (-), underscores (_), and forward slashes (/).
* `remark` - (Required, ForceNew) The remark for the topic. The remark must be 1-128 characters in length.
* `order_type` - (Required, ForceNew) The topic type. Valid values:
  * `0`: Normal topic
  * `1`: Ordered topic

## Attributes Reference

The following attributes are exported from the API response:

* `id` - The resource ID in the format of `{storeInstanceId:topic}`.
* `channel_id` - The channel ID.
* `channel_name` - The channel name.
* `create_time` - The creation time of the topic (timestamp in milliseconds).
* `independent_naming` - Indicates whether the independent namespace identifier is enabled.
* `namespace_id` - The namespace ID, which is the same as store_instance_id.
* `relation` - The relation identifier.
* `relation_name` - The relation name.
* `status` - The status identifier of the topic.
* `status_name` - The status name of the topic.
* `unit_flag` - The unit flag.
* `update_time` - The last update time of the topic (timestamp in milliseconds).