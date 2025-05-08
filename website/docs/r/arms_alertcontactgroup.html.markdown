---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_alertcontactgroup"
sidebar_current: "docs-Alibabacloudstack-arms-alertcontactgroup"
description: |- 
  Provides a arms Alertcontactgroup resource.
---

# alibabacloudstack_arms_alertcontactgroup
-> **NOTE:** Alias name has: `alibabacloudstack_arms_alert_contact_group`

Provides a arms Alertcontactgroup resource.

## Example Usage

Basic Usage

```terraform
variable "name" {
  default = "tf-testAccArmsAlertContactGroup1776884"
}

resource "alibabacloudstack_arms_alert_contact" "example" {
  alert_contact_name     = "example_value"
  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f6****"
  email                  = "someone@example.com"
  phone_num              = "1381111****"
}

resource "alibabacloudstack_arms_alert_contact_group" "default" {
  alert_contact_group_name = "${var.name}"
  contact_ids = [alibabacloudstack_arms_alert_contact.example.id]
}
```

## Argument Reference

The following arguments are supported:

* `alert_contact_group_name` - (Required) The name of the alert contact group. It must be unique within the specified Alibaba Cloud account and region.
* `contact_ids` - (Optional) A list of IDs for alert contacts that belong to this group. These IDs can be obtained from the `id` attribute of the `alibabacloudstack_arms_alert_contact` resource or manually created in the ARMS console.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alert contact group. This is automatically generated upon creation and can be used for importing the resource into Terraform.