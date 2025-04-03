---
subcategory: "VPC"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_vpc_ipv6_egress_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-vpc-ipv6-egress-rules"
description: |- 
  Provides a list of vpc ipv6egressrules owned by an Alibabacloudstack account.
---

# alibabacloudstack_vpc_ipv6_egress_rules
-> **NOTE:** Alias name has: `alibabacloudstack_vpc_ipv6_egressrules`

This data source provides a list of vpc ipv6egressrules in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_vpc_ipv6egressrules" "ids" {
  ipv6_gateway_id = "example_value"
  ids             = ["example_value-1", "example_value-2"]
}

output "vpc_ipv6_egress_rule_id_1" {
  value = data.alibabacloudstack_vpc_ipv6egressrules.ids.rules.0.id
}

data "alibabacloudstack_vpc_ipv6egressrules" "nameRegex" {
  ipv6_gateway_id = "example_value"
  name_regex      = "^my-Ipv6EgressRule"
}

output "vpc_ipv6_egress_rule_id_2" {
  value = data.alibabacloudstack_vpc_ipv6egressrules.nameRegex.rules.0.id
}

data "alibabacloudstack_vpc_ipv6egressrules" "status" {
  ipv6_gateway_id = "example_value"
  status          = "Available"
}

output "vpc_ipv6_egress_rule_id_3" {
  value = data.alibabacloudstack_vpc_ipv6egressrules.status.rules.0.id
}

data "alibabacloudstack_vpc_ipv6egressrules" "ipv6EgressRuleName" {
  ipv6_gateway_id       = "example_value"
  ipv6_egress_rule_name = "example_value"
}

output "vpc_ipv6_egress_rule_id_4" {
  value = data.alibabacloudstack_vpc_ipv6egressrules.ipv6EgressRuleName.rules.0.id
}
```

## Argument Reference

The following arguments are supported:

* `instance_id` - (Optional, ForceNew) The ID of the IPv6 address to which you want to apply the egress-only rule.
* `ids` - (Optional, ForceNew) A list of Ipv6 Egress Rule IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Ipv6 Egress Rule name.
* `ipv6_egress_rule_name` - (Optional, ForceNew) The name of the egress-only rule. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
* `ipv6_gateway_id` - (Required, ForceNew) The ID of the IPv6 gateway.
* `status` - (Optional, ForceNew) The status of the resource. Valid values: `Available`, `Deleting`, `Pending`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Ipv6 Egress Rule names.
* `rules` - A list of Vpc Ipv6 Egress Rules. Each element contains the following attributes:
  * `description` - The description of the egress-only rule. The description must be `2` to `256` characters in length. It cannot start with `http://` or `https://`.
  * `id` - The ID of the Ipv6 Egress Rule. The value formats as `<ipv6_gateway_id>:<ipv6_egress_rule_id>`.
  * `instance_id` - The ID of the instance to which the egress-only rule is applied.
  * `instance_type` - The type of instance to which you want to apply the egress-only rule. Valid values: `Ipv6Address`. `Ipv6Address` (default): an IPv6 address.
  * `ipv6_egress_rule_id` - The ID of the IPv6 EgressRule.
  * `ipv6_egress_rule_name` - The name of the egress-only rule. The name must be `2` to `128` characters in length, and can contain letters, digits, underscores (_), and hyphens (-). The name must start with a letter but cannot start with `http://` or `https://`.
  * `status` - The status of the resource. Valid values: `Available`, `Pending`, `Deleting`.
  * `ipv6_gateway_id` - The ID of the IPv6 gateway.