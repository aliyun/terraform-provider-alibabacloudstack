---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_gateway"
sidebar_current: "docs-alibabacloudstack-resource-vpc-ipv6-gateway"
description: |- 
  Provides a VPC Ipv6 Gateway resource.

---

# alibabacloudstack_vpc_ipv6_gateway
-> **NOTE:** Alias name has: `alibabacloudstack_vpc_ipv6gateway`

Provides a VPC Ipv6 Gateway resource.

For information about VPC Ipv6 Gateway and how to use it, see [What is Ipv6 Gateway](https://www.alibabacloud.com/help/doc-detail/102214.htm).



## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testaccvpcipv6gateway88979"
}

resource "alibabacloudstack_vpc" "default" {
  vpc_name    = var.name
  enable_ipv6 = "true"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = var.name
  vpc_id            = alibabacloudstack_vpc.default.id
  description       = var.name
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) The description of the IPv6 gateway. The description must be **2 to 256** characters in length. It cannot start with `http://` or `https://`.
* `ipv6_gateway_name` - (Optional) The name of the IPv6 gateway. The name must be **2 to 128** characters in length, and can contain letters, digits, underscores (`_`), and hyphens (`-`). The name must start with a letter but cannot start with `http://` or `https://`.
* `spec` - (Optional) The edition of the IPv6 gateway. Valid values: `Large`, `Medium`, and `Small`. 
  - `Small` (default): Free Edition.
  - `Medium`: Enterprise Edition.
  - `Large`: Enhanced Enterprise Edition.
  
  > **Note:** The throughput capacity of an IPv6 gateway varies based on the edition. For more information, see [Editions of IPv6 gateways](https://www.alibabacloud.com/help/doc-detail/98926.htm).
* `vpc_id` - (Required, ForceNew) The ID of the virtual private cloud (VPC) for which you want to create the IPv6 gateway.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The resource ID in Terraform of the Ipv6 Gateway.
* `status` - The status of the resource. Valid values: `Available`, `Pending`, and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 minute) Used when creating the Ipv6 Gateway.
* `update` - (Defaults to 1 minute) Used when updating the Ipv6 Gateway.
* `delete` - (Defaults to 5 minutes) Used when deleting the Ipv6 Gateway.

## Import

VPC Ipv6 Gateway can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_vpc_ipv6_gateway.example <id>
```