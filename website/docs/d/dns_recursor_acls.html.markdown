---
subcategory: "Universal DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_recursor_acls"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-recursor-acls"
description: |-
  Provides a list of DNS Recursor ACL policies for managing access control policies for cross-cloud domain name resolution.
---

# alibabacloudstack_dns_recursor_acls

This data source provides a list of DNS Recursor ACL policies for managing access control policies for cross-cloud domain name resolution.

## Example Usage

```hcl
variable "name" {
  default = "tfacc63456"
}

resource "alibabacloudstack_dns_line" "ipv4" {
  name         = "${var.name}ipv4"
  v4_addresses = ["192.168.0.1"]
}

resource "alibabacloudstack_dns_recursor_acl" "default" {
  name     = var.name
  remark   = var.name
  policy   = "ALLOW"
  line_ids = ["${alibabacloudstack_dns_line.ipv4.id}"]
}

data "alibabacloudstack_dns_recursor_acls" "default" {
  name_regex = alibabacloudstack_dns_recursor_acl.default.name
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of Recursor ACL IDs to query. If provided, only Recursor ACLs with IDs in the list will be returned.

* `name_regex` - (Optional) A regex string to filter Recursor ACL names. Only Recursor ACLs whose names match the regex will be returned.

## Attributes Reference

The following attributes are exported:

* `id` - (String) The unique identifier of the Recursor ACL.

* `create_timestamp` - (Integer) The creation timestamp of the Recursor ACL in seconds.

* `line_ids` - (List) A list of line IDs to which the ACL policy applies.

* `name` - (String) The name of the Recursor ACL.

* `policy` - (String) The recursive policy. Valid values are `ALLOW` (allow recursion) and `FORBID` (forbid recursion).

* `remark` - (String) The remark information of the Recursor ACL.

* `update_timestamp` - (Integer) The last update timestamp of the Recursor ACL in seconds.