---
subcategory: "OTS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ots_instances"
sidebar_current: "docs-Alibabacloudstack-datasource-ots-instances"
description: |- 
  Provides a list of ots instances owned by an Alibabacloudstack account.
---

# alibabacloudstack_ots_instances

This data source provides a list of ots instances in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_ots_instances" "instances_ds" {
  ids        = ["instance1", "instance2"]
  name_regex = "^my-instance-.*$"

  tags = {
    Environment = "Production"
    Owner      = "JohnDoe"
  }

  output_file = "instances.txt"
}

output "first_instance_id" {
  value = data.alibabacloudstack_ots_instances.instances_ds.instances.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of instance IDs. If specified, the data source will only return instances with these IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by instance name. This allows you to match instance names using regular expressions.
* `tags` - (Optional) A map of tags assigned to the instance. It must be in the format:
  ```terraform
  tags = {
    tagKey1 = "tagValue1",
    tagKey2 = "tagValue2"
  }
  ```

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance IDs matching the specified filters.
* `names` - A list of instance names matching the specified filters.
* `instances` - A list of instances. Each element contains the following attributes:
  * `id` - ID of the instance.
  * `name` - Instance name.
  * `status` - Instance status. Possible values: `Running`, `Disabled`, `Deleting`.
  * `write_capacity` - Reserve write throughput. The unit is CU (Capacity Unit). Only high-performance instances have this return value.
  * `read_capacity` - Reserve read throughput. The unit is CU (Capacity Unit). Only high-performance instances have this return value.
  * `cluster_type` - The cluster type of the instance. Possible values: `SSD`, `HYBRID`.
  * `create_time` - The create time of the instance.
  * `user_id` - The user ID associated with the instance.
  * `network` - The network type of the instance. Possible values: `NORMAL`, `VPC`, `VPC_CONSOLE`.
  * `description` - The description of the instance.
  * `entity_quota` - The instance quota indicating the maximum number of tables that can be created within this instance.
  * `tags` - The tags assigned to the instance.