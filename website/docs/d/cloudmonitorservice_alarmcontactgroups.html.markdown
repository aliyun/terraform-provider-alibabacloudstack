---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontactgroups"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-alarmcontactgroups"
description: |- 
  Provides a list of cloudmonitorservice alarmcontactgroups owned by an alibabacloudstack account.
---

# alibabacloudstack_cloudmonitorservice_alarmcontactgroups
-> **NOTE:** Alias name has: `alibabacloudstack_cms_alarm_contact_groups`

This data source provides a list of cloudmonitorservice alarmcontactgroups in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```hcl
data "alibabacloudstack_cloudmonitorservice_alarmcontactgroups" "example" {
  name_regex = "tf-testacc"
}

output "alarm_contact_groups" {
  value = data.alibabacloudstack_cloudmonitorservice_alarmcontactgroups.example.groups
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Alarm Contact Group IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Alarm Contact Group name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of Alarm Contact Group IDs.
* `names` - A list of Alarm Contact Group names.
* `groups` - A list of CloudMonitorService Alarm Contact Groups. Each element contains the following attributes:
  * `id` - The ID of the Alarm Contact Group.
  * `alarm_contact_group_name` - The name of the Alarm Contact Group.
  * `contacts` - A list of alarm contacts associated with this group.
    * `contact_id` - The ID of the contact.
    * `contact_name` - The name of the contact.
  * `describe` - The description of the Alarm Contact Group.
  * `enable_subscribed` - Indicates whether the alarm group subscribes to weekly reports. Valid values are `true` or `false`. If set to `true`, the alarm group will receive weekly summary reports.