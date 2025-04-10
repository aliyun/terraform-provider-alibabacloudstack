---
subcategory: "Application Real-Time Monitoring Service (ARMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_arms_dispatch_rule"
sidebar_current: "docs-alibabacloudstack-resource-arms-dispatch-rule"
description: |-
  编排应用实时监控服务(ARMS)告警分派规则
---

# alibabacloudstack_arms_dispatch_rule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_arms_dispatchrule`

使用Provider配置的凭证在指定的资源集下编排应用实时监控服务(ARMS)告警分派规则资源。

关于应用实时监控服务(ARMS)告警分派规则的更多信息以及如何使用它，请参见 [什么是告警分派规则](https://next.api.alibabacloud.com/document/ARMS/2019-08-08/CreateDispatchRule)。


## 示例用法

### 基础用法

```terraform
resource "alibabacloudstack_arms_alert_contact" "default" {
  alert_contact_name = "example_value"
  email              = "example_value@aaa.com"
}
resource "alibabacloudstack_arms_alert_contact_group" "default" {
  alert_contact_group_name = "example_value"
  contact_ids              = [alibabacloudstack_arms_alert_contact.default.id]
}

resource "alibabacloudstack_arms_dispatch_rule" "default" {
  dispatch_rule_name = "example_value"
  dispatch_type      = "CREATE_ALERT"
  group_rules {
    group_wait_time = 5
    group_interval  = 15
    repeat_interval = 100
    grouping_fields = [
    "alertname"]
  }
  label_match_expression_grid {
    label_match_expression_groups {
      label_match_expressions {
        key      = "_aliyun_arms_involvedObject_kind"
        value    = "app"
        operator = "eq"
      }
    }
  }

  notify_rules {
    notify_objects {
      notify_object_id = alibabacloudstack_arms_alert_contact.default.id
      notify_type      = "ARMS_CONTACT"
      name             = "example_value"
    }
    notify_objects {
      notify_object_id = alibabacloudstack_arms_alert_contact_group.default.id
      notify_type      = "ARMS_CONTACT_GROUP"
      name             = "example_value"
    }
    notify_channels = ["dingTalk", "wechat"]
  }
}
```

## 参数参考

支持以下参数：

* `dispatch_rule_name` - (必填) 分派规则的名称。
* `dispatch_type` - (可选) 告警处理方式。有效值：CREATE_ALERT：生成告警；DISCARD_ALERT：丢弃告警事件并生成无告警。
* `is_recover` - (可选) 是否发送恢复告警。有效值：true：发送告警；false：不发送告警。
* `group_rules` - (必填) 设置事件组。当 `dispatch_type = "DISCARD_ALERT"` 时，此参数将被忽略。
  * `group_wait_time` - (必填) 系统在发送第一个告警后等待的时间段（单位：秒）。在此时间段后，所有告警会以单个通知的形式发送给处理器。
  * `group_interval` - (必填) 系统在发送第一个告警后等待的时间段（单位：秒）。在此时间段后，所有告警会以单个通知的形式发送给处理器。
  * `grouping_fields` - (必填，List<String>) 用于分组事件的字段。具有相同字段内容的事件会被分配到一个组中。具有相同指定分组字段的告警会以单独的通知形式发送给处理器。
  * `repeat_interval` - (可选) 重复告警的静默期（单位：秒）。所有告警会在指定的时间间隔内重复发送，直到告警清除为止。最小值为61。默认为600。
  * `group_id` - (可选) 组规则的ID。
* `label_match_expression_grid` - (必填) 设置分派规则。
  * `label_match_expression_groups` - (必填) 设置分派规则。
    * `label_match_expressions` - (必填) 设置分派规则。
      * `key` - (必填) 分派规则标签的键。有效值：
        * `_aliyun_arms_userid`：用户ID
        * `_aliyun_arms_involvedObject_kind`：关联对象类型
        * `_aliyun_arms_involvedObject_id`：关联对象ID 
        * `_aliyun_arms_involvedObject_name`：关联对象名称
        * `_aliyun_arms_alert_name`：告警名称
        * `_aliyun_arms_alert_rule_id`：告警规则ID
        * `_aliyun_arms_alert_type`：告警类型
        * `_aliyun_arms_alert_level`：告警严重性
      * `value` - (必填) 标签的值。
      * `operator` - (必填) 分派规则中使用的运算符。有效值：
        * `eq`：等于。
        * `re`：匹配正则表达式。
* `notify_rules` - (必填) 设置通知规则。当 `dispatch_type = "DISCARD_ALERT"` 时，此参数将被忽略。
  * `notify_objects` - (必填) 设置通知对象。
    * `notify_object_id` - (必填) 联系人或联系组的ID。
    * `name` - (必填) 联系人或联系组的名称。
    * `notify_type` - (必填) 告警联系人的类型。有效值：`ARMS_CONTACT`：联系人；`ARMS_CONTACT_GROUP`：联系组。
  * `notify_channels` - (必填，List<String>) 通知方法。有效值：`dingTalk`、`sms`、`webhook`、`email` 和 `wechat`。


## 属性参考

导出以下属性：

* `id` - Terraform 中告警分派规则的资源ID。
* `status` - 告警分派规则的资源状态。

## 导入

应用实时监控服务(ARMS)告警联系人可以使用id导入，例如：

```shell
$ terraform import alibabacloudstack_arms_dispatch_rule.example <id>
```