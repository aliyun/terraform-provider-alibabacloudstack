---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnconnection"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnconnection"
description: |- 
  编排PN网关VPN连接
---

# alibabacloudstack_vpngateway_vpnconnection
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_vpn_connection`

使用Provider配置的凭证在指定的资源集编排PN网关VPN连接。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf-testaccVpnConnectionBaisc19277"
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

resource "alibabacloudstack_vpn_gateway" "foo" {
  name                 = "testAccVpnConfig_create"
  vpc_id               = alibabacloudstack_vpc.default.id
  bandwidth            = "10"
  enable_ssl           = true
  instance_charge_type = "PostPaid"
  description         = "test_create_description"
}

resource "alibabacloudstack_vpn_customer_gateway" "default" {
  name        = "${var.name}"
  ip_address  = "42.104.22.210"
  description = "testAccVpnConnectionDesc"
}

resource "alibabacloudstack_vpn_connection" "default" {
  vpn_connection_name = "${var.name}"
  vpn_gateway_id      = alibabacloudstack_vpn_gateway.foo.id
  customer_gateway_id = alibabacloudstack_vpn_customer_gateway.default.id
  local_subnet        = ["172.16.0.0/24", "172.16.1.0/24"]
  remote_subnet       = ["10.0.0.0/24", "10.0.1.0/24"]
  effect_immediately  = true
  ike_config {
    psk            = "tf-testvpn2"
    ike_version    = "ikev1"
    ike_mode       = "main"
    ike_enc_alg    = "des"
    ike_auth_alg   = "md5"
    ike_lifetime   = 86400
    ike_pfs        = "group1"
    ike_remote_id  = "testbob2"
    ike_local_id   = "testalice2"
  }
  ipsec_config {
    ipsec_pfs      = "group5"
    ipsec_enc_alg  = "des"
    ipsec_auth_alg = "md5"
    ipsec_lifetime = 8640
  }
}
```

## 参数参考

支持以下参数：
  * `customer_gateway_id` - (必填, 变更时重建) - 客户网关的 ID。
  * `vpn_gateway_id` - (必填, 变更时重建) - VPN 网关的 ID。
  * `name` - (选填) - IPsec 连接的名称，与 `vpn_connection_name` 功能相同。
  * `vpn_connection_name` - (选填) - IPsec 连接的名称。
  * `local_subnet` - (必填) - 虚拟私有云 (VPC) 的 CIDR 块。此参数用于第二阶段协商。
  * `remote_subnet` - (必填) - 本地数据中心的 CIDR 块。此参数用于第二阶段协商。
  * `effect_immediately` - (选填) - 是否立即启动 IPsec-VPN 协商。有效值：`true`(配置变更时触发重连)、`false`(有流量时触发重连)。默认值为 `true`。
  * `ike_config` - (选填) - 第一阶段协商的配置。
    
    * `psk` - (选填) - 用于验证 IPsec VPN 网关和客户网关之间身份的预共享密钥。
    * `ike_version` - (选填) - IKE 协议的版本。有效值：`ikev1` | `ikev2`。默认值：`ikev1`。
    * `ike_mode` - (选填) - IKE V1 的协商模式。有效值：`main`(主模式)| `aggressive`(积极模式)。默认值：`main`。
    * `ike_enc_alg` - (选填) - 第一阶段协商的加密算法。有效值：`aes` | `aes192` | `aes256` | `des` | `3des`。默认值：`aes`。
    * `ike_auth_alg` - (选填) - 第一阶段协商的身份验证算法。有效值：`md5` | `sha1` | `sha256` | `sha384` | `sha512`。默认值：`sha1`。
    * `ike_pfs` - (选填) - 第一阶段协商使用的 Diffie-Hellman 密钥交换算法。有效值：`group1` | `group2` | `group5` | `group14` | `group24`。默认值：`group2`。
    * `ike_lifetime` - (选填) - 第一阶段协商结果的 SA 生命周期。有效值范围为 [0, 86400]，单位为秒，默认值为 86400。
    * `ike_local_id` - (选填) - VPN 网关的标识符。
    * `ike_remote_id` - (选填) - 客户网关的标识符。
  * `ipsec_config` - (选填) - 第二阶段协商的配置。
    
    * `ipsec_enc_alg` - (选填) - 第二阶段协商的加密算法。有效值：`aes` | `aes192` | `aes256` | `des` | `3des`。默认值：`aes`。
    * `ipsec_auth_alg` - (选填) - 第二阶段协商的身份验证算法。有效值：`md5` | `sha1` | `sha256` | `sha384` | `sha512`。默认值：`sha1`。
    * `ipsec_pfs` - (选填) - 第二阶段协商使用的 Diffie-Hellman 密钥交换算法。有效值：`group1` | `group2` | `group5` | `group14` | `group24` | `disabled`。默认值：`group2`。
    * `ipsec_lifetime` - (选填) - 第二阶段协商结果的 SA 生命周期。有效值范围为 [0, 86400]，单位为秒，默认值为 86400。

## 属性参考

除了上述所有参数外，还导出了以下属性：
  * `id` - VPN 连接的 ID。
  * `name` - IPsec 连接的名称。
  * `vpn_connection_name` - IPsec 连接的名称。
  * `ike_config` - 第一阶段协商的配置。
    * `psk` - 预共享密钥。
    * `ike_version` - IKE 协议版本。
    * `ike_mode` - IKE 模式。
    * `ike_enc_alg` - 加密算法。
    * `ike_auth_alg` - 认证算法。
    * `ike_pfs` - Diffie-Hellman 组。
    * `ike_lifetime` - SA 生命周期。
    * `ike_local_id` - 本地标识符。
    * `ike_remote_id` - 远程标识符。
  * `ipsec_config` - 第二阶段协商的配置。
    * `ipsec_enc_alg` - 加密算法。
    * `ipsec_auth_alg` - 认证算法。
    * `ipsec_pfs` - Diffie-Hellman 组。
    * `ipsec_lifetime` - SA 生命周期。
  * `status` - 资源的状态。