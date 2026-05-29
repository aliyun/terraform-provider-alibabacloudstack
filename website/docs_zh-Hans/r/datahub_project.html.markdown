---
subcategory: "数据总线 DataHub"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_datahub_project"
sidebar_current: "docs-Alibabacloudstack-datahub-project"
description: |- 
  编排datahub项目
---

# alibabacloudstack_datahub_project

使用Provider配置的凭证在指定的资源集下编排datahub项目。

## 示例用法

### 基础用法：

```hcl
variable "name" {
    default = "tf_testacc_datahub_project37506"
}

resource "alibabacloudstack_datahub_project" "default" {
  name    = var.name
  comment = "This is a test project created via Terraform."
}
```

## 参数说明

支持以下参数：
  * `name` - (必填, 变更时重建) DataHub 项目的名称。其长度必须在 3 到 32 个字符之间，只允许字母、数字和下划线 (`_`)，且不区分大小写。
  * `comment` - (选填) 关于 DataHub 项目的简要说明或注释。最大长度为 255 个字符。默认值："project added by terraform"。
  * `vpc_ids` - (选填) 允许访问该 DataHub 项目的 VPC ID 列表。修改此参数不会触发资源重建。

## 属性说明

除了上述所有参数外，还导出了以下属性：
  * `id` - DataHub 项目的 ID。它与 `name` 相同。
  * `create_time` - DataHub 项目的创建时间，以字符串格式的 UNIX 时间戳表示。
  * `last_modify_time` - DataHub 项目的最后修改时间，以字符串格式的 UNIX 时间戳表示。