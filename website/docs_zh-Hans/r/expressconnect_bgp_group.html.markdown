---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_bgpgroup"
sidebar_current: "docs-Alibabacloudstack-expressconnect-bgpgroup"
description: |-
  Provides a expressconnect Bgpgroup resource.
---

# alibabacloudstack\_expressconnect\_bgpgroup

使用Provider配置的凭证在指定的资源集下编排高速通道虚拟边界路由器下的BGP组。

## 示例用法
```
resource "alibabacloudstack_expressconnect_bgp_group" "default" {
	bgp_group_name = "${var.name}"
	description =    "${var.name}"
	local_asn =      "65534"
	peer_asn =       "10"
	router_id =      "${alibabacloudstack_expressconnect_virtualborderrouter.default.id}$"
}
```

## 参数参考

支持以下参数：
  * `auth_key` - (选填) - BGP组使用的密钥。
  * `bgp_group_id` - (选填) - BGP组的ID。
  * `bgp_group_name` - (选填) - BGP组的名称。
  * `description` - (选填) - BGP组的描述。
  * `ip_version` - (选填) - IP版本
  * `local_asn` - (选填) - 本端AS号。
  * `peer_asn` - (必填) - 侧设备的AS号。
  * `region_id` - (选填) - BGP组所属的Region ID。
  * `router_id` - (必填) - VBR的ID。
  * `status` - (选填) - 代表资源状态的资源属性字段

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `bgp_group_id` - BGP组的ID。
  * `hold` - 等待BGP消息传入的保持时间。如果超过保持时间还没有消息传入，则认为BGP邻居断开了连接。
  * `ip_version` - IP版本
  * `is_fake` - AS号是否为假。
  * `keepalive` - 保活时间。
  * `region_id` - BGP组所属的Region ID。
  * `route_limit` - 路由限制。
  * `status` - 代表资源状态的资源属性字段
