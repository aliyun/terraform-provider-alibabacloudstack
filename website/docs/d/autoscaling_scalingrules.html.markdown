---
subcategory: "Auto Scaling (ESS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_scalingrules"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-scalingrules"
description: |- 
  Provides a list of autoscaling scalingrules owned by an AlibabaCloudStack account.
---

# alibabacloudstack_autoscaling_scalingrules
-> **NOTE:** Alias name has: `alibabacloudstack_ess_scaling_rules`

This data source provides a list of autoscaling scaling rules in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_scalingrules" "example" {
  scaling_group_id = "sg-1234567890abcdef"
  ids              = ["sr-abc123", "sr-def456"]
  name_regex       = "rule-name-*"
  type             = "SimpleScalingConfiguration"

  output_file = "scaling_rules_output.json"
}

output "first_scaling_rule_id" {
  value = data.alibabacloudstack_autoscaling_scalingrules.example.rules[0].id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional, ForceNew) The ID of the scaling group. This is used to filter scaling rules that belong to a specific scaling group.
* `name_regex` - (Optional, ForceNew) A regex string to apply as a filter on the names of the scaling rules. This allows you to retrieve only those rules whose names match the given regular expression.
* `ids` - (Optional) A list of scaling rule IDs. This can be used to filter the results to include only specific scaling rules by their IDs.
* `type` - (Optional, ForceNew) The type of the scaling rule. Valid values include:
  - `SimpleScalingConfiguration`: For simple scaling rules.
  - `TargetTrackingConfiguration`: For target tracking scaling rules.
  - `StepScalingConfiguration`: For step scaling rules.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of the names of all scaling rules that match the specified filters.
* `rules` - A list of scaling rules. Each element contains the following attributes:
  * `id` - The ID of the scaling rule.
  * `scaling_group_id` - The ID of the scaling group to which the scaling rule belongs.
  * `name` - The name of the scaling rule.
  * `type` - The type of the scaling rule. Possible values include `SimpleScalingConfiguration`, `TargetTrackingConfiguration`, and `StepScalingConfiguration`.
  * `cooldown` - The cooldown time for the scaling rule, during which no additional scaling activities will be triggered. This applies only to simple scaling rules. Value range: 0~86400 seconds.
  * `adjustment_type` - The adjustment type used in the scaling rule. Possible values are:
    - `QuantityChangeInCapacity`: Increases or decreases the number of ECS instances by a specified amount.
    - `PercentChangeInCapacity`: Increases or decreases the number of ECS instances by a specified percentage.
    - `TotalCapacity`: Sets the total number of ECS instances in the scaling group to a specified value.
  * `adjustment_value` - The adjustment value used in the scaling rule. This specifies the magnitude of the adjustment based on the `adjustment_type`.
  * `min_adjustment_magnitude` - The minimum adjustment magnitude for the scaling rule. This ensures that any scaling activity meets a minimum threshold.
  * `scaling_rule_ari` - The unique ARN (Aliyun Resource Name) of the scaling rule.