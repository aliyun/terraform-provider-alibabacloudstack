---
subcategory: "企业级分布式应用服务 EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_namespace"
sidebar_current: "docs-Alibabacloudstack-edas-namespace"
description: |-
  编排企业级分布式应用服务（EDAS）命名空间资源。
---

# alibabacloudstack_edas_namespace

使用Provider配置的凭证在指定的资源集下编排企业级分布式应用服务（Edas）命名空间。

## 示例用法

### 基础用法

```hcl
provider "alibabacloudstack" {
  region = var.region
}

variable "region" {
  default = "cn-hangzhou"
}

variable "name" {
  default = "tftestacc456"
}

resource "alibabacloudstack_edas_namespace" "default" {
  description          = var.name
  namespace_logical_id = "${var.region}:${var.name}"
  namespace_name       = var.name
}
```

## 参数说明

支持以下参数：

* `namespace_logical_id` - (必填, 变更时重建) 命名空间的逻辑ID。修改此参数会强制重新创建资源。对于自定义命名空间，格式为 `区域ID:命名空间标识符`，例如 `cn-beijing:tdy218`。对于默认命名空间，格式仅为 `区域ID`，例如 `cn-beijing`。
* `namespace_name` - (必填) 命名空间的名称。长度最多为63个字符。
* `description` - (可选) 命名空间的描述信息。长度最多为128个字符。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 命名空间在 Terraform 中的唯一标识符（ID）。
* `tenant_id` - 命名空间的租户ID。

## Timeouts

`timeouts` 块允许您为某些操作指定超时时间：

* `create` - （默认 1 分钟）创建命名空间时使用。
* `delete` - （默认 1 分钟）删除命名空间时使用。
* `update` - （默认 1 分钟）更新命名空间时使用。

## Import

EDAS 命名空间可以使用命名空间 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_edas_namespace.example 123456
```