---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_virtualborderrouters"
sidebar_current: "docs-Alibabacloudstack-datasource-expressconnect-virtualborderrouters"
description: |- 
  查询高速通道虚拟边界路由器
---

# alibabacloudstack_expressconnect_virtualborderrouters
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_express_connect_virtual_border_routers`

根据指定过滤条件列出当前凭证权限可以访问的高速通道虚拟边界路由器列表。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testAccExpressConnectVirtualBorderRoutersTest1519"
}

data "alibabacloudstack_express_connect_physical_connections" "nameRegex" {
  name_regex = "^preserved-NODELETING"
}

resource "alibabacloudstack_express_connect_virtual_border_router" "default" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = "pc-9wdbvb1hkf44szgqvgnor"
  virtual_border_router_name = var.name
  vlan_id                    = 960
  min_rx_interval            = 1000
  min_tx_interval            = 1000
  detect_multiplier          = 10
}

data "alibabacloudstack_express_connect_virtual_border_routers" "default" {
  ids = [alibabacloudstack_express_connect_virtual_border_router.default.id]
}

output "virtual_border_router_id" {
  value = data.alibabacloudstack_express_connect_virtual_border_routers.default.routers.0.id
}

data "alibabacloudstack_express_connect_virtual_border_routers" "byName" {
  name_regex = "^tf-testAccExpressConnectVirtualBorderRoutersTest1519"
}

output "virtual_border_router_id_by_name" {
  value = data.alibabacloudstack_express_connect_virtual_border_routers.byName.routers.0.id
}

data "alibabacloudstack_express_connect_virtual_border_routers" "filtered" {
  filter {
    key    = "physical_connection_id"
    values = ["pc-9wdbvb1hkf44szgqvgnor"]
  }
  filter {
    key    = "status"
    values = ["active"]
  }
}

output "filtered_virtual_border_router_id" {
  value = data.alibabacloudstack_express_connect_virtual_border_routers.filtered.routers.0.id
}
```

## 参数参考

以下参数是支持的：

* `filter` - (可选，变更时重建) 自定义过滤器块，用于筛选 VBR 资源。每个过滤器包含以下字段：
  * `key` - (必填) 要过滤的字段的键。
  * `values` - (必填) 给定字段接受的值集。
* `status` - (可选，变更时重建) 代表资源状态的资源属性字段。有效值：`active`, `deleting`, `recovering`, `terminated`, `terminating`, `unconfirmed`。
* `ids` - (可选，变更时重建) 虚拟边界路由器 ID 列表，用于筛选特定的 VBR 资源。
* `name_regex` - (可选，变更时重建) 用于按虚拟边界路由器名称筛选结果的正则表达式字符串。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 虚拟边界路由器名称列表。
* `routers` - Express Connect 虚拟边界路由器列表。每个元素包含以下属性：
  * `access_point_id` - 物理专线接入点的 ID。
  * `activation_time` - VBR 第一次激活的时间。
  * `circuit_code` - 运营商为物理专线提供的电路编码。
  * `cloud_box_instance_id` - 与 VBR 关联的云盒实例的 ID。
  * `create_time` - VBR 的创建时间。
  * `description` - VBR 的描述信息。
  * `detect_multiplier` - 检测时间倍数。即接收方允许发送方发送报文的最大连接丢包数，用来检测链路是否正常。取值范围：**3～10**。
  * `ecc_id` - 高速上云服务实例 ID。
  * `enable_ipv6` - 是否启用 IPv6。- **true**：开启。- **false**：关闭。
  * `id` - 虚拟边界路由器的 ID。
  * `local_gateway_ip` - VBR 实例的阿里云侧互联 IPv4 地址。
  * `local_ipv6_gateway_ip` - VBR 实例的阿里云侧互联 IPv6 地址。
  * `min_rx_interval` - 配置 BFD 报文的接收间隔，取值范围：**200～1000**，单位为 ms。
  * `min_tx_interval` - 配置 BFD 报文的发送间隔，取值范围：**200～1000**，单位为 ms。
  * `payment_vbr_expire_time` - VBR 的计费到期时间。
  * `peer_gateway_ip` - VBR 实例的客户侧互联 IPv4 地址。
  * `peer_ipv6_gateway_ip` - VBR 实例的客户侧互联 IPv6 地址。
  * `peering_ipv6_subnet_mask` - VBR 实例的阿里云侧互联 IPv6 与客户侧互联 IPv6 的子网掩码。
  * `peering_subnet_mask` - VBR 实例的阿里云侧互联 IPv4 和客户侧互联 IPv4 的子网掩码。
  * `physical_connection_business_status` - 物理专线业务状态。- **Normal**：正常。- **FinancialLocked**：欠费锁定。
  * `physical_connection_id` - VBR 所属的物理专线的 ID。
  * `physical_connection_owner_uid` - 物理专线所属的账号 ID。
  * `physical_connection_status` - 物理专线状态。- **Initial**：申请中。- **Approved**：审批通过。- **Allocating**：正在分配资源。- **Allocated**：接入施工中。- **Confirmed**：等待用户确认。- **Enabled**：已开通。- **Rejected**：申请被拒绝。- **Canceled**：已取消。- **Allocation Failed**：资源分配失败。- **Terminated**：已终止。
  * `recovery_time` - VBR 最近一次从 Terminated 状态恢复到 Active 状态的时间。
  * `route_table_id` - VBR 的路由表 ID。
  * `status` - 代表资源状态的资源属性字段。
  * `termination_time` - VBR 最近一次被终止的时间。
  * `type` - VBR 类型。
  * `virtual_border_router_id` - VBR 的 ID。
  * `virtual_border_router_name` - VBR 实例名称。
  * `vlan_id` - VBR 实例的 VLAN ID。
  * `vlan_interface_id` - VBR 的路由器接口的 ID。