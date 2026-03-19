---
subcategory: "Cloud DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_lines"
sidebar_current: "docs-Alibabacloudstack-datasource-dnslines"
description: |-
  Provides a list of DNS Global lines available to the user.
---

# alibabacloudstack_dns_lines

This data source queries the list of DNS Global lines in Alibaba Cloud.

## Example Usage

```hcl
variable "name" {
  default = "tfacc20899"
}

resource "alibabacloudstack_dns_line" "default" {
  name         = var.name
  v4_addresses = ["192.168.0.1"]
  v6_addresses = ["2020:148:2:28::", "2020:148:3:28::"]
}

data "alibabacloudstack_dns_lines" "default" {
  name_regex = alibabacloudstack_dns_line.default.name
}
```

## Argument Reference

The following arguments are used to filter query results:

* `name_regex` (String, Optional): Filters DNS Global line names by regular expression.
* `ids` (List of Strings, Optional): A list of DNS Global line IDs for exact matching of specified lines.

## Attributes Reference

The following attributes are exported:

* `id` (String): The unique identifier of the line.
* `create_timestamp` (Integer): The creation timestamp in seconds.
* `name` (String): The name of the line.
* `priority` (Integer): The priority of the line, where 1 is the highest priority and higher values indicate lower priority.
* `update_timestamp` (Integer): The update timestamp in seconds.
* `v4_addresses` (Set of Strings): A list of IPv4 addresses.
* `v6_addresses` (Set of Strings): A list of IPv6 addresses.