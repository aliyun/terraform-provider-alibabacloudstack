---
subcategory: "VPC"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_vswitches"
sidebar_current: "docs-apsarastack-datasource-vswitches"
description: |-
    Provides a list of VSwitch owned by an Apsarastack Cloud account.
---

# apsarastack\_vswitches

This data source provides a list of VSwitches owned by an Apsarastack Cloud account.

## Example Usage

```
resource "apsarastack_vpc" "vpc" {
  cidr_block = "172.16.0.0/16"
  name       = "${var.name}"
}

resource "apsarastack_vswitch" "vswitch" {
  name              = "${var.name}"
  cidr_block        = "172.16.0.0/24"
  vpc_id            = "${apsarastack_vpc.vpc.id}"
  availability_zone = "${var.availability_zone}"
}

data "apsarastack_vswitches" "default" {
  name_regex = "${apsarastack_vswitch.vswitch.name}"
}

output "vswitches" {
  value = data.apsarastack_vswitches.default.vswitches.*
}
```

## Argument Reference

The following arguments are supported:

* `cidr_block` - (Optional) Filter results by a specific CIDR block. For example: "172.16.0.0/12".
* `zone_id` - (Optional) The availability zone of the VSwitch.
* `name_regex` - (Optional) A regex string to filter results by name.
* `is_default` - (Optional, type: bool) Indicate whether the VSwitch is created by the system.
* `vpc_id` - (Optional) ID of the VPC that owns the VSwitch.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `ids` - (Optional) A list of VSwitch IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of VSwitch IDs.
* `names` - A list of VSwitch names.
* `vswitches` - A list of VSwitches. Each element contains the following attributes:
  * `id` - ID of the VSwitch.
  * `zone_id` - ID of the availability zone where the VSwitch is located.
  * `vpc_id` - ID of the VPC that owns the VSwitch.
  * `name` - Name of the VSwitch.
  * `instance_ids` - List of ECS instance IDs in the specified VSwitch.
  * `cidr_block` - CIDR block of the VSwitch.
  * `description` - Description of the VSwitch.
  * `is_default` - Whether the VSwitch is the default one in the region.
  * `creation_time` - Time of creation.
