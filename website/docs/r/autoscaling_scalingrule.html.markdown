---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingrule"
sidebar_current: "docs-Alibabacloudstack-autoscaling-scalingrule"
description: |- 
  Provides a autoscaling Scalingrule resource.
---

# alibabacloudstack_autoscaling_scalingrule
-> **NOTE:** Alias name has: `alibabacloudstack_autoscaling_rule` `alibabacloudstack_ess_scaling_rule`

Provides a autoscaling Scalingrule resource.

## Example Usage

```hcl
data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name = "${var.name}_vsw"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_ecs_securitygroup" "default" {
  name   = "${var.name}_sg"
  vpc_id = "${alibabacloudstack_vpc_vpc.default.id}"
}

resource "alibabacloudstack_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = "${alibabacloudstack_ecs_securitygroup.default.id}"
  cidr_ip           = "172.16.0.0/24"
}

data "alibabacloudstack_images" "default" {
  name_regex  = "^ubuntu_"
  most_recent = true
  owners      = "system"
}

data "alibabacloudstack_instance_types" "all" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
}

data "alibabacloudstack_instance_types" "any_n4" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

data "alibabacloudstack_instance_types" "default" {
  availability_zone = data.alibabacloudstack_zones.default.zones[0].id
  cpu_core_count       = 1
  memory_size          = 1
  instance_type_family = "ecs.n4"
  sorted_by            = "Memory"
}

locals {
  default_instance_type_id = try(element(sort(length(data.alibabacloudstack_instance_types.default.instance_types) > 0 ? data.alibabacloudstack_instance_types.default.ids : data.alibabacloudstack_instance_types.any_n4.ids), 0), sort(data.alibabacloudstack_instance_types.all.ids)[0])
}

resource "alibabacloudstack_ecs_instance" "default" {
  image_id             = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type        = "${local.default_instance_type_id}"
  system_disk_category = "${data.alibabacloudstack_zones.default.zones.0.available_disk_categories.0}"
  system_disk_size     = 20
  system_disk_name     = "test_sys_disk"
  security_groups      = [alibabacloudstack_ecs_securitygroup.default.id]
  instance_name        = "${var.name}_ecs"
  vswitch_id           = alibabacloudstack_vpc_vswitch.default.id
  zone_id    = data.alibabacloudstack_zones.default.zones[0].id
  is_outdated          = false
  lifecycle {
    ignore_changes = [
      instance_type
    ]
  }
}

variable "name" {
  default = "tf-testAccEssScalingRule-109553"
}

resource "alibabacloudstack_ess_scaling_group" "default" {
  min_size = 0
  max_size = 2
  default_cooldown = 20
  removal_policies = ["OldestInstance", "NewestInstance"]
  scaling_group_name = "${var.name}"
  vswitch_ids = ["${alibabacloudstack_vpc_vswitch.default.id}"]
}

resource "alibabacloudstack_ecs_deployment_set" "default" {
  strategy            = "Availability"
  domain              = "Default"
  granularity         = "Host"
  deployment_set_name = "example_value"
  description         = "example_value"
}

resource "alibabacloudstack_ess_scaling_configuration" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  image_id = "${data.alibabacloudstack_images.default.images.0.id}"
  instance_type = "ecs.e4.small"
  security_group_ids = [alibabacloudstack_ecs_securitygroup.default.id]
  force_delete = true
  active = true
  enable = true
  deployment_set_id = alibabacloudstack_ecs_deployment_set.default.id
}

resource "alibabacloudstack_autoscaling_scalingrule" "default" {
  scaling_group_id = "${alibabacloudstack_ess_scaling_group.default.id}"
  adjustment_type = "TotalCapacity"
  adjustment_value = "1"
  cooldown = 0
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Required, ForceNew) The ID of the scaling group.
* `adjustment_type` - (Required) Specifies the adjusting way in the scaling rule. Valid values:
  * `QuantityChangeInCapacity`: Increasing or decreasing the number of specified ECS instances.
  * `PercentChangeInCapacity`: Increasing or decreasing the specified proportion of ECS instances.
  * `TotalCapacity`: Setting the quantity of ECS instances in the current scaling group to a specified value.
* `adjustment_value` - (Required) Specifies the adjustment value in the scaling rule. The value range depends on the `adjustment_type`:
  * `QuantityChangeInCapacity`: (-500, 0] U (0, 500]
  * `PercentChangeInCapacity`: [-100, 0] U [0, 10000]
  * `TotalCapacity`: [0, 1000]
* `scaling_rule_name` - (Optional) Name shown for the scaling rule. It must be 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain numbers, underscores `_`, hyphens `-`, and decimal points `.`. If this parameter is not specified, the default value will be the `ScalingRuleId`.
* `cooldown` - (Optional) The cooldown time of the scaling rule. This parameter is applicable only to simple scaling rules. Value range: [0, 86400], in seconds. Default value is 0.
* `ari` - (Optional)  Field 'ari' is deprecated and will be removed in a future release. Please use new field 'scaling_rule_aris' instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the scaling rule.
* `ari` - The unique identifier of the scaling rule.
* `scaling_rule_aris` - The unique identifier list of the scaling rule.
* `scaling_rule_name` - Name shown for the scaling rule. It must be 2-64 characters (English or Chinese), starting with numbers, English letters or Chinese characters, and can contain numbers, underscores `_`, hyphens `-`, and decimal points `.`. If this parameter is not specified, the default value will be the `ScalingRuleId`.