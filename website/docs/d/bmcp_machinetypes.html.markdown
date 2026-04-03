---
subcategory: "BMCP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bmcp_machinetypes"
sidebar_current: "docs-Alibabacloudstack-datasource-bmcp-machinetypes"
description: |-
  Provides a list of BMCP (Bare Metal Compute Platform) machine types owned by an alibabacloudstack account.
---

# alibabacloudstack_bmcp_machinetypes

This data source provides a list of BMCP machine types in an AlibabacloudStack account according to the specified filters.

## Example Usage

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

### Combine filters

```hcl
data "alibabacloudstack_bmcp_machinetypes" "combined" {
  name_regex = "PG"
  arch_regex = "x86"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter resulting machine types by name.
* `arch_regex` - (Optional, ForceNew) A regex string to filter resulting machine types by CPU architecture.

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
  * `memory` - The memory size.
  * `disk` - The disk size.
  * `disk_type` - The disk type.
  * `gpu` - The GPU information.
  * `gpu_num` - The number of GPUs.
  * `gpu_manufacturer` - The GPU manufacturer.
  * `video_memory` - The video memory size.
  * `tflops_fp32` - The TFLOPS FP32 performance.
  * `network_card_type` - The network card type.
  * `network_card_num` - The number of network cards.
  * `rated_power` - The rated power.
  * `specification` - The specification of the machine type.
  * `unit_num` - The unit number.
  * `create_time` - The creation time.
  * `update_time` - The update time.
