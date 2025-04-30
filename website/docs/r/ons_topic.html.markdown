---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_topic"
sidebar_current: "docs-alibabacloudstack-resource-ons-topic"
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
* `topic` - (Required) Name of the topic. Two topics on a single instance cannot have the same name and the name cannot start with 'GID' or 'CID'. The length cannot exceed 64 characters.
* `message_type` - (Required) The type of the message.
* `remark` - (Required) This attribute is a concise description of topic. The length cannot exceed 128.
* `perm` - (Optional) This attribute is used to set the read-write mode for the topic.

## Attributes Reference

The following attributes are exported:

* `id` - Topic and InstanceID of the ONS Topic. The value is in format `Topic:InstanceID`.