---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_storageset"
sidebar_current: "docs-Alibabacloudstack-ecs-storageset"
description: |-
  Provides a ecs Storageset resource.
---

# alibabacloudstack_ecs_storageset
-> **NOTE:** Alias name has: `alibabacloudstack_ecs_ebs_storage_set`

Provides a ecs Storageset resource.

## Example Usage
```
variable "name" {
	default = "tf-testAcc_storage_set4148"
}
data "alibabacloudstack_zones" "default" {}



resource "alibabacloudstack_ecs_ebs_storage_set" "default" {
  storage_set_name = "tf-testAcc_storage_set4148"
  maxpartition_number = "2"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}
```

## Argument Reference

The following arguments are supported:
  * `storage_set_name` - (Required, ForceNew) - storage set name
  * `maxpartition_number` - (Optional, ForceNew) - The maximum number of partitions in the storage set.
  * `zone_id` - (Optional, ForceNew) - zone id
  * `storage_set_id` - (ForceNew) - storage set id

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
  * `storage_set_id` - storage set id
