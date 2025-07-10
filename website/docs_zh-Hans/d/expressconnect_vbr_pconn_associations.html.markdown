---
subcategory: "Express Connect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_vbrpconnassociations"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-vbrpconnassociations"
description: |-
  提供阿里云账号下拥有的expressconnect vbrpconnassociations列表。
---

# alibabacloudstack\_expressconnect\_vbrpconnassociations

此数据源提供根据指定过滤条件列出的阿里云账号下的expressconnect vbrpconnassociations资源列表。

## 示例用法
```
data "alibabacloudstack_expressconnect_vbrpconnassociations" "example" {
  vbr_id = "vbr-bp1d8yixxxxxxxxxxx"
}
```

## 参数参考
以下参数是支持的：
  * `vbr_id` - (必填) - VBR实例ID。
  * `physical_connection_id` - (选填) - 物理专线实例ID。
  * `vlan_id` - (选填) - VBR的VLAN ID，取值范围：**0～2999**。 > 只有物理专线的所有者可以指定该参数，同一条物理专线下的两个VBR的VLAN ID不能相同。

## Attributes Reference
除了上述参数外，还导出以下属性：
  * `vbr_pconn_associations` - VBR pconn association 列表。
    * `circuit_code` - 运营商为物理专线提供的电路编码。 > 只有物理专线的所有者可以指定该参数。
    * `enable_ipv6` - 是否启用IPv6。取值：- **true**：开启。- **false**（默认值）：关闭。
    * `local_gateway_ip` - VBR实例的阿里云侧互联IP。
    * `local_ipv6_gateway_ip` - VBR实例的阿里云侧互联IPv6地址。
    * `peer_gateway_ip` - VBR实例的客户侧互联IP。- 该属性只允许VBR所有者指定或修改。- 为物理专线所有者创建VBR实例时必填。
    * `peer_ipv6_gateway_ip` - VBR实例的客户侧互联IPv6地址。- 该属性只允许VBR所有者指定或者修改。- 为物理专线所有者创建VBR实例时必填。
    * `peering_ipv6_subnet_mask` - VBR实例的阿里云侧和客户侧互联IPv6的子网掩码。两个IPv6地址必须位于同一个子网中。
    * `peering_subnet_mask` - VBR实例的阿里云侧和客户侧互联IP的子网掩码。两个IP地址必须位于同一个子网中。
    * `physical_connection_id` - 物理专线实例ID。
    * `status` - 代表资源状态的资源属性字段
    * `vbr_id` - VBR实例ID。
    * `vlan_id` - VBR的VLAN ID，取值范围：**0～2999**。 > 只有物理专线的所有者可以指定该参数，同一条物理专线下的两个VBR的VLAN ID不能相同。
