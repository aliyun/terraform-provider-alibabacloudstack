---
subcategory: "Message Queuing Telemetry Transport"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-instances"
description: |-
  Query the list of Alibaba Cloud MQTT instances
---

# alibabacloudstack_mqtt_instances

Query the list of Alibaba Cloud MQTT instances, which can be used to retrieve information about created MQTT instances, including instance configuration, status, and access endpoints.

## Example Usage

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

## Argument Reference

The following arguments support filtering query results:

* `ids` (Optional): A list of MQTT instance IDs used for exact matching of instances to be queried.

* `name_regex` (Optional): A regular expression for instance names, used for fuzzy matching of instance names. For example, setting it to `"test.*"` will match all instance names starting with "test".

## Attributes Reference

The following attributes are exported:

* `id` (String): The data source ID, a unique identifier generated based on the query results.

* `ids` (List): A list of MQTT instance IDs retrieved from the query.

* `instances` (List): A list of detailed information for the retrieved MQTT instances, each instance contains the following attributes:
  * `create_time` (Integer): The creation time of the instance, Unix timestamp in milliseconds.
  * `endpoints` (Map): Access endpoint information of the instance, including internal and external access addresses.
  * `independent_naming` (Boolean): Whether independent namespace is used.
  * `instance_id` (String): The unique identifier of the MQTT instance.
  * `instance_name` (String): The name of the MQTT instance.
  * `instance_status` (Integer): The current status of the instance, 5 indicates running.
  * `instance_type` (Integer): The type of the instance, 0 indicates standard edition.
  * `max_conn` (Integer): The maximum number of connections supported by the instance.
  * `max_down_tps` (Integer): The maximum downstream TPS (transactions per second) supported by the instance.
  * `max_sub` (Integer): The maximum number of subscriptions supported by the instance.
  * `max_tps` (Integer): The maximum total TPS (transactions per second) supported by the instance.
  * `max_up_tps` (Integer): The maximum upstream TPS (transactions per second) supported by the instance.
  * `namespace_rules_type` (Boolean): Namespace rule type identifier.
  * `order_id` (String): The associated order ID.
  * `sp_instance_type` (Integer): Special instance type identifier.
  * `store_instance_id` (String): The associated storage instance ID.
  * `store_type` (Integer): Storage type, 1 indicates standard storage.

* `names` (List): A list of MQTT instance names retrieved from the query.