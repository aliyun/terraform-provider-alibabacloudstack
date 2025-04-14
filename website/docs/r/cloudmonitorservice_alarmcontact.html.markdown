---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontact"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-alarmcontact"
description: |- 
  Provides a cloudmonitorservice Alarmcontact resource.
---

# alibabacloudstack_cloudmonitorservice_alarmcontact
-> **NOTE:** Alias name has: `alibabacloudstack_cms_alarm_contact`

Provides a cloudmonitorservice Alarmcontact resource.

## Example Usage

Basic Usage:

```terraform
variable "name" {
    default = "tf-testacccloud_monitor_servicealarm_contact54950"
}

resource "alibabacloudstack_cloudmonitorservice_alarmcontact" "default" {
  alarm_contact_name = "Alice122"
  describe           = "报警联系人信息"
  channels_ali_im   = "leo"
  channels_ding_web_hook = "https://oapi.dingtalk.com/robot/send?access_token=7d49515e8ebf21106a80a9cc4bb3d2"
  channels_mail      = "alice@example.com"
  channels_sms       = "+1234567890"
}
```

Advanced Usage with lifecycle ignore_changes:

```terraform
resource "alibabacloudstack_cloudmonitorservice_alarmcontact" "example" {
  alarm_contact_name = "zhangsan"
  describe           = "For Test"
  channels_ali_im   = "zhangsan_trade_manager_id"
  channels_ding_web_hook = "https://oapi.dingtalk.com/robot/send?access_token=abcde12345"
  channels_mail      = "terraform.test.com"
  channels_sms       = "+0987654321"

  lifecycle {
    ignore_changes = [channels_mail]
  }
}
```

## Argument Reference

The following arguments are supported:

* `alarm_contact_name` - (Required, ForceNew) The name of the alarm contact. This name must be unique within your Alibaba Cloud account and cannot be modified after creation.
* `describe` - (Required) A brief description of the alarm contact. This helps in identifying the purpose or owner of the contact.
* `channels_ali_im` - (Optional) The TradeManager ID of the alarm contact. This is used for receiving notifications via TradeManager.
* `channels_ding_web_hook` - (Optional) The webhook URL of the DingTalk chatbot. This allows sending notifications to a DingTalk group.
* `channels_mail` - (Optional) The email address of the alarm contact. After you add or modify an email address, the recipient receives an email that contains an activation link. The system adds the recipient to the list of alarm contacts only after the recipient activates the email address.
* `channels_sms` - (Optional) The phone number of the alarm contact. Notifications will be sent via SMS to this number. Similar to email, an activation link may be required before the contact is added.
* `lang` - (Optional) The language type of the alarm. Valid values: `en`, `zh-cn`. Defaults to `zh-cn`.

-> **NOTE:** Specify at least one of the following alarm notification targets: `channels_ali_im`, `channels_ding_web_hook`, `channels_mail`, `channels_sms`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the alarm contact. It is the same as `alarm_contact_name`.
* `channels_aliim` - Deprecated. Use `channels_ali_im` instead.
* `channels_ali_im` - The TradeManager ID of the alarm contact.
* `channels_ding_web_hook` - The webhook URL of the DingTalk chatbot.
* `channels_mail` - The email address of the alarm contact.
* `channels_sms` - The phone number of the alarm contact.
* `describe` - A brief description of the alarm contact.

## Import

Alarm contact can be imported using the id, e.g.

```bash
$ terraform import alibabacloudstack_cloudmonitorservice_alarmcontact.example abc12345
```