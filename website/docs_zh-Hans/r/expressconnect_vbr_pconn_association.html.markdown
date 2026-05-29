---
subcategory: "高速通道"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbr_pconn_association"
sidebar_current: "docs-Alibabacloudstack-expressconnect-vbrpconnassociation"
description: |-
  Provides a expressconnect Vbrpconnassociation resource.
---

# alibabacloudstack\_expressconnect\_vbrpconnassociation

Provides a expressconnect Vbrpconnassociation resource.

## 示例用法

### 基本关联

```hcl
variable "name" {
  default = "tf-testaccexpressconnect-vbrpconn1073"
}

resource "alibabacloudstack_expressconnect_virtualborderrouter" "default" {
  physical_connection_id = alibabacloudstack_expressconnect_physicalconnection.default.id
  vlan_id                = 1073
  local_gateway_ip       = "10.0.0.1"
  peer_gateway_ip        = "10.0.0.2"
  peering_subnet_mask    = "255.255.255.252"
  name                   = var.name
}

resource "alibabacloudstack_expressconnect_vbr_pconn_association" "default" {
  physical_connection_id = alibabacloudstack_expressconnect_physicalconnection.default.id
  vbr_id                 = alibabacloudstack_expressconnect_virtualborderrouter.default.id
  vlan_id                = "1076"
  local_gateway_ip       = "10.100.0.1"
  peer_gateway_ip        = "10.100.0.2"
  peering_subnet_mask    = "255.255.255.0"
}
```

### 启用 IPv6 的关联

```hcl
resource "alibabacloudstack_expressconnect_vbr_pconn_association" "ipv6" {
  physical_connection_id   = alibabacloudstack_expressconnect_physicalconnection.default.id
  vbr_id                   = alibabacloudstack_expressconnect_virtualborderrouter.default.id
  vlan_id                  = "1076"
  local_gateway_ip         = "10.100.0.1"
  peer_gateway_ip          = "10.100.0.2"
  peering_subnet_mask      = "255.255.255.0"
  enable_ipv6              = true
  local_ipv6_gateway_ip    = "2408:4004:cc:500::1"
  peer_ipv6_gateway_ip     = "2408:4004:cc:500::2"
  peering_ipv6_subnet_mask = "2408:4004:cc:500::/56"
}
```

## 参数参考

支持以下参数：
  * `physical_connection_id` - (必填, 强制新建) - 物理专线实例ID。
  * `vbr_id` - (必填, 强制新建) - VBR实例ID。
  * `vlan_id` - (必填, 强制新建) - VBR 的 VLAN ID。取值范围：`0` 到 `2999`。只有物理专线的所有者可以指定该参数，同一条物理专线下的两个 VBR 的 VLAN ID 不能相同。
  * `local_gateway_ip` - (选填, 强制新建) - VBR实例的阿里云侧互联IP。
  * `peer_gateway_ip` - (选填, 强制新建) - VBR 实例的客户侧互联 IP。该属性只允许 VBR 所有者指定或修改。为物理专线所有者创建 VBR 实例时必填。
  * `peering_subnet_mask` - (选填, 强制新建) - VBR 实例的阿里云侧和客户侧互联 IP 的子网掩码。两个 IP 地址必须位于同一子网中。
  * `enable_ipv6` - (选填, 强制新建) - 是否启用 IPv6。取值：`true`（启用）、`false`（默认值，关闭）。
  * `local_ipv6_gateway_ip` - (选填, 强制新建) - VBR实例的阿里云侧互联IPv6地址。
  * `peer_ipv6_gateway_ip` - (选填, 强制新建) - VBR 实例的客户侧互联 IPv6 地址。该属性只允许 VBR 所有者指定或修改。为物理专线所有者创建 VBR 实例时必填。注意：当 `enable_ipv6` 为 `true` 时必填，且需与 `local_ipv6_gateway_ip` 和 `peering_ipv6_subnet_mask` 同时使用。
  * `peering_ipv6_subnet_mask` - (选填, 强制新建) - VBR 实例的阿里云侧和客户侧互联 IPv6 的子网掩码。两个 IPv6 地址必须位于同一子网中。注意：当 `enable_ipv6` 为 `true` 时必填，且需与 `local_ipv6_gateway_ip` 和 `peer_ipv6_gateway_ip` 同时使用。

## 属性参考

除了上述参数外，还导出以下属性：

  * `id` - 资源的 ID，格式为 `<physical_connection_id>:<vbr_id>`。
  * `circuit_code` - 运营商为物理专线提供的电路编码。只有物理专线的所有者可以指定该参数。
  * `status` - VBR-PCCN 关联的状态。可能的值包括 `Creating`、`Associated`、`UnAssociating`。
  * `vlan_id` - VBR 的 VLAN ID。

## Import

VBR-PCCN 关联可以使用 `physical_connection_id` 和 `vbr_id`（以冒号分隔）进行导入，例如：

```
$ terraform import alibabacloudstack_expressconnect_vbr_pconn_association.example pc-abc12345:vbr-xyz67890
```
