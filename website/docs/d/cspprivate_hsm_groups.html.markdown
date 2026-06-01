---
subcategory: "Cspprivate HSM"
layout: "alibabacloudstack"
page_title: "Alibabacloudstack: alibabacloudstack_cspprivate_hsm_groups"
sidebar_current: "docs-Alibabacloudstack-datasource-cspprivate-hsm-groups"
description: |-
  Provides a list of CSP Private HSM Groups.
---

# alibabacloudstack_cspprivate_hsm_groups

This data source provides the CSP Private HSM Groups available in ApsaraStack.

-> **NOTE:** Available in ApsaraStack.

## Example Usage

```hcl
data "alibabacloudstack_cspprivate_hsm_groups" "example" {
  ids = ["my-hsm-group"]
}

output "hsm_groups" {
  value = data.alibabacloudstack_cspprivate_hsm_groups.example.groups
}
```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of HSM group names to filter results by group name.

## Attributes Reference

The following attributes are exported:

* `ids` - A list of HSM group IDs.
* `groups` - A list of CSP Private HSM Groups. Each element contains the following attributes:
  * `id` - The ID of the HSM group (equivalent to group name).
  * `group_name` - The name of the HSM group.
  * `status` - The current status of the HSM group.
  * `unique_id` - The unique ID identifier of the HSM group.
  * `security_level_tag` - The security level tag of the HSM group.
  * `create_time` - The creation time of the HSM group, in ISO 8601 format.
  * `hsm_count` - The number of HSMs in the HSM group.
  * `vpc_id` - The ID of the VPC to which the HSM group belongs.
  * `update_time` - The last update time of the HSM group, in ISO 8601 format.
  * `zone_ids` - A list of zone IDs where the HSM group is located.