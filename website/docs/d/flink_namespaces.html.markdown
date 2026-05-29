---
subcategory: "Realtime Compute for Apache Flink"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_flink_namespaces"
sidebar_current: "docs-alibabacloudstack-datasource-flink-namespaces"
description: |-
  Provides a list of Flink namespaces.
---

# alibabacloudstack_flink_namespaces

This data source provides a list Flink namespaces on Alibabacloudstack Cloud.



## Example Usage

```
# Declare the data source
data "alibabacloudstack_flink_namespaces" "my_namespaces" {
  name_regex  = "my-namespace"
}

```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by namespace name.
* `ids` - (Optional) A string list to filter results by namespace name.
* `owner_id` - (Optional) Filter results by namespace Owner Id.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of matched Flink Registry namespaces. Its element is a namespace name.
* `names` - A list of namespace names.
* `namespaces` - A list of matched Flink Registry namespaces. Each element contains the following attributes:
  * `name` - Name of Flink Registry namespace. 
  * `cu` - Integer. Guaranteed Resources for Cpu.
  * `cpu_type` - The CPU type of the resource. Valid values: `intel`.
  * `owner_uid` - Owner ID for this Namesapce.