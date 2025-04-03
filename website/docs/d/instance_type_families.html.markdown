---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance_type_families"
sidebar_current: "docs-alibabacloudstack-datasource-instance-type-families"
description: |-
    Provides a list of ECS Instance Type Families to be used by the alibabacloudstack_instance resource.
---

# alibabacloudstack_instance_type_families

This data source provides the ECS instance type families of AlibabacloudStack.

## Example Usage

```
data "alibabacloudstack_instance_type_families" "default" {
  
}

output "first_instance_type_family_id" {
  value = "${data.alibabacloudstack_instance_type_families.default.instance_type_families.0.id}"
}

output "instance_ids" {
  value = "${data.alibabacloudstack_instance_type_families.default.ids}"
}
```

## Argument Reference

The following arguments are supported:

* `zone_id` - (Optional, ForceNew) The Zone to launch the instance.
* `generation` - (Optional) The generation of the instance type family,

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type family IDs.
* `families` - A list of image type families. Each element contains the following attributes:
  * `id` - ID of the instance type family.
  * `generation` - The generation of the instance type family.
  * `zone_ids` - A list of Zone to launch the instance.
* `families.zone_ids` - A list of Zone to launch the instance.