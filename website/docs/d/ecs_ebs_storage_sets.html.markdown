---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_ebs_storage_sets"
sidebar_current: "docs-Alibabacloudstack-datasource-ecs-ebs-storage-sets"
description: |- 
  Provides a list of ecs ebs storage sets
---

# alibabacloudstack_ecs_disks
-> **NOTE:** Alias name has: `alibabacloudstack_disks`

This data source provides a list of ECS EBS storage sets in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_ecs_storageset" "example" {
  storage_set_name = "example-storage-set"
  zone_id          = "cn-hangzhou-e"
}

output "storages" {
  value = data.alibabacloudstack_ecs_storageset.example.storages
}
```

## Argument Reference
The following arguments are supported:

* `storage_set_name` - (Optional) The name of the storage set to filter the results.
* `maxpartition_number` - (Optional) The maximum partition number of the storage set.
* `zone_id` - (Optional) The ID of the zone where the storage set is located.
* `storage_set_id` - (Optional) The ID of the storage set to filter the results.

## Attributes Reference
The following attributes are exported:

* `ids` - A list of IDs of the storage sets.
* `names` - A list of names of the storage sets.
* `storages` - A list of storage sets. Each element contains the following attributes:
    * `storage_set_id` - The unique identifier of the storage set.
    * `storage_set_name` - The name of the storage set.
    * `storage_set_partition_number` - The partition number of the storage set.