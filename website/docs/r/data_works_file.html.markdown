---
subcategory: "One-stop Big Data Development and Governance Platform"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_file"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-file"
description: |-
  Provides a DataWorks File resource.
---

# alibabacloudstack_data_works_file

Provides a DataWorks File resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf_testaccdataworksfile"
}

resource "alibabacloudstack_data_works_project" "default" {
  name           = var.name
  description    = "${var.name}_desc"
  task_auth_type = "PROJECT"
}

resource "alibabacloudstack_data_works_business" "default" {
  project_id = alibabacloudstack_data_works_project.default.id
  name       = var.name
  description = "${var.name}_desc"
}

resource "alibabacloudstack_data_works_folder" "default" {
  project_id    = alibabacloudstack_data_works_project.default.id
  business_name = alibabacloudstack_data_works_business.default.name
  engine_type   = "General"
  folder_path   = "folder_test"
}

data "alibabacloudstack_dataworks_file_types" "anyone" {
  name       = "Shell"
  project_id = alibabacloudstack_data_works_project.default.id
}

resource "alibabacloudstack_data_works_file" "default" {
  project_id = alibabacloudstack_data_works_project.default.id
  folder_id  = alibabacloudstack_data_works_folder.default.folder_id
  file_name  = var.name
  file_type  = data.alibabacloudstack_dataworks_file_types.anyone.file_types.0.node_type_id
  content    = <<-EOF
    #!/bin/bash
    echo "Hello World"
  EOF
}
```

## Argument Reference

The following arguments are supported:

* `project_id` - (Required, ForceNew) The ID of the DataWorks project where the file will be created.
* `file_name` - (Required) The name of the file.
* `file_type` - (Required) The type of the file. You can use the `alibabacloudstack_dataworks_file_types` data source to query available file types.
* `folder_id` - (Optional) The ID of the folder where the file will be located.
* `description` - (Optional) The description of the file.
* `content` - (Optional, Computed) The content of the file. For example, SQL scripts, shell scripts, etc.
* `file_id` - (Computed) The unique identifier of the file, returned by the API after creation.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The unique identifier of the resource. The value is formatted as `<project_id>:<file_id>`.
* `file_id` - The unique identifier of the file within the DataWorks project.

## Import

DataWorks File can be imported using the project_id and file_id separated by a colon, e.g.

```
$ terraform import alibabacloudstack_data_works_file.example 12345:67890
```
