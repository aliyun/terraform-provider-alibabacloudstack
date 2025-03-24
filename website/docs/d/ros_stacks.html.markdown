---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_stacks"
sidebar_current: "docs-Alibabacloudstack-datasource-ros-stacks"
description: |- 
  Provides a list of ros stacks owned by an Alibabacloudstack account.
---

# alibabacloudstack_ros_stacks

This data source provides a list of ROS Stacks in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_ros_stacks" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_stack_id" {
  value = data.alibabacloudstack_ros_stacks.example.stacks.0.id
}
```

Advanced Usage with Filters

```terraform
data "alibabacloudstack_ros_stacks" "example" {
  stack_name     = "my-stack"
  status         = "CREATE_COMPLETE"
  parent_stack_id = "parent-stack-id"
  show_nested_stack = true
  enable_details = true
  tags = jsonencode({
    Environment = "Production"
  })
}

output "stack_ids" {
  value = data.alibabacloudstack_ros_stacks.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `enable_details` - (Optional) Default to `false`. Set it to `true` can output more details about resource attributes.
* `ids` - (Optional, ForceNew) A list of Stack IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Stack name.
* `parent_stack_id` - (Optional, ForceNew) The ID of the parent stack.
* `show_nested_stack` - (Optional, ForceNew) Specifies whether to include nested stacks in the results. Valid values: `true`, `false`.
* `stack_name` - (Optional, ForceNew) The name of the stack. The name can be up to 255 characters in length, and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `status` - (Optional, ForceNew) The status of the stack. Valid values: `CREATE_COMPLETE`, `CREATE_FAILED`, `CREATE_IN_PROGRESS`, `DELETE_COMPLETE`, `DELETE_FAILED`, `DELETE_IN_PROGRESS`, `ROLLBACK_COMPLETE`, `ROLLBACK_FAILED`, `ROLLBACK_IN_PROGRESS`.
* `tags` - (Optional) Query the instance bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example `{"key1":"value1"}`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of Stack names.
* `stacks` - A list of Ros Stacks. Each element contains the following attributes:
  * `deletion_protection` - Specifies whether deletion protection is enabled for the stack. Valid values: `Enabled`, `Disabled`.
  * `description` - The description of the stack.
  * `disable_rollback` - Specifies whether to disable rollback of the stack when the stack fails to be created. Default value: `false`.
  * `drift_detection_time` - The time when the last successful drift detection was performed on the stack.
  * `id` - The ID of the stack.
  * `parent_stack_id` - The ID of the parent stack.
  * `ram_role_name` - The name of the RAM role. ROS assumes the RAM role to create the stack and uses credentials of the role to call the APIs of Alibaba Cloud services.
  * `root_stack_id` - The ID of the root stack. This parameter is returned when the specified stack is a nested stack.
  * `stack_drift_status` - The status of the stack on which the last successful drift detection was performed.
  * `stack_id` - The ID of the stack.
  * `stack_name` - The name of the stack.
  * `stack_policy_body` - The structure that contains the stack policy body.
  * `status` - Resource status.
  * `status_reason` - The reason why the stack is in a state.
  * `template_description` - The description of the template used to create the stack.
  * `timeout_in_minutes` - The timeout period that is specified for the stack creation request.
  * `parameters` - The list of parameters.
    * `parameter_key` - The key of the parameter.
    * `parameter_value` - The value of the parameter.
  * `tags` - Tags associated with the stack.
