---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_private_line"
sidebar_current: "docs-Alibabacloudstack-resource-dns-private-line"
description: |-
  Manage Alibaba Cloud DNS private lines
---

# alibabacloudstack_dns_private_line

Manages a DNS private line resource within Alibaba Cloud.

## Example Usage

### Basic Usage

```hcl
variable "name" {
  default = "tfacc32375"
}

resource "alibabacloudstack_dns_private_line" "default" {
  name = var.name
  v4_addresses = [
    "192.168.0.1"
  ]
  v6_addresses = [
    "2020:148:2:28::",
    "2020:148:3:28::"
  ]
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) The name of the DNS private line. The name is used to identify the line. For length and format restrictions, refer to the Alibaba Cloud API documentation.
* `v4_addresses` - (Optional) A list of IPv4 addresses. At least one of IPv4 or IPv6 addresses must be specified. Fuzzy query is supported for IPv4 addresses.
* `v6_addresses` - (Optional) A list of IPv6 addresses. At least one of IPv4 or IPv6 addresses must be specified. IPv6 addresses only support exact query.

## Attributes Reference

The following attributes are exported:

* `id` - The ID of the DNS private line.
* `priority` - The priority of the DNS private line. 1 indicates the highest priority, and a larger value indicates a lower priority. The system automatically assigns the priority value.