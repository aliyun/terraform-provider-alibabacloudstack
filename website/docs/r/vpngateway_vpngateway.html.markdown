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
  enable_ssl        = true
  instance_charge_type = "PostPaid"

  tags = {
    Created = "TF"
    Purpose = " acceptance_test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpn_gateway_name` - (Optional) The name of the VPN gateway.
* `vpc_id` - (Required, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `instance_charge_type` - (Optional, ForceNew) The billing method of the instance. Valid values:
  * **PrePaid**: Subscription.
  * **PostPaid**: Pay-As-You-Go.
  Default value: **PostPaid**.
* `period` - (Optional) Duration of purchase. This parameter is required when `instance_charge_type` is set to `PrePaid`. Valid values: [1-9, 12, 24, 36]. Default value: 1.
* `bandwidth` - (Required) The public network bandwidth of the VPN gateway. Unit: Mbps. Valid values for PostPaid instances: 10, 100, 200. Valid values for PrePaid instances: 5, 10, 20, 50, 100, 200.
* `enable_ipsec` - (Optional) Specifies whether to enable the IPsec-VPN feature. Default value: **true**.
* `enable_ssl` - (Optional) Specifies whether to enable the SSL-VPN feature. Default value: **false**.
* `ssl_connections` - (Optional) The maximum number of concurrent SSL-VPN connections. Valid values: 5, 10, 20, 50, 100, 200. Default value: **5**. This parameter takes effect only when `enable_ssl` is set to **true**.
* `description` - (Optional) The description of the VPN gateway.
* `vswitch_id` - (Optional, ForceNew) The ID of the vSwitch to which the VPN gateway belongs.
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The ID of the resource.
* `internet_ip` - The public IP address of the VPN gateway.
* `status` - The status of the resource. Valid values:
  * **Creating**: The resource is being created.
  * **Available**: The resource has been created and can be used normally.
  * **Deleting**: The resource is being deleted.
* `business_status` - The payment status of the VPN gateway. Valid values:
  * **Normal**: The resource is normal.
  * **Expired**: The resource has expired.
  * **LockDown**: The resource has been locked.