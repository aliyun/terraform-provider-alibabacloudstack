---
subcategory: "DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_project"
sidebar_current: "docs-Alibabacloudstack-datahub-project"
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
* `comment` - (Optional, ForceNew) A brief description or comment about the DataHub project. The maximum length is 255 characters.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the DataHub project. It is the same as its `name`.
* `create_time` - The creation time of the DataHub project. This is a human-readable string in the format `YYYY-MM-DD HH:mm:ss`.
* `last_modify_time` - The last modification time of the DataHub project. Initially, this value is the same as `create_time`. Like `create_time`, it is also a human-readable string in the format `YYYY-MM-DD HH:mm:ss`.

## Import

DataHub projects can be imported using their `name` or ID. For example:

```bash
$ terraform import alibabacloudstack_datahub_project.example tf_testacc_datahub_project
```