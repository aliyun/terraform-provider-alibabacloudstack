---
subcategory: "Bare Metal Computing Platform (BMCP)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bcmp_security_group_rule"
sidebar_current: "docs-Alibabacloudstack-resource-bcmp-security-group-rule"
description: |- 
  Provides a BCMP Security Group Rule resource.
---

# alibabacloudstack_bcmp_security_group_rule

Provides a BCMP Security Group Rule resource.

> **NOTE:** This resource can also be referred to by the following alias:
> - `alibabacloudstack_bcmp_security_group_rules` (deprecated)

## Example Usage

Basic Usage

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
  description       = "Allow SSH access"
}
```

Using Source Security Group

```hcl
resource "alibabacloudstack_bcmp_security_group" "new" {
  vpc_id = alibabacloudstack_vpc.default.id
  name   = "tf-test-sg-new"
}

resource "alibabacloudstack_bcmp_security_group_rule" "with_source_sg" {
  type                     = "ingress"
  ip_protocol              = "tcp"
  policy                   = "drop"
  port_range               = "22/22"
  priority                 = 100
  security_group_id        = alibabacloudstack_bcmp_security_group.default.id
  source_security_group_id = alibabacloudstack_bcmp_security_group.new.id
  description              = "Block SSH from other SG"
}
```

## Argument Reference

The following arguments are supported:

* `type` - (Required, ForceNew) Type of rule, `ingress` (inbound) or `egress` (outbound). Modifying this parameter will force a new resource to be created.
* `ip_protocol` - (Required) The IP protocol type. Valid values are `tcp`, `udp`, `icmp`, `gre`, and `all`.
* `policy` - (Optional) The policy of the rule. Valid values are `accept` and `drop`. Default is `accept`.
* `port_range` - (Required) The port range. For `tcp` and `udp`, the format is `start/end` (e.g., `22/22`). For `icmp`, `gre`, and `all`, use `-1/-1`.
* `priority` - (Optional) The priority of the rule. Valid range: 1-100. Default is 1.
* `security_group_id` - (Required) The ID of the security group.
* `cidr_ip` - (Optional) The CIDR IP address. Conflicts with `source_security_group_id`. Either `cidr_ip` or `source_security_group_id` must be specified.
* `source_security_group_id` - (Optional) The ID of the source security group. Conflicts with `cidr_ip`. Either `cidr_ip` or `source_security_group_id` must be specified.
* `description` - (Optional) The description of the rule.

> **NOTE:** Either `cidr_ip` or `source_security_group_id` must be specified.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the security group rule. The format is `<security_group_id>:<sgr_id>:<type>`.
* `type` - The type of the rule.
* `ip_protocol` - The IP protocol type.
* `policy` - The policy of the rule.
* `port_range` - The port range.
* `priority` - The priority of the rule.
* `security_group_id` - The ID of the security group.
* `cidr_ip` - The CIDR IP address.
* `source_security_group_id` - The ID of the source security group.
* `description` - The description of the rule.

## Import

BCMP Security Group Rule can be imported using the `<security_group_id>:<sgr_id>:<type>`, e.g.

```
$ terraform import alibabacloudstack_bcmp_security_group_rule.example sg-12345678:sgr-abcdef:ingress
```
