---
subcategory: "EasyAI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_evpc_evpc"
description: |- 
  Provides an EVPC resource.
---

# alibabacloudstack_evpc_evpc

Provides an EVPC (Elastic Virtual Private Cloud) resource.

## Example Usage

```hcl
resource "alibabacloudstack_evpc_evpc" "default" {
  evpc_name   = "my-evpc"
  description = "My EVPC instance"
}
```

## Argument Reference

The following arguments are supported:

* `evpc_name` - (Required) The name of the EVPC. It must be 2 to 128 characters in length.
* `description` - (Optional) The description of the EVPC. It must be 0 to 256 characters in length.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the EVPC.
* `evpc_id` - The ID of the EVPC.
* `evpc_name` - The name of the EVPC.
* `status` - The status of the EVPC.
* `description` - The description of the EVPC.
* `cidr` - The CIDR block of the EVPC.
* `tenant_id` - The tenant ID of the EVPC.
* `department` - The department ID of the EVPC.
* `department_name` - The department name of the EVPC.
* `region_id` - The region ID of the EVPC.
* `resource_group` - The resource group ID of the EVPC.
* `resource_group_name` - The resource group name of the EVPC.
* `cluster_id` - The cluster ID of the EVPC.
* `ascm_create_user` - The user who created the EVPC.
* `create_time` - The creation time of the EVPC.
* `update_time` - The last update time of the EVPC.
