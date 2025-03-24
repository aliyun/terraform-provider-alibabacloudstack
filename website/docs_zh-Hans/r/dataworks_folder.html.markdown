---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_folder"
sidebar_current: "docs-Alibabacloudstack-data-works-folder"
description: |- 
  编排Data Works文件夹
---

# alibabacloudstack_data_works_folder

使用Provider配置的凭证在指定的资源集下编排Data Works文件夹。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testaccdata_worksfolder99960"
}

resource "alibabacloudstack_data_works_folder" "example" {
  project_id   = "12345"
  folder_path  = "业务流程/test/folderMaxCompute/testcxt"
}
```

高级用法(指定 `folder_id` 和 `project_identifier`)

```terraform
variable "name" {
  default = "tf-testaccdata_worksfolder99960"
}

resource "alibabacloudstack_data_works_folder" "advanced" {
  project_id         = "12345"
  folder_path        = "业务流程/test/folderUserDefined/testcxt"
  folder_id          = "custom_folder_id_001"
  project_identifier = "test_project"
}
```

## 参数参考

支持以下参数：

* `folder_path` - (必填) 文件夹的路径。文件夹路径由四个部分组成：`业务流程/{业务流程名称}/{文件夹类型}/{目录名称}`。
  * 第一段必须是 `业务流程`。
  * 第二段必须是项目中已存在的业务流程名称。
  * 第三段必须是以下关键字之一：`folderDi`、`folderMaxCompute`、`folderGeneral`、`folderJdbc` 或 `folderUserDefined`。
  * 第四段是您指定的自定义目录名称。

* `project_id` - (必填，变更时重建) 要创建文件夹的数据工坊项目的ID。

* `project_identifier` - (选填)数据工坊项目的标识符(名称)。如果提供了 `project_id`，则此参数不是强制性的，但可以提供额外的清晰度。

* `folder_id` - (选填，变更时重建) 文件夹的唯一标识符。如果不指定，Terraform 将在创建期间自动生成一个。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 文件夹资源的唯一标识符。其值格式为 `<folder_id>:<project_id>`。

* `folder_id` - 数据工坊项目内文件夹的唯一标识符。