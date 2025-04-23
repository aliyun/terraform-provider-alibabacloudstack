---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_masterslaveservergroups"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-masterslaveservergroups"
description: |- 
  Provides a list of slb masterslaveservergroups owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_masterslaveservergroups
-> **NOTE:** Alias name has: `alibabacloudstack_slb_master_slave_server_groups`

This data source provides a list of SLB master slave server groups in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_masterslaveservergroups" "sample_ds" {
  load_balancer_id = "${alibabacloudstack_slb.sample_slb.id}"
  ids              = ["group1-id", "group2-id"]
  name_regex       = "group.*"
  output_file      = "output.txt"
}

output "first_slb_server_group_id" {
  value = "${data.alibabacloudstack_slb_masterslaveservergroups.sample_ds.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of the Server Load Balancer (SLB) instance.
* `ids` - (Optional) A list of master slave server group IDs to filter results.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by master slave server group name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of SLB master slave server groups names.
* `groups` - A list of SLB master slave server groups. Each element contains the following attributes:
  * `id` - The ID of the master slave server group.
  * `name` - The name of the master slave server group.
  * `servers` - ECS instances associated with the group. Each element contains the following attributes:
    * `instance_id` - The ID of the attached ECS instance.
    * `weight` - The weight associated with the ECS instance.
    * `port` - The port used by the master slave server group.
    * `server_type` - The server type of the attached ECS instance.
