---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_vservergroups"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-vservergroups"
description: |- 
  Provides a list of slb vservergroups owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_vservergroups
-> **NOTE:** Alias name has: `alibabacloudstack_slb_server_groups`

This data source provides a list of SLB VServer groups in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_vservergroups" "sample_ds" {
  load_balancer_id = "${alibabacloudstack_slb.sample_slb.id}"
  ids             = ["vsg-12345678", "vsg-abcdefg"]
  name_regex      = "^group-.*"

  output_file = "slb_vservergroups_output.txt"
}

output "first_slb_vserver_group_id" {
  value = data.alibabacloudstack_slb_vservergroups.sample_ds.slb_server_groups[0].id
}
```

## Argument Reference

The following arguments are supported:

* `load_balancer_id` - (Required) The ID of the SLB instance.
* `ids` - (Optional) A list of SLB VServer group IDs to filter results.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by VServer group name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB VServer group IDs.
* `names` - A list of SLB VServer group names.
* `slb_server_groups` - A list of SLB VServer groups. Each element contains the following attributes:
  * `id` - VServer group ID.
  * `name` - VServer group name.
  * `servers` - ECS instances associated with the group. Each element contains the following attributes:
    * `instance_id` - ID of the attached ECS instance.
    * `port` - Port number used by the back-end server.
    * `weight` - Weight associated with the ECS instance.