---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_projects"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-projects"
description: |-
  查询Max Compute Project
---

# alibabacloudstack_maxcompute_projects

根据指定过滤条件列出当前凭证权限可以访问的Max Compute Project。[什么是Project](https://www.alibabacloud.com/help/en/maxcompute/)

## 示例用法

```terraform
variable "name" {
  default = "tf_example_acc"
}

resource "alibabacloudstack_maxcompute_project" "default" {
  default_quota = "默认后付费Quota"
  project_name  = var.name
  comment       = var.name
  product_type  = "PayAsYouGo"
}

data "alibabacloudstack_maxcompute_projects" "default" {
  name_regex = alibabacloudstack_maxcompute_project.default.project_name
}

output "alibabacloudstack_maxcompute_project_example_id" {
  value = data.alibabacloudstack_maxcompute_projects.default.projects.0.project_name
}
```

## 参数说明

支持以下参数：
* `ids` - (可选，变更时重建) 项目ID列表，用于精确匹配特定的MaxCompute项目。
* `name` - (可选，变更时重建) 按项目名称过滤结果的字符串，支持模糊匹配。
* `name_regex` - (可选，必填) 用于通过项目名称筛选结果的正则表达式字符串，允许更灵活的过滤规则。

## 属性说明

除了上述列出的参数外，还导出以下属性：
* `ids` - 匹配到的项目ID列表，按照实际查询结果返回。
* `names` - 匹配到的项目的名称列表，与`ids`一一对应。
* `projects` - 匹配到的项目条目列表，每个元素包含以下属性：
  * `id` - 项目的唯一标识符（ID）。
  * `name` - 项目的名称，与创建时指定的名称一致。
  * `project_name` - 作为附加计算属性的项目名称，与`name`相同，便于引用。