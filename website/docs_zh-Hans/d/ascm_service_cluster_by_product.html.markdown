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

## 参数说明

以下是支持的参数：

* `ids` - (可选) 用于过滤结果的区域 ID 列表。通过指定此参数，可以筛选出特定区域的服务集群。
* `product_name` - (必填) 要检索的服务集群所属的产品名称。例如，`ecs` 表示弹性计算服务。
* `organization` - (可选) 要检索的服务集群所属的组织。如果指定此参数，则结果将限制在该组织下的服务集群。

## 属性说明

以下属性被导出：

* `ids` - 符合条件的区域 ID 列表。
* `region_list` - 区域列表。每个元素包含以下属性：
    * `region_id` - 区域的唯一标识符。
    * `region_type` - 区域的类型，描述该区域的分类或用途。
* `cluster_list` - 按某些标准分组的集群列表。具体包括：
    * `cluster_by_region` - 按区域分组的集群列表。此列表详细展示了每个区域内的服务集群信息。
