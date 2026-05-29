---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_kafka_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-datahub-kafka-groups"
description: |-
  Provides a list of DataHub Kafka Groups for a specified project.
---

# alibabacloudstack_datahub_kafka_groups

This data source provides a list of DataHub Kafka Groups for a specified project.

## Example Usage

```hcl
variable "name" {
  default = "tf_testacc_group5898772"
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
  project_name = alibabacloudstack_datahub_project.default.name
  comment      = "test group"
  group_name   = var.name
  topic_list = [
    "${alibabacloudstack_datahub_topic.default.0.name}",
    "${alibabacloudstack_datahub_topic.default.1.name}"
  ]
}

data "alibabacloudstack_datahub_kafka_groups" "default" {
  project_name = alibabacloudstack_datahub_project.default.name
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (String, Required) The name of the DataHub project to query.
* `ids` - (List, Optional) A list of Kafka Group IDs in the format "ProjectName:GroupName" to filter specific Kafka Groups.
* `name_regex` - (String, Optional) A regular expression to match Kafka Group names for filtering.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The data source ID, which is a hash value calculated from the list of matched Kafka Group IDs.
* `kafka_groups` - A list of matched Kafka Groups. Each element contains the following attributes:
  * `comment` - The description of the Kafka Group.
  * `create_time` - The creation time of the Kafka Group, represented as a Unix timestamp in milliseconds.
  * `creator` - The ID of the creator of the Kafka Group.
  * `group_name` - The name of the Kafka Group.
  * `last_modify_time` - The last modification time of the Kafka Group, represented as a Unix timestamp in milliseconds.
  * `topic_list` - A list of associated topic names.