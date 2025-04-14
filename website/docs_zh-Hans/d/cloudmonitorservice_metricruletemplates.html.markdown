---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricruletemplates"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-metricruletemplates"
description: |- 
  查询云监控告警模板
---

# alibabacloudstack_cloudmonitorservice_metricruletemplates
-> **NOTE:** 该资源等效别名有: `alibabacloudstack_cms_metric_rule_templates`

根据指定过滤条件列出当前凭证权限可以访问的云监控告警模板列表。

## 示例用法

### 基础用法：

```terraform
data "alibabacloudstack_cloudmonitorservice_metricruletemplates" "ids" {
  ids = ["example_value"]
}
output "cloudmonitorservice_metricruletemplate_id_1" {
  value = data.alibabacloudstack_cloudmonitorservice_metricruletemplates.ids.templates.0.id
}

data "alibabacloudstack_cloudmonitorservice_metricruletemplates" "nameRegex" {
  name_regex = "^my-MetricRuleTemplate"
}
output "cloudmonitorservice_metricruletemplate_id_2" {
  value = data.alibabacloudstack_cloudmonitorservice_metricruletemplates.nameRegex.templates.0.id
}

data "alibabacloudstack_cloudmonitorservice_metricruletemplates" "keyword" {
  keyword = "my-MetricRuleTemplate"
}
output "cloudmonitorservice_metricruletemplate_id_3" {
  value = data.alibabacloudstack_cloudmonitorservice_metricruletemplates.keyword.templates.0.id
}

data "alibabacloudstack_cloudmonitorservice_metricruletemplates" "template_id" {
  template_id = "example_value"
}
output "cloudmonitorservice_metricruletemplate_id_4" {
  value = data.alibabacloudstack_cloudmonitorservice_metricruletemplates.template_id.templates.0.id
}
```

## 参数说明

以下参数是支持的：

* `ids` - (可选，变更时重建) 告警模板ID列表。该属性用于匹配所需的告警模板ID。
* `keyword` - (可选，变更时重建) 用于过滤告警模板的关键字。可以根据模板名称进行模糊搜索。
* `name_regex` - (可选，变更时重建) 用于通过告警模板名称过滤结果的正则表达式字符串。
* `template_id` - (可选，变更时重建) 要检索的具体告警模板的ID。
* `is_default` - (可选) 是否过滤默认的告警模板。设置为`true`以仅包括默认模板。
* `history` - (可选) 是否显示应用到应用组的告警模板的历史记录。有效值：
  * `true`: 显示历史记录。
  * `false`(默认): 不显示历史记录。

## 属性说明

除了上述参数外，还导出以下属性：

* `names` - 告警模板名称列表。
* `templates` - Cms Metric Rule Templates 列表。每个元素包含以下属性：
  * `description` - 告警模板的描述信息。
  * `id` - 告警模板的ID。
  * `name` - 告警模板的名称。
  * `rest_version` - 告警模板的版本号。
