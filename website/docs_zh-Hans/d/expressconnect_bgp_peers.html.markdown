---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgppeers"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-bgppeers"
description: |-
  提供阿里云账号下拥有的expressconnect bgppeers列表。
---

# alibabacloudstack\_expressconnect\_bgppeers

此数据源提供根据指定过滤条件列出的阿里云账号下的expressconnect bgppeers资源列表。

## 示例用法

```
resource "alibabacloudstack_expressconnect_bgp_peer" "default" {
  bgp_group_id =   "${alibabacloudstack_expressconnect_bgp_group.default.id}"
  router_id =      "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  enable_bfd =      "true"
  peer_ip_address = "192.168.0.1"
  bfd_multi_hop =   "10"
}

data "alibabacloudstack_expressconnect_bgp_peers" "default" {
  
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - BGP邻居的ID列表。
  * `router_id` - (选填) - 路由器的ID。
  * `region_id` - (选填) - BGP组所属的地域ID。
  * `bgp_group_id` - (选填) - BGP组的ID。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `bgp_peers` - BGP邻居得列表。
    * `id` - BGP邻居的ID。
    * `auth_key` - BGP组的认证密钥。
    * `bfd_multi_hop` - 反射次数
    * `bgp_group_id` - BGP组的ID。
    * `bgp_peer_id` - BGP邻居的ID。
    * `bgp_peer_name` - BGP邻居的名称。
    * `bgp_status` - BGP的连接状态，包含以下状态：* creating：创建中。* working：使用中。* modifying：修改中。* deleting：删除中。* deleted：已删除。
    * `description` - BGP组的描述。 
    * `enable_bfd` - 是否开启了BFD协议。
    * `hold` - 保持时间。
    * `ip_version` - IP版本
    * `is_fake` - 是否启用了Fake AS号。
    * `keepalive` - 保活时间。
    * `local_asn` - 本地ASN号
    * `peer_asn` - BGP邻居的ASN。
    * `peer_ip_address` - BGP邻居的IP地址。
    * `region_id` - BGP组所属的地域ID。
    * `route_limit` - 路由限制。
    * `router_id` - 路由器的ID。
    * `status` - BGP邻居的状态
