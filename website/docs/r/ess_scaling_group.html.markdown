---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ess_scaling_group"
sidebar_current: "docs-alibabacloudstack-resource-ess-scaling-group"
description: |-
  Provides a ESS scaling group resource.
---

# alibabacloudstack_ess_scaling_group

Provides a ESS scaling group resource which is a collection of ECS instances with the same application scenarios.

It defines the maximum and minimum numbers of ECS instances in the group, and their associated Server Load Balancer instances, RDS instances, and other attributes.

-> **NOTE:** You can launch an ESS scaling group for a VPC network via specifying parameter `vswitch_ids`.

## Example Usage

```
variable "name" {
  default = "essscalinggroupconfig"
}

data "alibabacloudstack_zones" "default" {
  available_disk_category     = "cloud_efficiency"
  available_resource_creation = "VSwitch"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  cpu_core_count    = 2
  memory_size       = 4
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_18.*64"
  most_recent = true
  owners      = "system"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_security_group" "default" {
  name   = "${var.name}"
  vpc_id = "${alibabacloudstack_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

resource "alibabacloudstack_vswitch" "default2" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.1.0/24"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}-bar"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size           = 1
  max_size           = 1
  scaling_group_name = "${var.name}"
  default_cooldown   = 20
  vswitch_ids        = ["${alibabacloudstack_vswitch.default.id}", "${alibabacloudstack_vswitch.default2.id}"]
  removal_policies   = ["OldestInstance", "NewestInstance"]
}
```

## Argument Reference

The following arguments are supported:

* `min_size` - (Required) Minimum number of ECS instances in the scaling group. Value range: [0, 100].
* `max_size` - (Required) Maximum number of ECS instances in the scaling group. Value range: [0, 100].
* `scaling_group_name` - (Optional) Name shown for the scaling group, which must contain 2-40 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain numbers, underscores `_`, hyphens `-`, and decimal points `.`. If this parameter is not specified, the default value is ScalingGroupId.
* `default_cooldown` - (Optional) Default cool-down time (in seconds) of the scaling group. Value range: [0, 86400]. The default value is 300s.
* `vswitch_ids` - (Optional) List of virtual switch IDs in which the ecs instances to be launched.
* `removal_policies` - (Optional) RemovalPolicy is used to select the ECS instances you want to remove from the scaling group when multiple candidates for removal exist. Optional values:
    - OldestInstance: removes the first ECS instance attached to the scaling group.
    - NewestInstance: removes the first ECS instance attached to the scaling group.
    - OldestScalingConfiguration: removes the ECS instance with the oldest scaling configuration.
    - Default values: OldestScalingConfiguration and OldestInstance. You can enter up to two removal policies.
* `db_instance_ids` - (Optional) If an RDS instance is specified in the scaling group, the scaling group automatically attaches the Intranet IP addresses of its ECS instances to the RDS access whitelist.
    - The specified RDS instance must be in running status.
    - The specified RDS instance’s whitelist must have room for more IP addresses.
* `loadbalancer_ids` - (Optional) If a Server Load Balancer instance is specified in the scaling group, the scaling group automatically attaches its ECS instances to the Server Load Balancer instance.
    - The Server Load Balancer instance must be enabled.
    - At least one listener must be configured for each Server Load Balancer and it HealthCheck must be on. Otherwise, creation will fail (it may be useful to add a `depends_on` argument
      targeting your `alibabacloudstack_slb_listener` in order to make sure the listener with its HealthCheck configuration is ready before creating your scaling group).
    - The Server Load Balancer instance attached with VPC-type ECS instances cannot be attached to the scaling group.
    - The default weight of an ECS instance attached to the Server Load Balancer instance is 50.
* `multi_az_policy` - (Optional, ForceNew) Multi-AZ scaling group ECS instance expansion and contraction strategy. PRIORITY, BALANCE or COST_OPTIMIZED
-> **NOTE:** When detach loadbalancers, instances in group will be remove from loadbalancer's `Default Server Group`; On the contrary, When attach loadbalancers, instances in group will be added to loadbalancer's `Default Server Group`.

-> **NOTE:** When detach dbInstances, private ip of instances in group will be remove from dbInstance's `WhiteList`; On the contrary, When attach dbInstances, private ip of instances in group will be added to dbInstance's `WhiteList`.


## Attributes Reference

The following attributes are exported:

* `id` - The scaling group ID.
* `min_size` - The minimum number of ECS instances.
* `max_size` - The maximum number of ECS instances.
* `scaling_group_name` - The name of the scaling group.
* `default_cooldown` - The default cool-down of the scaling group.
* `removal_policies` - The removal policy used to select the ECS instance to remove from the scaling group.
* `db_instance_ids` - The db instances id which the ECS instance attached to.
* `loadbalancer_ids` - The slb instances id which the ECS instance attached to.
* `vswitch_ids` - The vswitches id in which the ECS instance launched.
