---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_regions_by_product"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-regions-by-product"
description: |-
    查询产品可用区域
---

# alibabacloudstack_ascm_regions_by_product

根据指定过滤条件列出当前凭证权限可以访问的区域列表

## 示例用法

```
data "alibabacloudstack_ascm_regions_by_product" "regions" {
  output_file = "product_regions"
  product_name = "ecs"
}
output "regions" {
  value = data.alibabacloudstack_ascm_regions_by_product.regions.*
}
```

## 参数说明

支持以下参数：

* `ids` - (可选) 区域 ID 列表。此参数用于通过指定的区域 ID 过滤结果。
* `product_name` - (必填) 服务名称。此参数用于通过指定的服务名称过滤结果，例如 `"ecs"` 表示弹性计算服务。
* `organization` - (可选) 组织名称。此参数用于通过指定的组织名称过滤结果。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `region_list` - 区域列表。每个元素包含以下属性：
    * `region_id` - 区域 ID。表示该区域的唯一标识符。
    * `region_type` - 区域类型。表示该区域的类型，可能用于区分不同类型的部署环境或地理位置。
