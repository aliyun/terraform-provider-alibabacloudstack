---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_topics"
sidebar_current: "docs-Alibabacloudstack-datasource-datahub-topics"
description: |-
  Provides a list of DataHub Topics.
---

# alibabacloudstack_datahub_topics

This data source provides the DataHub Topics available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_datahub_topics" "example" {
  project_name = "my-project"
  name_regex   = "^test-.*"
}

output "datahub_topics" {
  value = data.alibabacloudstack_datahub_topics.example.topics
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the DataHub project.
* `name_regex` - (Optional, ForceNew) A regex string to filter topics by topic name.
* `names` - (Optional, Computed) A list of topic names to filter results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of topic IDs.
* `names` - A list of topic names.
* `topics` - A list of DataHub Topics. Each element contains the following attributes:
  * `id` - The ID of the topic, in the format `{project_name}:{topic_name}`.
  * `name` - The name of the topic.
  * `shard_count` - The number of shards.
  * `life_cycle` - The lifecycle of the topic in days.
  * `comment` - The comment for the topic.
  * `record_type` - The record type of the topic.
  * `create_time` - The creation time of the topic.
