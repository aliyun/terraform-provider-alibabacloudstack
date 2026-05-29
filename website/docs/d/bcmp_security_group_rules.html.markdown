---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_security_group_rules"
sidebar_current: "docs-Alibabacloudstack-datasource-bcmp-security-group-rules"
description: |- 
  Provides a list of BCMP security group rules.
---

# alibabacloudstack_bcmp_security_group_rules

This data source provides a list of BCMP security group rules in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
resource "alibabacloudstack_vpc" "default" {
  name       = "tf-test-vpc"
  cidr_block = "172.16.0.0/12"
}

resource "alibabacloudstack_bcmp_security_group" "default" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = "tf-test-sg"
}

resource "alibabacloudstack_bcmp_security_group_rule" "default" {
  type              = "ingress"
  ip_protocol       = "tcp"
  policy            = "accept"
  port_range        = "22/22"
  priority          = 1
  security_group_id = alibabacloudstack_bcmp_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
  description       = "test"
}

# Retrieve all rules for a security group
data "alibabacloudstack_bcmp_security_group_rules" "default" {
  security_group_id = alibabacloudstack_bcmp_security_group.default.id
}

output "rules" {
  value = data.alibabacloudstack_bcmp_security_group_rules.default.rules
}
```

## Argument Reference

The following arguments are supported:

* `security_group_id` - (Required) The ID of the security group.
* `type` - (Optional) The type of rule to filter. Valid values are `ingress` and `egress`.
* `ip_protocol` - (Optional) The IP protocol type to filter. Valid values are `tcp`, `udp`, `icmp`, `gre`, and `all`.
* `policy` - (Optional) The policy to filter. Valid values are `accept` and `drop`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `rules` - A list of security group rules. Each element contains the following attributes:
  * `type` - The type of the rule.
  * `ip_protocol` - The IP protocol type.
  * `port_range` - The port range.
  * `cidr_ip` - The CIDR IP address.
  * `source_security_group_id` - The ID of the source security group.
  * `policy` - The policy of the rule.
  * `priority` - The priority of the rule.
  * `description` - The description of the rule.
