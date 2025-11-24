---
subcategory: "MQTT"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_instance"
sidebar_current: "docs-Alibabacloudstack-datasource-mqtt-instance"
description: |-
    Provides a MQTT instance resource
---
# alibabacloudstack_mqtt_instance

Creates and manages MQTT instances in the specified region using the credentials configured in the Provider.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tftestacc729"
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
  max_conn          = "1000"
  max_up_tps        = "500"
  max_down_tps      = "1000"
  remark            = "test"
  store_instance_id = alibabacloudstack_ons_instance.default.id
  max_sub           = "10000"
  instance_name     = var.name
  cluster_name      = "mqtt4Private"
}
```

## Argument Reference

The following arguments are supported:

* `instance_name` - (Required) The name of the MQTT instance. The name must be 1 to 128 characters in length and cannot start with `http://` or `https://`.
* `cluster_name` - (Required, Forces new resource when changed) The name of the cluster. Default value is "mqtt4Private", indicating the use of a private cluster.
* `independent_naming` - (Optional) Whether to use an independent namespace. Default value is `true`, indicating that independent namespace is enabled.
* `max_conn` - (Optional) The maximum number of connections. Valid values: 1000 to 100000.
* `max_down_tps` - (Optional) The maximum downstream TPS (transactions per second). Valid values: 1000 to 100000.
* `max_sub` - (Optional) The maximum number of subscriptions. Default value is 1000000.
* `max_up_tps` - (Optional) The maximum upstream TPS (transactions per second). Valid values: 500 to 20000.
* `remark` - (Optional) The remark for the instance. The description must be 1 to 256 characters in length.
* `store_instance_id` - (Optional) The storage instance ID. Used to bind an existing Message Queue for RocketMQ instance.
* `store_type` - (Optional) The storage type. Default value is 1, indicating standard storage type.

## Attributes Reference

The following attributes are exported from the API response:

* `instance_id` - The ID of the MQTT instance.
* `endpoints` - The endpoint information of the instance, including internal and external access addresses, such as `tcpInternalEndpoint`, `tcpInternetEndpoint`, etc.
* `namespace_id` - The namespace ID, which is usually the same as the instance ID.