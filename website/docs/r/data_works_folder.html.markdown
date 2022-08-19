---
subcategory: "Data Works"
layout: "apsarastack"
page_title: "ApsaraStack: apsarastack_data_works_folder"
sidebar_current: "docs-apsarastack-resource-data-works-folder"
description: |-
  Provides a ApsaraStack Data Works Folder resource.
---

# apsarastack\_data\_works\_folder

Provides a Data Works Folder resource.

For information about Data Works Folder and how to use it, see [What is Folder](https://help.aliyun.com/document_detail/173940.html).

## Example Usage

Basic Usage

```terraform
resource "apsarastack_data_works_folder" "example" {
  project_id  = "320687"
  folder_path = "Business Flow/tfTestAcc/folderDi/tftest1"
}
```

## Argument Reference

The following arguments are supported:

* `folder_path` - (Required) Folder Path. The folder path composed with for part: `Business Flow/{Business Flow Name}/[folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined]/{Directory Name}`. The first segment of path must be `Business Flow`, and sencond segment of path must be a Business Flow Name within the project. The third part of path must be one of those keywords:`folderDi|folderMaxCompute|folderGeneral|folderJdbc|folderUserDefined`. Then the finial part of folder path can be specified in yourself.
* `project_id` - (Required, ForceNew) The ID of the project.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Folder. The value formats as `<folder_id>:<$.ProjectId>`.

## Import

Data Works Folder can be imported using the id, e.g.

```
$ terraform import apsarastack_data_works_folder.example <folder_id>:<$.ProjectId>
```
