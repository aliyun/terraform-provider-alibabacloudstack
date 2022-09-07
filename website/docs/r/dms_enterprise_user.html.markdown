---
subcategory: "DMS Enterprise"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_dms_enterprise_user"
sidebar_current: "docs-alibabacloudstack-resource-dms-enterprise-user"
description: |-
  Provides a DMS Enterprise User resource.
---

# alibabacloudstack\_dms\_enterprise\_user

Provides a DMS Enterprise User resource. For information about Alidms Enterprise User and how to use it, see [What is Resource Alidms Enterprise User](https://www.alibabacloud.com/help/doc-detail/98001.htm).

## Example Usage

```terraform
resource "alibabacloudstack_dms_enterprise_user" "example" {
  uid        = "uid"
  user_name  = "tf-test"
  role_names = ["DBA"]
  mobile     = "1591066xxxx"
}
```

## Argument Reference

The following arguments are supported:

* `tid` - (Optional) The tenant ID. 
* `uid` - (Required, ForceNew) The Alibaba Cloud unique ID (UID) of the user to add.
* `status` - (Optional) The state of DMS Enterprise User. Valid values: `NORMAL`, `DISABLE`.
* `role_names` - (Optional) The roles that the user plays.
* `nick_name` - (Optional) It has been deprecated and use `user_name` instead.
* `user_name` - (Optional) The nickname of the user.
* `mobile` - (Optional) The DingTalk number or mobile number of the user.
* `max_result_count` - (Optional) Query the maximum number of rows on the day.
* `max_execute_count` - (Optional) Maximum number of inquiries on the day.
                         
## Attributes Reference

The following attributes are exported:

* `id` - The Alibaba Cloud unique ID of the user. The value is same as the UID.
* `mobile` - The DingTalk number or mobile number of the user.
* `nick_name` - The nickname of the user.
* `role_names` - The list of roles that the user plays.
* `status` - The state of DMS Enterprise User.

## Import

DMS Enterprise User can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_dms_enterprise_user.example 24356xxx
```
