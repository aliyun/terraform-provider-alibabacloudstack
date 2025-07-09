---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbrpconnassociation"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbrpconnassociation"
description: |-
  Provides a expressconnect Vbrpconnassociation resource.
---

# alibabacloudstack\_expressconnect\_vbrpconnassociation

Provides a expressconnect Vbrpconnassociation resource.

## 示例用法
```
variable "name" {
  default = "tf-testaccexpressconnect-vbrpconn1073"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
	physical_connection_id =     ""
	vlan_id =                    1073
	local_gateway_ip =           "10.0.0.1"
	peer_gateway_ip =            "10.0.0.2"
	peering_subnet_mask =        "255.255.255.252"
	virtual_border_router_name = "${var.name}"
	enable_ipv6              = true
	local_ipv6_gateway_ip = "2408:4004:cc:400::1"
	peer_ipv6_gateway_ip= "2408:4004:cc:400::2"
	peering_ipv6_subnet_mask= "2408:4004:cc:400::/56"
}



resource "alibabacloudstack_expressconnect_vbr_pconn_association" "default" {
  vlan_id = "1076"
  peering_subnet_mask = "255.255.255.0"
  physical_connection_id = ""
  local_gateway_ip = "10.100.0.1"
  local_ipv6_gateway_ip = "2408:4004:cc:500::1"
  vbr_id = "${alibabacloudstack_express_connect_virtual_border_router.default.id}"
  peering_ipv6_subnet_mask = "2408:4004:cc:500::/56"
  peer_gateway_ip = "10.100.0.2"
  enable_ipv6 = "true"
  peer_ipv6_gateway_ip = "2408:4004:cc:500::2"
}
```

## 参数参考

支持以下参数：
  * `physical_connection_id` - (必填, 强制新建) - 物理专线实例ID。
  * `vbr_id` - (必填, 强制新建) - VBR实例ID。
  * `vlan_id` - (必填, 强制新建) - VBR的VLAN ID，取值范围：**0～2999**。 > 只有物理专线的所有者可以指定该参数，同一条物理专线下的两个VBR的VLAN ID不能相同。
  * `local_gateway_ip` - (选填, 强制新建) - VBR实例的阿里云侧互联IP。
  * `peer_gateway_ip` - (选填, 强制新建) - VBR实例的客户侧互联IP。- 该属性只允许VBR所有者指定或修改。- 为物理专线所有者创建VBR实例时必填。
  * `peering_subnet_mask` - (选填, 强制新建) - VBR实例的阿里云侧和客户侧互联IP的子网掩码。两个IP地址必须位于同一个子网中。
  * `enable_ipv6` - (选填, 强制新建) - 是否启用IPv6。取值：- **true**：开启。- **false**（默认值）：关闭。
  * `local_ipv6_gateway_ip` - (选填, 强制新建) - VBR实例的阿里云侧互联IPv6地址。
  * `peer_ipv6_gateway_ip` - (选填, 强制新建) - VBR实例的客户侧互联IPv6地址。- 该属性只允许VBR所有者指定或者修改。- 为物理专线所有者创建VBR实例时必填。
  * `peering_ipv6_subnet_mask` - (选填, 强制新建) - VBR实例的阿里云侧和客户侧互联IPv6的子网掩码。两个IPv6地址必须位于同一个子网中。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `circuit_code` - 运营商为物理专线提供的电路编码。 > 只有物理专线的所有者可以指定该参数。
  * `status` - 代表资源状态的资源属性字段
