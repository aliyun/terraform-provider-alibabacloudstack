---
subcategory: "RocketMQ (ONS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ons_topics"
sidebar_current: "docs-alibabacloudstack-datasource-ons-topics"
description: |-
    Provides a list of ons topics available to the user.
---

# alibabacloudstack_ons_topics

This data source provides a list of ONS Topics in an Apsara Stack Cloud account according to the specified filters.



## Example Usage

```
variable "name" {
  default = "onsInstanceName"
}

variable "topic" {
  default = "onsTopicDatasourceName"
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

data "alibabacloudstack_ons_topics" "topics_ds" {
 instance_id = alibabacloudstack_ons_topic.topic.instance_id
  output_file = "topics.txt"
}

output "first_topic_name" {
   value = data.alibabacloudstack_ons_topics.topics_ds.*
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Required) ID of the ONS Instance that owns the topics.
* `name_regex` - (Optional) A regex string to filter results by the topic name. 


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of topic names.
* `topics` - A list of topics. Each element contains the following attributes:
  * `id` - ID of the topics.
  * `instance_id` - ID of the intance to which the topic belongs.
  * `topic` - The name of the topic.
  * `owner` - The ID of the topic owner, which is the Apsara Stack Cloud UID.
  * `relation` - The relation ID. 
  * `relation_name` - The name of the relation, for example, owner, publishable, subscribable, and publishable and subscribable.
  * `message_type` - The type of the message.
  * `independent_naming` - Indicates whether namespaces are available. 
  * `create_time` - Time of creation.
  * `remark` - Remark of the topic.
  * `computed_attribute` - Description of the computed attribute.