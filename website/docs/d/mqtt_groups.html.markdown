---
subcategory: "Message Queuing Telemetry Transport"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-groups"
description: |-
  Query AlibabacloudStack MQTT Group resources.
---

# alibabacloudstack_mqtt_groups

Query AlibabacloudStack MQTT Group resources.

## Example Usage

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

## Argument Reference

The following arguments are supported:

* `group_id` (Required, ForceNew): MQTT Group ID, in the format "GID_xxx".
* `instance_id` (Required, ForceNew): MQTT instance ID. 

## Attributes Reference

The following attributes are exported:

* `id` (String): Resource ID, in the format "{instance_id}:{group_id}".
* `channel_name` (String): Channel name.
* `create_time` (Integer): Creation time of the Group, in Unix timestamp format.
* `independent_naming` (Boolean): Whether independent namespace is enabled.
* `region_name` (String): Region name.
* `update_time` (Integer): Last update time of the Group, in Unix timestamp format.