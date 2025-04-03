---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_natgateways"
sidebar_current: "docs-Alibabacloudstack-datasource-natgateway-natgateways"
description: |- 
  Provides a list of natgateway natgateways owned by an alibabacloudstack account.
---

# alibabacloudstack_natgateway_natgateways
-> **NOTE:** Alias name has: `alibabacloudstack_nat_gateways`

This data source provides a list of natgateway natgateways in an Alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
variable "name" {
  default = "natGatewaysDatasource"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "foo" {
  name       = var.name
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_nat_gateway" "foo" {
  vpc_id        = alibabacloudstack_vpc.foo.id
  specification = "Small"
  name          = var.name
}

data "alibabacloudstack_natgateway_natgateways" "foo" {
  vpc_id     = alibabacloudstack_vpc.foo.id
  name_regex = alibabacloudstack_nat_gateway.foo.name
  ids        = [alibabacloudstack_nat_gateway.foo.id]
}

output "nat_gateways" {
  value = data.alibabacloudstack_natgateway_natgateways.foo.gateways
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Optional, ForceNew) The ID of the VPC where the NAT gateway is deployed.
* `ids` - (Optional) A list of NAT Gateway IDs. If specified, the data source will only return results that match these IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter NAT gateways by their names. This allows for more flexible filtering based on naming conventions.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of names of all the matched NAT Gateways.
* `gateways` - A list of NAT Gateways. Each element contains the following attributes:
  * `id` - The ID of the NAT Gateway.
  * `name` - The name of the NAT Gateway.
  * `description` - The description of the NAT Gateway.
  * `creation_time` - The time when the NAT Gateway was created.
  * `spec` - The specification of the NAT Gateway (e.g., Small, Medium).
  * `status` - The status of the NAT Gateway (e.g., Available, Pending, Deleting).
  * `snat_table_id` - The ID of the SNAT table associated with the NAT Gateway.
  * `forward_table_id` - The ID of the forward table associated with the NAT Gateway.