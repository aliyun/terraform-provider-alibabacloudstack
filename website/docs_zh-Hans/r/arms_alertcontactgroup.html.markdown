---
subcategory: "ARMS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_alertcontactgroup"
sidebar_current: "docs-Alibabacloudstack-arms-alertcontactgroup"
description: |- 
  编排应用实时监控服务(ARMS)警报联系人组
---

# alibabacloudstack_arms_alertcontactgroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_arms_alert_contact_group`

使用Provider配置的凭证在指定的资源集下编排应用实时监控服务(ARMS)警报联系人组。

## 示例用法

### 基础用法

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
  alert_contact_group_name = var.name
  contact_ids = [alibabacloudstack_arms_alert_contact.example.id]
}
```

## 参数参考

支持以下参数：

* `alert_contact_group_name` - (必填) 告警联系组的名称。它在指定的阿里云账户和区域中必须唯一。
* `contact_ids` - (可选) 属于此组的告警联系人的 ID 列表。这些 ID 可以从 `alibabacloudstack_arms_alert_contact` 资源的 `id` 属性或在 ARMS 控制台中手动创建获取。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - 告警联系组的 ID。这是在创建时自动生成的，可以用于将资源导入 Terraform。
```