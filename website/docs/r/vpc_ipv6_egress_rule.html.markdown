---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_egress_rule"
sidebar_current: "docs-alibabacloudstack-resource-vpc-ipv6-egress-rule"
description: |-
  Provides a Alibabacloudstack VPC Ipv6 Egress Rule resource.
---

# alibabacloudstack\_vpc\_ipv6\_egress\_rule

Provides a VPC Ipv6 Egress Rule resource.

For information about VPC Ipv6 Egress Rule and how to use it, see [What is Ipv6 Egress Rule](https://www.alibabacloud.com/help/doc-detail/102200.htm).

-> **NOTE:** Available in v1.142.0+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_vpc" "default" {
  vpc_name    = "example_value"
  enable_ipv6 = "true"
}

resource "alibabacloudstack_vpc_ipv6_gateway" "example" {
  ipv6_gateway_name = "example_value"
  vpc_id            = alibabacloudstack_vpc.default.id
}

data "alibabacloudstack_instances" "default" {
  name_regex = "ecs_with_ipv6_address"
  status     = "Running"
}

data "alibabacloudstack_vpc_ipv6_addresses" "default" {
  associated_instance_id = data.alibabacloudstack_instances.default.instances.0.id
  status                 = "Available"
}

resource "alibabacloudstack_vpc_ipv6_egress_rule" "example" {
  instance_id           = data.alibabacloudstack_vpc_ipv6_addresses.default.ids.0
  ipv6_egress_rule_name = "example_value"
  description           = "example_value"
  ipv6_gateway_id       = alibabacloudstack_vpc_ipv6_gateway.example.id
  instance_type         = "Ipv6Address"
}

```

## Argument Reference

The following arguments are supported:

* `description` - (Optional, ForceNew) The description of the egress-only rule. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
* `instance_id` - (Required, ForceNew) The ID of the IPv6 address to which you want to apply the egress-only rule.
* `instance_type` - (Optional, Computed, ForceNew) The type of instance to which you want to apply the egress-only rule. Valid values: `Ipv6Address`. `Ipv6Address` (default): an IPv6 address.
* `ipv6_egress_rule_name` - (Optional, ForceNew) The name of the egress-only rule. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID in terraform of Ipv6 Egress Rule. The value formats as `<ipv6_gateway_id>:<ipv6_egress_rule_id>`.
* `status` - The status of the resource. Valid values: `Available`, `Pending` and `Deleting`.

### Timeouts

The `timeouts` block allows you to specify [timeouts](https://www.terraform.io/docs/configuration-0-11/resources.html#timeouts) for certain actions:

* `create` - (Defaults to 1 mins) Used when create the Ipv6 Egress Rule.
* `delete` - (Defaults to 1 mins) Used when delete the Ipv6 Egress Rule.

## Import

VPC Ipv6 Egress Rule can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_vpc_ipv6_egress_rule.example <ipv6_gateway_id>:<ipv6_egress_rule_id>
```
