---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_machinetypes"
description: |-
  Provides a list of BMCP (Bare Metal Compute Platform) machine types owned by an alibabacloudstack account.
---

# alibabacloudstack_bmcp_machinetypes

This data source provides a list of BMCP machine types in an AlibabacloudStack account according to the specified filters.

## Example Usage

### Query all machine types

```hcl
data "alibabacloudstack_bmcp_machinetypes" "default" {}
```

### Filter by name regex

```hcl
data "alibabacloudstack_bmcp_machinetypes" "name_filtered" {
  name_regex = "PG"
}
```

### Filter by architecture regex

```hcl
data "alibabacloudstack_bmcp_machinetypes" "arch_filtered" {
  arch_regex = "x86"
}
```

### Filter by minimum standard instance count

```hcl
data "alibabacloudstack_bmcp_machinetypes" "min_count_filtered" {
  min_standard_instance_count = 1
}
```

### Filter by maximum standard instance count

```hcl
data "alibabacloudstack_bmcp_machinetypes" "max_count_filtered" {
  max_standard_instance_count = 100
}
```

### Filter by standard instance count range

```hcl
data "alibabacloudstack_bmcp_machinetypes" "range_filtered" {
  min_standard_instance_count = 1
  max_standard_instance_count = 100
}
```

### Combine filters

```hcl
data "alibabacloudstack_bmcp_machinetypes" "combined" {
  name_regex                  = "PG"
  arch_regex                  = "x86"
  min_standard_instance_count = 1
}
```

### Dynamic filtering with pre-query

```hcl
# First query all machine types to get real data
data "alibabacloudstack_bmcp_machinetypes" "all" {
}

# Use dynamic values from real data for filtering
data "alibabacloudstack_bmcp_machinetypes" "filtered" {
  name_regex = data.alibabacloudstack_bmcp_machinetypes.all.machinetypes.0.name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting machine types by name.
* `arch_regex` - (Optional, ForceNew) A regex string to filter resulting machine types by CPU architecture.
* `min_standard_instance_count` - (Optional, ForceNew) The minimum standard instance count filter. Only machine types with standard instance count greater than or equal to this value will be returned.
* `max_standard_instance_count` - (Optional, ForceNew) The maximum standard instance count filter. Only machine types with standard instance count less than or equal to this value will be returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of machine type IDs.
* `machinetypes` - A list of machine types. Each element contains the following attributes:
  * `id` - The ID of the machine type.
  * `name` - The name of the machine type.
  * `description` - The description of the machine type.
  * `deploy_type` - The deploy type of the machine type.
  * `manufacturer` - The manufacturer of the machine type.
  * `cpu_arch` - The CPU architecture of the machine type.
  * `cpu_model` - The CPU model of the machine type.
  * `cpu_manufacturer` - The CPU manufacturer of the machine type.
  * `cpu_number` - The number of CPUs.
  * `memory` - The memory size (in GB).
  * `disk` - The disk size.
  * `disk_type` - The disk type.
  * `gpu` - The GPU model.
  * `gpu_num` - The number of GPUs.
  * `gpu_manufacturer` - The GPU manufacturer.
  * `video_memory` - The video memory size (in GB).
  * `tflops_fp32` - The TFLOPS FP32 performance.
  * `network_card_type` - The network card type.
  * `network_card_num` - The number of network cards.
  * `rated_power` - The rated power.
  * `specification` - The specification of the machine type.
  * `unit_num` - The unit number.
  * `standard_instance_count` - The count of standard instances available for this machine type.
  * `default_fmin` - The default FMin value.
  * `create_time` - The creation time.
  * `update_time` - The update time.
