---
subcategory: "ECS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ecs_hpccluster"
sidebar_current: "docs-Alibabacloudstack-ecs-hpccluster"
description: |- 
  编排云服务器（Ecs）高性能计算集群（HPC）
---

# alibabacloudstack_ecs_hpccluster
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_ecs_hpc_cluster`

使用Provider配置的凭证在指定的资源集下编排云服务器（Ecs）高性能计算集群（HPC）。

## 示例用法

### 基础用法

```terraform
variable "name" {
    default = "tf-testaccecshpc_cluster26697"
}

resource "alibabacloudstack_ecs_hpccluster" "example" {
  name        = var.name
  description = "For Terraform Test"
}
```

## 参数说明

支持以下参数：

* `name` - (必填) ECS HPC集群的名称。此名称在同一区域内必须唯一，可以由大写/小写字母、数字、连字符(-)和下划线(_)组成，长度不得超过128个字符。
* `description` - (可选) ECS HPC集群的描述。该描述可以包含大写/小写字母、数字、句点(.)、冒号(:)、下划线(_)、连字符(-)和 at 符号(@)，长度不得超过256个字符。

## 属性说明

除了上述参数外，还导出以下属性：

* `id` - ECS HPC集群的ID。此属性与`name`参数相同，可用于在Terraform配置的其他部分中引用此资源。
* `description` - (输出) ECS HPC集群的描述信息。