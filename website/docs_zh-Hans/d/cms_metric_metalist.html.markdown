---
subcategory: "Cloud Monitor"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_metric_metalist"
sidebar_current: "docs-alibabacloudstack-datasource-cms-metric-metalist"
description: |-
    查询云监控指标
---

# alibabacloudstack_cms_metric_metalist

根据指定过滤条件列出当前凭证权限可以访问的云监控指标列表。

## 示例用法

### 基础用法

```
data "alibabacloudstack_cms_metric_metalist" "default" {
  namespace="acs_slb_dashboard"
}

output "metric_metalist" {
  value = data.alibabacloudstack_cms_metric_metalist.default
}
```

## 参数说明

以下参数被支持：

* `namespace` - (必填，变更时重建) 服务的命名空间。您可以调用操作以获取命名空间。
* `resources` - (可选) cms metriclist 列表。此参数用于进一步筛选符合条件的指标资源。

## 属性参考

以下属性被导出：

* `resources` - 云监控指标列表。每个元素包含以下属性：
    * `metric_name` - 指标的名称。
    * `periods` - 指标的统计周期（单位为秒）。例如，60 表示 1 分钟，300 表示 5 分钟。
    * `description` - 指标的描述信息，通常用于解释该指标的具体含义。
    * `dimensions` - 指标的维度。多个维度用逗号分隔。维度是用于区分不同实例或资源的关键字，例如 `instanceId` 或 `deviceName`。
    * `labels` - 指标的标签。值为 JSON 数组字符串。数组可以包括重复的标签名。示例值：`[{"name":"Tag name","value":"Tag value"}]`。可用的标签名如下：
        * `metricCategory`：指标类别。
        * `alertEnable`：指示是否启用告警。
        * `alertUnit`：告警中的指标单位。
        * `unitFactor`：指标单位转换因子。
        * `minAlertPeriod`：触发新告警的最短时间间隔。
        * `productCategory`：服务类别。
    * `unit` - 指标的单位，例如 `Bytes`, `Percent`, `Count` 等。
    * `statistics` - 指标的统计方法。多个统计方法用逗号(,)分隔，例如，`Average,Minimum,Maximum`。常见的统计方法包括：
        * `Average`：平均值。
        * `Minimum`：最小值。
        * `Maximum`：最大值。
        * `Sum`：总和。
        * `SampleCount`：样本数量。
    * `namespace` - 监控服务的命名空间，用于标识该指标所属的服务或产品。