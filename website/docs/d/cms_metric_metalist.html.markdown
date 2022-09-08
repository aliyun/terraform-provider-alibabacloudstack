---
subcategory: "Cloud Monitor"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_metric_metalist"
sidebar_current: "docs-alibabacloudstack-datasource-cms-metric-metalist"
description: |-
    Provides a Metalist owned by an Alibabacloudstack Cloud account.
---

# alibabacloudstack\_cms\_project\_metalist

Provides a Metalist of project  owned by an Alibabacloudstack Cloud account.

## Example Usage

Basic Usage

```
data "alibabacloudstack_cms_metric_metalist" "default" {
  namespace="acs_slb_dashboard"
}

output "metric_metalist" {
  value = data.alibabacloudstack_cms_metric_metalist.default
}
```

## Argument Reference

The following arguments are supported:

* `namespace` - (Required, ForceNew) The namespace of the service. You can call the  operation to obtain namespaces. 

## Attributes Reference

The following attributes are exported:

* `resources` - A list of cms metriclist. Each element contains the following attributes:
    * `metric_name` - The name of the metric.
    * `periods` -     The statistical period of the metric.
    * `description` - The description of the metric. 
    * `dimensions` - The dimensions of the metric. Multiple dimensions are separated with commas.
    * `labels` - The tags of the metric. The value is a JSON array string. The array can include repeated tag names. Sample value: [{"name":"Tag name","value":"Tag value"}] . 
                 The available tag names are as follows: metricCategory: 
          the category of the metrics. alertEnable: indicates whether the alert is enabled. alertUnit: the unit of the metric in the alert. unitFactor: the factor for metric unit conversion. minAlertPeriod: the minimum time interval to raise a new alert. productCategory: the category of the service.
    * `unit` - The unit of the metric. 
    * `statistics` - The statistical method of the metric. Multiple statistical methods are separated with commas (,), for example, Average,Minimum,Maximum.
    * `namespace` - The namespace of the monitored service.

  