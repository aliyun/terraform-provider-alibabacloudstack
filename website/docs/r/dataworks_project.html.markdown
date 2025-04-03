---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_project"
sidebar_current: "docs-Alibabacloudstack-data-works-project"
description: |- 
  Provides a data works Project resource.
---

# alibabacloudstack_data_works_project

Provides a data works Project resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_testacc97890"
}

resource "alibabacloudstack_data_works_project" "default" {
  project_name   = var.name
  task_auth_type = "PROJECT"
}
```

## Argument Reference

The following arguments are supported:

* `project_name` - (Required) The name of the DataWorks project. This is also referred to as the workspace name. It must be unique within your AlibabaCloudStack account and follow naming conventions.
* `task_auth_type` - (Optional) The task authorization type for the DataWorks project. Valid values include:
  * `PROJECT`: Authorization is scoped to the entire project.
  * `CUSTOM`: Customized authorization settings can be applied. If not specified, the default value is `PROJECT`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `project_id` - The unique identifier (ID) of the DataWorks project. This attribute is automatically generated after the project is created and can be used to reference the project in other resources.
```

### Explanation of Changes

1. **Example Usage**: Updated the example to use a Terraform variable (`var.name`) for better reusability and clarity. This aligns with best practices for writing reusable Terraform configurations.
   
2. **Argument Reference**:
   - Added detailed descriptions for `project_name` and `task_auth_type`.
   - Clarified valid values for `task_auth_type` and provided explanations for each option.
   - Ensured that the argument descriptions are consistent with Terraform documentation standards.

3. **Attributes Reference**:
   - Provided a clear explanation of the `project_id` attribute, emphasizing that it is computed and automatically generated upon resource creation.
   - Ensured the attribute description aligns with the expected behavior of the resource.

This updated documentation provides a comprehensive guide for using the `alibabacloudstack_dataworks_project` resource, making it easier for users to understand and implement in their Terraform configurations.