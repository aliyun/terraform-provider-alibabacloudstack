---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_vpngateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-vpngateways"
description: |- 
  Provides a list of vpngateway vpngateways owned by an alibabacloudstack account.
---

# alibabacloudstack_vpngateway_vpngateways
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_gateways`

This data source provides a list of vpngateway vpngateways in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpngateway_vpngateways" "example" {
  vpc_id          = "vpc-example-id"
  ids             = ["vgw-example-id1", "vgw-example-id2"]
  status          = "Active"
  business_status = "Normal"
  name_regex      = "testAcc*"
  output_file     = "/tmp/vpns"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of VPN gateway IDs. This will be used as a filter condition.
* `vpc_id` - (Optional, ForceNew) The ID of the VPC to which the VPN gateway belongs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by the name of the VPN gateway.
* `status` - (Optional, ForceNew) The status of the resource. Valid values include: "Init", "Provisioning", "Active", "Updating", "Deleting".
* `business_status` - (Optional, ForceNew) The payment status of the VPN gateway. Valid values include: "Normal", "FinancialLocked".

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of the matched vpngateways.
* `gateways` - A list of vpngateway resources. Each element contains the following attributes:
  * `id` - The ID of the vpngateway.
  * `vpc_id` - The ID of the VPC to which the vpngateway belongs.
  * `internet_ip` - The public IP address of the vpngateway.
  * `create_time` - The time when the vpngateway was created.
  * `end_time` - The expiration time of the vpngateway.
  * `specification` - The specification of the vpngateway.
  * `name` - The name of the vpngateway.
  * `description` - The description of the vpngateway.
  * `status` - The status of the vpngateway.
  * `business_status` - The payment status of the vpngateway.
  * `instance_charge_type` - The charge type of the vpngateway.
  * `enable_ipsec` - Indicates whether the IPsec-VPN function is enabled.
  * `enable_ssl` - Indicates whether the SSL-VPN function is enabled.
  * `ssl_connections` - The maximum number of SSL-VPN client connections allowed.