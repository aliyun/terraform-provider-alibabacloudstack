---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_ha"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbr_ha"
description: |-
  高速通道边界路由器快速倒换组
---

# alibabacloudstack_vbr_ha

使用Provider配置的凭证在指定的资源集下编排高速通道虚拟边界路由器下的快速倒换组。

## 示例用法
```
variable "name" {
  default = "tf-testaccexpressconnect-vbrha1926"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default1" {
	physical_connection_id =     ""
	vlan_id =                    1926
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_1"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default2" {
	physical_connection_id =     ""
	vlan_id =                    1929
	local_gateway_ip =           "10.0.1.1"
	peer_gateway_ip =            "10.0.1.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}_2"
}

resource "alibabacloudstack_expressconnect_vbr_ha" "default" {
  name = "tf-testaccexpressconnect-vbrha1926"
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default1.id}"
  peer_vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default2.id}"
  description = "tf-testaccexpressconnect-vbrha1926"
}
```

## 参数参考

支持以下参数：
  * `name` - (必填，变更时重建) - 快速倒换组名称。
  * `description` - (选填，变更时重建) - 快速倒换组的描述。 
  * `vbr_id` - (必填，变更时重建) - 高速通道虚拟边界路由器ID。
  * `peer_vbr_id` - (必填，变更时重建) - 对端高速通道虚拟边界路由器ID。

## 属性参考

除了上述所有参数外，无输出只读属性。