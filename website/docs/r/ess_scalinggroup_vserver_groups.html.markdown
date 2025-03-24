---
subcategory: "ESS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_scalinggroup_vserver_groups"
sidebar_current: "docs-Alibabacloudstack-ess-scalinggroup-vserver-groups"
description: |-
  Provides a ESS Scaling Group VServer Groups resource.
---

# alibabacloudstack_ess_scalinggroup_vserver_groups

Attaches/Detaches vserver groups to a specified scaling group.

-> **NOTE:** The load balancer of which vserver groups belongs to must be in `active` status.

-> **NOTE:** If scaling group's network type is `VPC`, the vserver groups must be in the same `VPC`.
 
-> **NOTE:** A scaling group can have at most 5 vserver groups attached by default.

-> **NOTE:** Vserver groups and the default group of loadbalancer share the same backend server quota.

-> **NOTE:** When attach vserver groups to scaling group, existing ECS instances will be added to vserver groups; Instead, ECS instances will be removed from vserver group when detach.

-> **NOTE:** Detach action will be executed before attach action.

-> **NOTE:** Vserver group is defined uniquely by `loadbalancer_id`, `vserver_group_id`, `port`.

-> **NOTE:** Modifing `weight` attribute means detach vserver group first and then, attach with new weight parameter.

## Example Usage

```hcl
resource "alibabacloudstack_ess_scalinggroup_vserver_groups" "default" {
  scaling_group_id = "your_scaling_group_id"
  vserver_groups {
    loadbalancer_id = "your_loadbalancer_id"
    vserver_attributes {
      vserver_group_id = "your_vserver_group_id"
      port             = 80
      weight           = 100
    }
  }
  vserver_groups {
    loadbalancer_id = "another_loadbalancer_id"
    vserver_attributes {
      vserver_group_id = "another_vserver_group_id"
      port             = 8080
      weight           = 200
    }
  }
  force = true
}
```

### Argument Reference
The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) - The ID of the scaling group.
* `vserver_groups` - (Required) - A set of VServer groups to be attached to the scaling group.
  * `loadbalancer_id` - (Required) - The ID of the load balancer.
  * `vserver_attributes` - (Required) - A set of VServer attributes.
    * `vserver_group_id` - (Required) - The ID of the VServer group.
    * `port` - (Required) - The port number for the VServer group.
    * `weight` - (Required) - The weight of the VServer group.
* `force` - (Optional) - Whether to force the attachment or detachment of VServer groups. Default is true.


### Attributes Reference
The following attributes are exported in addition to the arguments listed above:

* `id` - (Required, ForceNew) The ESS vserver groups attachment resource ID.