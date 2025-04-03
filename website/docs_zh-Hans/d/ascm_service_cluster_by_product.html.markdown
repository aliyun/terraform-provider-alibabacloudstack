---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ascm_service_cluster_by_product"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-service-cluster-by-product"
description: |-
    查询服务集群列表
---

# alibabacloudstack_ascm_service_cluster_by_product

根据指定过滤条件列出当前凭证权限可以访问的服务集群的列表。

## 示例用法

```hcl
data "alibabacloudstack_ascm_regions_by_product" "example" {
  product_name = "ecs"
}

output "regions" {
  value = data.alibabacloudstack_ascm_regions_by_product.example.region_list
}
```

## 参数参考

以下是支持的参数：

* `ids` - (可选) 用于过滤结果的区域 ID 列表。
* `product_name` - (必填) 要检索区域的产品名称。
* `organization` - (可选) 要检索区域的组织。

## 属性参考

以下属性被导出：

* `ids` - 区域的 ID 列表。
* `region_list` - 区域列表。每个元素包含以下属性：
    * `region_id` - 区域的唯一标识符。
    * `region_type` - 区域的类型。