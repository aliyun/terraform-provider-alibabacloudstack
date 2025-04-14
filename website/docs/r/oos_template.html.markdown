---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_template"
sidebar_current: "docs-Alibabacloudstack-oos-template"
description: |- 
  Provides a oos Template resource.
---

# alibabacloudstack_oos_template

Provides a OOS Template resource. For information about Alibaba Cloud OOS Template and how to use it, see [What is Resource Alibaba Cloud OOS Template](https://www.alibabacloud.com/help/doc-detail/120761.htm).

## Example Usage

```terraform
variable "name" {
    default = "tf-testaccoostemplate93918"
}

resource "alibabacloudstack_oos_template" "default" {
  content       = <<EOF
  {
    "FormatVersion": "OOS-2019-06-01",
    "Description": "Update Describe instances of given status",
    "Parameters": {
      "Status": {
        "Type": "String",
        "Description": "(Required) The status of the Ecs instance."
      }
    },
    "Tasks": [
      {
        "Properties": {
          "Parameters": {
            "Status": "{{ Status }}"
          },
          "API": "DescribeInstances",
          "Service": "Ecs"
        },
        "Name": "foo",
        "Action": "ACS::ExecuteApi"
      }
    ]
  }
  EOF
  template_name = var.name
  version_name  = "v1.0"
  tags = {
    "Created" = "TF",
    "For"     = "acceptance Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `content` - (Required) The content of the template. The template must be in the JSON or YAML format. Maximum size: 64 KB. This field defines the structure and logic of the template, including parameters, tasks, and other configurations.
* `auto_delete_executions` - (Optional) When deleting a template, whether to delete its related executions. Default value is `false`.
* `template_name` - (Required, ForceNew) The name of the template. The template name can be up to 200 characters in length. The name can contain letters, digits, hyphens (-), and underscores (_). It cannot start with `ALIYUN`, `ACS`, `ALIBABA`, or `ALICLOUD`.
* `version_name` - (Optional) The name of the template version. This allows you to manage different versions of the same template.
* `tags` - (Optional) A mapping of tags to assign to the resource. Tags help in organizing and categorizing your resources.
* `description` - (Required) Description of the template.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource. It is the same as `template_name`.
* `created_by` - The creator of the template. Indicates who created the template.
* `created_date` - The time when the template was created. This is useful for tracking the lifecycle of the template.
* `description` - The description of the template. Provides a brief summary of what the template does.
* `has_trigger` - Indicates whether the template has been triggered successfully. This attribute helps in understanding if the template execution was initiated.
* `share_type` - The sharing type of the template. Templates created by users are set to `Private`. Common templates provided by OOS are set to `Public`.
* `template_format` - The format of the template. The system automatically identifies whether the template is in JSON or YAML format.
* `template_id` - The unique identifier of the OOS Template. Useful for referencing the template in other API calls.
* `template_type` - The type of the OOS Template. `Automation` means the implementation of Alibaba Cloud API template, while `Package` represents a template for installing software.
* `template_version` - The version of the OOS Template. Helps in managing different iterations of the same template.
* `updated_by` - The user who last updated the template. Useful for auditing purposes.
* `updated_date` - The time when the template was last updated. This helps in tracking changes made to the template over time.
* `template_name` - (Computed) Name of the template.