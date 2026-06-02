---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_projects"
description: |-
  Provides a list of DataHub Projects.
---

# alibabacloudstack\_datahub\_projects

This data source provides a list of DataHub Projects in an Alibaba Cloud account.

## Example Usage

```hcl
data "alibabacloudstack_datahub_projects" "example" {
  name_regex = "^my-Project"
}

output "first_project_id" {
  value = data.alibabacloudstack_datahub_projects.example.projects.0.id
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter results by project name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of project IDs.
* `projects` - A list of DataHub projects. Each element contains the following attributes:
  * `id` - The ID of the project (same as the project name).
  * `name` - The name of the project.
  * `comment` - The comment or description of the project.
  * `create_time` - The creation time of the project (Unix timestamp).
  * `last_modify_time` - The last modification time of the project (Unix timestamp).
