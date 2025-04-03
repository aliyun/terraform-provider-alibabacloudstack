---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpngateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpngateways"
description: |- 
  查询VPN网关实例
---

# alibabacloudstack_vpngateway_vpngateways
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_gateways`

根据指定过滤条件列出当前凭证权限可以访问的VPN网关实例列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccVpcGatewayConfig6349"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  vpc_name   = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
  name                = "${var.name}"
  vpc_id             = "${alibabacloudstack_vswitch.default.vpc_id}"
  bandwidth          = "10"
  enable_ssl         = true
  instance_charge_type = "PostPaid"
  description        = "${var.name}"
  vswitch_id         = "${alibabacloudstack_vswitch.default.id}"
}

data "alibabacloudstack_vpngateway_vpngateways" "default" {
  ids             = ["${alibabacloudstack_vpn_gateway.default.id}"]
  vpc_id          = "${alibabacloudstack_vpc.default.id}"
  status          = "Active"
  business_status = "Normal"
  name_regex      = "^${var.name}$"
  output_file     = "/tmp/vpns.txt"
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选) VPN网关ID列表。这将作为过滤条件使用。
* `vpc_id` - (可选，变更时重建) VPN网关所属的VPC的ID。
* `name_regex` - (可选，变更时重建) 用于按名称筛选结果的正则表达式字符串。
* `status` - (可选，变更时重建) 资源状态。有效值包括："Init"、"Provisioning"、"Active"、"Updating"、"Deleting"。
* `business_status` - (可选，变更时重建) VPN网关的支付状态。有效值包括："Normal"、"FinancialLocked"。

## 属性参考

除了上述参数外，还导出以下属性：

* `names` - 匹配的VPN网关名称列表。
* `gateways` - 匹配的VPN网关列表。每个元素包含以下属性：
  * `id` - VPN网关的ID。
  * `vpc_id` - 所属VPC的ID。
  * `internet_ip` - 公共IP地址。
  * `create_time` - 创建时间。
  * `end_time` - 到期时间。
  * `specification` - 规格。
  * `name` - 名称。
  * `description` - 描述。
  * `status` - 状态。
  * `business_status` - 支付状态。
  * `instance_charge_type` - 计费类型。可能的值为："PrePaid"(预付费)、"PostPaid"(后付费)。
  * `enable_ipsec` - 是否启用IPsec-VPN功能。
  * `enable_ssl` - 是否启用SSL-VPN功能。
  * `ssl_connections` - 允许的最大SSL-VPN客户端连接数。