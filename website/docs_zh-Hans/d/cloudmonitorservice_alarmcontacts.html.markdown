---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontacts"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-alarmcontacts"
description: |- 
  查询云监控报警联系人
---

# alibabacloudstack_cloudmonitorservice_alarmcontacts
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarm_contacts`

根据指定过滤条件列出当前凭证权限可以访问的云监控报警联系人列表。

## 示例用法

### 基础用法

```terraform
data "alibabacloudstack_cloudmonitorservice_alarmcontacts" "example" {
  ids = ["tf-testAccCmsAlarmContact"]
}

output "first-contact" {
  value = data.alibabacloudstack_cloudmonitorservice_alarmcontacts.example.contacts[0]
}
```

## 参数参考

以下参数是支持的：

* `ids` - (可选, 变更时重建) 报警联系人的ID列表。用于通过报警联系人ID进行精确过滤。
* `name_regex` - (可选, 变更时重建) 用于通过报警联系人名称过滤结果的正则表达式字符串。例如，可以使用 `^test-` 来匹配所有以 `test-` 开头的报警联系人。
* `chanel_type` - (可选, 变更时重建) 报警通知方式。可以通过 `Email`, `DingWebHook` 或其他方法发送报警通知。
* `chanel_value` - (可选, 变更时重建) 报警通知目标，例如电子邮件地址或钉钉机器人的webhook URL。

-> **注意：** 至少需要指定以下一个报警通知目标：电话号码、电子邮件地址、钉钉群机器人webhook URL 和 TradeManager ID。

## 属性参考

除了上述参数外，还导出以下属性：

* `contacts` - 报警联系人列表。每个元素包含以下属性：
    * `id` - 报警联系人的唯一标识符。
    * `alarm_contact_name` - 报警联系人的名称。
    * `channels_aliim` - 报警联系人的TradeManager ID。
    * `channels_ding_web_hook` - 钉钉群机器人的webhook URL。
    * `channels_mail` - 报警联系人的电子邮件地址。
    * `channels_sms` - 报警联系人的电话号码。
    * `describe` - 报警联系人的描述信息。
    * `contact_groups` - 添加该报警联系人的告警组列表。
    * `channels_state_aliim` - 表示TradeManager ID是否有效(布尔值)。
    * `channels_state_ding_web_hook` - 表示钉钉群机器人是否正常(布尔值)。
    * `channels_state_mail` - 电子邮件地址的状态(布尔值)。
    * `channels_status_sms` - 电话号码的状态(布尔值)。
    * `Lang` - 报警的语言类型(如 `zh` 表示中文，`en` 表示英文)。