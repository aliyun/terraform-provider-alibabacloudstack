---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_environment_services_by_product"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-environment-services"
description: |-
    查询环境服务
---

# alibabacloudstack_ascm_environment_services_by_product

根据指定过滤条件列出当前凭证权限可以访问的环境服务列表

## 示例用法

```
data "alibabacloudstack_ascm_environment_services_by_product" "default" {
}
output "envser" {
  value = data.alibabacloudstack_ascm_environment_services_by_product.default.*
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 环境服务ID列表。用于过滤特定的环境服务。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `result` - 环境服务列表。每个元素包含以下属性：
  * `id` - 环境服务的唯一标识符。
  * `name` - 环境服务的名称。
  * `product` - 环境服务所属的产品类型。
  * `status` - 环境服务的状态。
  * `region` - 环境服务所在的区域。
