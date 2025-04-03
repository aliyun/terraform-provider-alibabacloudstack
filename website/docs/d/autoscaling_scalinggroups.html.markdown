---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalinggroups"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalinggroups"
description: |- 
  Provides a list of autoscaling scalinggroups owned by an AlibabaCloudStack account.
---

# alibabacloudstack_autoscaling_scalinggroups
-> **NOTE:** Alias name has: `alibabacloudstack_ess_scaling_groups`

This data source provides a list of autoscaling scalinggroups in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_scalinggroups" "example" {
  name_regex = "scaling_group_name"
  ids        = ["scaling_group_id1", "scaling_group_id2"]
}

output "first_scaling_group_id" {
  value = "${data.alibabacloudstack_autoscaling_scalinggroups.example.groups.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter resulting scaling groups by name.
* `ids` - (Optional) A list of scaling group IDs.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of scaling group IDs.
* `names` - A list of scaling group names.
* `groups` - A list of scaling groups. Each element contains the following attributes:
  * `id` - ID of the scaling group.
  * `name` - Name of the scaling group.
  * `active_scaling_configuration` - The active scaling configuration for the scaling group.
  * `region_id` - Region ID the scaling group belongs to.
  * `min_size` - The minimum number of ECS instances in the scaling group. When the number of ECS instances in the scaling group is less than this value, Auto Scaling automatically creates ECS instances.
  * `max_size` - The maximum number of ECS instances in the scaling group. When the number of ECS instances in the scaling group exceeds this value, Auto Scaling automatically removes ECS instances. The value range of `max_size` depends on the usage quota and can be viewed in the quota center. For example, if the quota allows up to 2000 instances per scaling group, `max_size` ranges from 0 to 2000.
  * `cooldown_time` - Default cooldown time of the scaling group, in seconds. During this period, no scaling activities are triggered.
  * `removal_policies` - Removal policy used to select the ECS instance to remove from the scaling group when reducing the number of instances.
  * `load_balancer_ids` - A list of SLB (Server Load Balancer) instance IDs that the ECS instances in the scaling group are attached to.
  * `db_instance_ids` - A list of RDS (Relational Database Service) instance IDs that the ECS instances in the scaling group are attached to.
  * `vswitch_ids` - A list of VSwitch IDs in which the ECS instances are launched.
  * `lifecycle_state` - Lifecycle state of the scaling group. Possible values include `Active`, `Inactive`, etc.
  * `total_capacity` - Total number of ECS instances in the scaling group.
  * `active_capacity` - Number of ECS instances that have successfully joined the scaling group and are running normally.
  * `pending_capacity` - Number of ECS instances that are in the process of joining the scaling group but have not completed their configurations yet.
  * `removing_capacity` - Number of ECS instances that are being removed from the scaling group.
  * `creation_time` - Creation time of the scaling group, in ISO 8601 format.