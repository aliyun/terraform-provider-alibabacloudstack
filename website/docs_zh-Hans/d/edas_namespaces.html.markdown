---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_namespaces"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-namespaces"
description: |- 
   查询企业级分布式应用服务命名空间
---

# alibabacloudstack_edas_namespaces

根据指定过滤条件列出当前凭证权限可以访问的企业级分布式应用服务命名空间列表。

## 示例用法

以下示例展示了如何使用 `alibabacloudstack_edas_namespaces` 数据源来获取 EDAS 命名空间的详细信息。

### 根据命名空间 ID 获取命名空间

```terraform
variable "name" {	
	default = "tf-testAccNamespace-102"
}

variable "logical_id" {
  default = ":tftest102"
}

resource "alibabacloudstack_edas_namespace" "default" {
	description = var.name
	namespace_logical_id = var.logical_id
	namespace_name = var.name
}

data "alibabacloudstack_edas_namespaces" "by_ids" {
	ids = [alibabacloudstack_edas_namespace.default.id]
}

output "namespace_by_ids" {
  value = data.alibabacloudstack_edas_namespaces.by_ids.namespaces[0].namespace_name
}
```

### 根据命名空间名称正则表达式获取命名空间

```terraform
data "alibabacloudstack_edas_namespaces" "by_name_regex" {
	name_regex = "^tf-testAccNamespace"
}

output "namespace_by_name_regex" {
  value = data.alibabacloudstack_edas_namespaces.by_name_regex.namespaces[0].namespace_name
}
```

## 参数说明

以下参数是支持的：

* `ids` - （可选，变更时重建）命名空间ID列表。用于通过特定的命名空间ID过滤结果。
* `name_regex` - （可选，变更时重建）用于按命名空间名称筛选结果的正则表达式字符串。当你想查找匹配特定命名模式的命名空间时，这会非常有用。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 命名空间名称列表。
* `namespaces` - Edas命名空间列表。每个元素包含以下属性：
  * `description` - 命名空间的描述。它是对命名空间的简要说明，长度不得超过128个字符。
  * `id` - 资源ID，它在Terraform中唯一标识命名空间。
  * `namespace_id` - 由企业分布式应用服务（EDAS）生成的命名空间唯一ID。此ID由EDAS内部使用来管理命名空间。
  * `namespace_logical_id` - 命名空间的逻辑ID。**注意：** 创建命名空间后，逻辑ID不能更改。其格式为`物理区域ID:逻辑区域标识符`。
  * `namespace_name` - 命名空间的名称。这是用户定义的命名空间名称，在阿里云账户内必须是唯一的。
  * `belong_region` - 命名空间所属的物理区域ID。这表示命名空间所在的地理位置。
  * `user_id` - 命名空间所属的阿里云账户ID。这有助于识别命名空间的所有者。