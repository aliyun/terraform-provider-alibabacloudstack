---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_group"
sidebar_current: "docs-alibabacloudstack-resource-dns-group"
description: |-
  Provides a DNS Group resource.
---

# alibabacloudstack_dns_group
-> **NOTE:** Alias name has: `alibabacloudstack_alidns_domaingroup`

Provides a DNS Group resource.

## Example Usage

```
# Add a new Domain group.
resource "alibabacloudstack_dns_group" "group" {
  name = "testgroup"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) Name of the domain group.    

* `child` - (Optional) The child property of the domain group. 

* `propreties` - (Required) The properties of the domain group.

## Attributes Reference

The following attributes are exported:

* `id` - The group id.
* `name` - The group name.
* `propreties` - The properties of the domain group.