---
subcategory: "Virtual Private Cloud (VPC)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpngateway"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpngateway"
description: |- 
  编排VPN网关实例
---

# alibabacloudstack_vpngateway_vpngateway
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_gateway`

使用 Provider 配置的凭证在指定的资源集编排 VPN 网关实例。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccvpn_gatewayvpn_gateway87805"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
  enable_details = true
}

resource "alibabacloudstack_vpc_vpc" "default" {
  vpc_name   = "${var.name}_vpc"
  cidr_block = "172.16.0.0/16"
}

resource "alibabacloudstack_vpc_vswitch" "default" {
  name       = "${var.name}_vsw"
  vpc_id     = "${alibabacloudstack_vpc_vpc.default.id}"
  cidr_block = "172.16.0.0/24"
  zone_id    = "${data.alibabacloudstack_zones.default.zones.0.id}"
}

resource "alibabacloudstack_vpngateway_vpngateway" "default" {
  description        = "test_vpn"
  vpn_gateway_name  = "test_vpn"
  bandwidth         = "10"
  vswitch_id        = "${alibabacloudstack_vpc_vswitch.default.id}"
  vpc_id            = "${alibabacloudstack_vpc_vpc.default.id}"
  ssl_vpn           = true
  instance_charge_type = "PostPaid"

  tags = {
    Created = "TF"
    Purpose = "acceptance_test"
  }
}
```

## 参数说明

支持以下参数：

* `vpc_id` - (必填，变更时重建) VPN 网关所属的 VPC 的 ID。
* `bandwidth` - (必填) VPN 网关的公网带宽。单位：Mbps。有效值：5, 10, 20, 50, 100, 200, 500, 1000。
* `vpn_gateway_name` - (可选) VPN 网关的名称。长度为 1~128 个字符。
* `instance_charge_type` - (可选，变更时重建) 实例的计费方式。有效值：`PrePaid`（包年包月）、`PostPaid`（按量付费）。默认值：`PostPaid`。
* `period` - (可选) 购买时长。有效值：1~9、12、24、36。默认值：1。
* `vswitch_id` - (可选，变更时重建) VPN 网关所属的交换机的 ID。
* `ipsec_vpn` - (可选) 是否开启 IPsec-VPN 功能。
* `ssl_vpn` - (可选) 是否开启 SSL-VPN 功能。
* `ssl_max_connections` - (可选) 最大 SSL-VPN 并发连接数。默认值：5。
* `description` - (可选) VPN 网关的描述信息。长度为 2~256 个字符。
* `tags` - (可选) 要分配给资源的标签映射。
* `name` - (可选，已废弃) 此参数已废弃，请使用 `vpn_gateway_name`。
* `enable_ipsec` - (可选，已废弃) 此参数已废弃，请使用 `ipsec_vpn`。
* `enable_ssl` - (可选，已废弃) 此参数已废弃，请使用 `ssl_vpn`。
* `ssl_connections` - (可选，已废弃) 此参数已废弃，请使用 `ssl_max_connections`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - VPN 网关的 ID。
* `internet_ip` - VPN 网关的公网 IP 地址。
* `status` - VPN 网关的状态。有效值：`provisioning`、`init`、`active`。
* `business_status` - VPN 网关的付费状态。
* `vpn_gateway_name` - VPN 网关的名称。
* `ipsec_vpn` - 是否开启了 IPsec-VPN 功能。
* `ssl_vpn` - 是否开启了 SSL-VPN 功能。
* `ssl_max_connections` - 最大 SSL-VPN 并发连接数。

## Import

VPN Gateway 可以使用 VpnGatewayId 导入，例如：

```
$ terraform import alibabacloudstack_vpngateway_vpngateway.example vgw-xxxxxxxxx
```
