---
subcategory: "DataWorks"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_data_works_project"
sidebar_current: "docs-Alibabacloudstack-data-works-project"
description: |- 
  编排Data Works项目
---

# alibabacloudstack_data_works_project

使用Provider配置的凭证在指定的资源集下编排Data Works项目。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf_testacc97890"
}

resource "alibabacloudstack_data_works_project" "default" {
  project_name   = var.name
  task_auth_type = "PROJECT"
}
```

## 参数参考

支持以下参数：

* `project_name` - (必填) data_works项目的名称，也被称为工作区名称。必须是唯一的，并遵循命名约定。
* `task_auth_type` - (可选) data_works项目的任务授权类型。有效值包括：
  * `PROJECT`: 授权范围为整个项目。
  * `CUSTOM`: 可以应用自定义授权设置。如果不指定，默认值为 `PROJECT`。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `project_id` - data_works项目的唯一标识符(ID)。此属性是在项目创建后自动生成的，可以用于在其他资源中引用该项目。
