---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_kafka_group"
sidebar_current: "docs-Alibabacloudstack-resource-datahub-kafka-group"
description: |-
  Manages DataHub Kafka group resources
---

# alibabacloudstack_datahub_kafka_group

Manages DataHub Kafka group resources, used for creating and managing Kafka consumer groups in DataHub projects.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tf_testacc_datahub_group21719"
}
resource "alibabacloudstack_datahub_project" "default" {
  comment = "test"
  name    = var.name
}

resource "alibabacloudstack_datahub_topic" "default" {
  count        = 2
  name         = "${var.name}_${count.index}"
  comment      = "test"
  record_type  = "BLOB"
  project_name = alibabacloudstack_datahub_project.default.name
}

resource "alibabacloudstack_datahub_kafka_group" "default" {
  topic_list = [
    "${alibabacloudstack_datahub_topic.default.0.name}"
  ]
  project_name = alibabacloudstack_datahub_project.default.name
  comment      = "test group"
  group_name   = "tf_testacc_datahub_group21719"
}
```

## Argument Reference

The following arguments are supported:

* `comment` - (Required) The description of the Kafka group. Cannot be modified after creation.
* `group_name` - (Required, ForceNew) The name of the Kafka group. Must comply with DataHub naming conventions, with a length limit of 1-128 characters.
* `project_name` - (Required, ForceNew) The name of the DataHub project. Must already exist and comply with DataHub project naming conventions.
* `topic_list` - (Optional) The list of topics bound to this Kafka group. At least one topic must be specified, and topic names must comply with DataHub topic naming conventions. When updating, the `UpdateTopicsForKafkaGroup` API is called to change the binding relationship.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in the format `{ProjectName}:{GroupName}`.
* `create_time` - The creation time of the Kafka group, in Unix timestamp (milliseconds).