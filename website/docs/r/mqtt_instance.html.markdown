---
subcategory: "Message Queuing Telemetry Transport"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_mqtt_instance"
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
* `cluster_name` - (Optional, ForceNew) The name of the cluster. Default value is "mqtt4Private", indicating the use of a private cluster. Changing this parameter forces a new resource to be created.
* `independent_naming` - (Optional) Whether to use an independent namespace. Default value is `true`, indicating that independent namespace is enabled.
* `max_conn` - (Optional) The maximum number of connections. Valid values: 1000 to 100000.
* `max_down_tps` - (Optional) The maximum downstream TPS (transactions per second). Valid values: 1000 to 100000.
* `max_sub` - (Optional) The maximum number of subscriptions. Default value is 1000000.
* `max_up_tps` - (Optional) The maximum upstream TPS (transactions per second). Valid values: 500 to 20000.
* `remark` - (Optional) The remark for the instance. The description must be 1 to 256 characters in length.
* `store_instance_id` - (Optional, Computed) The storage instance ID. Used to bind an existing Message Queue for RocketMQ instance. This attribute is read from the API response.
* `store_type` - (Optional) The storage type. Valid values: 1. Default value is 1, indicating standard storage type.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the resource.
* `instance_id` - The ID of the MQTT instance.
* `endpoints` - The endpoint information of the instance, including internal and external access addresses, such as `tcpInternalEndpoint`, `tcpInternetEndpoint`, etc. This is a map of strings.
* `namespace_id` - The namespace ID, which is usually the same as the instance ID.
* `store_instance_id` - The storage instance ID bound to this MQTT instance.
* `instance_name` - The name of the MQTT instance.
* `max_conn` - The maximum number of connections.
* `max_sub` - The maximum number of subscriptions.
* `max_up_tps` - The maximum upstream TPS.
* `max_down_tps` - The maximum downstream TPS.
* `independent_naming` - Whether independent namespace is enabled.
* `remark` - The remark for the instance.
* `store_type` - The storage type.
* `cluster_name` - The cluster name. Note: This attribute is not returned by the Read API and will be ignored during import verification.

## Import

MQTT Instance can be imported using the instance ID, e.g.

```
$ terraform import alibabacloudstack_mqtt_instance.example mqtt-instance-id-12345
```