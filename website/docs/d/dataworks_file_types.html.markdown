---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dataworks_file_types"
description: |-
  Provides a list of DataWorks File Types.
---

# alibabacloudstack_dataworks_file_types

This data source provides the DataWorks File Types available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_dataworks_file_types" "example" {
  project_id = 12345
  name       = "Shell"
}

output "file_types" {
  value = data.alibabacloudstack_dataworks_file_types.example.file_types
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required) The ID of the DataWorks project.
* `name` - (Optional) The keyword used to filter file types by name.
* `ids` - (Optional, Computed) A list of file type IDs to filter results.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of file type IDs.
* `names` - A list of file type names.
* `file_types` - A list of DataWorks File Types. Each element contains the following attributes:
  * `node_type_name` - The display name of the node type.
  * `node_type_id` - The numeric identifier of the node type.
