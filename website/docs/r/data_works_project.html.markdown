---
subcategory: "Data Works"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_data_works_project"
sidebar_current: "docs-apsarastack-resource-data-works-project"
description: |- Provides a ApsaraStack Data Works Project resource.
---

# apsarastack\_data\_works\_project

Provides a Data Works Project resource.

For information about Data Works Project and how to use it,

## Example Usage

Basic Usage

```terraform
resource "apsarastack_data_works_project" "default" {
  project_name = "tf_testacc46774"
  task_auth_type = "PROJECT"
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Computed) The ID of the project.
* `project_name` - (Required) The name of the project.
* `task_auth_type` - (Optional) The task auth type of the project, default value is PROJECT.

