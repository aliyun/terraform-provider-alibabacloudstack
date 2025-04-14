---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontact"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-alarmcontact"
description: |- 
  云监控服务（CMS）报警联系人
---

# alibabacloudstack_cloudmonitorservice_alarmcontact
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarm_contact`

使用Provider配置的凭证在指定的资源集下编排云监控服务（CMS）报警联系人。

## 示例用法

### 基础用法：

```terraform
variable "name" {
    default = "tf-testacccloud_monitor_servicealarm_contact54950"
}

resource "alibabacloudstack_cloudmonitorservice_alarmcontact" "default" {
  alarm_contact_name = "Alice122"
  describe           = "报警联系人信息"
  channels_ali_im   = "leo"
  channels_ding_web_hook = "https://oapi.dingtalk.com/robot/send?access_token=abcde12345"
  channels_mail      = "alice@example.com"
  channels_sms       = "+1234567890"
}
```

高级用法，包含生命周期忽略更改：

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

## 参数说明

支持以下参数：

* `alarm_contact_name` - (必填，变更时重建) 报警联系人的名称。该名称在您的阿里云账户内必须唯一，并且创建后无法修改。
* `describe` - (必填) 报警联系人的简要描述。这有助于识别联系人的目的或所有者。
* `channels_ali_im` - (选填) 报警联系人的TradeManager ID。用于通过TradeManager接收通知。
* `channels_ding_web_hook` - (选填) 钉钉群机器人的webhook URL。允许将通知发送到钉钉群。
* `channels_mail` - (选填) 报警联系人的电子邮件地址。添加或修改电子邮件地址后，接收方会收到一封包含激活链接的电子邮件。只有在接收方激活电子邮件地址后，系统才会将其添加到报警联系人列表中。
* `channels_sms` - (选填) 报警联系人的电话号码。通知将通过短信发送到此号码。与电子邮件类似，在联系人被添加之前可能需要激活链接。
* `lang` - (选填) 报警的语言类型。有效值：`en`、`zh-cn`。默认为`zh-cn`。

-> **注意：** 至少需要指定以下报警通知目标之一：`channels_ali_im`、`channels_ding_web_hook`、`channels_mail`、`channels_sms`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 报警联系人的ID。它与`alarm_contact_name`相同。
* `channels_aliim` - 已弃用。请改用`channels_ali_im`。
* `channels_ali_im` - 报警联系人的TradeManager ID。
* `channels_ding_web_hook` - 钉钉群机器人的webhook URL。
* `channels_mail` - 报警联系人的电子邮件地址。
* `channels_sms` - 报警联系人的电话号码。
* `describe` - 报警联系人的简要描述。

## 导入

报警联系人可以通过ID导入，例如：

```bash
$ terraform import alibabacloudstack_cloudmonitorservice_alarmcontact.example abc12345
```