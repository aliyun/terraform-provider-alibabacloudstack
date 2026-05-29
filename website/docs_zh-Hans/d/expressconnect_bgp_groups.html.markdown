---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgp_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-bgp-groups"
description: |-
  提供阿里云账号下拥有的expressconnect bgpgroups列表。
---

# alibabacloudstack\_expressconnect\_bgpgroups

此数据源提供根据指定过滤条件列出的阿里云账号下的expressconnect bgpgroups资源列表。

## 示例用法
```
resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      "65534"
	peer_asn =       "10"
	router_id =      "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}$"
}

data "alibabacloudstack_expressconnect_bgp_groups" "default" {
	router_id = "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}"
}
```

## 参数参考
以下参数是支持的：
  * `ids` - (选填) - BGP组ID列表。
  * `router_id` - (必填) - VBR的ID。
  * `region_id` - (选填) - BGP组所属的Region ID。
  * `name_regex` - (选填) - 用于过滤结果，支持BGP组的名称。
  * `description_regex` - (选填) - 用于过滤结果，支持BGP组的描述。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `bgp_groups` - BGP组列表。
    * `id` - BGP组的ID。
    * `auth_key` - BGP组使用的密钥。
    * `bgp_group_id` - BGP组的ID。
    * `bgp_group_name` - BGP组的名称。
    * `description` - BGP组的描述。
    * `hold` - 等待BGP消息传入的保持时间。如果超过保持时间还没有消息传入，则认为BGP邻居断开了连接。
    * `ip_version` - IP版本
    * `is_fake` - AS号是否为假。
    * `keepalive` - 保活时间。
    * `local_asn` - 本端AS号。
    * `peer_asn` - 侧设备的AS号。
    * `route_limit` - 路由限制。
    * `router_id` - VBR的ID。
    * `region_id` - BGP组所属的Region ID。
    * `status` - 代表资源状态的资源属性字段
