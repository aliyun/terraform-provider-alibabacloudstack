---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_forward_domain"
description: |-
  Global DNS Forward Domain
---

# alibabacloudstack_dns_forward_domain

This resource creates a global DNS forward domain in the specified resource set using the credentials configured in the provider.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc14020.test."
}

resource "alibabacloudstack_dns_forward_domain" "default" {
  name         = var.name
  remark       = "test1234"
  forward_mode = "FORWARD_FIRST"
  forwarders = [
    "192.168.101.1"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the forward domain. Must end with a dot, for example, `"example.com."`.
* `forward_mode` - (Required) The forwarding mode. Valid values:
  * `FORWARD_FIRST`: Forward First mode - If forwarding fails, it degrades to internet recursion.
  * `FORWARD_ONLY`: Forward Only mode (recommended) - All requests are forwarded only (no recursion).
* `forwarders` - (Required) A list of forwarder IP addresses. For example: `["192.168.101.1"]`.
* `remark` - (Optional) Remark information describing the purpose of the forward domain.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the forward domain.
* `caller_uid` - A system parameter identifying the caller who created the resource.
* `create_timestamp` - The creation timestamp in seconds.
* `update_timestamp` - The update timestamp in seconds.

## Import

DNS Forward Domain can be imported using the resource ID, e.g.

```
$ terraform import alibabacloudstack_dns_forward_domain.example <resource_id>
```