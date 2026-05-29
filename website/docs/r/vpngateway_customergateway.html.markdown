---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_customergateway"
sidebar_current: "docs-Alibabacloudstack-vpngateway-customergateway"
description: |-
  Provides a vpngateway Customergateway resource.
---

# alibabacloudstack_vpngateway_customergateway
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_customer_gateway`

Provides a vpngateway Customergateway resource.

## Example Usage

Basic Usage

```hcl
resource "alibabacloudstack_vpngateway_customergateway" "default" {
  ip_address           = "1.1.1.1"
  customer_gateway_name = "example-customer-gateway"
  description          = "This is a test customer gateway."
}
```

## Argument Reference

The following arguments are supported:

* `ip_address` - (Required, ForceNew) The IP address of the customer gateway. This must be a valid IP address. If you plan to create an IPsec connection of the Internet type, enter a public IP address. If you plan to create an IPsec connection of the VPC type, enter a private IP address. Modifying this argument forces a new resource to be created.
* `customer_gateway_name` - (Optional) The name of the customer gateway. The name must be 2 to 128 characters in length, and can contain letters, digits, periods (.), underscores (_), and hyphens (-). It must start with a letter but cannot start with http:// or https://. Conflicts with `name`.
* `name` - (Deprecated) The name of the customer gateway. This field is deprecated and will be removed in a future release. Please use `customer_gateway_name` instead. Conflicts with `customer_gateway_name`.
* `description` - (Optional) The description of the customer gateway. The description must be 2 to 256 characters in length, and must start with a letter but cannot start with http:// or https://.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the customer gateway instance.
* `customer_gateway_name` - The name of the customer gateway.

## Import

VPN Customer Gateway can be imported using the customer gateway ID, e.g.

```
$ terraform import alibabacloudstack_vpngateway_customergateway.example cgw-bp1jrawp82av6bws9****
```
