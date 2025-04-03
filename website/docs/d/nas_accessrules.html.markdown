---
subcategory: "NAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_nas_accessrules"
sidebar_current: "docs-Alibabacloudstack-datasource-nas-accessrules"
description: |- 
  Provides a list of nas accessrules owned by an alibabacloudstack account.
---

# alibabacloudstack_nas_accessrules
-> **NOTE:** Alias name has: `alibabacloudstack_nas_access_rules`

This data source provides a list of NAS Access Rules in an AlibabaCloudStack account according to the specified filters.

## Example Usage

```terraform
data "alibabacloudstack_nas_accessrules" "foo" {
  access_group_name = "tf-testAccAccessGroupsdatasource"
  source_cidr_ip    = "168.1.1.0/16"
  rw_access         = "RDWR"
  user_access       = "no_squash"
}

output "alibabacloudstack_nas_accessrules_id" {
  value = "${data.alibabacloudstack_nas_accessrules.foo.rules.0.access_rule_id}"
}
```

## Argument Reference

The following arguments are supported:

* `access_group_name` - (Required, ForceNew) The name of the access group. This is a required field and cannot be modified after creation.
* `source_cidr_ip` - (Optional) The CIDR block that specifies the IP range for the access rule.
* `rw_access` - (Optional) The read/write access type for the access rule. Valid values include:
  * `RDONLY`: Read-only access.
  * `RDWR`: Read-write access.
* `user_access` - (Optional) The user access type for the access rule. Valid values include:
  * `no_squash`: No root squashing.
  * `root_squash`: Root squashing.
  * `all_squash`: All squashing.
* `ids` - (Optional) A list of NAS Access Rule IDs to filter results.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of NAS Access Rule IDs.
* `rules` - A list of NAS Access Rules. Each element contains the following attributes:
  * `source_cidr_ip` - The CIDR block that specifies the IP range for the access rule.
  * `priority` - The priority of the access rule.
  * `access_rule_id` - The ID of the access rule.
  * `user_access` - The user access type for the access rule.
  * `rw_access` - The read/write access type for the access rule.
