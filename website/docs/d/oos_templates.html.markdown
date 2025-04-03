---
subcategory: "OOS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_oos_templates"
sidebar_current: "docs-Alibabacloudstack-datasource-oos-templates"
description: |- 
  Provides a list of oos templates owned by an Alibabacloudstack account.
---

# alibabacloudstack_oos_templates

This data source provides a list of OOS Templates in an Alibaba Cloud account according to the specified filters.

## Example Usage

```hcl
# Declare the data source

data "alibabacloudstack_oos_templates" "example" {
  name_regex = "test"
  tags = {
    "Created" = "TF"
    "For"     = "template Test"
  }
  share_type = "Private"
  has_trigger = false
}

output "first_template_name" {
  value = data.alibabacloudstack_oos_templates.example.templates.0.template_name
}
```

## Argument Reference

The following arguments are supported:

* `name_regex` - (Optional) A regex string to filter results by the `template_name`.
* `category` - (Optional) The category of the template.
* `created_by` - (Optional) The creator of the template.
* `created_date` - (Optional) The creation time of the template, less than or equal to the specified time. Format: `YYYY-MM-DDThh:mm:ssZ`.
* `created_date_after` - (Optional) The creation time of the template, greater than or equal to the specified time. Format: `YYYY-MM-DDThh:mm:ssZ`.
* `has_trigger` - (Optional) Whether the template has been triggered successfully.
* `share_type` - (Optional) The sharing type of the template. Valid values: `Private`, `Public`.
* `sort_field` - (Optional) The field used for sorting. Valid values: `TotalExecutionCount`, `Popularity`, `TemplateName`, and `CreatedDate`. Default: `TotalExecutionCount`.
* `sort_order` - (Optional) The order of sorting. Valid values: `Ascending`, `Descending`. Default: `Descending`.
* `template_format` - (Optional) The format of the template. Valid values: `JSON`, `YAML`.
* `template_type` - (Optional) The type of the OOS Template.
* `ids` - (Optional) A list of OOS Template IDs (`template_name`).
* `tags` - (Optional) A mapping of tags assigned to the resource.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `ids` - A list of OOS Template IDs. Each element in the list is the same as `template_name`.
* `names` - (Available in v1.114.0+) A list of OOS Template names.
* `templates` - A list of OOS Templates. Each element contains the following attributes:
  * `id` - The ID of the OOS Template. Same as `template_name`.
  * `template_name` - The name of the OOS Template.
  * `description` - The description of the OOS Template.
  * `template_id` - The ID of the OOS Template resource.
  * `template_version` - The version of the OOS Template.
  * `updated_by` - The user who last updated the template.
  * `updated_date` - The time when the template was last updated.
  * `category` - The category of the template.
  * `created_by` - The creator of the template.
  * `created_date` - The creation time of the template.
  * `has_trigger` - Whether the template has been triggered successfully.
  * `share_type` - The sharing type of the template. Valid values: `Private`, `Public`.
  * `tags` - A mapping of tags assigned to the resource.
  * `template_format` - The format of the template. Valid values: `JSON`, `YAML`.
  * `template_type` - The type of the OOS Template.