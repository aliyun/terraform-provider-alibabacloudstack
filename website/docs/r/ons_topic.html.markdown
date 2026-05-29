---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_topic"
sidebar_current: "docs-Alibabacloudstack-resource-ons-topic"
description: |-
  Provides a alibabacloudstack ONS Topic resource.
---

# alibabacloudstack_ons_topic

Provides an ONS topic resource.


## Example Usage

Basic Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicName"
}

resource "alibabacloudstack_ons_instance" "default" {
  tps_receive_max = 500
  tps_send_max = 500
  topic_capacity = 50
  cluster = "cluster1"
  independent_naming = "true"
  name = var.name
  remark = "Ons Instance"
}

resource "alibabacloudstack_ons_topic" "default" {
  topic = var.topic
  instance_id = alibabacloudstack_ons_instance.default.id
  message_type = 0
  remark = "dafault_ons_topic_remark"
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional) ID of the ONS Instance that owns the topics. 
* `topic` - (Required) Name of the topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length must be between 1 and 128 characters.
* `message_type` - (Required) The type of the message. Modifying this argument forces the creation of a new resource.
* `remark` - (Required) A concise description of the topic. The length must be between 1 and 128 characters.

## Attributes Reference

The following attributes are exported:

* `id` - Topic and InstanceID of the ONS Topic. The value is in format `Topic:InstanceID`.

## Import

ONS Topic can be imported using the topic and instance_id, e.g.

```
$ terraform import alibabacloudstack_ons_topic.example tf-topic12345:mq-instance-abc
```