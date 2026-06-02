---
subcategory: "Bare-Metal Management Service (BMS)"
layout: "alibabacloudstack"
page_title: "Data Source: alibabacloudstack_bms_machinetypes"
subcategory: "Bare-Metal Management Service (BMS)"
description: |-
    Query the available machine types of BMS.
---

# alibabacloudstack_bms_machinetypes

This data source provides the BMS Machine Types available to the user.

## Example Usage

```hcl
data "alibabacloudstack_bms_machinetypes" "default" {
  name_regex = "PG"
}

output "first_machine_type" {
  value = data.alibabacloudstack_bms_machinetypes.default.machinetypes.0.name
}
```

## Arguments Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string used to filter results by machine type name.
* `arch_regex` - (Optional) A regex string used to filter results by CPU architecture.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of machine type IDs.
* `machinetypes` - A list of machine type objects. Each element contains the following attributes:
  * `id` - The ID of the machine type.
  * `name` - The name of the machine type.
  * `description` - The description of the machine type.
  * `deploy_type` - The deployment type. Valid values: `bmcp`, `bmcp_managed`, `bmcp_no_clone`, `base`, `ehpc`, `ehpc_managed`, `aspeed`.
  * `manufacturer` - The manufacturer of the machine.
  * `cpu_arch` - The CPU architecture (e.g., x86, ARM).
  * `cpu_model` - The CPU model.
  * `cpu_manufacturer` - The CPU manufacturer.
  * `cpu_number` - The number of CPU cores.
  * `memory` - The memory size in GB.
  * `disk` - The disk size in GB.
  * `disk_type` - The disk type.
  * `gpu` - The GPU model.
  * `gpu_num` - The number of GPUs.
  * `gpu_manufacturer` - The GPU manufacturer.
  * `video_memory` - The video memory size in GB.
  * `tflops_fp32` - The TFLOPS FP32 performance metric.
  * `network_card_type` - The network card type.
  * `network_card_num` - The number of network cards.
  * `rated_power` - The rated power consumption.
  * `specification` - The specification details.
  * `unit_num` - The unit number.
  * `create_time` - The creation time.
  * `update_time` - The last update time.
