---
subcategory: "ASCM"
layout: "alibabacloudstack"
page_title: "AlibabacloudStack: alibabacloudstack_ascm_ram_policies_for_user"
sidebar_current: "docs-alibabacloudstack-datasource-ascm-ram-policies-for-user"
description: |-
    Provides a list of ram policy of the user.
---

# alibabacloudstack\_ascm_ram_policies_for_user

This data source provides the ram policy for user of the current Apsara Stack Cloud user.

## Example Usage

```
data "alibabacloudstack_ascm_ram_policies_for_user" "default" {
  login_name = "test_admin"
}
output "ramPolicy" {
  value = data.alibabacloudstack_ascm_ram_policies_for_user.default.*
}

```

## Argument Reference

The following arguments are supported:

* `ids` - (Optional) A list of ram policy IDs.
* `name_regex` - (Optional) A regex string to filter results by login name of the user.
* `output_file` - (Optional) File name where to save data source results (after running `terraform plan`).

## Attributes Reference

The following attributes are exported in addition to the arguments listed above:

* `policies` - A list of policies. Each element contains the following attributes:
  * `policy_name` - Policy name.
  * `description` - Description about the policy.
  * `attach_date` -  Creation Date of ram policy.
  * `policy_type` - Type of the policy.
  * `default_version` - Default version.
  * `policy_document` - Policy Document.
     
     
