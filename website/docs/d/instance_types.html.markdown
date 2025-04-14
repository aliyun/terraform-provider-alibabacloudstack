---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance_types"
sidebar_current: "docs-alibabacloudstack-datasource-instance-types"
description: |-
    Provides a list of ECS Instance Types to be used by the alibabacloudstack_instance resource.
---

# alibabacloudstack_instance_types

This data source provides the ECS instance types of ApsarStack.

~> **NOTE:** By default, only the upgraded instance types are returned. If you want to get outdated instance types, you must set `is_outdated` to true.

## Example Usage

```
# Declare the data source
data "alibabacloudstack_instance_types" "types_ds" {
  cpu_core_count = 1
  memory_size    = 2
}

output "instance_types"{
  value=data.alibabacloudstack_instance_types.types_ds.*
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The zone where instance types are supported.
* `cpu_core_count` - (Optional) Filter the results to a specific number of cpu cores.
* `cpu_type` - (Optional) Filter the results to a specific cpu type. Optional Values: `intel`, `hg`, `kp`, `ft`.
* `memory_size` - (Optional) Filter the results to a specific memory size in GB.
* `sorted_by` - (Optional, ForceNew) Sort mode, valid values: `CPU`, `Memory`, `Price`.
* `instance_type_family` - (Optional) Filter the results based on their family name. For example: 'ecs.n4'.
* `eni_amount` - (Optional) Filter the result whose network interface number is no more than `eni_amount`.
* `kubernetes_node_role` - (Optional) Filter the result which is used to create a kubernetes cluster Optional Values: `Master` and `Worker`.
* `ids` - (Optional) A list of instance type IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of image types. Each element contains the following attributes:
  * `id` - ID of the instance type.
  * `price` - The price of instance type.
  * `cpu_core_count` - Number of CPU cores.
  * `memory_size` - Size of memory, measured in GB.
  * `family` - The instance type family.
  * `availability_zones` - List of availability zones that support the instance type.
  * `burstable_instance` - The burstable instance attribution:
    * `initial_credit` - The initial CPU credit of a burstable instance.
    * `baseline_credit` - The compute performance benchmark CPU credit of a burstable instance.
  * `eni_amount` - The maximum number of network interfaces that an instance type can be attached to.
  * `local_storage` - Local storage of an instance type:
    * `capacity` - The capacity of a local storage in GB.
    * `amount` - The number of local storage devices that an instance has been attached to.
    * `category` - The category of local storage that an instance has been attached to.
  * `cpu_type` - Type of CPU.
  * `instance_type_family` - The instance type family.