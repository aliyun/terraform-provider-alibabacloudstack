---
subcategory: "CloudMonitorService"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cloudmonitorservice_metricruletemplates"
sidebar_current: "docs-Alibabacloudstack-datasource-cloudmonitorservice-metricruletemplates"
description: |- 
  Provides a list of cloudmonitorservice metricruletemplates owned by an alibabacloudstack account.
---

# alibabacloudstack_cloudmonitorservice_metricruletemplates
-> **NOTE:** Alias name has: `alibabacloudstack_cms_metric_rule_templates`

This data source provides a list of cloudmonitorservice metricruletemplates in an alibabacloudstack account according to the specified filters.

## Example Usage

Basic Usage

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
  keyword = "^my-MetricRuleTemplate"
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

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew) A list of Metric Rule Template IDs. The attribute is used to match against the IDs of the desired templates.
* `keyword` - (Optional, ForceNew) The keyword for filtering alert templates. You can perform fuzzy search based on the template name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Metric Rule Template name.
* `template_id` - (Optional, ForceNew) The ID of the specific alert template you want to retrieve.
* `is_default` - (Optional) Specifies whether to filter for default Metric Rule Templates. Set to `true` to include only default templates.
* `history` - (Optional) Whether to display the history of alarm templates applied to application groups. Valid values:
  * `true`: Display the history.
  * `false` (default): Do not display the history.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `names` - A list of Metric Rule Template names.
* `templates` - A list of Cms Metric Rule Templates. Each element contains the following attributes:
  * `description` - The description of the alert template.
  * `id` - The ID of the Metric Rule Template.
  * `name` - Name of the Metric Rule Template.
  * `rest_version` - The version of the alert template.