---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_alarmcontactgroup"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-alarmcontactgroup"
description: |- 
  编排云监控服务（CMS）报警联系人组
---

# alibabacloudstack_cloudmonitorservice_alarmcontactgroup
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarm_contact_group`

使用Provider配置的凭证在指定的资源集下编排云监控服务（CMS）报警联系人组。

## 示例用法

```hcl
variable "name" {
    default = "tf-testacccloud_monitor_servicealarm_contact_group59875"
}

resource "alibabacloudstack_cloudmonitorservice_alarmcontactgroup" "default" {
  alarm_contact_group_name = var.name
  describe                = "This is a test description for the alarm contact group."
  contacts               = ["Contact1", "Contact2"]
  enable_subscribed      = true
}
```

## 参数说明

支持以下参数：

* `alarm_contact_group_name` - (必填, 变更时重建) 报警联系组的名称。此名称在您的账户内必须唯一，并且创建后无法修改。
* `contacts` - (选填) 属于该报警联系组的联系人名称列表。当报警触发时，这些联系人将收到通知。
* `describe` - (必填) 报警联系组的描述信息。这提供了关于该组目的的额外上下文或详细信息。
* `enable_subscribed` - (选填) 指示是否为报警联系组启用每周订阅通知。默认值为 `false`。

## 属性说明

除了上述所有参数外，还导出了以下属性：

* `id` - 报警联系组的ID。它由阿里云自动生成，可用于资源识别。
* `enable_subscribed` - 指示是否为报警联系组启用了每周订阅通知。此属性反映了当前设置的状态。