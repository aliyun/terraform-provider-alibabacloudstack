---
subcategory: "VPNGateway"
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
variable "name" {
    default = "tf-testaccvpn_gatewaycustomer_gateway50970"
}

resource "alibabacloudstack_vpngateway_customergateway" "default" {
  ip_address           = "1.1.1.1" # 必须是一个有效的公共 IP 地址
  customer_gateway_name = var.name # 自定义客户网关名称
  description         = "This is a test customer gateway." # 客户网关的描述信息
}
```

## 参数说明

支持以下参数：

* `ip_address` - (必填，变更时重建) 客户网关的 IP 地址。这必须是一个有效的公共 IP 地址。
* `customer_gateway_name` - (可选) 客户网关的名称。默认为 null。如果不指定，Terraform 将自动生成一个唯一的名称用于客户网关。
* `description` - (可选) 客户网关的描述信息。默认为 null。

> **注意**：如果未提供 `customer_gateway_name`，系统将生成一个默认名称。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `id` - 客户网关实例的唯一标识符(ID)。
* `customer_gateway_name` - 实际设置或生成的客户网关名称。
* `status` - 客户网关的状态，例如 `Available` 或 `Deleting` 等。
* `name` - 如果在参数中未明确指定 `customer_gateway_name`，此字段将返回生成的默认名称。