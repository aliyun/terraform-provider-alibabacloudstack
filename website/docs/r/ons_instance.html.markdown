---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_instance"
sidebar_current: "docs-alibabacloudstack-resource-ons-instance"
description: |-
  Provides a alibabacloudstack ONS Instance resource.
---

# alibabacloudstack_ons_instance

Provides an ONS instance resource.

## Example Usage

Basic Usage

```
resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "Ons_Apsara_instance"
  remark = "Ons Instance"
}

output "inst" {
  value = alibabacloudstack_ons_instance.default.*
}
```

## Argument Reference

The following arguments are supported:


* `name` - (Required)Two instances on a single account in the same region cannot have the same name. The length must be 3 to 64 characters. Chinese characters, English letters digits and hyphen are allowed.
* `tps_receive_max` - (Required)This attribute is used to set the message receiving transactions per second (TPS) of the topic during a certain period of time.
* `tps_send_max` - (Required)This attribute is used to set the message sending transactions per second (TPS) of the topic during a certain period of time.
* `topic_capacity` - (Required)This attribute is used to set the topic capacity.
* `independent_naming` - (Required)This attribute is used to define an independent name or not. It takes only bool value.
* `cluster` - (Required)This attribute is a used to add cluster name.
* `remark` - (Optional)This attribute is a concise description of instance. The length cannot exceed 128.
* `instance_type` - (Required)  This attribute specifies the type of the instance.
* `instance_status` - (Required)  This attribute specifies the status of the instance.

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `instance_type` - The edition of instance. 1 represents the postPaid edition, and 2 represents the platinum edition.
* `instance_status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.
* `create_time` - The create time of the reousrce.
* `name` -  This attribute indicates the name of the instance.
* `tps_receive_max` -  This attribute indicates the maximum receive TPS of the instance.
* `tps_send_max` -  This attribute indicates the maximum send TPS of the instance.
* `topic_capacity` -  This attribute indicates the topic capacity of the instance.
* `independent_naming` -  This attribute indicates whether the instance has independent naming.
* `cluster` -  This attribute indicates the cluster associated with the instance.
* `remark` -  This attribute indicates the remark of the instance.