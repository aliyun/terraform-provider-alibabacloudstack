---
subcategory: "DMSEnterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dmsenterprise_user"
sidebar_current: "docs-Alibabacloudstack-dmsenterprise-user"
description: |- 
  Provides a dmsenterprise User resource.
---

# alibabacloudstack_dmsenterprise_user
-> **NOTE:** Alias name has: `alibabacloudstack_dms_enterprise_user`

Provides a dmsenterprise User resource.

## Example Usage

```terraform
variable "name" {
    default = "tf-testaccdms_enterpriseuser93463"
}

resource "alibabacloudstack_dmsenterprise_user" "default" {
  mobile      = "11111111111"
  uid         = "265530631068325049"
  user_name   = "rdktest"
  role_names  = ["DBA"]
  status      = "NORMAL"
  max_execute_count = 100
  max_result_count  = 500
}
```

## Argument Reference

The following arguments are supported:

* `uid` - (Required, ForceNew) The Alibaba Cloud unique ID (UID) of the user to add. This field cannot be modified after creation.
* `user_name` - (Optional) The nickname of the user.
* `mobile` - (Optional) The DingTalk number or mobile number of the user.
* `role_names` - (Optional) The roles that the user plays. For example: `["DBA"]`.
* `status` - (Optional) The state of the DMS Enterprise User. Valid values: `NORMAL`, `DISABLE`.
* `max_execute_count` - (Optional) The maximum number of queries allowed for the user on the day.
* `max_result_count` - (Optional) The maximum number of rows that can be queried by the user on the day.
* `tid` - (Optional) The tenant ID. This is the ID of the tenant displayed in the upper right corner of the system. For more information, see [view tenant information](https://www.alibabacloud.com/help/doc-detail/181330.htm).
* `nick_name` - (Optional) The deprecated nickname of the user. It is recommended to use `user_name` instead.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `id` - The Alibaba Cloud unique ID (UID) of the user. The value is the same as the UID.
* `nick_name` - The nickname of the user (deprecated, use `user_name` instead).
* `role_names` - The list of roles that the user plays.
* `status` - The state of the DMS Enterprise User.
* `mobile` - The DingTalk number or mobile number of the user.

## Import

DMS Enterprise User can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_dmsenterprise_user.example <uid>
```