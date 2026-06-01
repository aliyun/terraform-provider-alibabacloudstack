---
subcategory: "一站式大数据开发治理平台"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_file"
sidebar_current: "docs-Alibabacloudstack-resource-data-works-file"
description: |-
  编排Data Works文件
---

# alibabacloudstack_data_works_file

使用Provider配置的凭证在指定的数据工坊项目中编排Data Works文件。

## 示例用法

### 基础用法

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

## 参数说明

支持以下参数：

* `project_id` - (必填，变更时重建) 要创建文件的数据工坊项目的ID。

* `file_name` - (必填) 文件名称。

* `file_type` - (必填) 文件类型。您可以使用 `alibabacloudstack_dataworks_file_types` 数据源查询可用的文件类型。

* `folder_id` - (选填) 文件所在文件夹的ID。

* `description` - (选填) 文件的描述信息。

* `content` - (选填，由系统返回) 文件内容。例如SQL脚本、Shell脚本等。

* `file_id` - (由系统返回) 文件的唯一标识符，由API创建后返回。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的唯一标识符。其值格式为 `<project_id>:<file_id>`。

* `file_id` - 数据工坊项目内文件的唯一标识符。

## Import

DataWorks 文件可以使用 project_id 和 file_id 进行导入，两者用冒号分隔，例如：

```
$ terraform import alibabacloudstack_data_works_file.example 12345:67890
```
