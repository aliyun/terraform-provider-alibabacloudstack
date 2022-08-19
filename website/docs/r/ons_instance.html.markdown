---
subcategory: "RocketMQ"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_ons_instance"
sidebar_current: "docs-apsarastack-resource-ons-instance"
description: |-
  Provides a apsarastack ONS Instance resource.
---

# apsarastack\_ons\_instance

Provides an ONS instance resource.

## Example Usage

Basic Usage

```
resource "apsarastack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = "Ons_Apsara_instance"
  remark = "Ons Instance"
}

output "inst" {
  value = apsarastack_ons_instance.default.*
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

## Attributes Reference

The following attributes are exported:

* `id` - The `key` of the resource supplied above.
* `instance_type` - The edition of instance. 1 represents the postPaid edition, and 2 represents the platinum edition.
* `instance_status` - The status of instance. 1 represents the platinum edition instance is in deployment. 2 represents the postpaid edition instance are overdue. 5 represents the postpaid or platinum edition instance is in service. 7 represents the platinum version instance is in upgrade and the service is available.
* `release_time` - Platinum edition instance expiration time.


