---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnconnections"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpnconnections"
description: |- 
  查询VPN网关VPN连接
---

# alibabacloudstack_vpngateway_vpnconnections
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_connections`

根据指定过滤条件列出当前凭证权限可以访问的VPN网关VPN连接列表。

## 示例用法

```hcl
variable "name" {
  default = "tf-testAccVpnConnDataResource2865"
}

resource "alibabacloudstack_vpc" "default" {
  cidr_block = "172.16.0.0/12"
  name       = "${var.name}"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name             = "${var.name}"
}

resource "alibabacloudstack_vpn_gateway" "default" {
  name                = "${var.name}"
  vpc_id             = "${alibabacloudstack_vswitch.default.vpc_id}"
  bandwidth          = "10"
  enable_ssl         = true
  instance_charge_type = "PostPaid"
  description        = "test_create_description"
}

resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name         = "${var.name}"
  ip_address   = "41.104.22.229"
  description = "${var.name}"
}

resource "alibabacloudstack_vpn_connection" "default" {
  name               = "${var.name}"
  vpn_gateway_id    = "${alibabacloudstack_vpn_gateway.default.id}"
  customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
  local_subnet      = ["172.16.1.0/24"]
  remote_subnet     = ["10.4.0.0/24"]
  effect_immediately = true
  ike_config {
    ike_auth_alg = "sha1"
    ike_enc_alg  = "3des"
    ike_version  = "ikev2"
    ike_mode     = "aggressive"
    ike_lifetime = 8640
    psk          = "tf-testvpn1"
    ike_pfs      = "group2"
    ike_remote_id = "testbob1"
    ike_local_id = "testalice1"
  }
  
  ipsec_config {
    ipsec_pfs        = "group2"
    ipsec_enc_alg    = "aes"
    ipsec_auth_alg   = "sha1"
    ipsec_lifetime   = 86400
  }
}

data "alibabacloudstack_vpn_connections" "default" {
  ids                 = ["${alibabacloudstack_vpn_connection.default.id}"]
  vpn_gateway_id      = "${alibabacloudstack_vpn_gateway.default.id}"
  customer_gateway_id = "${alibabacloudstack_vpn_customer_gateway.default.id}"
  name_regex         = "${var.name}"
  output_file        = "vpn_connections_output.txt"
}
```

## 参数说明

以下参数是支持的：

* `ids` - （可选）要查询的 IPsec-VPN 连接的 ID 列表。如果指定了此参数，仅返回指定的资源。
* `vpn_gateway_id` - （可选，变更时重建）与 IPsec-VPN 连接关联的 VPN 网关的 ID。如果要按 VPN 网关过滤连接，这是必填的参数。
* `customer_gateway_id` - （可选，变更时重建）与 IPsec-VPN 连接关联的客户网关的 ID。如果要按客户网关过滤连接，这是必填的参数。
* `name_regex` - （可选，变更时重建）用于基于 IPsec-VPN 连接名称过滤结果的正则表达式。

## 属性说明

除了上述所有参数外，还导出以下属性：

* `ids` - 匹配的 IPsec-VPN 连接的 ID 列表。
* `names` - 匹配的 IPsec-VPN 连接的名称列表。
* `connections` - IPsec-VPN 连接详细信息列表。每个元素包含以下属性：
  * `id` - IPsec-VPN 连接的唯一标识符。
  * `customer_gateway_id` - 与此 IPsec-VPN 连接关联的客户网关的 ID。
  * `vpn_gateway_id` - 与此 IPsec-VPN 连接关联的 VPN 网关的 ID。
  * `name` - IPsec-VPN 连接的名称。
  * `local_subnet` - 此 IPsec-VPN 连接附加到的虚拟私有云 (VPC) 的 CIDR 块。
  * `remote_subnet` - 通过此 IPsec-VPN 连接连接的本地数据中心或其他网络的 CIDR 块。
  * `create_time` - 创建 IPsec-VPN 连接的时间。
  * `effect_immediately` - 表示在更新配置后是否立即启动 IPsec-VPN 协商。有效值：`true` 或 `false`。
  * `status` - 当前的 IPsec-VPN 连接状态。可能的值包括：
    * `ike_sa_not_established` - IKE 安全关联 (SA) 尚未建立。
    * `ike_sa_established` - IKE SA 已建立。
    * `ipsec_sa_not_established` - IPsec SA 尚未建立。
    * `ipsec_sa_established` - IPsec SA 已建立。
  * `ike_config` - 第一阶段协商的配置。它包括以下属性：
    * `psk` - 用于 IPsec-VPN 网关和客户网关之间身份验证的预共享密钥。
    * `ike_version` - IKE 协议的版本。例如，`ikev1` 或 `ikev2`。
    * `ike_mode` - IKE 第一阶段的协商模式。例如，`main` 或 `aggressive`。
    * `ike_enc_alg` - IKE 第一阶段协商中使用的加密算法。例如，`aes-128-cbc`、`aes-256-cbc`。
    * `ike_auth_alg` - IKE 第一阶段协商中使用的身份验证算法。例如，`sha1`、`sha256`。
    * `ike_pfs` - IKE 第一阶段协商中使用的 Diffie-Hellman 密钥交换算法。例如，`group2`、`group14`。
    * `ike_lifetime` - IKE 第一阶段协商结果的安全关联 (SA) 的生命周期。以秒为单位测量。
    * `ike_local_id` - 本地 (VPN 网关) 端点的身份识别。
    * `ike_remote_id` - 远程 (客户网关) 端点的身份识别。
  * `ipsec_config` - 第二阶段协商的配置。它包括以下属性：
    * `ipsec_enc_alg` - IPsec 第二阶段协商中使用的加密算法。例如，`aes-128-cbc`、`aes-256-cbc`。
    * `ipsec_auth_alg` - IPsec 第二阶段协商中使用的身份验证算法。例如，`hmac-sha1-96`、`hmac-sha256-128`。
    * `ipsec_pfs` - IPsec 第二阶段协商中使用的 Diffie-Hellman 密钥交换算法。例如，`group2`、`group14`。
    * `ipsec_lifetime` - IPsec 第二阶段协商结果的安全关联 (SA) 的生命周期。以秒为单位测量。