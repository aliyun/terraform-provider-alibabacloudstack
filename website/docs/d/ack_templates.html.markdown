---
subcategory: "Container Service for Kubernetes"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_ack_templates"
description: |-
  Provides a list of Ack Templates available to the user.
---

# alibabacloudstack_ack_templates

This data source provides a list of ACK Templates according to the specified filters.

## Example Usage

```hcl
# Declare the data source
data "alibabacloudstack_ack_templates" "example" {
  name_regex = "my-template"
}
```
## Argument Reference
The following arguments are supported:

* `ids` - (Optional) A list of Template IDs to filter results.
* `name_regex` - (Optional) A regex string to filter results by * template name.
* `description_regex` - (Optional) A regex string to filter results by template description.
* `template_type` - (Optional) Filter results by template type.
Attributes Reference
The following attributes are exported:

* `ids` - A list of Template IDs.
* `templates` - A list of templates. Each element contains the following attributes:
* `template` - The template content.
* `name` - The name of the template.
* `description` - The description of the template.
* `template_type` - The type of the template.
* `template_id` - The ID of the template.
* `template_with_hist_id` - The template ID with history.
* `template_hash_code_version` - The hash code version of the template.
* `created` - The creation time of the template.
* `acl` - The ACL of the template.
* `version` - The version of the template.
* `tags` - The tags of the template.
* `ali_uid` - The Alibaba Cloud UID.
* `updated` - The last update time of the template.