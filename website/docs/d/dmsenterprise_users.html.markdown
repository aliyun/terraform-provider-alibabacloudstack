---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_users"
sidebar_current: "docs-Alibabacloudstack-datasource-dmsenterprise-users"
description: |- 
  Provides a list of dmsenterprise users owned by an alibabacloudstack account.
---

# alibabacloudstack_dmsenterprise_users
-> **NOTE:** Alias name has: `alibabacloudstack_dms_enterprise_users`

This data source provides a list of dmsenterprise users in an Alibabacloudstack account according to the specified filters.

## Example Usage

```terraform
# Declare the data source
data "alibabacloudstack_dmsenterprise_users" "example" {
  name_regex = "user-.*"
  role       = "USER"
  status     = "NORMAL"
  tid        = "1234567890"
}

output "first_user_id" {
  value = "${data.alibabacloudstack_dmsenterprise_users.example.users.0.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional, ForceNew) A regex string to filter the results by the DMS Enterprise User nick_name.
* `role` - (Optional, ForceNew) The role of the user to query.
* `search_key` - (Optional, ForceNew) The keyword used to query users.
* `status` - (Optional, ForceNew) The status of the user.
* `tid` - (Optional, ForceNew) The ID of the tenant in DMS Enterprise. This is the ID of the tenant displayed in the upper right corner of the system. For more information, see [view tenant information](~~ 181330 ~~).
* `ids` - (Optional, ForceNew) A list of DMS Enterprise User IDs (UID).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of DMS Enterprise User IDs (UID).
* `names` - A list of DMS Enterprise User names.
* `users` - A list of DMS Enterprise Users. Each element contains the following attributes:
  * `mobile` - The DingTalk number or mobile number of the user.
  * `nick_name` - The nickname of the user.
  * `user_name` - The nickname of the user.
  * `parent_uid` - The Alibaba Cloud unique ID (UID) of the parent account if the user corresponds to a Resource Access Management (RAM) user.
  * `role_ids` - The list of IDs of the roles that the user plays.
  * `role_names` - The list of names of the roles that the user plays.
  * `status` - The status of the user.
  * `id` - The Alibaba Cloud unique ID (UID) of the user.
  * `uid` - Alias of `id`.
  * `user_id` - The ID of the user.