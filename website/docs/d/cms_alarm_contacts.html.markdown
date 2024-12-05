---
subcategory: "Cloud Monitor Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_alarm_contacts"
sidebar_current: "docs-alibabacloudstack-resource-cms-alarm-contacts"
description: |-
  Provides a list of alarm contact owned by an Alibaba Cloud account.
---

# alibabacloudstack\_cms\_alarm\_contacts

Provides a list of alarm contact owned by an Alibaba Cloud account.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_cms_alarm_contacts" "example" {
  ids = ["tf-testAccCmsAlarmContact"]
}
output "first-contact" {
  value = data.alibabacloudstack_cms_alarm_contacts.this.contacts
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of alarm contact IDs. 
* `name_regex` - (Optional, ForceNew) A regex string to filter results by alarm contact name. 
* `chanel_type` - (Optional, ForceNew)  The alarm notification method. Alarm notifications can be sent by using `Email` or `DingWebHook`.
* `chanel_value` - (Optional, ForceNew)  The alarm notification target.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`). 

-> **NOTE:** Specify at least one of the following alarm notification targets: phone number, email address, webhook URL of the DingTalk chatbot, and TradeManager ID.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of alarm contact IDs.
* `names` - A list of alarm contact names.
* `contacts` - A list of alarm contacts. Each element contains the following attributes:
    * `id` - The ID of the alarm contact.
    * `alarm_contact_name` - The name of the alarm contact.
    * `channels_aliim` - The TradeManager ID of the alarm contact.
    * `channels_ding_web_hook` - The webhook URL of the DingTalk chatbot.
    * `channels_mail` - The email address of the alarm contact. 
    * `channels_sms` - The phone number of the alarm contact.
    * `describe` - The description of the alarm contact.
    * `contact_groups` - The alert groups to which the alarm contact is added.
    * `channels_state_aliim` - Indicates whether the TradeManager ID is valid.
    * `channels_state_ding_web_hook` - Indicates whether the DingTalk chatbot is normal.
    * `channels_state_mail` - The status of the email address.
    * `channels_status_sms` - The status of the phone number.
    * `Lang` - The language type of the alarm.
