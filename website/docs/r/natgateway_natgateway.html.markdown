---
subcategory: "NATGateway"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_natgateway_natgateway"
sidebar_current: "docs-Alibabacloudstack-natgateway-natgateway"
description: |- 
  Provides a natgateway Natgateway resource.
---

# alibabacloudstack_natgateway_natgateway
-> **NOTE:** Alias name has: `alibabacloudstack_nat_gateway`

Provides a natgateway Natgateway resource.

## Example Usage

```hcl
variable "name" {
  default = "tf-testAccNatGatewayConfig13663"
}

data "alibabacloudstack_zones" "default" {
  available_resource_creation = "VSwitch"
}

resource "alibabacloudstack_vpc" "default" {
  name       = "${var.name}"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_vswitch" "default" {
  vpc_id            = "${alibabacloudstack_vpc.default.id}"
  cidr_block        = "172.16.0.0/21"
  availability_zone = "${data.alibabacloudstack_zones.default.zones.0.id}"
  name              = "${var.name}"
}

resource "alibabacloudstack_nat_gateway" "default" {
  vpc_id             = "${alibabacloudstack_vswitch.default.vpc_id}"
  specification     = "Small"
  nat_gateway_name  = "${var.name}"
  description       = "This is a test NAT Gateway"
  bandwidth_packages = [
    {
      ip_count         = 2
      bandwidth        = 10
      zone            = "${data.alibabacloudstack_zones.default.zones.0.id}"
      public_ip_addresses = ["10.0.0.1", "10.0.0.2"]
    }
  ]
  tags = {
    CreatedBy = "Terraform"
    Env       = "Test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` - (Required, ForceNew) The ID of the VPC where the NAT gateway is deployed.
* `specification` - (Optional) The specification of the NAT gateway. Valid values are `Small`, `Middle`, and `Large`. Default to `Small`.
* `nat_gateway_name` - (Optional) The name of the NAT gateway. The value can have a string of 2 to 128 characters, must contain only alphanumeric characters or hyphens, such as "-", ".", "_", and must not begin or end with a hyphen, and must not begin with http:// or https://. Defaults to null.
* `description` - (Optional) The description of the NAT gateway. This description can have a string of 2 to 256 characters, and it cannot begin with http:// or https://. Defaults to null.
* `bandwidth_packages` - (Optional) A list of bandwidth packages for the NAT gateway. Only supports NAT gateways created before 00:00 on November 4, 2017.
  * `ip_count` - (Required) The number of IP addresses in the current bandwidth package. Its value range is from 1 to 50.
  * `bandwidth` - (Required) The bandwidth value of the current bandwidth package. Its value range is from 5 to 5000 Mbps.
  * `zone` - (Optional) The Availability Zone for the current bandwidth package. If this value is not specified, Terraform will set a random AZ.
  * `public_ip_addresses` - The public IP addresses for the bandwidth package. The count of public IPs equals `ip_count`, and multiple IPs are separated by commas, such as "10.0.0.1,10.0.0.2".
* `tags` - (Optional, Map) A mapping of tags to assign to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the NAT gateway.
* `name` - The name of the NAT gateway.
* `description` - The description of the NAT gateway.
* `specification` - The specification of the NAT gateway.
* `vpc_id` - The VPC ID for the NAT gateway.
* `bandwidth_package_ids` - A list of IDs of the bandwidth packages, separated by commas.
* `snat_table_ids` - The ID of the SNAT table automatically created by the NAT gateway.
* `forward_table_ids` - The ID of the Destination Network Address Translation (DNAT) table automatically created by the NAT gateway.