---
subcategory: "Bastion Host"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_bastionhost_instances"
sidebar_current: "docs-alibabacloudstack-resource-bastionhost_instances"
description: |-
  Provides a list of Bastionhost instances in Alibaba Cloud Stack based on the provided filters.
---


## Example Usage
The following example retrieves Bastionhost instances by description regex and writes the output to a file:
```hcl
data "alibabacloudstack_bastionhost_instances" "example" {
  description_regex = "^example"
  output_file       = "output.json"
}

output "instances" {
  value = data.alibabacloudstack_bastionhost_instances.example.ids
}
```
## Argument Reference
The following arguments are supported:

* `description_regex` - (Optional) A regular expression to filter Bastionhost instances by their description.
* `output_file` - (Optional) File name to save the results to in JSON format.
* `Deprecated`: This field is deprecated and will be removed in version 3.19.0. Use the local_file provider instead.
* `ids` - (Optional, ForceNew) A list of Bastionhost instance IDs to filter the results.
* `tags` - (Optional) A map of tags to filter Bastionhost instances by.
## Attributes Reference
In addition to all the above arguments, the following attributes are exported:

* `ids` - A list of Bastionhost instance IDs matching the criteria.
* `descriptions` - A list of descriptions corresponding to the matched instances.
* `instances` - A list of Bastionhost instances with their detailed attributes. Each entry contains:
* `id` - The ID of the Bastionhost instance.
* `description` - Description of the instance.
* `user_vswitch_id` - VSwitch ID associated with the instance.
* `private_domain` - Private domain or intranet endpoint of the instance.
* `public_domain` - Public domain or internet endpoint of the instance (if available).
* `instance_status` - Current status of the instance.
* `license_code` - License code for the instance.
* `public_network_access` - Boolean indicating whether public network access is enabled.
* `security_group_ids` - List of security group IDs associated with the instance.
* `tags` - Tags assigned to the instance.


