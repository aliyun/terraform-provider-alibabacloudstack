---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_zones"
sidebar_current: "docs-alibabacloudstack-datasource-slb-zones"
description: |-
    查询负载均衡(SLB)可用区
---

# alibabacloudstack_slb_zones

根据指定过滤条件列出当前凭证权限可以访问的负载均衡(SLB)可用区列表。


## 示例用法

```
# 声明数据源
data "alibabacloudstack_slb_zones" "zones_ids" {}

output "slb_zones" {
  value = data.alibabacloudstack_slb_zones.zones_ids.*
}
```

## 参数说明

支持以下参数：

* `enable_details` - (可选) 默认为 `false`，仅在 `zones` 块中输出 `id`。将其设置为 `true` 可以输出更多详细信息。

## 属性说明

除了上述列出的参数外，还导出以下属性：

* `ids` - 区域 ID 列表。
* `zones` - 可用区列表。每个元素包含以下属性：
  * `id` - 区域的 ID。
  * `slb_slave_zone_ids` - SLB 主可用区中的从可用区 ID 列表。
  * `local_name` - 次要区域的名称。
  * `computed_attribute_example` - 计算属性示例（如果适用）。