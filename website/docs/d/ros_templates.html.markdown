---
subcategory: "ROS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ros_templates"
sidebar_current: "docs-Alibabacloudstack-datasource-ros-templates"
description: |- 
  Provides a list of ros templates owned by an Alibabacloudstack account.
---

# alibabacloudstack_ros_templates

This data source provides a list of ROS Templates in an Alibaba Cloud Stack account according to the specified filters.

## Example Usage

Basic Usage

```terraform
data "alibabacloudstack_ros_templates" "example" {
  ids        = ["example_value"]
  name_regex = "the_resource_name"
}

output "first_ros_template_id" {
  value = data.alibabacloudstack_ros_templates.example.templates.0.id
}
```

Advanced Usage with Tags

```terraform
data "alibabacloudstack_ros_templates" "example_with_tags" {
  template_name = "example_template"
  tags = jsonencode({
    Environment = "Production"
    Owner      = "TeamA"
  })
}

output "ros_template_ids_with_tags" {
  value = data.alibabacloudstack_ros_templates.example_with_tags.ids
}
```

## Argument Reference

The following arguments are supported:

* `share_type` - (Optional, ForceNew) The share type of the ROS Template. Valid values: `Private`, `Shared`.
* `ids` - (Optional, ForceNew) A list of Template IDs.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by Template name.
* `template_name` - (Optional, ForceNew) The name of the ROS Template. The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
* `enable_details` - (Optional) Default to `false`. Set it to `true` to output more details about resource attributes.
* `tags` - (Optional) Query the resource bound to the tag. The format of the incoming value is `json` string, including `TagKey` and `TagValue`. `TagKey` cannot be null, and `TagValue` can be empty. Format example: `{"key1":"value1"}`.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `names` - A list of Template names.
* `templates` - A list of ROS Templates. Each element contains the following attributes:
  * `change_set_id` - The ID of the change set associated with the template.
  * `description` - The description of the template. The description can be up to 256 characters in length.
  * `id` - The ID of the Template.
  * `share_type` - The share type of the template (`Private` or `Shared`).
  * `stack_group_name` - The name of the stack group associated with the template. The name must be unique in a region and can be up to 255 characters in length.
  * `stack_id` - The ID of the stack associated with the template.
  * `tags` - Tags associated with the template.
    * `tag_key` - The key of tag N of the resource.
    * `tag_value` - The value of tag N of the resource.
  * `template_body` - The structure that contains the template body. The template body must be 1 to 524,288 bytes in length. If the length of the template body is longer than required, we recommend adding parameters to the HTTP POST request body to avoid request failures due to excessive length of URLs. You must specify one of the `TemplateBody` and `TemplateURL` parameters, but you cannot specify both of them.
  * `template_id` - The ID of the template.
  * `template_name` - The name of the template. The name can be up to 255 characters in length and can contain digits, letters, hyphens (-), and underscores (_). It must start with a digit or letter.
  * `template_version` - The version of the template.
