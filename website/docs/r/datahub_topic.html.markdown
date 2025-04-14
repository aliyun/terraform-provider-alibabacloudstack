---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_topic"
sidebar_current: "docs-Alibabacloudstack-datahub-topic"
description: |- 
  Provides a datahub Topic resource.
---

# alibabacloudstack_datahub_topic

Provides a datahub Topic resource.

## Example Usage

### Basic Usage

- **BLOB Topic**

```hcl
resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = "tf_testacc_datahub_project"
}

resource "alibabacloudstack_datahub_topic" "blob_example" {
  name         = "tf_testacc_datahub_blob_topic"
  project_name = alibabacloudstack_datahub_project.default.name
  record_type  = "BLOB"
  shard_count  = 3
  life_cycle   = 7
  comment      = "created by terraform"
}
```

- **TUPLE Topic**

```hcl
resource "alibabacloudstack_datahub_project" "default" {
    comment = "test"
    name = "tf_testacc_datahub_project"
}

resource "alibabacloudstack_datahub_topic" "tuple_example" {
  name         = "tf_testacc_datahub_tuple_topic"
  project_name = alibabacloudstack_datahub_project.default.name
  record_type  = "TUPLE"
  record_schema = {
    bigint_field    = "BIGINT"
    timestamp_field = "TIMESTAMP"
    string_field    = "STRING"
    double_field    = "DOUBLE"
    boolean_field   = "BOOLEAN"
  }
  shard_count = 3
  life_cycle  = 7
  comment     = "created by terraform"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required, ForceNew) The name of the DataHub project that this topic belongs to. It is case-insensitive and cannot exceed 128 characters.
* `name` - (Required, ForceNew) The name of the DataHub topic. Its length is limited to 1-128 characters and only letters, digits, and underscores (`_`) are allowed. It is case-insensitive.
* `shard_count` - (Optional, ForceNew) The number of shards this topic contains. The permitted range of values is [1, 10]. The default value is 1.
* `life_cycle` - (Optional) The retention period for the topic's data in days. The permitted range of values is [1, 7]. The default value is 3.
* `record_type` - (Optional, ForceNew) The type of the topic. It must be one of `BLOB` or `TUPLE`. For `BLOB` topics, data will be organized as binary and encoded by BASE64. For `TUPLE` topics, data has a fixed schema. The default value is `TUPLE` with a schema `{STRING}`.
* `record_schema` - (Optional, ForceNew) Schema of this topic, required only for `TUPLE` topics. Supported data types (case-insensitive) are:
  - `BIGINT`
  - `STRING`
  - `BOOLEAN`
  - `DOUBLE`
  - `TIMESTAMP`
* `comment` - (Optional) Comment for the DataHub topic. It cannot exceed 255 characters.
* `create_time` - (Optional, Computed) The creation time of the DataHub topic.
* `last_modify_time` - (Optional, Computed) The last modification time of the DataHub topic. Initially, it is the same as the `create_time`.

**Note:** Currently, the `life_cycle` field cannot be modified and will be supported in the future.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DataHub topic. It is composed of the project name and the topic name, formatted as `<project_name>:<name>`.
* `create_time` - The creation time of the DataHub topic. It is a human-readable string rather than a 64-bit UTC timestamp.
* `last_modify_time` - The last modification time of the DataHub topic. Initially, it is the same as the `create_time`. It is also a human-readable string rather than a 64-bit UTC timestamp.