---
subcategory: "ExpressConnect"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_expressconnect_virtualborderrouter"
sidebar_current: "docs-Alibabacloudstack-expressconnect-virtualborderrouter"
description: |- 
  编排高速通道虚拟边界路由器
---

# alibabacloudstack_expressconnect_virtualborderrouter
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_express_connect_virtual_border_router`

使用Provider配置的凭证在指定的资源集下编排高速通道虚拟边界路由器。

## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_express_connect_physical_connections" "nameRegex" {
  name_regex = "^my-PhysicalConnection"
}

resource "alibabacloudstack_expressconnect_virtualborderrouter" "example" {
  local_gateway_ip           = "10.0.0.1"
  peer_gateway_ip            = "10.0.0.2"
  peering_subnet_mask        = "255.255.255.252"
  physical_connection_id     = data.alibabacloudstack_express_connect_physical_connections.nameRegex.connections.0.id
  virtual_border_router_name = "example_value"
  vlan_id                    = 1
  min_rx_interval            = 1000
  min_tx_interval           = 1000
  detect_multiplier         = 10
}
```

## 参数说明

支持以下参数：

* `associated_physical_connections` - (选填) 关联的物理专线信息。
* `bandwidth` - (选填) VBR实例的带宽。
* `circuit_code` - (选填) 运营商为物理专线提供的电路编码。
* `description` - (选填) VBR的描述信息。长度为2到256个字符，必须以字母或汉字开头，不能以`http://`或`https://`开头。
* `detect_multiplier` - (选填) 检测时间倍数。即接收方允许发送方发送报文的最大连接丢包数，用来检测链路是否正常。取值范围：**3～10**。
* `enable_ipv6` - (选填) 是否启用IPv6。有效值：
  - `true`: 开启。
  - `false`: 关闭。
* `local_gateway_ip` - (必填) VBR实例的阿里云侧互联IPv4地址。
* `local_ipv6_gateway_ip` - (选填) VBR实例的阿里云侧互联IPv6地址。
* `min_rx_interval` - (选填) 配置BFD报文的接收间隔。取值范围：**200～1000**，单位为ms。
* `min_tx_interval` - (选填) 配置BFD报文的发送间隔。取值范围：**200～1000**，单位为ms。
* `peer_gateway_ip` - (必填) VBR实例的客户侧互联IPv4地址。
* `peer_ipv6_gateway_ip` - (选填) VBR实例的客户侧互联IPv6地址。
* `peering_ipv6_subnet_mask` - (选填) VBR实例的阿里云侧互联IPv6与客户侧互联IPv6的子网掩码。
* `peering_subnet_mask` - (必填) VBR实例的阿里云侧互联IPv4和客户侧互联IPv4的子网掩码。
* `physical_connection_id` - (必填，变更时重建) VBR所属的物理专线的ID。
* `status` - (选填) 资源状态。有效值：
  - `active`
  - `deleting`
  - `recovering`
  - `terminated`
  - `terminating`
  - `unconfirmed`
* `vbr_owner_id` - (选填) VBR实例所有者的账号ID。默认为登录的阿里云账号ID。
* `virtual_border_router_name` - (选填) VBR实例名称。长度为2到128个字符，必须以字母或汉字开头，可以包含数字、下划线(`_`)和连字符(`-`)，但不能以`http://`或`https://`开头。
* `vlan_id` - (必填) VBR实例的VLAN ID。有效范围：**0 到 2999**。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform中的虚拟边界路由器资源ID。
* `route_table_id` - VBR的路由表ID。
* `detect_multiplier` - 检测时间倍数。即接收方允许发送方发送报文的最大连接丢包数，用来检测链路是否正常。取值范围：**3～10**。
* `enable_ipv6` - 是否启用IPv6。有效值：
  - `true`: 开启。
  - `false`: 关闭。
* `min_rx_interval` - 配置BFD报文的接收间隔。取值范围：**200～1000**，单位为ms。
* `min_tx_interval` - 配置BFD报文的发送间隔。取值范围：**200～1000**，单位为ms。
* `status` - 资源状态。

### 超时设置

`timeouts` 块允许您为某些操作指定 [超时时间](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts)：

* `update` - (默认为2分钟)用于更新虚拟边界路由器时。

## 导入

可以通过id导入Express Connect虚拟边界路由器，例如：

```bash
$ terraform import alibabacloudstack_expressconnect_virtualborderrouter.example <id>
```