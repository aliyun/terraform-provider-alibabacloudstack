---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpnconnection"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpnconnection"
description: |- 
  Provides a vpngateway Vpnconnection resource.
---

# alibabacloudstack_vpngateway_vpnconnection
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_connection`

Provides a vpngateway Vpnconnection resource.

## Example Usage

Basic Usage

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

## Argument Reference

The following arguments are supported:

* `customer_gateway_id` - (Required, ForceNew) The ID of the customer gateway.
* `vpn_gateway_id` - (Required, ForceNew) The ID of the VPN gateway.
* `vpn_connection_name` - (Optional) The name of the IPsec-VPN connection.
* `local_subnet` - (Required, Type:Set) The CIDR block of the virtual private cloud (VPC). This parameter is used for phase-two negotiation.
* `remote_subnet` - (Required, Type:Set) The CIDR block of the on-premises data center. This parameter is used for phase-two negotiation.
* `effect_immediately` - (Optional) Indicates whether IPsec-VPN negotiations are initiated immediately. Valid values: `true`, `false`.
* `ike_config` - (Optional) The configuration of Phase 1 negotiations. The `ike_config` mapping supports the following:
  * `psk` - (Optional) Used for authentication between the IPsec VPN gateway and the customer gateway.
  * `ike_version` - (Optional) The version of the IKE protocol. Valid value: `ikev1` | `ikev2`. Default value: `ikev1`.
  * `ike_mode` - (Optional) The negotiation mode of IKE V1. Valid value: `main` (main mode) | `aggressive` (aggressive mode). Default value: `main`.
  * `ike_enc_alg` - (Optional) The encryption algorithm of phase-one negotiation. Valid value: `aes` | `aes192` | `aes256` | `des` | `3des`. Default value: `aes`.
  * `ike_auth_alg` - (Optional) The authentication algorithm of phase-one negotiation. Valid value: `md5` | `sha1` | `sha256` | `sha384` | `sha512`. Default value: `sha1`.
  * `ike_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-one negotiation. Valid value: `group1` | `group2` | `group5` | `group14` | `group24`. Default value: `group2`.
  * `ike_lifetime` - (Optional) The SA lifecycle as the result of phase-one negotiation. The valid value of n is [0, 86400], the unit is second and the default value is 86400.
  * `ike_local_id` - (Optional) The identification of the VPN gateway.
  * `ike_remote_id` - (Optional) The identification of the customer gateway.
* `ipsec_config` - (Optional) The configurations of phase-two negotiation. The `ipsec_config` mapping supports the following:
  * `ipsec_enc_alg` - (Optional) The encryption algorithm of phase-two negotiation. Valid value: `aes` | `aes192` | `aes256` | `des` | `3des`. Default value: `aes`.
  * `ipsec_auth_alg` - (Optional) The authentication algorithm of phase-two negotiation. Valid value: `md5` | `sha1` | `sha256` | `sha384` | `sha512`. Default value: `sha1`.
  * `ipsec_pfs` - (Optional) The Diffie-Hellman key exchange algorithm used by phase-two negotiation. Valid value: `group1` | `group2` | `group5` | `group14` | `group24` | `disabled`. Default value: `group2`.
  * `ipsec_lifetime` - (Optional) The SA lifecycle as the result of phase-two negotiation. The valid value is [0, 86400], the unit is second and the default value is 86400.
* `name` - (Optional) The name of the IPsec-VPN connection.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the VPN connection.
* `status` - The status of the resource.
* `ike_config` - The configuration of Phase 1 negotiations.
* `ipsec_config` - IPsec configuration.
* `vpn_connection_name` - The name of the IPsec-VPN connection.
* `name` - The name of the IPsec-VPN connection.

## Import

VPN connection can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_vpn_connection.example vco-abc123456
```