---
subcategory: "Server Load Balancer (SLB)"
layout: "apsarastack"
page_title: "Apsarastack: apsarastack_slb_server_groups"
sidebar_current: "docs-apsarastack-datasource-slb-server_groups"
description: |-
    Provides a list of VServer groups related to a server load balancer to the user.
---

# apsarastack\_slb_server_groups

This data source provides the VServer groups related to a server load balancer.

## Example Usage

```
data "apsarastack_slb_server_groups" "sample_ds" {
  load_balancer_id = "${apsarastack_slb.sample_slb.id}"
}

output "first_slb_server_group_id" {
  value = "${data.apsarastack_slb_server_groups.sample_ds.slb_server_groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) ID of the SLB.
* `ids` - (Optional) A list of VServer group IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by VServer group name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB VServer groups IDs.
* `names` - A list of SLB VServer groups names.
* `slb_server_groups` - A list of SLB VServer groups. Each element contains the following attributes:
  * `id` - VServer group ID.
  * `name` - VServer group name.
  * `servers` - ECS instances associated to the group. Each element contains the following attributes:
    * `instance_id` - ID of the attached ECS instance.
    * `port`- Port number used by the back-end server.
    * `weight` - Weight associated to the ECS instance.
