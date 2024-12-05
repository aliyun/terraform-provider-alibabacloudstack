---
subcategory: "Cloud Monitor Service"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cms_metric_rule_templates"
sidebar_current: "docs-alibabacloudstack-datasource-cms-metric-rule-templates"
description: |- 
    Provides a list of Cms Metric Rule Templates to the user.
---

# alibabacloudstack\_cms\_metric\_rule\_templates

This data source provides the Cms Metric Rule Templates of the current Alibaba Cloud user.


## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_cms_metric_rule_templates" "ids" {
  ids = ["example_value"]
}
output "cms_metric_rule_template_id_1" {
  value = data.alibabacloudstack_cms_metric_rule_templates.ids.templates.0.id
}

data "alibabacloudstack_cms_metric_rule_templates" "nameRegex" {
  name_regex = "^my-MetricRuleTemplate"
}
output "cms_metric_rule_template_id_2" {
  value = data.alibabacloudstack_cms_metric_rule_templates.nameRegex.templates.0.id
}

data "alibabacloudstack_cms_metric_rule_templates" "keyword" {
  keyword = "^my-MetricRuleTemplate"
}
output "cms_metric_rule_template_id_3" {
  value = data.alibabacloudstack_cms_metric_rule_templates.nameRegex.templates.0.id
}

data "alibabacloudstack_cms_metric_rule_templates" "template_id" {
  template_id = "example_value"
}
output "cms_metric_rule_template_id_4" {
  value = data.alibabacloudstack_cms_metric_rule_templates.nameRegex.templates.0.id
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional, ForceNew, Computed)  A list of Metric Rule Template IDs.
* `keyword` - (Optional, ForceNew) The name of the alert template. You can perform fuzzy search based on the template name.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Metric Rule Template name.
* `template_id` - (Optional, ForceNew) The ID of the alert template.
* `is_default` - (Optional) Is the default Metric Rule Template.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).
* `history` - (Optional) Whether to display the history of alarm templates applied to application groups. Value:
  * True: Display. 
  * False (default): Not displayed

## Argument Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Metric Rule Template names.
* `templates` - A list of Cms Metric Rule Templates. Each element contains the following attributes:
  * `description` - The description of the alert template.
  * `id` - The ID of the Metric Rule Template.
  * `name` - Name of the Metric Rule Template.
  * `rest_version` - The version of the alert template.