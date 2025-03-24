---
subcategory: "EDAS"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_edas_deploygroups"
sidebar_current: "docs-Alibabacloudstack-datasource-edas-deploygroups"
description: |- 
  Provides a list of edas deploygroups owned by an alibabacloudstack account.
---

# alibabacloudstack_edas_deploygroups
-> **NOTE:** Alias name has: `alibabacloudstack_edas_deploy_groups`

This data source provides a list of edas deploygroups in an alibabacloudstack account according to the specified filters.

## Example Usage

```hcl
data "alibabacloudstack_edas_deploygroups" "example" {
  app_id   = "your_app_id"
  name_regex = "group-.*"

  output_file = "deploygroups.txt"
}

output "first_group_name" {
  value = data.alibabacloudstack_edas_deploygroups.example.groups.0.group_name
}
```

## Argument Reference

The following arguments are supported:

* `app_id` - (Required, ForceNew) The ID of the EDAS application for which you want to retrieve the deploy groups.
* `name_regex` - (Optional, ForceNew) A regex string to filter results by the deploy group name.

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `ids` - A list of deploy group IDs.
* `names` - A list of deploy group names.
* `groups` - A list of deploy groups. Each element contains the following attributes:
  * `group_id` - The ID of the instance group.
  * `group_name` - The name of the instance group. The length cannot exceed 64 characters.
  * `group_type` - The type of the instance group. Valid values:
    - `0`: Default Grouping.
    - `1`: Grayscale is not enabled for traffic management.
    - `2`: Traffic Management Enable Grayscale.
  * `create_time` - The timestamp of the creation time.
  * `update_time` - The timestamp of the update time.
  * `app_id` - The ID of the application that you want to deploy.
  * `cluster_id` - The ID of the cluster that you want to create the application.
  * `package_version_id` - The version ID of the group deployment package.
  * `app_version_id` - The version ID of the application deployment record.