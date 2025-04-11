---
subcategory: "Quick BI"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_quick_bi_workspace"
sidebar_current: "docs-alibabacloudstack-resource-quick-bi-workspace"
description: |-
  编排Quick BI工作区
---

# alibabacloudstack_quick_bi_workspace

使用Provider配置的凭证在指定的资源集编排Quick BI工作区。

## 示例用法

### 基础用法

```terraform

resource "alibabacloudstack_quick_bi_workspace" "default" {
  workspace_name = "example_value"
  workspace_desc = "example_value"
  use_comment = "false"
  allow_share = "false"
  allow_publish = "false"
}

```

## 参数说明

以下是支持的参数：

* `workspace_name` - (必填) 工作区名称。
* `workspace_desc` - (可选) 工作区描述。
* `use_comment` - (可选) 在创建数据集时是否使用表注释（对应偏好设置）。有效值：`true` 和 `false`。
* `allow_share` - (可选) 是否允许报告共享（对应功能权限 - 作品可以授权）。有效值：`false`、`true`。
* `allow_publish` - (可选) 是否允许报告公开（对应功能权限 - 作品可以公开）。有效值：`false`、`true`。

## 属性说明

以下属性将导出：

* `workspace_id` - 工作区在 Terraform 中的资源 ID。

## 导入

Quick BI 工作区可以使用 id 导入，例如

```bash
$ terraform import alibabacloudstack_quick_bi_workspace.example <id>
```