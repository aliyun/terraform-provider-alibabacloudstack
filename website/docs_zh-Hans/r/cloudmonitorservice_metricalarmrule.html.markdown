---
subcategory: "Cloud Monitor Service (CMS)"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricalarmrule"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-metricalarmrule"
description: |- 
  云监控服务（CMS）报警监控项规则
---

# alibabacloudstack_cloudmonitorservice_metricalarmrule
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_alarm`

使用Provider配置的凭证在指定的资源集下编排云监控服务（CMS）报警监控项规则。

## 示例用法

### 基础用法

```hcl
variable "name" {
  default = "tf_testacc_cmsalarm4649479"
}

resource "alibabacloudstack_slb" "basic" {
  name = "${var.name}"
}

resource "alibabacloudstack_cloudmonitorservice_metricalarmrule" "default" {
  rule_name    = "${var.name}"
  namespace    = "acs_slb_dashboard"
  metric_name  = "ActiveConnection"
  dimensions   = {
    instanceId = alibabacloudstack_slb.basic.id
  }
  period       = 300
  escalations_critical {
    statistics        = "Average"
    comparison_operator = "<="
    threshold         = 35
    times             = 2
  }
  enabled          = true
  contact_groups   = ["test-group"]
  effective_interval = "0:00-2:00"
  silence_time     = 86400
  webhook          = ""
}
```

## 参数参考

支持以下参数：

* `rule_name` - (必填) 报警规则的名称。
* `namespace` - (必填，变更时重建) 数据命名空间用于区分不同的产品。例如，ECS为`acs_ecs_dashboard`，RDS为`acs_rds_dashboard`。
* `metric_name` - (必填，变更时重建) 监控项名称。例如，ECS的`CPUUtilization`或SLB的`ActiveConnection`。更多详细信息，请参阅[指标参考](https://www.alibabacloud.com/help/doc-detail/28619.htm)。
* `dimensions` - (必填，变更时重建) 与报警规则关联的资源映射，如“instanceId”、“device”和“port”。每个键的值是一个字符串，并使用逗号分隔多个项。更多信息，请参见[指标参考](https://www.alibabacloud.com/help/doc-detail/28619.htm)。
* `period` - (可选) 统计周期(秒)。有效值取决于指标。默认值为`300`。
* `escalations_critical` - (可选) 定义关键报警设置的配置块(如下文所述)。每个报警规则只能定义一个级别。
  * `statistics` - 关键级别报警统计方法。它必须与为指标定义的内容一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
  * `comparison_operator` - 关键级别报警比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
  * `threshold` - 关键级别报警阈值，目前必须是数值。
  * `times` - 关键级别报警重试次数。默认为`3`。
* `escalations_warn` - (可选) 定义警告报警设置的配置块(如下文所述)。每个报警规则只能定义一个级别。
  * `statistics` - 警告级别报警统计方法。它必须与为指标定义的内容一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
  * `comparison_operator` - 警告级别报警比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
  * `threshold` - 警告级别报警阈值，目前必须是数值。
  * `times` - 警告级别报警重试次数。默认为`3`。
* `escalations_info` - (可选) 定义信息报警设置的配置块(如下文所述)。每个报警规则只能定义一个级别。
  * `statistics` - 信息级别报警统计方法。它必须与为指标定义的内容一致。有效值：["Average", "Minimum", "Maximum"]。默认为"Average"。
  * `comparison_operator` - 信息级别报警比较运算符。有效值：["<=", "<", ">", ">=", "==", "!="]。默认为"=="。
  * `threshold` - 信息级别报警阈值，目前必须是数值。
  * `times` - 信息级别报警重试次数。默认为`3`。
* `contact_groups` - (必填) 报警规则的联系组列表，必须已在控制台上创建。报警通知将发送到报警联系组中的报警联系人。
* `effective_interval` - (可选) 报警规则生效的时间间隔。格式为"hh:mm-hh:mm"，如"0:00-4:00"。默认为"00:00-23:59"。
* `silence_time` - (可选) 在报警状态下的通知静默期，单位为秒。有效范围：[300, 86400]。默认为`86400`。
* `enabled` - (可选) 是否启用报警规则。默认为`true`。
* `webhook` - (可选) 当报警触发时应调用的Webhook。目前仅支持HTTP协议。默认为空字符串。

> **注意**：每个资源支持创建以下三个级别之一：关键、警告或信息。

## 属性参考

除了上述所有参数外，还导出以下属性：

* `rule_id` - 报警规则的唯一ID。
* `status` - 报警规则的当前状态(`true`表示启用，`false`表示禁用)。
