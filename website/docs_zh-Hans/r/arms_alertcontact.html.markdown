---
subcategory: "ARMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_alertcontact"
sidebar_current: "docs-Alibabacloudstack-arms-alertcontact"
description: |- 
  编排应用实时监控服务(ARMS)警报联系人
---

# alibabacloudstack_arms_alertcontact
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_arms_alert_contact`

使用Provider配置的凭证在指定的资源集下编排应用实时监控服务(ARMS)警报联系人。

## 示例用法

```hcl
variable "name" {
    default = "tf-testaccarmsalert_contact84949"
}

resource "alibabacloudstack_arms_alertcontact" "default" {
  alert_contact_name     = var.name
  ding_robot_webhook_url = "https://oapi.dingtalk.com/robot/send?access_token=91f2f6****"
  email                  = "someone@example.com"
  phone_num              = "1381111****"
  system_noc            = true
}
```

## 参数说明

支持以下参数：

* `alert_contact_name` - (必填) 警报联系人的名称。这是您阿里巴巴云账户内警报联系人的唯一标识符。
* `ding_robot_webhook_url` - (可选) 钉钉机器人的Webhook URL。有关如何获取URL的更多信息，请参阅[配置钉钉机器人发送告警通知](https://www.alibabacloud.com/help/en/doc-detail/106247.htm)。必须指定以下参数中的至少一个：`phone_num`、`email` 和 `ding_robot_webhook_url`。
* `email` - (可选) 警报联系人的电子邮件地址。必须指定以下参数中的至少一个：`phone_num`、`email` 和 `ding_robot_webhook_url`。
* `phone_num` - (可选) 警报联系人的手机号码。必须指定以下参数中的至少一个：`phone_num`、`email` 和 `ding_robot_webhook_url`。
* `system_noc` - (可选) 指定警报联系人是否接收系统通知。有效值：
  * `true`: 警报联系人接收系统通知。
  * `false`: 警报联系人不接收系统通知。默认值：`false`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 警报联系人的ID。它由阿里云自动生成，并唯一标识该警报联系人。
* `alert_contact_name` - 警报联系人的名称。这是您阿里巴巴云账户内警报联系人的唯一标识符。