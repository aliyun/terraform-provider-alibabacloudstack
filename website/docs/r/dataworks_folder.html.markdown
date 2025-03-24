---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_folder"
sidebar_current: "docs-Alibabacloudstack-data-works-folder"
description: |- 
  Provides a data works Folder resource.
---

# alibabacloudstack_data_works_folder

Provides a dataworks Folder resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccdataworksfolder99960"
}

resource "alibabacloudstack_data_works_folder" "example" {
  project_id   = "12345"
  folder_path  = "业务流程/test/folderMaxCompute/testcxt"
}
```

## Argument Reference

The following arguments are supported:

* `folder_path` - (Required) The path of the folder. The folder path is composed of four parts: `业务流程/{Business Flow Name}/{Folder Type}/{Directory Name}`. 
  * The first segment must be `业务流程`.
  * The second segment must be the name of an existing Business Flow within the project.
  * The third segment must be one of the following keywords: `folderDi`, `folderMaxCompute`, `folderGeneral`, `folderJdbc`, or `folderUserDefined`.
  * The fourth segment is the custom directory name you specify.
* `project_id` - (Required, ForceNew) The ID of the DataWorks project where the folder will be created.
* `project_identifier` - (Optional) The identifier (name) of the DataWorks project. This can be used to provide additional clarity but is not mandatory if `project_id` is provided.
* `folder_id` - (Optional, ForceNew) The unique identifier for the folder. If not specified, Terraform will automatically generate one during creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the folder resource. The value is formatted as `<folder_id>:<project_id>`.
* `folder_id` - The unique identifier of the folder within the DataWorks project.
```

### Explanation of Changes:
1. **Example Usage**: Added a `variable` block to demonstrate how users might define reusable variables in their Terraform configuration. Also updated the example to align with the new argument descriptions.
2. **Argument Reference**:
   - Clarified the structure and components of the `folder_path`.
   - Marked `project_id` as `(Required, ForceNew)` since it is essential and cannot be changed after creation.
   - Included `project_identifier` as an optional field with a brief explanation.
   - Added `folder_id` as an optional field that can be force-new.
3. **Attributes Reference**:
   - Described the `id` attribute format explicitly.
   - Confirmed that `folder_id` is exported as part of the resource attributes.