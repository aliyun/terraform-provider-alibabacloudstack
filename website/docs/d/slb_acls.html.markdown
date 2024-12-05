---
subcategory: "Server Load Balancer (SLB)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_slb_acls"
sidebar_current: "docs-alibabacloudstack-datasource-slb-acls"
description: |-
    Provides a list of server load balancer acls (access control lists) to the user.
---

# alibabacloudstack\_slb_acls

This data source provides the acls in the region.

## Example Usage

```
data "alibabacloudstack_slb_acls" "sample_ds" {
}

output "first_slb_acl_id" {
  value = "${data.alibabacloudstack_slb_acls.sample_ds.acls.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of acls IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by acl name.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `tags` - (Optional) A mapping of tags to assign to the resource.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of SLB acls IDs.
* `names` - A list of SLB acls names.
* `acls` - A list of SLB  acls. Each element contains the following attributes:
  * `id` - Acl ID.
  * `name` - Acl name.
  * `ip_version` - The IP Version of access control list is the type of its entry (IP addresses or CIDR blocks). It values ipv4/ipv6.
  * `entry_list` - A list of entry (IP addresses or CIDR blocks).
    * `entry`   - An IP addresses or CIDR blocks.
    * `comment` - the comment of the entry.
  * `related_listeners` - A list of listener are attached by the acl.
    * `load_balancer_id` - the id of load balancer instance, the listener belongs to.
    * `frontend_port` - the listener port.
    * `protocol`      - the listener protocol (such as tcp/udp/http/https, etc).
    * `acl_type`      - the type of acl (such as white/black).
  * `tags` - A mapping of tags to assign to the resource.
