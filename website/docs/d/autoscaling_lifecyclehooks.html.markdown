---
subcategory: "AutoScaling"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_autoscaling_lifecyclehooks"
sidebar_current: "docs-Alibabacloudstack-datasource-autoscaling-lifecyclehooks"
description: |- 
  Provides a list of autoscaling lifecyclehooks owned by an alibabacloudstack account.
---

# alibabacloudstack_autoscaling_lifecyclehooks
-> **NOTE:** Alias name has: `alibabacloudstack_ess_lifecycle_hooks`

This data source provides a list of autoscaling lifecyclehooks in an AlibabacloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_autoscaling_lifecyclehooks" "example" {
  scaling_group_id = "sg-1234567890abcdefg"
  name_regex       = "lifecyclehook-*"

  ids = ["lh-1234567890abcdefg"]
}

output "first_lifecycle_hook_id" {
  value = data.alibabacloudstack_autoscaling_lifecyclehooks.example.hooks.0.id
}
```

## Argument Reference

The following arguments are supported:

* `scaling_group_id` - (Optional) The ID of the scaling group to which the lifecycle hooks belong.
* `name_regex` - (Optional) A regex string used to filter the results by lifecycle hook name.
* `ids` - (Optional) A list of lifecycle hook IDs used to filter the results.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `hooks` - A list of autoscaling lifecycle hooks. Each element contains the following attributes:
  * `id` - The ID of the lifecycle hook.
  * `name` - The name of the lifecycle hook.
  * `scaling_group_id` - The ID of the scaling group to which the lifecycle hook belongs.
  * `default_result` - Defines the action the Auto Scaling group should take when the lifecycle hook timeout elapses. It can be either `CONTINUE` or `ABANDON`.
  * `heartbeat_timeout` - Defines the amount of time, in seconds, that can elapse before the lifecycle hook times out. When the lifecycle hook times out, Auto Scaling performs the action defined in the `default_result` parameter.
  * `lifecycle_transition` - The type of scaling activity associated with the lifecycle hook. Possible values include `INSTANCE_LAUNCHING`, `INSTANCE_TERMINATING`, etc.
  * `notification_arn` - The ARN of the notification target that Auto Scaling will notify when an instance enters the wait state due to the lifecycle hook.
  * `notification_metadata` - Additional information that you want to include when Auto Scaling sends a message to the notification target.

* `ids` - A list of lifecycle hook IDs.
* `names` - A list of lifecycle hook names.