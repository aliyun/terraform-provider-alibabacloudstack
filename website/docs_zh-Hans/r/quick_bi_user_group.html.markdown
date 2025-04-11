---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_user_group"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-user-group"
description: |-
  编排Quick BI 用户组
---

# alibabacloudstack_quick_bi_user_group

使用Provider配置的凭证在指定的资源集编排Quick BI 用户组。


## 示例用法

### 基础用法

```terraform

resource "alibabacloudstack_quick_bi_user_group" "example" {
  user_group_name = "example_value"
  user_group_description = "example_value"
  parent_user_group_id = "-1"
}

```

## 参数说明

支持以下参数：

* `user_group_name` - (必填) 用户组名称。
* `user_group_description` - (必填) 用户组描述。
* `parent_user_group_id` - (必填) 父用户组 ID。您可以将新用户组添加到此分组中。当您输入 `-1` 时，新创建的用户组将被添加到根目录下。
* `user_group_id` - (可选) 用户组 ID。用于唯一标识用户组。

## 属性说明

导出以下属性：

* `user_group_id` - 用户组在 Terraform 中的资源 ID，用于唯一标识该用户组。

## 导入

Quick BI 用户组可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_quick_bi_user_group.example <id>
```