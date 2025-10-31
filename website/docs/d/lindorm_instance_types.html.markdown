---
subcategory: "Lindorm"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_lindorm_instance_types"
description: |-
  Provides a list of Lindorm Instance Types to the user.
---

# alibabacloudstack\_lindorm\_instance\_types

This data source provides Lindorm instance types available to the user.


## Example Usage

### Basic Usage

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
}

output "first_instance_type_id" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types.0.id
}
```

### Filter by Engine Type

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  engine_type = "lindorm"
}

output "lindorm_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

### Filter by CPU and Memory

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  cpu    = 4
  memory = 8
}

output "filtered_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

### Sort by CPU

```terraform
data "alibabacloudstack_lindorm_instance_types" "example" {
  sorted_by = "CPU"
}

output "sorted_instance_types" {
  value = data.alibabacloudstack_lindorm_instance_types.example.instance_types
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of instance type IDs.
* `engine_type` - (Optional) The type of the engine. Valid values: `lindorm`, `tsdb`, `solr`, `lts`.
* `cpu` - (Optional) The number of CPUs.
* `memory` - (Optional) The memory size in GB.
* `sorted_by` - (Optional) Sorting field. Valid values: `CPU`, `Memory`.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of instance type IDs.
* `instance_types` - A list of Lindorm instance types. Each element contains the following attributes:
  * `id` - The ID of the instance type.
  * `cpu` - The number of CPUs.
  * `memory` - The memory size in GB.
  * `rate` - The rate.
  * `name` - The name of the instance type.
