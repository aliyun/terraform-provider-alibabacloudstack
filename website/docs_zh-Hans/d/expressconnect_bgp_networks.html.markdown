---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgpnetworks"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-bgpnetworks"
description: |-
  提供阿里云账号下拥有的expressconnect bgpnetworks列表。
---

# alibabacloudstack\_expressconnect\_bgpnetworks

此数据源提供根据指定过滤条件列出的阿里云账号下的expressconnect bgpnetworks资源列表。

## 示例用法
```
data "alibabacloudstack_expressconnect_bgpnetworks" "default" {
  router_id = "vbr-bp1d8yqvswjk4qr6i***"
  dst_cidr_block = "192.168.0.0/16"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - 已宣告的BGP网络的id列表。
  * `router_id` - (必填, 强制新建) - 路由器的ID。
  * `dst_cidr_block` - (选填) - 已宣告的BGP网络。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `bgp_networks` - 已宣告的BGP网络列表。
    * `id` - 已宣告的BGP网络的id。
    * `dst_cidr_block` - 已宣告的BGP网络。
    * `router_id` - 路由器的ID。
    * `status` - 已宣告的BGP网络状态。
