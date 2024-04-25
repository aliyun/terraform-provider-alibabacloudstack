---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_instance_types"
sidebar_current: "docs-alibabacloudstack-datasource-instance-types"
description: |-
    Provides a list of ECS Instance Types to be used by the alibabacloudstack_instance resource.
---

# alibabacloudstack\_instance\_types

This data source provides the ECS instance types of ApsarStack.

~> **NOTE:** By default, only the upgraded instance types are returned. If you want to get outdated instance types, you must set `is_outdated` to true.

~> **NOTE:** If one instance type is sold out, it will not be exported.

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
* `gpu_amount` - (Optional) The GPU amount of an instance type.
* `gpu_spec` - (Optional) The GPU spec of an instance type.
* `instance_type_family` - (Optional) Filter the results based on their family name. For example: 'ecs.n4'.
* `eni_amount` - (Optional) Filter the result whose network interface number is no more than `eni_amount`.
* `kubernetes_node_role` - (Optional) Filter the result which is used to create a kubernetes cluster Optional Values: `Master` and `Worker`.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of image types. Each element contains the following attributes:
  * `id` - ID of the instance type.
  * `cpu_core_count` - Number of CPU cores.
  * `memory_size` - Size of memory, measured in GB.
  * `family` - The instance type family.
  * `availability_zones` - List of availability zones that support the instance type.
  * `gpu` - The GPU attribution of an instance type:
    * `amount` - The amount of GPU of an instance type.
    * `category` - The category of GPU of an instance type.
  * `burstable_instance` - The burstable instance attribution:
    * `initial_credit` - The initial CPU credit of a burstable instance.
    * `baseline_credit` - The compute performance benchmark CPU credit of a burstable instance.
  * `eni_amount` - The maximum number of network interfaces that an instance type can be attached to.
  * `local_storage` - Local storage of an instance type:
    * `capacity` - The capacity of a local storage in GB.
    * `amount` - The number of local storage devices that an instance has been attached to.
    * `category` - The category of local storage that an instance has been attached to.
