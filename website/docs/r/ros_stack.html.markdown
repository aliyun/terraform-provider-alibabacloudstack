---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_stack"
sidebar_current: "docs-Alibabacloudstack-ros-stack"
description: |- 
  Provides a ROS Stack resource.
---

# alibabacloudstack_ros_stack

Provides a ROS Stack resource.

For information about ROS Stack and how to use it, see [What is Stack](https://www.alibabacloud.com/help/en/doc-detail/132086.htm).



## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_ros_stack" "example" {
  stack_name = "tf-testaccstack"

  template_body = <<EOF
    {
      "ROSTemplateFormatVersion": "2015-09-01",
      "Parameters": {
        "VpcName": {
          "Type": "String"
        },
        "InstanceType": {
          "Type": "String"
        }
      }
    }
    EOF

  stack_policy_body = <<EOF
    {
      "Statement": [
        {
          "Action": "Update:Delete",
          "Resource": "*",
          "Effect": "Allow",
          "Principal": "*"
        }
      ]
    }
    EOF

  tags = {
    Created = "TF"
    For     = "ROS"
  }

  parameters {
    parameter_key   = "VpcName"
    parameter_value = "MyVpc"
  }

  parameters {
    parameter_key   = "InstanceType"
    parameter_value = "ecs.t5-lc2m1.nano"
  }

  timeout_in_minutes = 90
}
```

## Argument Reference

The following arguments are supported:

* `create_option` - (Optional, ForceNew) Specifies whether to delete the stack after the stack is created. Default value: `KeepStackOnCreationComplete`. Valid values:
  * `KeepStackOnCreationComplete`: Retains the stack and all of its resources after the stack is created.
  * `AbandonStackOnCreationComplete`: Deletes the stack but retains all of its resources after the stack is created. This ensures that the maximum number of stacks allowed to be created is not reached. If the stack fails to be created, the stack is retained.
  * `AbandonStackOnCreationRollbackComplete`: Deletes the stack after rollback on stack creation failure is complete. This ensures that the maximum number of stacks allowed to be created is not reached. If the stack is created or the rollback fails to complete, the stack is retained.
* `deletion_protection` - (Optional, ForceNew) Specifies whether deletion protection is enabled for the stack. Valid values:
  * `Enabled`: Deletion protection is enabled for the stack.
  * `Disabled`: Deletion protection is disabled for the stack. You can delete the stack by using the Resource Orchestration Service (ROS) console or by calling the DeleteStack operation.
  
  > Deletion protection of a nested stack is the same as that of its root stack.
* `disable_rollback` - (Optional) Specifies whether to disable rollback of the stack when the stack fails to be created. Default value: `false`.
* `notification_urls` - (Optional) The callback URLs for receiving stack event notifications. Only HTTP POST is supported. Maximum of 5 URLs can be specified.
* `parameters` - (Optional) The list of parameters. Each parameter supports the following:
  * `parameter_key` - (Required) The key of the parameter.
  * `parameter_value` - (Required) The value of the parameter.
* `ram_role_name` - (Optional) The name of the RAM role. ROS assumes the specified RAM role to create the stack and call API operations by using the credentials of the role.
* `replacement_option` - (Optional) Whether to use replacement update. When the resource attribute does not support modification update, you can use replacement update to change the resource attribute. The replacement update will delete the resource and recreate the resource. The physical ID of the new resource will change. Valid values:
  * `Enabled`: Allows replacement updates.
  * `Disabled` (default): Replacement updates are not allowed.
  
  > The priority of modifying updates is higher than that of replacing updates.
* `retain_all_resources` - (Optional) Specifies whether to retain all resources in the stack during deletion.
* `retain_resources` - (Optional) Specifies whether to retain specific resources in the stack during deletion.
* `stack_name` - (Required, ForceNew) The name of the stack. The name can be up to 255 characters in length, and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `stack_policy_body` - (Optional) The structure that contains the stack policy body. The stack policy body must be 1 to 16,384 bytes in length.
* `stack_policy_during_update_body` - (Optional) Temporarily overrides the structure of the resource stack policy body. The length is 1~16,384 bytes. If you want to update protected resources, specify a temporary overwrite resource stack policy during the update. If no resource stack policy is specified, the current policy associated with the resource stack will be used.
* `stack_policy_during_update_url` - (Optional) The location of the file that updates the resource stack policy. The URL must point to the storage space located on the Web server (HTTP or HTTPS) or Alibaba Cloud OSS (for example, `oss:// ros/stack-policy/demo`, `oss:// ros/stack-policy/demo? RegionId = cn-hangzhou`). The maximum file value of the policy is 16,384 bytes.
* `stack_policy_url` - (Optional) The location of the file that contains the resource stack policy. The URL must point to the storage space located on the Web server (HTTP or HTTPS) or Alibaba Cloud OSS (for example, `oss:// ros/stack-policy/demo`, `oss:// ros/stack-policy/demo? RegionId = cn-hangzhou`). The maximum length of the policy file is 16,384 bytes.
* `template_body` - (Optional) The structure of the template body. The length is 1~524,288 bytes. If the length is long, we recommend that you pass the parameters in the request Body by using HTTP POST + Body Param to avoid request failure due to a long URL.
* `template_url` - (Optional) The location of the file that contains the template body. The URL must point to the storage space located on the Web server (HTTP or HTTPS) or Alibaba Cloud OSS (for example, `oss:// ros/template/demo`, `oss:// ros/template/demo? RegionId = cn-hangzhou`). The maximum length of the template is 524,288 bytes.
* `template_version` - (Optional) The version of the template.
* `timeout_in_minutes` - (Optional) The timeout period that is specified for the stack creation request. Default value: `60`.
* `use_previous_parameters` - (Optional) Whether the unpassed parameter uses the last passed value. Valid values:
  * `true`: The unpassed parameter uses the last passed value.
  * `false`: Unpassed parameters do not use the last passed value.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID in Terraform of Stack. Value as `stack_id`.
* `status` - The status of the stack.