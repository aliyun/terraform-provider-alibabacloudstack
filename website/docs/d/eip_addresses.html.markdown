---
subcategory: "EIP"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_eip_addresses"
sidebar_current: "docs-Alibabacloudstack-datasource-eip-addresses"
description: |- 
  Provides a list of eip addresses owned by an alibabacloudstack account.
---

# alibabacloudstack_eip_addresses
-> **NOTE:** Alias name has: `alibabacloudstack_eips`

This data source provides a list of eip addresses in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_eip_addresses" "example" {
  ids        = ["eip-12345678"]
  ip_addresses = ["192.168.0.1"]

  output_file = "eips_output.txt"
}

output "first_eip_id" {
  value = data.alibabacloudstack_eip_addresses.example.eips.0.id
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of EIP IDs. If specified, the data source will return only those EIPs whose IDs match the ones provided.
* `ip_addresses` - (Optional) A list of EIP public IP addresses. If specified, the data source will return only those EIPs whose IP addresses match the ones provided.


## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of EIP IDs matching the specified filters.
* `names` - A list of EIP names corresponding to the filtered EIPs. 
* `eips` - A list of EIPs. Each element contains the following attributes:
  * `id` - The ID of the EIP.
  * `status` - The status of the EIP. Possible values include: `Associating`, `Unassociating`, `InUse`, and `Available`.
  * `ip_address` - The public IP address of the EIP.
  * `bandwidth` - The maximum internet bandwidth (in Mbps) of the EIP.
  * `instance_id` - The ID of the instance that the EIP is currently bound to.
  * `instance_type` - The type of instance that the EIP is bound to.