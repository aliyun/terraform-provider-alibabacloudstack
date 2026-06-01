---
subcategory: "Message Queuing Telemetry Transport"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_topics"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-topics"
description: |-
  Provides a list of MQTT Topics.
---

# alibabacloudstack_mqtt_topics

This data source provides the MQTT Topics available in ApsaraStack.

-> **NOTE:** Available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_mqtt_topics" "example" {
  store_instance_id = "MQ_INST_xxx"
  name_regex        = "^my-topic.*"
}

output "mqtt_topics" {
  value = data.alibabacloudstack_mqtt_topics.example.topics
}
```

## Argument Reference

The following arguments are supported:

* `store_instance_id` - (Required) The ID of the MQTT instance.
* `ids` - (Optional) A list of topic IDs to filter results.
* `name_regex` - (Optional) A regex string to filter topics by topic name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of topic IDs.
* `topics` - A list of MQTT Topics. Each element contains the following attributes:
  * `store_instance_id` - The ID of the MQTT instance.
  * `topic` - The topic name.
  * `remark` - The remark for the topic.
  * `order_type` - The order type. 1 indicates Professional Edition, 0 indicates Basic Edition.
  * `independent_naming` - Whether independent naming is enabled.
  * `update_time` - The last update timestamp in milliseconds.
  * `relation` - The resource relation identifier.
  * `relation_name` - The resource relation name.
  * `create_time` - The creation timestamp in milliseconds.
  * `namespace_id` - The namespace ID.
  * `unit_flag` - Whether it is a unitized deployment identifier.
  * `status_name` - The resource status name.
  * `channel_name` - The channel name.
  * `channel_id` - The channel ID.
  * `status` - The resource status code. 0 indicates running normally.