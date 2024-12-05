---
subcategory: "Max Compute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_projects"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-projects"
description: |-
  Provides a datasource of Max Compute Project owned by an Alibaba Cloud account.
---

# alibabacloudstack_maxcompute_projects

This data source provides Max Compute Project available to the user.[What is Project](https://www.alibabacloud.com/help/en/maxcompute/)

## Example Usage

```terraform
variable "name" {
  default = "tf_example_acc"
}

resource "alibabacloudstack_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

data "alibabacloudstack_maxcompute_projects" "default" {
  name_regex = alibabacloudstack_maxcompute_project.default.project_name
}

output "alibabacloudstack_maxcompute_project_example_id" {
  value = data.alibabacloudstack_maxcompute_projects.default.projects.0.project_name
}
```

## Argument Reference

The following arguments are supported:
* `ids` - (Optional, ForceNew, Computed) A list of Project IDs.
* `name` - (Optional, ForceNew) A string to filter results by Project name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:
* `ids` - A list of Project IDs.
* `names` - A list of name of Projects.
* `projects` - A list of Project Entries. Each element contains the following attributes:
  * `id` - ID of the Project.
  * `name` - Name of the Project.
