---
subcategory: "Data Works"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_project"
sidebar_current: "docs-alibabacloudstack-resource-data-works-project"
description: |- 
  Provides a AlibabacloudStack Data Works Project resource.
---

# alibabacloudstack\_data\_works\_project

Provides a Data Works Project resource.

For information about Data Works Project and how to use it,

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_project" "default" {
  project_name = "tf_testacc46774"
  task_auth_type = "PROJECT"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Computed) The ID of the project.
* `project_name` - (Required) The name of the project.
* `task_auth_type` - (Optional) The task auth type of the project, default value is PROJECT.

