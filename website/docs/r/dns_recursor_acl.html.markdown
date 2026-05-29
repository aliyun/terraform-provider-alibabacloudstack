---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_recursor_acl"
sidebar_current: "docs-Alibabacloudstack-dns-recursor-acl"
description: |-
    Configures DNS recursive ACL policies in the specified resource set using the credentials configured in the provider.
---

# alibabacloudstack_dns_recursor_acl

This resource configures DNS recursive ACL policies in the specified resource set using the credentials configured in the provider.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc67879"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name         = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_line" "ipv6" {
  name         = "${var.name}ipv6"
  v6_addresses = ["2020:148:2:28::"]
}

resource "alibabacloudstack_dns_recursor_acl" "default" {
  line_ids = [
    "${alibabacloudstack_dns_line.ipv4.id}",
    "${alibabacloudstack_dns_line.ipv6.id}"
  ]
  name   = var.name
  remark = var.name
  policy = "ALLOW"
}
```

## Argument Reference

The following arguments are supported:

* `line_ids` - (Required) A list of request line IDs. At least one line ID must be specified.
* `name` - (Required) The name of the ACL policy.
* `policy` - (Required) The recursive policy. Valid values: `ALLOW` (allow recursion) or `FORBID` (forbid recursion).
* `remark` - (Optional) The remark information.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the ACL policy.

## Import

DNS Recursor ACL can be imported using the ACL ID, e.g.

```
$ terraform import alibabacloudstack_dns_recursor_acl.example acl-12345678
```