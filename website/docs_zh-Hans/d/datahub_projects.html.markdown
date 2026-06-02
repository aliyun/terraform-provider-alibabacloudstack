---
subcategory: "数据总线 DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_projects"
sidebar_current: "docs-Alibabacloudstack-datasource-datahub-projects"
description: |-
  提供DataHub项目列表。
---

# alibabacloudstack\_datahub\_projects

此数据源提供阿里云账号下的DataHub项目列表。

## 示例用法

```hcl
data "alibabacloudstack_datahub_projects" "example" {
  name_regex = "^my-Project"
}

output "first_project_id" {
  value = data.alibabacloudstack_datahub_projects.example.projects.0.id
}
```

## 参数说明

支持以下参数：

* `name_regex` - (可选, 变更后重建) 用于按项目名称过滤结果的正则表达式字符串。

## 属性说明

导出以下属性：

* `ids` - 项目ID列表。
* `projects` - DataHub项目列表。每个元素包含以下属性：
  * `id` - 项目ID（与项目名称相同）。
  * `name` - 项目名称。
  * `comment` - 项目注释或描述。
  * `create_time` - 项目创建时间（Unix时间戳）。
  * `last_modify_time` - 项目最后修改时间（Unix时间戳）。
