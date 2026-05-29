---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_project"
sidebar_current: "docs-Alibabacloudstack-resource-datahub-project"
description: |- 
  Provides a datahub Project resource.
---

# alibabacloudstack_datahub_project

Provides a datahub Project resource.

## Example Usage

Basic Usage:

```hcl
variable "name" {
    default = "tf_testacc_datahub_project"
}

resource "alibabacloudstack_datahub_project" "default" {
  name    = var.name
  comment = "This project is created using Terraform for testing purposes."
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required, ForceNew) The name of the DataHub project. Its length must be between 3 and 32 characters. Only letters, digits, and underscores (`_`) are allowed. It is case-insensitive.
* `comment` - (Optional) A brief description or comment about the DataHub project. The maximum length is 255 characters. Default: "project added by terraform".
* `vpc_ids` - (Optional) A list of VPC IDs that are allowed to access the DataHub project. Changes to this parameter will not trigger resource recreation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DataHub project. It is the same as its `name`.
* `create_time` - The creation time of the DataHub project, represented as a UNIX timestamp in string format.
* `last_modify_time` - The last modification time of the DataHub project, represented as a UNIX timestamp in string format.

## Import

DataHub projects can be imported using their `name` or ID. For example:

```bash
$ terraform import alibabacloudstack_datahub_project.example tf_testacc_datahub_project
```