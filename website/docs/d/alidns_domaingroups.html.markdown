---
subcategory: "DNS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dns_groups"
sidebar_current: "docs-alibabacloudstack-datasource-dns-groups"
description: |-
    Provides a list of groups available to the dns.
---

# alibabacloudstack_dns_groups
-> **NOTE:** Alias name has: `alibabacloudstack_alidns_domaingroups`

This data source provides a list of DNS Domain Groups in an Alibabacloudstack Cloud account according to the specified filters.

## Example Usage

```
data "alibabacloudstack_dns_groups" "groups_ds" {
  name_regex  = "^y[A-Za-z]+"
}

output "first_group_name" {
  value = "${data.alibabacloudstack_dns_groups.groups_ds.groups.0.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by group name. 
* `ids` - (Optional) A list of group IDs.
* `names` - (Optional)  A list of group names.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of group IDs. 
* `names` - A list of group names.
* `groups` - A list of groups. Each element contains the following attributes:
  * `group_id` - Id of the group.
  * `group_name` - Name of the group.
* `name_regex` - A regex string to filter results by group name. 