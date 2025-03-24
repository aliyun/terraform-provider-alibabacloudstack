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

## 参数参考

支持以下参数：

* `ids` - (可选) 区域 ID 列表。
* `product_name` - (必填) 通过指定服务名称过滤结果。
* `organization` - (可选) 通过指定组织名称过滤结果。

## 属性参考

除了上述列出的参数外，还导出以下属性：

* `region_list` - 区域列表。每个元素包含以下属性：
    * `region_id` - 区域 ID。
    * `region_type` - 区域类型。 
