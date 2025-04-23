---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_accesscontrollists"
sidebar_current: "docs-Alibabacloudstack-datasource-slb-accesscontrollists"
description: |- 
  Provides a list of slb accesscontrollists owned by an alibabacloudstack account.
---

# alibabacloudstack_slb_accesscontrollists
-> **NOTE:** Alias name has: `alibabacloudstack_slb_acls`

This data source provides a list of SLB Access Control Lists (ACLs) in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_slb_accesscontrollists" "sample_ds" {
  ids        = ["acl-12345678"]
  name_regex = "^my-acl-.*$"
  tags       = {
    Environment = "Production"
  }
}

output "first_slb_acl_id" {
  value = "${data.alibabacloudstack_slb_accesscontrollists.sample_ds.acls.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ACL IDs to filter results. This can be useful if you know the IDs of specific ACLs you want to retrieve.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by ACL name. This allows you to match ACL names based on a pattern.
* `tags` - (Optional) A mapping of tags to filter the ACLs by. Only ACLs with matching tags will be returned.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB ACL IDs.
* `names` - A list of SLB ACL names.
* `acls` - A list of SLB ACLs. Each element contains the following attributes:
  * `id` - The unique ID of the ACL.
  * `name` - The name of the ACL.
  * `ip_version` - The IP Version of the access control list, which determines the type of its entries (IPv4 or IPv6). Possible values are `ipv4` or `ipv6`.
  * `entry_list` - A list of entries (IP addresses or CIDR blocks) associated with the ACL. Each entry contains:
    * `entry` - An IP address or CIDR block.
    * `comment` - A comment associated with the entry.
  * `related_listeners` - A list of listeners that are attached to this ACL. Each listener contains:
    * `load_balancer_id` - The ID of the load balancer instance to which the listener belongs.
    * `frontend_port` - The port number of the listener.
    * `protocol` - The protocol used by the listener (e.g., TCP, UDP, HTTP, HTTPS, etc.).
    * `acl_type` - The type of ACL applied to the listener (e.g., white/black).
  * `tags` - A mapping of tags assigned to the ACL.