---
subcategory: "Elastic Compute Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_security_group_rule"
sidebar_current: "docs-alibabacloudstack-resource-security-group-rule"
description: |-
  Provides a Alibabacloudstack Security Group Rule resource.
---

# alibabacloudstack_security_group_rule

Provides a security group rule resource.
Represents a single `ingress` or `egress` group rule, which can be added to external Security Groups.

-> **NOTE:**  `nic_type` should set to `intranet` when security group type is `vpc` or specifying the `source_security_group_id`. In this situation it does not distinguish between intranet and internet, the rule is effective on them both.

> **NOTE:** This resource can also be referred to by the following alias:
> - `alibabacloudstack_ecs_securitygrouprule`


## Example Usage

Basic Usage

```
resource "alibabacloudstack_vpc" "vpc" {
  cidr_block = "10.1.0.0/21"
}

resource "alibabacloudstack_security_group" "group" {
  vpc_id = "${alibabacloudstack_vpc.vpc.id}"
}

resource "alibabacloudstack_security_group_rule" "allow_all_tcp" {
  type              = "ingress"
  ip_protocol       = "tcp"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "1/65535"
  priority          = 1
  security_group_id = "${alibabacloudstack_security_group.default.id}"
  cidr_ip           = "0.0.0.0/0"
}
```
## Argument Reference

The following arguments are supported:

* `type` - (Required, ForceNew) The type of rule being created. Valid options are `ingress` (inbound) or `egress` (outbound).
* `ip_protocol` - (Required, ForceNew) The protocol. Can be `tcp`, `udp`, `icmp`, `gre` or `all`.
* `port_range` - (Required, ForceNew) The range of port numbers relevant to the IP protocol. When the protocol is tcp or udp, each side port number range from 1 to 65535 and '-1/-1' will be invalid. For example, `1/200` means that the range of the port numbers is 1-200. Other protocols' 'port_range' can only be "-1/-1", and other values will be invalid.
* `security_group_id` - (Required, ForceNew) The security group to apply this rule to.
* `cidr_ip` - (Optional, ForceNew) The target IPv4 CIDR block. Conflicts with `ipv6_cidr_ip` and `source_security_group_id`.
* `ipv6_cidr_ip` - (Optional, ForceNew) The target IPv6 CIDR block. Conflicts with `cidr_ip` and `source_security_group_id`.
* `source_security_group_id` - (Optional, ForceNew) The target security group ID within the same region. Conflicts with `cidr_ip` and `ipv6_cidr_ip`. If this field is specified, the `nic_type` must be `intranet`.
* `nic_type` - (Optional, ForceNew) Network type, can be either `internet` or `intranet`, the default value is `intranet`.
* `policy` - (Optional, ForceNew) Authorization policy, can be either `accept` or `drop`, the default value is `accept`.
* `priority` - (Optional, ForceNew) Authorization policy priority, with parameter values: `1-100`, default value: 1.
* `description` - (Optional) The description of the security group rule. The description can be up to 1 to 512 characters in length.
* `source_group_owner_account` - (Optional, Deprecated) The Alibaba Cloud user account Id of the target security group when security groups are authorized across accounts. This parameter is invalid in apsarastack and is scheduled for removal in version 3.19.0.

-> **NOTE:** One of the `source_security_group_id`, `cidr_ip`, or `ipv6_cidr_ip` must be set.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the security group rule, formatted as `<security_group_id>:<type>:<ip_protocol>:<port_range>:<nic_type>:<cidr_ip>:<policy>:<priority>`.
* `type` - The type of rule, `ingress` or `egress`.
* `port_range` - The range of port numbers.
* `ip_protocol` - The protocol of the security group rule.
* `nic_type` - Indicates the network type, either `internet` or `intranet`.
* `policy` - The authorization policy, `accept` or `drop`.
* `priority` - The priority of the rule.
* `description` - The description of the security group rule.

## Import

Security Group Rule can be imported using the composite ID, e.g.

```
$ terraform import alibabacloudstack_security_group_rule.example sg-12345678:ingress:tcp:22/22:intranet:10.0.0.0/8:accept:1
```