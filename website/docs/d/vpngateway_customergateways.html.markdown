---
subcategory: "VPNGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpngateway_customergateways"
sidebar_current: "docs-Alibabacloudstack-datasource-vpngateway-customergateways"
description: |- 
  Provides a list of vpngateway customergateways owned by an alibabacloudstack account.
---

# alibabacloudstack_vpngateway_customergateways
-> **NOTE:** Alias name has: `alibabacloudstack_vpn_customer_gateways`

This data source provides a list of vpngateway customergateways in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_vpngateway_customergateways" "example" {
  name_regex = "example-cgw-*"
  ids        = ["cgw-12345678", "cgw-87654321"]
  output_file = "./customergateways_output.txt"
}

output "customergateway_ids" {
  value = data.alibabacloudstack_vpngateway_customergateways.example.ids
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of customer gateway IDs. If specified, the data source will return matching customer gateways.
* `name_regex` - (Optional, ForceNew) A regex string used to filter customer gateways by their names. This allows you to match specific patterns in the names of the customer gateways.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of names of the matched customer gateways.
* `gateways` - A list of customer gateway objects. Each object contains the following attributes:
  * `id` - The ID of the customer gateway.
  * `name` - The name of the customer gateway.
  * `ip_address` - The IP address of the customer gateway.
  * `description` - The description of the customer gateway.
  * `create_time` - The time when the customer gateway was created, formatted as an ISO8601 string (e.g., `2023-09-01T12:00:00Z`).

### Notes

- The `ids`, `names`, and `gateways` attributes provide different ways to access the same information about the customer gateways. You can choose the one that best fits your use case.
```