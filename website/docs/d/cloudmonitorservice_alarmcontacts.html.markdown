---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontacts"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-alarmcontacts"
description: |- 
  Provides a list of cloudmonitorservice alarmcontacts in an Alibabacloudstack account according to the specified filters.

---

# alibabacloudstack_cloudmonitorservice_alarmcontacts
-> **NOTE:** Alias name has: `alibabacloudstack_cms_alarm_contacts`

This data source provides a list of cloudmonitorservice alarmcontacts in an Alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_cloudmonitorservice_alarmcontacts" "example" {
  ids = ["tf-testAccCmsAlarmContact"]
}

output "first-contact" {
  value = data.alibabacloudstack_cloudmonitorservice_alarmcontacts.example.contacts[0]
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of alarm contact IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by alarm contact name.
* `chanel_type` - (Optional, ForceNew) The alarm notification method. Alarm notifications can be sent by using `Email`, `DingWebHook`, or other methods.
* `chanel_value` - (Optional, ForceNew) The alarm notification target, such as email address or webhook URL.
* `names` - (Optional) A list of alarm contact names.

-> **NOTE:** Specify at least one of the following alarm notification targets: phone number, email address, webhook URL of the DingTalk chatbot, and TradeManager ID.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

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