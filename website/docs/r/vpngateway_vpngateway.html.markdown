---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpngateway"
sidebar_current: "docs-Alibabacloudstack-vpngateway-vpngateway"
description: |- 
  Provides a vpngateway Vpngateway resource.
---

# alibabacloudstack_vpngateway_vpngateway
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_gateway`

Provides a vpngateway Vpngateway resource.

## Example Usage

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
    Purpose = " acceptance_test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `bandwidth` - (Required) The public network bandwidth of the VPN gateway. Unit: Mbps. Valid values: 5, 10, 20, 50, 100, 200, 500, 1000.
* `vpn_gateway_name` - (Optional) The name of the VPN gateway. The name must be 1 to 128 characters in length.
* `instance_charge_type` - (Optional, ForceNew) The billing method of the instance. Valid values: `PrePaid`, `PostPaid`. Default value: `PostPaid`.
* `period` - (Optional) Duration of purchase. Valid values: 1 to 9, 12, 24, 36. Default value: 1.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch to which the VPN gateway belongs.
* `ipsec_vpn` - (Optional) Specifies whether to enable the IPsec-VPN feature.
* `ssl_vpn` - (Optional) Specifies whether to enable the SSL-VPN feature.
* `ssl_max_connections` - (Optional) The maximum number of concurrent SSL-VPN connections. Default value: 5.
* `description` - (Optional) The description of the VPN gateway. The description must be 2 to 256 characters in length.
* `tags` - (Optional) A mapping of tags to assign to the resource.
* `name` - (Optional, Deprecated) This parameter is deprecated. Use `vpn_gateway_name` instead.
* `enable_ipsec` - (Optional, Deprecated) This parameter is deprecated. Use `ipsec_vpn` instead.
* `enable_ssl` - (Optional, Deprecated) This parameter is deprecated. Use `ssl_vpn` instead.
* `ssl_connections` - (Optional, Deprecated) This parameter is deprecated. Use `ssl_max_connections` instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the VPN gateway.
* `internet_ip` - The public IP address of the VPN gateway.
* `status` - The status of the VPN gateway. Valid values: `provisioning`, `init`, `active`.
* `business_status` - The payment status of the VPN gateway.
* `vpn_gateway_name` - The name of the VPN gateway.
* `ipsec_vpn` - Indicates whether the IPsec-VPN feature is enabled.
* `ssl_vpn` - Indicates whether the SSL-VPN feature is enabled.
* `ssl_max_connections` - The maximum number of concurrent SSL-VPN connections.

## Import

VPN Gateway can be imported using the VpnGatewayId, e.g.

```
$ terraform import alibabacloudstack_vpngateway_vpngateway.example vgw-xxxxxxxxx
```
