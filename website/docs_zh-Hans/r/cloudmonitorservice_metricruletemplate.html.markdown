---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricruletemplate"
sidebar_current: "docs-Alibabacloudstack-cloudmonitorservice-metricruletemplate"
description: |- 
  云监控服务（CMS）报警监控项模板
---
# alibabacloudstack_cloudmonitorservice_metricruletemplate
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_metric_rule_template`

使用Provider配置的凭证在指定的资源集下编排云监控服务（CMS）报警监控项模板。

## 示例用法

### 基础用法

```terraform
variable "name" {
  default = "tf-testacccloudmonitorservicemetricruletemplate38509"
}

resource "alibabacloudstack_cms_monitor_group" "default" {
  monitor_group_name = var.name
}

resource "alibabacloudstack_cms_metric_rule_template" "default" {
  metric_rule_template_name = var.name
  description              = var.name
  enable                   = true
  notify_level             = 4
  silence_time             = 86400
  group_id                 = alibabacloudstack_cms_monitor_group.default.id

  alert_templates {
    category    = "ecs"
    metric_name = "cpu_total"
    namespace   = "acs_ecs_dashboard"
    rule_name   = "tf_testAcc_new"

    escalations {
      critical {
        comparison_operator = "GreaterThanThreshold"
        statistics          = "Average"
        threshold           = "90"
        times               = "3"
      }
    }
  }

  enable_start_time = "0"
  enable_end_time   = "23"
  webhook           = "https://www.example.com"
}
```

## 参数参考

支持以下参数：

* `alert_templates` - (可选) 描述时序指标报警模板规则的信息。
  * `category` - (必填) 服务名称的缩写。有效值包括：`ecs`, `rds`, `ads`, `slb`, `vpc`, `apigateway`, `cdn`, `cs`, `dcdn`, `ddos`, `eip`, `elasticsearch`, `emr`, `ess`, `hbase`, `iot_edge`, `kvstore_sharding`, `kvstore_splitrw`, `kvstore_standard`, `memcache`, `mns`, `mongodb`, `mongodb_cluster`, `mongodb_sharding`, `mq_topic`, `ocs`, `opensearch`, `oss`, `polardb`, `petadata`, `scdn`, `sharebandwidthpackages`, `sls`, `vpn`。
  * `escalations` - (可选) 基于警报级别触发条件的信息。
    * `critical` - (可选) 触发关键级警报的条件。
      * `comparison_operator` - (可选) 关键级警报的阈值比较运算符。有效值：`GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`。
      * `statistics` - (可选) 关键级警报的统计聚合方法。
      * `threshold` - (可选) 关键级警报的阈值。
      * `times` - (可选) 在触发关键级警报之前连续测量的次数。
    * `info` - (可选) 触发信息级警报的条件。
      * `comparison_operator` - (可选) 信息级警报的阈值比较运算符。有效值：`GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`。
      * `statistics` - (可选) 信息级警报的统计聚合方法。
      * `threshold` - (可选) 信息级警报的阈值。
      * `times` - (可选) 在触发信息级警报之前连续测量的次数。
    * `warn` - (可选) 触发警告级警报的条件。
      * `comparison_operator` - (可选) 警告级警报的阈值比较运算符。有效值：`GreaterThanOrEqualToThreshold`, `GreaterThanThreshold`, `LessThanOrEqualToThreshold`, `LessThanThreshold`, `NotEqualToThreshold`, `GreaterThanYesterday`, `LessThanYesterday`, `GreaterThanLastWeek`, `LessThanLastWeek`, `GreaterThanLastPeriod`, `LessThanLastPeriod`。
      * `statistics` - (可选) 警告级警报的统计聚合方法。
      * `threshold` - (可选) 警告级警报的阈值。
      * `times` - (可选) 在触发警告级警报之前连续测量的次数。
  * `metric_name` - (必填) 指标名称。更多信息，请参见 [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) 或 [附录1：指标](https://www.alibabacloud.com/help/doc-detail/28619.htm)。
  * `namespace` - (必填) 服务命名空间。更多信息，请参见 [DescribeMetricMetaList](https://www.alibabacloud.com/help/doc-detail/98846.htm) 或 [附录1：指标](https://www.alibabacloud.com/help/doc-detail/28619.htm)。
  * `rule_name` - (必填) 警报规则名称。
  * `webhook` - (可选) 当根据警报规则触发警报时，发送POST请求的回调URL。
* `apply_mode` - (可选) 模板应用方式。取值：
  * `GROUP_INSTANCE_FIRST`: 应用组实例优先。当应用报警模板时，以应用组实例优先。如果应用组中不存在该实例，则忽略模板中的规则。
  * `ALARM_TEMPLATE_FIRST`: 报警模板实例优先。当应用报警模板时，不管应用组中是否存在该实例，都创建报警规则。
* `description` - (可选) 警报模板的描述信息。
* `enable` - (可选) 是否将警报模板应用于应用组。有效值：`true` 或 `false`。默认值：`false`。
* `enable_end_time` - (可选) 警报生效的结束时间。取值范围：`00`~`23`，表示`00:59`到`23:59`。
* `enable_start_time` - (可选) 警报生效的开始时间。取值范围：`00`~`23`，表示`00:00`到`23:00`。
* `group_id` - (可选) 应用组ID。
* `metric_rule_template_name` - (必填，强制更新) 警报模板名称。
* `notify_level` - (可选) 报警通知模式。有效值：
  * `2`: 电话 + 短信 + 邮件 + 王旺 + 钉钉机器人。
  * `3`: 短信 + 邮件 + 王旺 + 钉钉机器人。
  * `4`: 王旺 + 钉钉机器人。
* `rest_version` - (可选) 警报模板的版本。
* `silence_time` - (可选) 通道静默周期。单位：秒。默认值：`86400`。当监控数据持续超过警报规则阈值时，每个静默周期内仅发送一次警报通知。
* `webhook` - (可选) 当发生警报时，指定的URL地址被回调，并发送POST请求。
* `overwrite` - (可选) 是否覆盖警报模板。有效值：`true` 或 `false`。默认值：`true`。

## 属性参考

除了上述所有参数外，还导出了以下属性：

* `id` - Terraform 中 Metric Rule Template 的资源 ID。
* `rest_version` - 警报模板的版本。