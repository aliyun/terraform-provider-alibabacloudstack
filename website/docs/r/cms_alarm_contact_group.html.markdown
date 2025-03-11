---
subcategory: "Cloud Monitor"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_alarm_contact_group"
sidebar_current: "docs-alibabacloudstack-resource-cms-alarm-contact-group"
description: |-
  Provides a Alibabacloudstack CMS Alarm Contact Group resource.
---

# alibabacloudstack\_cms\_alarm\_contact\_group

Provides a CMS Alarm Contact Group resource.

For information about CMS Alarm Contact Group and how to use it, see [What is Alarm Contact Group](https://www.alibabacloud.com/help/en/doc-detail/114929.htm).

-> **NOTE:** Available in v1.101.0+.

## Example Usage

Basic Usage

```terraform
resource "alibabacloudstack_cms_alarm_contact_group" "example" {
  alarm_contact_group_name = "tf-test"
}
```

## Argument Reference

The following arguments are supported:

* `alarm_contact_group_name` - (Required, ForceNew) The name of the alarm group.
* `contacts` - (Optional) The name of the alert contact.
* `describe` - (Optional) The description of the alert group.
* `enable_subscribed` - (Optional) Whether to open weekly subscription.

## Attributes Reference

The following attributes are exported:

* `id` - The resource ID of Alarm Contact Group.

## Import

CMS Alarm Contact Group can be imported using the id, e.g.

```
$ terraform import alibabacloudstack_cms_alarm_contact_group.example tf-testacc123
```
