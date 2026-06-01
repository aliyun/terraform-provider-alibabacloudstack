---
subcategory: "Alibaba Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_forward_domains"
sidebar_current: "docs-Alibabacloudstack-datasource-dns-forward-domains"
description: |-
  Provides a list of forwarding domains available to the user.
---

# alibabacloudstack_dns_forward_domains

> Query the list of Alibaba Cloud DNS global forwarding domains

## Example Usage

```hcl
variable "name" {
  default = "tfacc6273812765840923665.test."
}

resource "alibabacloudstack_dns_forward_domain" "default" {
  name         = var.name
  remark       = "Created by Terraform"
  forward_mode = "FORWARD_FIRST"
  forwarders   = ["192.168.101.1"]
}

data "alibabacloudstack_dns_forward_domains" "default" {
  name_regex = alibabacloudstack_dns_forward_domain.default.name
}
```

## Argument Reference

The following arguments support result filtering:

* `forward_mode` - (String, Optional) Forwarding mode. Valid values are `"FORWARD_FIRST"` or `"FORWARD_ONLY"`. Value description:
  * `FORWARD_ONLY`: Full forwarding mode (recommended) - all requests are forwarded only (no recursion).
  * `FORWARD_FIRST`: Priority forwarding mode - if forwarding fails, it falls back to internet recursion.
* `ids` - (List, Optional) List of forwarding domain IDs for filtering results. Only forwarding domains with IDs in the list will be returned.
* `name` - (String, Optional) Forwarding domain name for filtering results. Matches forwarding domains containing the specified name.
* `name_regex` - (String, Optional) Regular expression for the forwarding domain name for filtering results. Only forwarding domains with names matching the regular expression will be returned.

## Attributes Reference

The following attributes are exported:

* `forward_domains` - (List) List of matching forwarding domains. Each element contains:
  * `id` - (String) Forwarding domain ID.
  * `caller_uid` - (String) Caller UID.
  * `create_timestamp` - (Integer) Creation timestamp in seconds.
  * `forward_mode` - (String) Forwarding mode. Valid values:
    * `FORWARD_ONLY`: Full forwarding mode (recommended) - all requests are forwarded only (no recursion).
    * `FORWARD_FIRST`: Priority forwarding mode - if forwarding fails, it falls back to internet recursion.
  * `forwarders` - (Set) List of forwarder IP addresses.
  * `name` - (String) Forwarding domain name.
  * `remark` - (String) Remark.
  * `update_timestamp` - (Integer) Update timestamp in seconds.
* `ids` - (List) List of matching forwarding domain IDs.