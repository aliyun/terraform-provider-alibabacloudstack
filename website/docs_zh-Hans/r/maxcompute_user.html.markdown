---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_user"
sidebar_current: "docs-alibabacloudstack-resource-maxcompute-user"
description: |-
  集编排Max Compute用户
---

# alibabacloudstack_maxcompute_user

使用Provider配置的凭证在指定的资源集编排Max Compute用户。  
它类似于传统数据库中的 Database 或 Schema 的概念，设置了 MaxCompute 多用户隔离和访问控制的边界。

## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_maxcompute_user" "example" {
  user_name             = "%s"
  description           = "TestAccAlibabacloudStackMaxcomputeUser"
  lifecycle {
    ignore_changes = [
      organization_id,
    ]
  }
}
```

## 参数说明

以下是支持的参数：  
* `user_name` - (必填，变更时重建) 您要创建的用户的名称。  
* `description` - (必填，变更时重建) 您要创建的用户的描述。  
* `organization_id` - (可选) 组织的 ID。  
* `organization_name` - (可选) 组织的名称。

## 属性说明

以下是资源提供的属性：  
* `id` - 用户的 ID。  
* `user_id` - `id` 的别名。  
* `user_pk` - 用户的主键（PK）。  
* `user_type` - 用户的类型。  
* `user_name` - (计算后返回) 用户的名称。  
* `description` - (计算后返回) 用户的描述。  
* `organization_name` - 组织的名称。

## 导入

MaxCompute 项目可以使用 *name* 或 ID 导入，例如：

```bash
$ terraform import alibabacloudstack_maxcompute_cu.example tf_maxcompute_cu
```