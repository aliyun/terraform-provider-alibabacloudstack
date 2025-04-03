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

## 参数参考
支持以下参数：

* `ids` - (可选) 用于过滤结果的 CU ID 列表。
* `name_regex` - (可选) 通过名称过滤 CU 的正则表达式模式。
* `cluster_name` - (可选) 用于过滤 CU 的集群名称。

## 属性参考
导出以下属性：

* `ids` - 集群的 ID 列表。
* `cus` - CU 列表。每个元素包含以下属性：
    * `id` - 集群的唯一标识符。
    * `cu_name` - CU 的名称。
    * `cu_num` - CU 的数量。
    * `cluster_name` - 集群的名称。