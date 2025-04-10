---
subcategory: "MaxCompute"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_maxcompute_cus"
sidebar_current: "docs-alibabacloudstack-datasource-maxcompute-cus"
description: |-
  查询Max Compute Cus
---

# alibabacloudstack_maxcompute_cus

根据指定过滤条件列出当前凭证权限可以访问的Max Compute Cus


## 示例用法

```hcl
data "alibabacloudstack_maxcompute_cus" "example" {
  name_regex = "example-cu"
}

output "cus" {
  value = data.alibabacloudstack_maxcompute_cus.example.cus
}
```

## 参数说明
支持以下参数：

* `ids` - (可选) 用于过滤结果的 CU ID 列表。此参数可以帮助用户通过指定的 CU ID 列表筛选出特定的计算单元。
* `name_regex` - (可选) 通过名称过滤 CU 的正则表达式模式。此参数允许用户使用正则表达式匹配 CU 名称，从而实现更灵活的筛选。
* `cluster_name` - (可选) 用于过滤 CU 的集群名称。此参数可以通过指定集群名称来筛选属于该集群的 CU。

## 属性说明
导出以下属性：

* `ids` - 集群的 ID 列表。此列表包含所有匹配过滤条件的 CU 所属集群的唯一标识符。
* `cus` - CU 列表。每个元素包含以下属性：
    * `id` - 集群的唯一标识符。此字段表示 CU 所属集群的 ID。
    * `cu_name` - CU 的名称。此字段表示计算单元的名称。
    * `cu_num` - CU 的数量。此字段表示该 CU 的实例数量或规模。
    * `cluster_name` - 集群的名称。此字段表示 CU 所属的集群名称。