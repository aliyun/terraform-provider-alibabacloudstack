---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_ebs_storage_sets"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-ebs-storage-sets"
description: |- 
  Provides a list of ecs ebs storage sets
---

# alibabacloudstack_ecs_ebs_storage_sets

This data source provides a list of ECS EBS storage sets in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_ebs_storage_sets" "example" {
  storage_set_name = "example-storage-set"
  zone_id          = "cn-hangzhou-e"
}

output "storages" {
  value = data.alibabacloudstack_ecs_ebs_storage_sets.example.storages
}
```

## Argument Reference
The following arguments are supported:

* `storage_set_name` - (Optional) The name of the storage set to filter the results.
* `storage_set_id` - (Optional) The ID of the storage set to filter the results.
* `zone_id` - (Optional) The ID of the zone where the storage set is located.
* `shared` - (Optional) Whether to query resources shared from other organizations. If set to true, shared resources will be included in the results.
* `output_file` - (Optional, Deprecated) The 'output_file' field has been deprecated and is scheduled for removal in version 3.19.0. To write content to a file, use the 'local_file' provider instead.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of storage set IDs.
* `names` - A list of storage set names.
* `storages` - A list of storage sets. Each element contains the following attributes:
    * `storage_set_id` - The unique identifier of the storage set.
    * `storage_set_name` - The name of the storage set.
    * `storage_set_partition_number` - The partition number of the storage set.
    * `zone_id` - The ID of the zone where the storage set is located.