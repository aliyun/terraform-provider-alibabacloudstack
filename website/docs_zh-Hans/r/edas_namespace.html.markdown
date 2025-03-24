---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_namespace"
sidebar_current: "docs-Alibabacloudstack-edas-namespace"
description: |- 
  编排企业级分布式应用服务（Edas）命名空间
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

variable "logical_id" {
  default = "${var.region}:tftestacc456"
}

resource "alibabacloudstack_edas_namespace" "default" {
  debug_enable         = false
  description          = var.name
  namespace_logical_id = var.logical_id
  namespace_name       = var.name
}
```

## 参数参考

支持以下参数：

* `description` - (可选) 命名空间的描述信息。长度最多为128个字符。
* `namespace_logical_id` - (必填, 变更时重建) 命名空间的逻辑ID。  
  - 对于自定义命名空间，格式为 `区域ID:命名空间标识符`，例如 `cn-beijing:tdy218`。
  - 对于默认命名空间，格式仅为 `区域ID`，例如 `cn-beijing`。
* `namespace_name` - (必填)  命名空间的名称。长度最多为63个字符。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 命名空间在Terraform中的唯一标识符(ID)。