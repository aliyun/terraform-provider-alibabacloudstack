---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricalarmrules"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-metricalarmrules"
description: |- 
  查询云监控告警规则
---

# alibabacloudstack_cloudmonitorservice_metricalarmrules
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarms`

根据指定过滤条件列出当前凭证权限可以访问的云监控告警规则列表。

## 示例用法

```hcl
variable "name" {
  default = "tf_testacc_cmsalarm18358"
}

resource "alibabacloudstack_slb" "basic" {
  name          = var.name
}

resource "alibabacloudstack_cms_alarm" "default" {
  name                = var.name
  project             = "acs_slb_dashboard"
  metric              = "ActiveConnection"
  dimensions = {
    instanceId = alibabacloudstack_slb.basic.id
  }
  escalations_critical {
    statistics         = "Average"
    comparison_operator = "<="
    threshold          = 35
    times              = 2
  }
  enabled             = true
  contact_groups      = ["test-group"]
  effective_interval  = "0:00-2:00"
  
  lifecycle {
    ignore_changes = [
      dimensions,
      period,
    ]
  }
}

data "alibabacloudstack_cloudmonitorservice_metricalarmrules" "example" {
  name_regex = alibabacloudstack_cms_alarm.default.rule_name
}

output "metricalarmrules_list" {
  value = data.alibabacloudstack_cloudmonitorservice_metricalarmrules.example.alarms
}
```

## 参数参考

以下参数是支持的：

* `name_regex` - (可选) 用于按规则名称过滤告警规则的正则表达式字符串。该参数可以帮助用户根据告警规则名称进行精确或模糊匹配。

## 属性参考

除了上述参数外，还导出以下属性：

* `alarms` - 告警规则列表。每个元素包含以下属性：
  * `group_name` - 应用程序组的名称。
  * `metric_name` - 与告警规则关联的指标名称。
  * `no_effective_interval` - 告警规则无效的时间段。
  * `silence_time` - 在告警状态下的通知静默期，以秒为单位。有效值范围：[300, 86400]。默认为86400。
  * `contact_groups` - 告警规则的联系人组列表，这些联系人组必须已在控制台上创建。
  * `mail_subject` - 告警通知电子邮件的主题。
  * `source_type` - 告警规则的类型。有效值：
    - `METRIC`：时间序列指标的告警规则。
    - `EVENT`：事件类型指标的告警规则。这种类型在早期版本中使用，现已弃用。
  * `rule_id` - 告警规则的ID。
  * `period` - 索引查询周期，必须与为指标定义的一致。默认为300，以秒为单位。
  * `dimensions` - 与告警规则关联的资源映射，例如“instanceId”、“device”和“port”。每个键的值是一个字符串，并使用逗号分隔多个项。有关更多信息，请参阅 [Metrics Reference](https://www.alibabacloud.com/help/doc-detail/28619.htm)。
  * `effective_interval` - 生效告警规则的时间间隔。其格式为"hh:mm-hh:mm"，例如"0:00-4:00"。默认为"00:00-23:59"。
  * `namespace` - 监控服务的命名空间。
  * `enable_state` - 指示是否启用了告警规则。有效值：
    - `true`：表示已启用告警规则。
    - `false`：表示已禁用告警规则。如果未在请求中指定此参数，则默认返回所有启用和禁用的规则。
  * `webhook` - 回调URL。
  * `resources` - 与告警规则关联的资源。
  * `rule_name` - 告警规则名称。
  * `escalations` - 触发不同级别告警的条件。
    * `critical_comparison_operator` - 严重级别告警的比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
    * `critical_times` - 严重级别告警的重试次数。默认为3。
    * `critical_statistics` - 严重级别告警的统计方法。它必须与为指标定义的一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
    * `critical_threshold` - 严重级别告警的阈值，必须是数值。
    * `info_comparison_operator` - 信息级别告警的比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
    * `info_times` - 信息级别告警的重试次数。默认为3。
    * `info_statistics` - 信息级别告警的统计方法。它必须与为指标定义的一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
    * `info_threshold` - 信息级别告警的阈值，必须是数值。
    * `warn_comparison_operator` - 警告级别告警的比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
    * `warn_times` - 警告级别告警的重试次数。默认为3。
    * `warn_statistics` - 警告级别告警的统计方法。它必须与为指标定义的一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
    * `warn_threshold` - 警告级别告警的阈值，必须是数值。