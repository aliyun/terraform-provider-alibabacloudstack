---
subcategory: "专有网络 VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_customergateway"
sidebar_current: "docs-Alibabacloudstack-vpngateway-customergateway"
description: |-
  编排VPN网关客户网关
---

# alibabacloudstack_vpngateway_customergateway
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_customer_gateway`

使用Provider配置的凭证在指定的资源集编排VPN网关客户网关。

## 示例用法

以下是一个完整的示例，展示如何创建一个自定义客户网关资源：

```hcl
resource "alibabacloudstack_vpngateway_customergateway" "default" {
  ip_address           = "1.1.1.1"
  customer_gateway_name = "example-customer-gateway"
  description          = "This is a test customer gateway."
}
```

## 参数说明

支持以下参数：

* `ip_address` - (必填，变更时重建) 客户网关的 IP 地址。如果计划创建公网类型的 IPsec 连接，请输入公网 IP 地址；如果计划创建 VPC 类型的 IPsec 连接，请输入私网 IP 地址。修改此参数会强制重新创建资源。
* `customer_gateway_name` - (可选) 客户网关的名称。名称长度为 2~128 个字符，可以包含字母、数字、半角句号（.）、下划线（_）和短横线（-）。必须以字母或中文开头，不能以 http:// 或 https:// 开头。与 `name` 参数互斥。
* `name` - (已弃用) 客户网关的名称。该参数已弃用，将在后续版本中移除。请使用 `customer_gateway_name` 替代。与 `customer_gateway_name` 参数互斥。
* `description` - (可选) 客户网关的描述信息。描述信息长度为 2~256 个字符，必须以字母或中文开头，不能以 http:// 或 https:// 开头。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 客户网关实例的唯一标识符(ID)。
* `customer_gateway_name` - 实际设置的客户网关名称。

## Import

VPN 客户网关可以使用客户网关 ID 进行导入，例如：

```
$ terraform import alibabacloudstack_vpngateway_customergateway.example cgw-bp1jrawp82av6bws9****
```
