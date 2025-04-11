---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpngateway"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpngateway"
description: |- 
  编排VPN网关实例
---

# alibabacloudstack_vpngateway_vpngateway
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_gateway`

使用Provider配置的凭证在指定的资源集编排VPN网关实例。

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
  enable_ssl        = true
  instance_charge_type = "PostPaid"

  tags = {
    Created = "TF"
    Purpose = "acceptance_test"
  }
}
```

## 参数说明

支持以下参数：

* `vpn_gateway_name` - (可选) VPN 网关的名称。
* `vpc_id` - (必填，变更时重建) VPN 网关所属的 VPC 的 ID。
* `instance_charge_type` - (可选，变更时重建) 实例的计费方式。有效值：
  * **PrePaid**：包年包月。
  * **PostPaid**：按量付费。
  默认值：**PostPaid**。
* `period` - (可选) 购买时长。当 `instance_charge_type` 设置为 `PrePaid` 时，此参数是必填的。有效值：[1-9, 12, 24, 36]。默认值：1。
* `bandwidth` - (必填) VPN 网关的公网带宽。单位：Mbps。对于 PostPaid 实例的有效值：10, 100, 200。对于 PrePaid 实例的有效值：5, 10, 20, 50, 100, 200。
* `enable_ipsec` - (可选) 是否启用 IPsec-VPN 功能。默认值：**true**。
* `enable_ssl` - (可选) 是否启用 SSL-VPN 功能。默认值：**false**。
* `ssl_connections` - (可选) SSL-VPN 并发连接的最大数量。有效值：5, 10, 20, 50, 100, 200。默认值：**5**。仅当 `enable_ssl` 设置为 **true** 时此参数生效。
* `description` - (可选) VPN 网关的描述信息。
* `vswitch_id` - (可选，变更时重建) VPN 网关所属的交换机的 ID。
* `tags` - (可选) 要分配给资源的标签映射。
* `ipsec_vpn` - (可选) 是否开启 IPsec-VPN 功能。
* `ssl_vpn` - (可选) 是否开启 SSL-VPN 功能。
* `ssl_max_connections` - (可选) 最大 SSL-VPN 并发连接用户数的规格。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 资源的 ID。
* `internet_ip` - 公网 IP 地址。
* `status` - 资源的状态。有效值：
  * **Creating**：资源正在创建中。
  * **Available**：资源已创建并可以正常使用。
  * **Deleting**：资源正在删除中。
* `business_status` - VPN 网关的支付状态。有效值：
  * **Normal**：资源正常。
  * **Expired**：资源已过期。
  * **LockDown**：资源已被锁定。
* `name` - VPN 网关的名称。