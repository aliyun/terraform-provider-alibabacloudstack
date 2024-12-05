---
subcategory: "Data Works"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_data_works_folder"
sidebar_current: "docs-alibabacloudstack-resource-data-works-folder"
description: |-
  Provides a AlibabacloudStack Data Works Folder resource.
---

# alibabacloudstack\_data\_works\_folder

Provides a Data Works Folder resource.

For information about Data Works Folder and how to use it, see [What is Folder](https://help.aliyun.com/document_detail/173940.html).

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_data_works_folder" "example" {
  project_id  = "320687"
  folder_path = "Business Flow/tfTestAcc/folderDi/tftest1"
}
```

## Argument Reference

The following arguments are supported:

* `folder_path` - (Required) Folder Path. The folder path composed with for part: `Business Flow/{Business Flow Name}/[folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined]/{Directory Name}`. The first segment of path must be `Business Flow`, and sencond segment of path must be a Business Flow Name within the project. The third part of path must be one of those keywords:`folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined`. Then the finial part of folder path can be specified in yourself.
* `project_id` - (Required, ForceNew) The ID of the project.
* `project_identifier` - (Optional) The name of the project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Folder. The value formats as `<folder_id>:<$.ProjectId>`.
* `folder_id` - The resource ID of Folder.

## Import

Data Works Folder can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_data_works_folder.example <folder_id>:<$.ProjectId>
```
